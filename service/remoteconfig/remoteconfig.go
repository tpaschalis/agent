package remoteconfig

import (
	"context"
	"fmt"
	"hash"
	"hash/fnv"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"time"

	"connectrpc.com/connect"
	"github.com/go-kit/log"
	agentv1 "github.com/grafana/agent/api/gen/proto/go/agent/v1"
	"github.com/grafana/agent/api/gen/proto/go/agent/v1/agentv1connect"
	"github.com/grafana/agent/component"
	"github.com/grafana/agent/component/common/config"
	"github.com/grafana/agent/internal/agentseed"
	"github.com/grafana/agent/pkg/flow/logging/level"
	"github.com/grafana/agent/service"
	commonconfig "github.com/prometheus/common/config"
)

var fnvHash hash.Hash32 = fnv.New32()

func getHash(in string) string {
	fnvHash.Write([]byte(in))
	defer fnvHash.Reset()

	return fmt.Sprintf("%x", fnvHash.Sum(nil))
}

// TODO: Add metrics for every case.
// TODO: Add jitter for requests.
// TODO: Handle backoff from API in response to 429 or Retry-After headers.

// Service implements a service for remote configuration.
type Service struct {
	opts                  Options
	args                  Arguments
	mod                   component.Module
	asClient              agentv1connect.AgentServiceClient
	ticker                *time.Ticker
	dataPath              string
	getConfigRequest      *connect.Request[agentv1.GetConfigRequest]
	currentConfigHash     string
	lastSuccesfulContents []byte
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
	ID               string                   `river:"id,attr,optional"`
	Metadata         map[string]string        `river:"metadata,attr,optional"`
	PollFrequency    time.Duration            `river:"poll_frequency,attr,optional"`
	HTTPClientConfig *config.HTTPClientConfig `river:",squash"`
}

// GetDefaultArguments populates the default values for the Arguments struct.
func GetDefaultArguments() Arguments {
	return Arguments{
		ID:               agentseed.Get().UID,
		Metadata:         make(map[string]string),
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
	// We're on the initial start-up of the service.
	// Let's try to read from the API and load the response contents.
	// If either the request itself or loading its contents fails, try to read
	// from the on-disk cache on a best-effort basis.
	var gcr *connect.Response[agentv1.GetConfigResponse]
	var getErr, loadErr error

	gcr, getErr = s.asClient.GetConfig(ctx, s.getConfigRequest)
	if getErr == nil {
		// Reading from the API succeeded, let's try to load the contents.
		loadErr = s.mod.LoadConfig([]byte(gcr.Msg.GetContent()), nil)
		if loadErr != nil {
			level.Info(s.opts.Logger).Log("msg", "could not load the API response contents")
		} else {
			s.lastSuccesfulContents = []byte(gcr.Msg.GetContent())
		}
	}
	if getErr != nil || loadErr != nil {
		// Either the API call or loading its contents failed, let's try the
		// on-disk cache.
		b, err := os.ReadFile(s.dataPath)
		if err != nil {
			level.Info(s.opts.Logger).Log("msg", "could not read from the on-disk cache")
		}
		if len(b) > 0 {
			loadErr = s.mod.LoadConfig(b, nil)
			if err != nil {
				level.Info(s.opts.Logger).Log("msg", "could not load the on-disk cache contents")
			} else {
				s.lastSuccesfulContents = b
				level.Info(s.opts.Logger).Log("msg", "loaded with on-disk cache contents successfully")
			}
		}
	}

	go func() {
		if err := s.mod.Run(ctx); err != nil {
			level.Info(s.opts.Logger).Log("msg", "remote_configuration module exited with", "err", err)
		}
	}()

	for {
		select {
		case <-s.ticker.C:
			gcr, err := s.asClient.GetConfig(ctx, s.getConfigRequest)
			if err != nil {
				level.Error(s.opts.Logger).Log("msg", "failed to fetch remote_configuration", "err", err)
				continue
			}

			newConfig := []byte(gcr.Msg.Content)

			newConfigHash := getHash(gcr.Msg.Content)
			if s.currentConfigHash == newConfigHash {
				continue
			}

			// The polling loop got a new configuration to be loaded,
			// attempt to load it.
			err = s.mod.LoadConfig(newConfig, nil)
			if err != nil {
				level.Error(s.opts.Logger).Log("msg", "failed to load fetched remote_configuration", "err", err)

				// Since it failed, try to reload the last successful
				// configuration instead.
				s.mod.LoadConfig(s.lastSuccesfulContents, nil)
				continue
			}

			// If successful, flush to disk and keep a copy.
			err = os.WriteFile(s.dataPath, newConfig, 0750)
			if err != nil {
				level.Error(s.opts.Logger).Log("msg", "failed to flush remote_configuration contents the on-disk cache", "err", err)
			}
			s.lastSuccesfulContents = newConfig
		case <-ctx.Done():
			s.ticker.Stop()
			return nil
		default:
		}
	}
}

// Update implements [service.Service] and applies settings.
func (s *Service) Update(newConfig any) error {
	newArgs := newConfig.(Arguments)

	if newArgs.URL == "" {
		// TODO: We either never set the block on the first place,
		// or recently removed it. Make sure we stop everything before
		// returning.
		return nil
	}

	// TODO: Should we also include the username, or some other property?
	s.dataPath = filepath.Join(s.opts.StoragePath, ServiceName, getHash(newArgs.URL))
	s.ticker.Reset(newArgs.PollFrequency)

	if !reflect.DeepEqual(s.args.HTTPClientConfig, newArgs.HTTPClientConfig) {
		httpClient, err := commonconfig.NewClientFromConfig(*newArgs.HTTPClientConfig.Convert(), "remoteconfig")
		if err != nil {
			return err
		}
		s.asClient = agentv1connect.NewAgentServiceClient(
			httpClient,
			newArgs.URL,
		)
	}

	// TODO: Is this ok to reuse since it contains some kind of 'state' field?
	// TODO: Wire in Agent ID, Metadata from the Arguments.
	s.getConfigRequest = connect.NewRequest(&agentv1.GetConfigRequest{
		Id:       newArgs.ID,
		Metadata: newArgs.Metadata,
	})

	s.args = newArgs

	return nil
}

// SetModule sets up the module used to create and run pipelines fetched from a
// remote configuration endpoint.
// TODO: This is likely not the best option, but used as a starting
// point; we should find a better way of passing or creating a module from the
// root Flow controller, likely after Modules v2.0 is finished.
func (s *Service) SetModule(mod component.Module) {
	s.mod = mod
}
