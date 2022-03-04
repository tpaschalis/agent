package ebpf

import (
	"context"
	"fmt"

	"github.com/go-kit/log"
	"github.com/grafana/agent/pkg/integrations/v2"
)

type ebpfHandler struct{}

func (e ebpfHandler) RunIntegration(ctx context.Context) error {
	fmt.Println("Running epbf handler!")
	return nil
}

type Config struct{}

// var DefaultConfig = Config{}
// func (c *Config) UnmarshalYAML(unmarshal func(interface{}) error) error {
// 	*c = DefaultConfig
// 	type plain Config
// 	return unmarshal((*plain)(c))
// }

func (c *Config) ApplyDefaults(globals integrations.Globals) error {
	return nil
}

func (c *Config) Identifier(globals integrations.Globals) (string, error) {
	return globals.AgentIdentifier, nil
}

func (c *Config) Name() string { return "ebpf" }

func (c *Config) NewIntegration(l log.Logger, globals integrations.Globals) (integrations.Integration, error) {
	ebpf := ebpfHandler{}

	return ebpf, nil
}

func init() {
	integrations.Register(&Config{}, integrations.TypeSingleton)
}
