package file

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/go-kit/log/level"
	"github.com/grafana/agent/component"
	"github.com/grafana/agent/component/discovery"
	"github.com/grafana/loki/clients/pkg/promtail/api"
	"github.com/grafana/loki/clients/pkg/promtail/positions"
	"github.com/prometheus/common/model"
)

func init() {
	component.Register(component.Registration{
		Name: "loki.source.file",
		Args: Arguments{},

		Build: func(opts component.Options, args component.Arguments) (component.Component, error) {
			return New(opts, args.(Arguments))
		},
	})
}

const (
	pathLabel = "__path__"
)

// Arguments holds values which are used to configure the loki.source.file
// component.
type Arguments struct {
	Targets   []discovery.Target `river:"targets,attr,optional"`
	ForwardTo []chan api.Entry   `river:"forward_to,attr,optional"`
}

// DefaultArguments defines the default settings for log scraping.
var DefaultArguments = Arguments{}

// UnmarshalRiver implements river.Unmarshaler.
func (arg *Arguments) UnmarshalRiver(f func(interface{}) error) error {
	*arg = DefaultArguments

	type args Arguments
	return f((*args)(arg))
}

var (
	_ component.Component = (*Component)(nil)
)

// Component implements the loki.source.file component.
type Component struct {
	opts component.Options

	metrics *metrics

	mut     sync.RWMutex
	args    Arguments
	posFile positions.Positions
	readers map[string][]Reader
}

// New creates a new loki.source.file component.
func New(o component.Options, args Arguments) (*Component, error) {
	err := os.Mkdir(o.DataPath, 0750)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}
	positionsFile, err := positions.New(o.Logger, positions.Config{
		SyncPeriod:        10 * time.Second,
		PositionsFile:     filepath.Join(o.DataPath, "positions.yml"),
		IgnoreInvalidYaml: false,
		ReadOnly:          false,
	})
	if err != nil {
		return nil, err
	}

	c := &Component{
		opts:    o,
		metrics: newMetrics(o.Registerer),
		readers: make(map[string][]Reader),
		posFile: positionsFile,
	}

	// Call to Update() to set the targets and receivers once at the start.
	if err := c.Update(args); err != nil {
		return nil, err
	}

	return c, nil
}

// Run implements component.Component.
func (c *Component) Run(ctx context.Context) error {
	defer func() {
		level.Info(c.opts.Logger).Log("msg", "loki.source.file component shutting down, stopping readers")
		for _, rs := range c.readers {
			for _, r := range rs {
				r.Stop()
			}
		}
	}()

	<-ctx.Done()
	return nil
}

// Update implements component.Component.
func (c *Component) Update(args component.Arguments) error {
	newArgs := args.(Arguments)

	c.mut.Lock()
	defer c.mut.Unlock()
	c.args = newArgs

	if len(newArgs.Targets) == 0 {
		level.Debug(c.opts.Logger).Log("msg", "no files targets were passed, nothing will be tailed")
		return nil
	}

	// If any of the previous paths had at least one reader stopped because of
	// errors, remove all readers from the running list. They will be restarted
	// in the following loop if the path is still on the list of targets.
	c.pruneStoppedReaders()

	// 	t.reportSize(matches) ??

	var paths []string
	for _, target := range newArgs.Targets {
		path := target[pathLabel]
		paths = append(paths, path)
		var labels = make(model.LabelSet)
		for k, v := range target {
			if strings.HasPrefix(k, "__") {
				continue
			}

			labels[model.LabelName(k)] = model.LabelValue(v)
		}

		for _, receiver := range newArgs.ForwardTo {
			handler := api.AddLabelsMiddleware(labels).Wrap(api.NewEntryHandler(receiver, func() {}))

			reader, err := c.startTailing(path, handler, labels)
			if err != nil {
				continue // TODO (@tpaschalis) return err maybe?
			}

			c.readers[path] = append(c.readers[path], reader)
		}
	}

	// Stop tailing any files which are no longer in our targets.
	toStopTailing := toStopTailing(paths, c.readers)
	c.stopTailingAndRemovePosition(toStopTailing)

	return nil
}

