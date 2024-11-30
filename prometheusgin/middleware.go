// prometheusgin/middleware.go

package prometheusgin

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type MetricType int

const (
	CounterIncrement MetricType = iota
	GaugeIncrement
	GaugeDecrement
	HistogramObserve
)

type MetricUpdate struct {
	Type      MetricType
	Name      string
	Labels    map[string]string
	Value     float64
	MetricPtr interface{}
}

func PrometheusMiddleware(reg *MetricRegistry) gin.HandlerFunc {
	updateChan := make(chan MetricUpdate, 1000)
	go processMetricUpdates(reg, updateChan)

	return func(c *gin.Context) {
		start := time.Now()

		updateChan <- MetricUpdate{
			Type:   GaugeIncrement,
			Name:   "active_requests",
			Labels: map[string]string{"method": c.Request.Method, "path": c.FullPath()},
			Value:  1,
		}

		c.Next()

		duration := time.Since(start).Seconds()

		updateChan <- MetricUpdate{
			Type:   GaugeDecrement,
			Name:   "active_requests",
			Labels: map[string]string{"method": c.Request.Method, "path": c.FullPath()},
			Value:  -1,
		}

		updateChan <- MetricUpdate{
			Type:   CounterIncrement,
			Name:   "http_requests_total",
			Labels: map[string]string{"method": c.Request.Method, "path": c.FullPath()},
			Value:  1,
		}

		if c.Writer.Status() >= 400 {
			updateChan <- MetricUpdate{
				Type:   CounterIncrement,
				Name:   "http_errors_total",
				Labels: map[string]string{"method": c.Request.Method, "path": c.FullPath()},
				Value:  1,
			}
		}

		updateChan <- MetricUpdate{
			Type:   HistogramObserve,
			Name:   "http_latency_seconds_total",
			Labels: map[string]string{"method": c.Request.Method, "path": c.FullPath()},
			Value:  duration,
		}

		log.Printf("Request %s %s - Status: %d, Duration: %f seconds", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)
	}
}

func PrometheusMiddlewareWithSize(reg *MetricRegistry, size int) gin.HandlerFunc {
	updateChan := make(chan MetricUpdate, size)
	go processMetricUpdates(reg, updateChan)

	return func(c *gin.Context) {
		start := time.Now()

		updateChan <- MetricUpdate{
			Type:   GaugeIncrement,
			Name:   "active_requests",
			Labels: map[string]string{"method": c.Request.Method, "path": c.FullPath()},
			Value:  1,
		}

		c.Next()

		duration := time.Since(start).Seconds()

		updateChan <- MetricUpdate{
			Type:   GaugeDecrement,
			Name:   "active_requests",
			Labels: map[string]string{"method": c.Request.Method, "path": c.FullPath()},
			Value:  -1,
		}

		updateChan <- MetricUpdate{
			Type:   CounterIncrement,
			Name:   "http_requests_total",
			Labels: map[string]string{"method": c.Request.Method, "path": c.FullPath()},
			Value:  1,
		}

		if c.Writer.Status() >= 400 {
			updateChan <- MetricUpdate{
				Type:   CounterIncrement,
				Name:   "http_errors_total",
				Labels: map[string]string{"method": c.Request.Method, "path": c.FullPath()},
				Value:  1,
			}
		}

		updateChan <- MetricUpdate{
			Type:   HistogramObserve,
			Name:   "http_latency_seconds_total",
			Labels: map[string]string{"method": c.Request.Method, "path": c.FullPath()},
			Value:  duration,
		}

		log.Printf("Request %s %s - Status: %d, Duration: %f seconds", c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration)
	}
}

func processMetricUpdates(reg *MetricRegistry, updates chan MetricUpdate) {
	for update := range updates {
		switch update.Type {
		case CounterIncrement:
			counter := getOrCreateCounter(reg, update.Name, "A counter metric", update.Labels)
			counter.Add(update.Value)
		case GaugeIncrement:
			gauge := getOrCreateGauge(reg, update.Name, "A gauge metric", update.Labels)
			gauge.Add(update.Value)
		case GaugeDecrement:
			gauge := getOrCreateGauge(reg, update.Name, "A gauge metric", update.Labels)
			gauge.Add(update.Value)
		case HistogramObserve:
			histogram := getOrCreateHistogram(reg, update.Name, "A histogram metric", []float64{0.1, 0.3, 1.2, 5.0}, update.Labels)
			histogram.Observe(update.Value)
		}
	}
}

func getOrCreateCounter(reg *MetricRegistry, name, help string, labels map[string]string) *Counter {
	reg.mu.RLock()
	metricsList, exists := reg.metrics[name]
	reg.mu.RUnlock()
	if exists && len(metricsList) > 0 {
		if counter, ok := metricsList[0].(*Counter); ok {
			return counter
		}
	}
	counter := NewCounter(name, help, labels)
	reg.Register(counter)
	return counter
}

func getOrCreateGauge(reg *MetricRegistry, name, help string, labels map[string]string) *Gauge {
	reg.mu.RLock()
	metricsList, exists := reg.metrics[name]
	reg.mu.RUnlock()
	if exists && len(metricsList) > 0 {
		if gauge, ok := metricsList[0].(*Gauge); ok {
			return gauge
		}
	}
	gauge := NewGauge(name, help, labels)
	reg.Register(gauge)
	return gauge
}

func getOrCreateHistogram(reg *MetricRegistry, name, help string, buckets []float64, labels map[string]string) *Histogram {
	reg.mu.RLock()
	metricsList, exists := reg.metrics[name]
	reg.mu.RUnlock()
	if exists && len(metricsList) > 0 {
		if histogram, ok := metricsList[0].(*Histogram); ok {
			return histogram
		}
	}
	histogram := NewHistogram(name, help, buckets, labels)
	reg.Register(histogram)
	return histogram
}
