// prometheusgin/registry.go

package prometheusgin

import (
	"strings"
	"sync"
)

type MetricRegistry struct {
	metrics map[string][]Metric
	mu      sync.RWMutex
}

func NewMetricRegistry() *MetricRegistry {
	return &MetricRegistry{
		metrics: make(map[string][]Metric),
	}
}

func (r *MetricRegistry) Register(metric Metric) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.metrics[getMetricName(metric)] = append(r.metrics[getMetricName(metric)], metric)
}

func (r *MetricRegistry) ExportAll() string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var sb strings.Builder
	for _, metricsList := range r.metrics {
		for _, metric := range metricsList {
			sb.WriteString(metric.Export())
		}
	}
	return sb.String()
}

func getMetricName(metric Metric) string {
	switch m := metric.(type) {
	case *Counter:
		return m.name
	case *Gauge:
		return m.name
	case *Histogram:
		return m.name
	case *Summary:
		return m.name
	case *Info:
		return m.name
	case *Stateset:
		return m.name
	case *Untyped:
		return m.name
	default:
		return "unknown_metric"
	}
}

func (r *MetricRegistry) GetCounter(name string) *Counter {
	r.mu.RLock()
	defer r.mu.RUnlock()
	metricsList, exists := r.metrics[name]
	if !exists {
		return nil
	}
	for _, metric := range metricsList {
		if counter, ok := metric.(*Counter); ok {
			return counter
		}
	}
	return nil
}

func (r *MetricRegistry) GetGauge(name string) *Gauge {
	r.mu.RLock()
	defer r.mu.RUnlock()
	metricsList, exists := r.metrics[name]
	if !exists {
		return nil
	}
	for _, metric := range metricsList {
		if gauge, ok := metric.(*Gauge); ok {
			return gauge
		}
	}
	return nil
}

func (r *MetricRegistry) GetHistogram(name string) *Histogram {
	r.mu.RLock()
	defer r.mu.RUnlock()
	metricsList, exists := r.metrics[name]
	if !exists {
		return nil
	}
	for _, metric := range metricsList {
		if histogram, ok := metric.(*Histogram); ok {
			return histogram
		}
	}
	return nil
}

func (r *MetricRegistry) GetSummary(name string) *Summary {
	r.mu.RLock()
	defer r.mu.RUnlock()
	metricsList, exists := r.metrics[name]
	if !exists {
		return nil
	}
	for _, metric := range metricsList {
		if summary, ok := metric.(*Summary); ok {
			return summary
		}
	}
	return nil
}

func (r *MetricRegistry) GetInfo(name string) *Info {
	r.mu.RLock()
	defer r.mu.RUnlock()
	metricsList, exists := r.metrics[name]
	if !exists {
		return nil
	}
	for _, metric := range metricsList {
		if info, ok := metric.(*Info); ok {
			return info
		}
	}
	return nil
}

func (r *MetricRegistry) GetStateset(name string) *Stateset {
	r.mu.RLock()
	defer r.mu.RUnlock()
	metricsList, exists := r.metrics[name]
	if !exists {
		return nil
	}
	for _, metric := range metricsList {
		if stateset, ok := metric.(*Stateset); ok {
			return stateset
		}
	}
	return nil
}

func (r *MetricRegistry) GetUntyped(name string) *Untyped {
	r.mu.RLock()
	defer r.mu.RUnlock()
	metricsList, exists := r.metrics[name]
	if !exists {
		return nil
	}
	for _, metric := range metricsList {
		if untyped, ok := metric.(*Untyped); ok {
			return untyped
		}
	}
	return nil
}
