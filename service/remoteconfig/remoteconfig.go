package remoteconfig

import (
	"context"
	"fmt"
	"hash"
	"hash/fnv"
	"math"
	"os"
	"path/filepath"
	"time"

	"github.com/go-kit/log"
	"github.com/grafana/agent/api/gen/proto/go/agent/v1/agentv1connect"
	"github.com/grafana/agent/component"
	"github.com/grafana/agent/component/common/config"
	"github.com/grafana/agent/service"
	commonconfig "github.com/prometheus/common/config"
)

var fnvHash hash.Hash32 = fnv.New32()

func getHash(in string) string {
	fnvHash.Write([]byte(in))
	defer fnvHash.Reset()

	return fmt.Sprintf("%x", fnvHash.Sum(nil))
}

// Service implements a service for remote configuration.
type Service struct {
	mod component.Module

	asClient agentv1connect.AgentServiceClient

	opts Options

	dataPath string
	ticker   *time.Ticker
}

// ServiceName defines the name used for the remote config service.
const ServiceName = "remote_configuration"

// Options are used to configure the remote config service. Options are
// constant for the lifetime of the remote config service.
type Options struct {
	Logger      log.Logger // Where to send logs.
	StoragePath string     // Where to cache configuration on-disk.
}

// Arguments holds runtime settings for the remote config service.
type Arguments struct {
	URL              string                   `river:"url,attr,optional"`
	PollFrequency    time.Duration            `river:"poll_frequency,attr,optional"`
	HTTPClientConfig *config.HTTPClientConfig `river:",squash"`
}

// GetDefaultArguments populates the default values for the Arguments struct.
func GetDefaultArguments() Arguments {
	return Arguments{
		PollFrequency:    1 * time.Minute,
		HTTPClientConfig: config.CloneDefaultHTTPClientConfig(),
	}
}

// SetToDefault implements river.Defaulter.
func (a *Arguments) SetToDefault() {
	*a = GetDefaultArguments()
}

// Validate implements river.Validator.
func (a *Arguments) Validate() error {
	// We must explicitly Validate because HTTPClientConfig is squashed and it
	// won't run otherwise
	if a.HTTPClientConfig != nil {
		return a.HTTPClientConfig.Validate()
	}

	return nil
}

// Data includes information associated with the remote configuration service.
type Data struct {
}

// New returns a new instance of the remote config service.
func New(opts Options) (*Service, error) {
	basePath := filepath.Join(opts.StoragePath, ServiceName)
	err := os.MkdirAll(basePath, 0750)
	if err != nil {
		return nil, err
	}

	return &Service{
		opts:   opts,
		ticker: time.NewTicker(math.MaxInt64),
	}, nil
}

// Data returns an instance of [Data]. Calls to Data are cachable by the
// caller.
//
// Data must only be called after parsing command-line flags.
func (s *Service) Data() any {
	return map[string]string{}
}

// Definition returns the definition of the remote configuration service.
func (s *Service) Definition() service.Definition {
	return service.Definition{
		Name:       ServiceName,
		ConfigType: Arguments{},
		DependsOn:  nil, // remoteconfig has no dependencies.
	}
}

var _ service.Service = (*Service)(nil)

// Run implements [service.Service] and starts the remote configuration
// service. It will run until the provided context is canceled or there is a
// fatal error.
func (s *Service) Run(ctx context.Context, host service.Host) error {
	return nil
}

// Update implements [service.Service] and applies settings.
func (s *Service) Update(newConfig any) error {
	newArgs := newConfig.(Arguments)

	if newArgs.URL == "" {
		// TODO(@tpaschalis) We either never set the block on the first place,
		// or recently removed it. Make sure we stop everything before
		// returning.
		return nil
	}

	s.ticker.Reset(newArgs.PollFrequency)
	s.dataPath = filepath.Join(s.opts.StoragePath, ServiceName, getHash(newArgs.URL))

	httpClient, err := commonconfig.NewClientFromConfig(*newArgs.HTTPClientConfig.Convert(), "remoteconfig")
	if err != nil {
		return err
	}

	s.asClient = agentv1connect.NewAgentServiceClient(
		httpClient,
		newArgs.URL,
	)

	return nil
}

// SetModule sets up the module used to create and run pipelines fetched from a
// remote configuration endpoint.
// TODO(@tpaschalis) This is likely not the best option, but used as a starting
// point; we should find a better way of passing or creating a module from the
// root Flow controller.
func (s *Service) SetModule(mod component.Module) {
	s.mod = mod
}
