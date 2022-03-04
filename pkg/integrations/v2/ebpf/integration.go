package ebpf

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"
	"time"

	ebpf_config "github.com/cloudflare/ebpf_exporter/config"
	"github.com/cloudflare/ebpf_exporter/exporter"
	"gopkg.in/yaml.v2"

	"github.com/go-kit/log"
	"github.com/gorilla/mux"
	"github.com/grafana/agent/pkg/integrations/v2"
	"github.com/grafana/agent/pkg/integrations/v2/autoscrape"
	"github.com/grafana/agent/pkg/integrations/v2/common"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/model"
	"github.com/prometheus/prometheus/config"
	"github.com/prometheus/prometheus/discovery"
	"github.com/prometheus/prometheus/discovery/targetgroup"
)

type Config struct {
	common  common.MetricsConfig
	globals integrations.Globals
}

// var DefaultConfig = Config{}
// func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	*c = DefaultConfig
// 	type plain Config
// 	return unmarshal((*plain)(c))
// }

func (c *Config) ApplyDefaults(globals integrations.Globals) error {
	c.common.ApplyDefaults(globals.SubsystemOpts.Metrics.Autoscrape)
	return nil
}

func (c *Config) Identifier(globals integrations.Globals) (string, error) {
	return "ebpf", nil
}

func (c *Config) Name() string { return "ebpf" }

func (c *Config) NewIntegration(l log.Logger, globals integrations.Globals) (integrations.Integration, error) {
	c.globals = globals

	ebpf := &ebpfHandler{}
	ebpf.cfg = c

	return ebpf, nil
}

func init() {
	integrations.Register(&Config{}, integrations.TypeSingleton)
}

type ebpfHandler struct {
	cfg *Config
}

// RunIntegration implements the Integration interface.
func (e *ebpfHandler) RunIntegration(ctx context.Context) error {
	fmt.Println("Running epbf handler!")
	time.Sleep(5 * time.Minute)
	<-ctx.Done()
	fmt.Println("Exiting from ebpf handler...")
	return nil
}

// Handler implements the HTTPIntegration interface.
func (e *ebpfHandler) Handler(prefix string) (http.Handler, error) {
	r := mux.NewRouter()
	r.Handle(path.Join(prefix, "metrics"), createHandler())
	return r, nil
}

// TODO this might be better integrated INTO ebpfHandler, OR as multiple handlers, one per each 'program'.
// We'll see how it goes anyway.
func createHandler() http.HandlerFunc {

	bioyaml := `programs:
  # See:
  # * https://github.com/iovisor/bcc/blob/master/tools/biolatency.py
  # * https://github.com/iovisor/bcc/blob/master/tools/biolatency_example.txt
  #
  # See also: bio-tracepoints.yaml
  - name: bio
    metrics:
      histograms:
        - name: bio_latency_seconds
          help: Block IO latency histogram
          table: io_latency
          bucket_type: exp2
          bucket_min: 0
          bucket_max: 26
          bucket_multiplier: 0.000001 # microseconds to seconds
          labels:
            - name: device
              size: 32
              decoders:
                - name: string
            - name: operation
              size: 8
              decoders:
                - name: uint
                - name: static_map
                  static_map:
                    1: read
                    2: write
            - name: bucket
              size: 8
              decoders:
                - name: uint
        - name: bio_size_bytes
          help: Block IO size histogram with kibibyte buckets
          table: io_size
          bucket_type: exp2
          bucket_min: 0
          bucket_max: 15
          bucket_multiplier: 1024 # kibibytes to bytes
          labels:
            - name: device
              size: 32
              decoders:
                - name: string
            - name: operation
              size: 8
              decoders:
                - name: uint
                - name: static_map
                  static_map:
                    1: read
                    2: write
            - name: bucket
              size: 8
              decoders:
                - name: uint
    kprobes:
      # Remove blk_start_request if you're running Linux 5.3+, or better yet
      # use tracepoint based code that depends on stable kernel ABI.
	  # blk_start_request: trace_req_start
      blk_mq_start_request: trace_req_start
	  # blk_account_io_completion: trace_req_completion
    code: |
      #include <linux/blkdev.h>
      #include <linux/blk_types.h>

      typedef struct disk_key {
          char disk[32];
          u8 op;
          u64 slot;
      } disk_key_t;

      // Max number of disks we expect to see on the host
      const u8 max_disks = 255;

      // 27 buckets for latency, max range is 33.6s .. 67.1s
      const u8 max_latency_slot = 26;

      // 16 buckets per disk in kib, max range is 16mib .. 32mib
      const u8 max_size_slot = 15;

      // Hash to temporily hold the start time of each bio request, max 10k in-flight by default
      BPF_HASH(start, struct request *);

      // Histograms to record latencies
      BPF_HISTOGRAM(io_latency, disk_key_t, (max_latency_slot + 2) * max_disks);

      // Histograms to record sizes
      BPF_HISTOGRAM(io_size, disk_key_t, (max_size_slot + 2) * max_disks);

      // Record start time of a request
      int trace_req_start(struct pt_regs *ctx, struct request *req) {
          u64 ts = bpf_ktime_get_ns();
          start.update(&req, &ts);

          return 0;
      }

      // Calculate request duration and store in appropriate histogram bucket
      int trace_req_completion(struct pt_regs *ctx, struct request *req, unsigned int bytes) {
          u64 *tsp, delta;

          // Fetch timestamp and calculate delta
          tsp = start.lookup(&req);
          if (tsp == 0) {
              return 0; // missed issue
          }

          // There are write request with zero length on sector zero,
          // which do not seem to be real writes to device.
          if (req->__sector == 0 && req->__data_len == 0) {
            return 0;
          }

          // Disk that received the request
          struct gendisk *disk = req->rq_disk;

          // Delta in nanoseconds
          delta = bpf_ktime_get_ns() - *tsp;

          // Convert to microseconds
          delta /= 1000;

          // Latency histogram key
          u64 latency_slot = bpf_log2l(delta);

          // Cap latency bucket at max value
          if (latency_slot > max_latency_slot) {
              latency_slot = max_latency_slot;
          }

          disk_key_t latency_key = { .slot = latency_slot };
          bpf_probe_read(&latency_key.disk, sizeof(latency_key.disk), &disk->disk_name);

          // Size in kibibytes
          u64 size_kib = bytes / 1024;

          // Request size histogram key
          u64 size_slot = bpf_log2(size_kib);

          // Cap latency bucket at max value
          if (size_slot > max_size_slot) {
              size_slot = max_size_slot;
          }

          disk_key_t size_key = { .slot = size_slot };
          bpf_probe_read(&size_key.disk, sizeof(size_key.disk), &disk->disk_name);

          if ((req->cmd_flags & REQ_OP_MASK) == REQ_OP_WRITE) {
              latency_key.op = 2;
              size_key.op    = 2;
          } else {
              latency_key.op = 1;
              size_key.op    = 1;
          }

          io_latency.increment(latency_key);
          io_size.increment(size_key);

          // Increment sum keys
          latency_key.slot = max_latency_slot + 1;
          io_latency.increment(latency_key, delta);
          size_key.slot = max_size_slot + 1;
          io_size.increment(size_key, size_kib);

          start.delete(&req);

          return 0;
      }`
	config := ebpf_config.Config{}
	dec := yaml.NewDecoder(strings.NewReader(bioyaml))
	dec.Decode(&config)
	e, err := exporter.New(config)
	if err != nil {
		panic(err)
	}

	err = e.Attach()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Starting with %d programs found in the config\n", len(config.Programs))

	// err = prometheus.Register(version.NewCollector("ebpf_exporter"))
	// if err != nil {
	// 	log.Fatalf("Error registering version collector: %s", err)
	// }

	// err = prometheus.Register(e)
	// if err != nil {
	// 	log.Fatalf("Error registering exporter: %s", err)
	// }

	return func(w http.ResponseWriter, r *http.Request) {

		registry := prometheus.NewRegistry()

		registry.MustRegister(e)
		h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)

		ln := []byte(`an_important_metric_total{method="GET",code="200"}  3`)
		w.Write(ln)
		return
	}
}

