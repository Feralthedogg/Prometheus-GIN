package prometheusgin

import (
	"fmt"
	"strings"
	"sync"
)

type Gauge struct {
	name        string
	help        string
	value       float64
	labels      map[string]string
	mu          sync.Mutex
	labelString string
}

func NewGauge(name, help string, labels map[string]string) *Gauge {
	labelStr := formatLabels(labels)
	return &Gauge{
		name:        name,
		help:        help,
		labels:      labels,
		labelString: labelStr,
	}
}

func (g *Gauge) Set(v float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value = v
}

func (g *Gauge) Inc() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value++
}

func (g *Gauge) Dec() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value--
}

func (g *Gauge) Add(v float64) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.value += v
}

func (g *Gauge) Export() string {
	g.mu.Lock()
	defer g.mu.Unlock()
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# HELP %s %s\n", g.name, g.help))
	sb.WriteString(fmt.Sprintf("# TYPE %s gauge\n", g.name))
	if g.labelString != "" {
		sb.WriteString(fmt.Sprintf("%s{%s} %s\n", g.name, g.labelString, formatFloat(g.value)))
	} else {
		sb.WriteString(fmt.Sprintf("%s %s\n", g.name, formatFloat(g.value)))
	}
	return sb.String()
}
