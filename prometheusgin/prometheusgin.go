// prometheusgin/prometheusgin.go

package prometheusgin

import (
	"github.com/gin-gonic/gin"
)

type PrometheusGin struct {
	registry *MetricRegistry
	engine   *gin.Engine
}

func NewPrometheusGin() *PrometheusGin {
	reg := NewMetricRegistry()
	engine := gin.Default()
	return &PrometheusGin{
		registry: reg,
		engine:   engine,
	}
}

func (pg *PrometheusGin) UseMetricsMiddleware() {
	pg.engine.Use(PrometheusMiddleware(pg.registry))
}

func (pg *PrometheusGin) RegisterCounter(name, help string, labels map[string]string) *Counter {
	counter := NewCounter(name, help, labels)
	pg.registry.Register(counter)
	return counter
}

func (pg *PrometheusGin) RegisterGauge(name, help string, labels map[string]string) *Gauge {
	gauge := NewGauge(name, help, labels)
	pg.registry.Register(gauge)
	return gauge
}

func (pg *PrometheusGin) RegisterHistogram(name, help string, buckets []float64, labels map[string]string) *Histogram {
	histogram := NewHistogram(name, help, buckets, labels)
	pg.registry.Register(histogram)
	return histogram
}

func (pg *PrometheusGin) RegisterSummary(name, help string, quantiles []float64, labels map[string]string) *Summary {
	summary := NewSummary(name, help, quantiles, labels)
	pg.registry.Register(summary)
	return summary
}

func (pg *PrometheusGin) RegisterInfo(name, help string, info string, labels map[string]string) *Info {
	infoMetric := NewInfo(name, help, info, labels)
	pg.registry.Register(infoMetric)
	return infoMetric
}

func (pg *PrometheusGin) RegisterStateset(name, help string, state string, labels map[string]string) *Stateset {
	stateset := NewStateset(name, help, state, labels)
	pg.registry.Register(stateset)
	return stateset
}

func (pg *PrometheusGin) RegisterUntyped(name, help string, labels map[string]string) *Untyped {
	untyped := NewUntyped(name, help, labels)
	pg.registry.Register(untyped)
	return untyped
}

func (pg *PrometheusGin) MetricsHandler(path string) {
}

func (pg *PrometheusGin) Engine() *gin.Engine {
	return pg.engine
}

func (pg *PrometheusGin) Registry() *MetricRegistry {
	return pg.registry
}

func (pg *PrometheusGin) Run(addr string) error {
	return pg.engine.Run(addr)
}