func toStopTailing(newMatches []string, existingTailers map[string][]Reader) []string {
	// Make a set of all existing tails
	existingTails := make(map[string]struct{}, len(existingTailers))
	for file := range existingTailers {
		existingTails[file] = struct{}{}
	}

	// Make a set of what we are about to start tailing
	newTails := make(map[string]struct{}, len(newMatches))
	for _, p := range newMatches {
		newTails[p] = struct{}{}
	}

	// Find the tails in our existing which are not in the new, these need to be stopped!
	ts := missing(newTails, existingTails)
	ta := make([]string, len(ts))
	i := 0
	for t := range ts {
		ta[i] = t
		i++
	}
	return ta
}

// Returns the elements from set b which are missing from set a
func missing(as map[string]struct{}, bs map[string]struct{}) map[string]struct{} {
	c := map[string]struct{}{}
	for a := range bs {
		if _, ok := as[a]; !ok {
			c[a] = struct{}{}
		}
	}
	return c
}

// stopTailingAndRemovePosition will stop all readers for the given paths and
// remove their positions entries reader. This should be called when a file no longer
// exists and you want to remove all traces of it.
func (c *Component) stopTailingAndRemovePosition(ps []string) {
	for _, p := range ps {
		if readers, ok := c.readers[p]; ok {
			for _, reader := range readers {
				reader.Stop()
			}
			c.posFile.Remove(p)
			delete(c.readers, p)
		}
	}
}

// pruneStoppedReaders removes all readers from any paths which have at least
// one reader that has stopped running. This allows them to be restarted if
// there were errors.
func (c *Component) pruneStoppedReaders() {
	toRemove := make([]string, 0, len(c.readers))
	for path, readers := range c.readers {
		for _, reader := range readers {
			if !reader.IsRunning() {
				toRemove = append(toRemove, path)
			}
		}
	}
	for _, tr := range toRemove {
		delete(c.readers, tr)
	}
}

// startTailing starts and returns a reader for the given path. For most files,
// this will be a tailer implementation. If the file suffix alludes to it being
// a compressed file, then a decompressor will be started instead.
func (c *Component) startTailing(path string, handler api.EntryHandler, labels model.LabelSet) (Reader, error) {
	fi, err := os.Stat(path)
	if err != nil {
		level.Error(c.opts.Logger).Log("msg", "failed to tail file, stat failed", "error", err, "filename", path)
		c.metrics.totalBytes.DeleteLabelValues(path)
		return nil, fmt.Errorf("failed to stat path %s", path)
	}

	if fi.IsDir() {
		level.Info(c.opts.Logger).Log("msg", "failed to tail file", "error", "file is a directory", "filename", path)
		c.metrics.totalBytes.DeleteLabelValues(path)
		return nil, fmt.Errorf("failed to tail file, it was a directory %s", path)
	}

	var reader Reader
	if isCompressed(path) {
		level.Debug(c.opts.Logger).Log("msg", "reading from compressed file", "filename", path)
		decompressor, err := newDecompressor(
			c.metrics,
			c.opts.Logger,
			handler,
			c.posFile,
			path,
			"",
		)
		if err != nil {
			level.Error(c.opts.Logger).Log("msg", "failed to start decompressor", "error", err, "filename", path)
			return nil, fmt.Errorf("failed to start decompressor %s", err)
		}
		reader = decompressor
	} else {
		level.Debug(c.opts.Logger).Log("msg", "tailing new file", "filename", path)
		tailer, err := newTailer(
			c.metrics,
			c.opts.Logger,
			handler,
			c.posFile,
			path,
			"",
		)
		if err != nil {
			level.Error(c.opts.Logger).Log("msg", "failed to start tailer", "error", err, "filename", path)
			return nil, fmt.Errorf("failed to start tailer %s", err)
		}
		reader = tailer
	}

	return reader, nil
}
