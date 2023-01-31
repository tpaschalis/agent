package metric

import (
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
)

// DefaultHistogramConfig sets the defaults for a Histogram.
var DefaultHistogramConfig = HistogramConfig{
	IdleDuration: 5 * time.Minute,
}

// HistogramConfig defines a histogram metric whose values are bucketed.
type HistogramConfig struct {
	// Shared fields
	Name         string        `river:"name,attr"`
	Description  string        `river:"description,attr,optional"`
	Source       string        `river:"source,attr,optional"`
	Prefix       string        `river:"prefix,attr,optional"`
	IdleDuration time.Duration `river:"max_idle_duration,attr,optional"`
	Value        string        `river:"value,attr,optional"`

	// Histogram-specific fields
	Buckets []float64 `river:"buckets,attr"`
}

// UnmarshalRiver implements the unmarshaller
func (h *HistogramConfig) UnmarshalRiver(f func(v interface{}) error) error {
	*h = DefaultHistogramConfig
	type histogram HistogramConfig
	err := f((*histogram)(h))
	if err != nil {
		return err
	}

	if h.IdleDuration < 1*time.Second {
		return fmt.Errorf("idle duration must be greater than 1s")
	}

	if h.Source == "" {
		h.Source = h.Name
	}
	return nil
}

// Histograms is a vector of histograms for a each log stream.
type Histograms struct {
	*metricVec
	Cfg *HistogramConfig
}

// NewHistograms creates a new histogram vec.
func NewHistograms(name string, config *HistogramConfig) (*Histograms, error) {
	return &Histograms{
		metricVec: newMetricVec(func(labels map[string]string) prometheus.Metric {
			return &expiringHistogram{prometheus.NewHistogram(prometheus.HistogramOpts{
				Help:        config.Description,
				Name:        name,
				ConstLabels: labels,
				Buckets:     config.Buckets,
			}),
				0,
			}
		}, int64(config.IdleDuration.Seconds())),
		Cfg: config,
	}, nil
}

// With returns the histogram associated with a stream labelset.
func (h *Histograms) With(labels model.LabelSet) prometheus.Histogram {
	return h.metricVec.With(labels).(prometheus.Histogram)
}

type expiringHistogram struct {
	prometheus.Histogram
	lastModSec int64
}

// Observe adds a single observation to the histogram.
func (h *expiringHistogram) Observe(val float64) {
	h.Histogram.Observe(val)
	h.lastModSec = time.Now().Unix()
}

// HasExpired implements Expirable
func (h *expiringHistogram) HasExpired(currentTimeSec int64, maxAgeSec int64) bool {
	return currentTimeSec-h.lastModSec >= maxAgeSec
}
