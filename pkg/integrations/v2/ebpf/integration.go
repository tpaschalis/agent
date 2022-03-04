package ebpf

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"time"

	ebpf_config "github.com/cloudflare/ebpf_exporter/config"
	"github.com/cloudflare/ebpf_exporter/exporter"

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

	config := ebpf_config.Config{}
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
