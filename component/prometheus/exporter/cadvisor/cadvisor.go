package cadvisor

import (
	"time"

	"github.com/grafana/agent/component"
	"github.com/grafana/agent/component/prometheus/exporter"
	"github.com/grafana/agent/pkg/integrations"
	"github.com/grafana/agent/pkg/integrations/cadvisor"
)

func init() {
	component.Register(component.Registration{
		Name:          "prometheus.exporter.cadvisor",
		Args:          Arguments{},
		Exports:       exporter.Exports{},
		NeedsServices: exporter.RequiredServices(),
		Build:         exporter.New(createExporter, "cadvisor"),
	})
}

func createExporter(opts component.Options, args component.Arguments) (integrations.Integration, error) {
	a := args.(Arguments)
	return a.Convert().NewIntegration(opts.Logger)
}

// func customizeTarget(baseTarget discovery.Target, args component.Arguments) []discovery.Target {
// 	a := args.(Arguments)
// 	target := baseTarget
// 	url, err := url.Parse(a.APIURL)
// 	if err != nil {
// 		return []discovery.Target{target}
// 	}
// 	target["instance"] = url.Host
// 	return []discovery.Target{target}
// }

// DefaultArguments holds non-zero default options for Arguments when it is
// unmarshaled from river.
var DefaultArguments = Arguments{
	StoreContainerLabels: true,
	ResctrlInterval:      0,
	StorageDuration:      2 * time.Minute,

	ContainerdHost:      "/run/containerd/containerd.sock",
	ContainerdNamespace: "k8s.io",

	// TODO(@tpaschalis) Do we need the default cert/key/ca since tls is disabled by default?
	DockerHost:    "unix:///var/run/docker.sock",
	UseDockerTLS:  false,
	DockerTLSCert: "cert.pem",
	DockerTLSKey:  "key.pem",
	DockerTLSCA:   "ca.pem",

	DockerOnly: false,
}

// Arguments configures the prometheus.exporter.cadvisor component.
type Arguments struct {
	StoreContainerLabels       bool          `river:"store_container_labels,attr,optional"`
	AllowlistedContainerLabels []string      `river:"allowlisted_container_labels,attr,optional"`
	EnvMetadataAllowlist       []string      `river:"env_metadata_allowlist,attr,optional"`
	RawCgroupPrefixAllowlist   []string      `river:"raw_cgroup_prefix_allowlist,attr,optional"`
	PerfEventsConfig           string        `river:"perf_events_config,attr,optional"`
	ResctrlInterval            time.Duration `river:"resctrl_interval,attr,optional"`
	DisabledMetrics            []string      `river:"disabled_metrics,attr,optional"`
	EnabledMetrics             []string      `river:"enabled_metrics,attr,optional"`
	StorageDuration            time.Duration `river:"storage_duration,attr,optional"`
	ContainerdHost             string        `river:"containerd_host,attr,optional"`
	ContainerdNamespace        string        `river:"containerd_namespace,attr,optional"`
	DockerHost                 string        `river:"docker_host,attr,optional"`
	UseDockerTLS               bool          `river:"use_docker_tls,attr,optional"`
	DockerTLSCert              string        `river:"docker_tls_cert,attr,optional"`
	DockerTLSKey               string        `river:"docker_tls_key,attr,optional"`
	DockerTLSCA                string        `river:"docker_tls_ca,attr,optional"`
	DockerOnly                 bool          `river:"docker_only,attr,optional"`
}

// SetToDefault implements river.Defaulter.
func (a *Arguments) SetToDefault() {
	*a = DefaultArguments
}

// Convert returns the upstream-compatible configuration struct.
func (a *Arguments) Convert() *cadvisor.Config {
	return &cadvisor.Config{
		StoreContainerLabels:       a.StoreContainerLabels,
		AllowlistedContainerLabels: a.AllowlistedContainerLabels,
		EnvMetadataAllowlist:       a.EnvMetadataAllowlist,
		RawCgroupPrefixAllowlist:   a.RawCgroupPrefixAllowlist,
		PerfEventsConfig:           a.PerfEventsConfig,
		ResctrlInterval:            int64(a.ResctrlInterval), // TODO(@tpaschalis) This is so that the cadvisor package can re-cast back to time.Duration. Can we make it use time.Duration directly instead?
		DisabledMetrics:            a.DisabledMetrics,
		EnabledMetrics:             a.EnabledMetrics,
		StorageDuration:            a.StorageDuration,
		Containerd:                 a.ContainerdHost,
		ContainerdNamespace:        a.ContainerdNamespace,
		Docker:                     a.DockerHost,
		DockerTLS:                  a.UseDockerTLS,
		DockerTLSCert:              a.DockerTLSCert,
		DockerTLSKey:               a.DockerTLSKey,
		DockerTLSCA:                a.DockerTLSCA,
		DockerOnly:                 a.DockerOnly,
	}
}
