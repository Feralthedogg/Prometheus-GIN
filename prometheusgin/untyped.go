package prometheusgin

import (
	"fmt"
	"strings"
	"sync"
)

type Untyped struct {
	name        string
	help        string
	value       float64
	labels      map[string]string
	mu          sync.Mutex
	labelString string
}

func NewUntyped(name, help string, labels map[string]string) *Untyped {
	labelStr := formatLabels(labels)
	return &Untyped{
		name:        name,
		help:        help,
		labels:      labels,
		labelString: labelStr,
	}
}

func (u *Untyped) Set(v float64) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.value = v
}

func (u *Untyped) Export() string {
	u.mu.Lock()
	defer u.mu.Unlock()
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# HELP %s %s\n", u.name, u.help))
	sb.WriteString(fmt.Sprintf("# TYPE %s untyped\n", u.name))
	if u.labelString != "" {
		sb.WriteString(fmt.Sprintf("%s{%s} %s\n", u.name, u.labelString, formatFloat(u.value)))
	} else {
		sb.WriteString(fmt.Sprintf("%s %s\n", u.name, formatFloat(u.value)))
	}
	return sb.String()
}
