package prometheusgin

import (
	"fmt"
	"strings"
	"sync"
)

type Counter struct {
	name        string
	help        string
	value       float64
	labels      map[string]string
	mu          sync.Mutex
	labelString string
}

func NewCounter(name, help string, labels map[string]string) *Counter {
	labelStr := formatLabels(labels)
	return &Counter{
		name:        name,
		help:        help,
		labels:      labels,
		labelString: labelStr,
	}
}

func (c *Counter) Inc() {
	c.Add(1)
}

func (c *Counter) Add(v float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.value += v
}

func (c *Counter) Export() string {
	c.mu.Lock()
	defer c.mu.Unlock()
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# HELP %s %s\n", c.name, c.help))
	sb.WriteString(fmt.Sprintf("# TYPE %s counter\n", c.name))
	if c.labelString != "" {
		sb.WriteString(fmt.Sprintf("%s{%s} %s\n", c.name, c.labelString, formatFloat(c.value)))
	} else {
		sb.WriteString(fmt.Sprintf("%s %s\n", c.name, formatFloat(c.value)))
	}
	return sb.String()
}
