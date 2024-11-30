package prometheusgin

import (
	"fmt"
	"strings"
	"sync"
)

type Info struct {
	name        string
	help        string
	info        string
	labels      map[string]string
	mu          sync.Mutex
	labelString string
}

func NewInfo(name, help string, info string, labels map[string]string) *Info {
	labelStr := formatLabels(labels)
	return &Info{
		name:        name,
		help:        help,
		info:        info,
		labels:      labels,
		labelString: labelStr,
	}
}

func (i *Info) SetInfo(v string) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.info = v
}

func (i *Info) Export() string {
	i.mu.Lock()
	defer i.mu.Unlock()
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# HELP %s %s\n", i.name, i.help))
	sb.WriteString(fmt.Sprintf("# TYPE %s info\n", i.name))
	if i.labelString != "" {
		sb.WriteString(fmt.Sprintf("%s{%s} \"%s\"\n", i.name, i.labelString, i.info))
	} else {
		sb.WriteString(fmt.Sprintf("%s \"%s\"\n", i.name, i.info))
	}
	return sb.String()
}