// Targets implements the MetricsIntegration interface.
func (e *ebpfHandler) Targets(ep integrations.Endpoint) []*targetgroup.Group {
	integrationNameValue := model.LabelValue("integrations/ebpf")

	group := &targetgroup.Group{
		Labels: model.LabelSet{
			model.InstanceLabel: model.LabelValue("my_m1_macbook"),
			model.JobLabel:      integrationNameValue,
			"agent_hostname":    model.LabelValue(e.cfg.globals.AgentIdentifier),

			// Meta labels that can be used during SD.
			"__meta_agent_integration_name":       model.LabelValue("ebpf"),
			"__meta_agent_integration_instance":   model.LabelValue("my_m1_macbook"),
			"__meta_agent_integration_autoscrape": model.LabelValue("1"),
		},
		Source: fmt.Sprintf("%s/%s", "ebpf", "my_m1_macbook"),
	}

	for _, lbl := range e.cfg.common.ExtraLabels {
		group.Labels[model.LabelName(lbl.Name)] = model.LabelValue(lbl.Value)
	}

	group.Targets = append(group.Targets, model.LabelSet{
		model.AddressLabel:     model.LabelValue(ep.Host),
		model.MetricsPathLabel: model.LabelValue(path.Join(ep.Prefix, "/metrics")),
	})

	return []*targetgroup.Group{group}
}

// ScrapeConfigs implements the MetricsIntegration interface.
func (e *ebpfHandler) ScrapeConfigs(sd discovery.Configs) []*autoscrape.ScrapeConfig {
	name := "ebpf"

	cfg := config.DefaultScrapeConfig
	cfg.JobName = fmt.Sprintf("%s", name)
	cfg.Scheme = e.cfg.globals.AgentBaseURL.Scheme
	cfg.HTTPClientConfig = e.cfg.globals.SubsystemOpts.ClientConfig
	cfg.ServiceDiscoveryConfigs = sd
	cfg.ScrapeInterval = e.cfg.common.Autoscrape.ScrapeInterval
	cfg.ScrapeTimeout = e.cfg.common.Autoscrape.ScrapeTimeout
	cfg.RelabelConfigs = e.cfg.common.Autoscrape.RelabelConfigs
	cfg.MetricRelabelConfigs = e.cfg.common.Autoscrape.MetricRelabelConfigs

	return []*autoscrape.ScrapeConfig{{
		Instance: e.cfg.common.Autoscrape.MetricsInstance,
		Config:   cfg,
	}}
}
