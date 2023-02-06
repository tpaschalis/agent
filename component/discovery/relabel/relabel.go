package relabel

import (
	"context"
	"reflect"
	"sync"

	"github.com/go-kit/log/level"
	"github.com/grafana/agent/component"
	flow_relabel "github.com/grafana/agent/component/common/relabel"
	"github.com/grafana/agent/component/discovery"
	"github.com/grafana/agent/pkg/river"
	lru "github.com/hashicorp/golang-lru"
	"github.com/prometheus/prometheus/model/labels"
	"github.com/prometheus/prometheus/model/relabel"
)

func init() {
	component.Register(component.Registration{
		Name:    "discovery.relabel",
		Args:    Arguments{},
		Exports: Exports{},

		Build: func(opts component.Options, args component.Arguments) (component.Component, error) {
			return New(opts, args.(Arguments))
		},
	})
}

// Arguments holds values which are used to configure the discovery.relabel component.
type Arguments struct {
	// Targets contains the input 'targets' passed by a service discovery component.
	Targets []discovery.Target `river:"targets,attr"`

	// The relabelling rules to apply to the each target's label set.
	RelabelConfigs []*flow_relabel.Config `river:"rule,block,optional"`

	// The maximum number of items to hold in the component's LRU cache.
	MaxCacheSize int `river:"max_cache_size,attr,optional"`
}

// DefaultArguments provides the default arguments for the loki.relabel
// component.
var DefaultArguments = Arguments{
	MaxCacheSize: 10_000,
}

var _ river.Unmarshaler = (*Arguments)(nil)

// UnmarshalRiver implements river.Unmarshaler.
func (a *Arguments) UnmarshalRiver(f func(interface{}) error) error {
	*a = DefaultArguments

	type arguments Arguments
	return f((*arguments)(a))
}

// Exports holds values which are exported by the discovery.relabel component.
type Exports struct {
	Output []discovery.Target `river:"output,attr"`
	Rules  flow_relabel.Rules `river:"rules,attr"`
}

// Component implements the discovery.relabel component.
type Component struct {
	opts component.Options

	mut          sync.RWMutex
	rcs          []*relabel.Config
	cache        *lru.Cache
	maxCacheSize int
}

var _ component.Component = (*Component)(nil)

// New creates a new discovery.relabel component.
func New(o component.Options, args Arguments) (*Component, error) {
	cache, err := lru.New(args.MaxCacheSize)
	if err != nil {
		return nil, err
	}

	c := &Component{
		opts:         o,
		cache:        cache,
		maxCacheSize: args.MaxCacheSize,
	}

	// Call to Update() to set the output once at the start
	if err := c.Update(args); err != nil {
		return nil, err
	}

	return c, nil
}

// Run implements component.Component.
func (c *Component) Run(ctx context.Context) error {
	<-ctx.Done()
	return nil
}

// Update implements component.Component.
func (c *Component) Update(args component.Arguments) error {
	c.mut.Lock()
	defer c.mut.Unlock()

	newArgs := args.(Arguments)

	targets := make([]discovery.Target, 0, len(newArgs.Targets))
	relabelConfigs := flow_relabel.ComponentToPromRelabelConfigs(newArgs.RelabelConfigs)
	if relabelingChanged(c.rcs, relabelConfigs) {
		level.Debug(c.opts.Logger).Log("msg", "received new relabel configs, purging cache")
		c.cache.Purge()
		c.rcs = relabelConfigs
		// c.metrics.cacheSize.Set(0)
	}
	if newArgs.MaxCacheSize != c.maxCacheSize {
		evicted := c.cache.Resize(newArgs.MaxCacheSize)
		if evicted > 0 {
			level.Debug(c.opts.Logger).Log("msg", "resizing the cache lead to evicting of items", "len_items_evicted", evicted)
		}
	}

	for _, t := range newArgs.Targets {
		lset := c.relabel(t)
		if lset != nil {
			targets = append(targets, lset)
		}
	}

	c.opts.OnStateChange(Exports{
		Output: targets,
		Rules:  newArgs.RelabelConfigs,
	})

	return nil
}

func componentMapToPromLabels(ls discovery.Target) labels.Labels {
	res := make([]labels.Label, 0, len(ls))
	for k, v := range ls {
		res = append(res, labels.Label{Name: k, Value: v})
	}

	return res
}

func promLabelsToComponent(ls labels.Labels) discovery.Target {
	res := make(map[string]string, len(ls))
	for _, l := range ls {
		res[l.Name] = l.Value
	}

	return res
}

func relabelingChanged(prev, next []*relabel.Config) bool {
	if len(prev) != len(next) {
		return true
	}
	for i := range prev {
		if !reflect.DeepEqual(prev[i], next[i]) {
			return true
		}
	}
	return false
}

type cacheItem struct {
	original  labels.Labels
	relabeled discovery.Target
}

func (c *Component) relabel(tgt discovery.Target) discovery.Target {
	lset := labels.FromMap(tgt)
	hash := lset.Hash()

	// Let's look in the cache for the hash of the entry's labels.
	val, found := c.cache.Get(hash)

	// We've seen this hash before; let's see if we've already relabeled this
	// specific entry before and can return early, or if it's a collision.
	if found {
		for _, ci := range val.([]cacheItem) {
			if labels.Equal(lset, ci.original) {
				// c.metrics.cacheHits.Inc()
				return ci.relabeled
			}
		}
	}

	// Seems like it's either a new entry or a hash collision.
	// c.metrics.cacheMisses.Inc()
	relabeled := c.process(lset)

	// In case it's a new hash, initialize it as a new cacheItem.
	// If it was a collision, append the result to the cached slice.
	if !found {
		val = []cacheItem{{lset, relabeled}}
	} else {
		val = append(val.([]cacheItem), cacheItem{lset, relabeled})
	}

	c.cache.Add(hash, val)
	// c.metrics.cacheSize.Set(float64(c.cache.Len()))

	return relabeled
}

func (c *Component) process(lset labels.Labels) discovery.Target {
	res := relabel.Process(lset, c.rcs...)
	return res.Map()
}
