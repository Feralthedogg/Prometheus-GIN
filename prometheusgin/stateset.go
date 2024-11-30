package prometheusgin

import (
	"fmt"
	"strings"
	"sync"
)

type Stateset struct {
	name        string
	help        string
	state       string
	labels      map[string]string
	mu          sync.Mutex
	labelString string
}

func NewStateset(name, help string, state string, labels map[string]string) *Stateset {
	labelStr := formatLabels(labels)
	return &Stateset{
		name:        name,
		help:        help,
		state:       state,
		labels:      labels,
		labelString: labelStr,
	}
}

func (s *Stateset) SetState(v string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.state = v
}

func (s *Stateset) Export() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# HELP %s %s\n", s.name, s.help))
	sb.WriteString(fmt.Sprintf("# TYPE %s stateset\n", s.name))
	if s.labelString != "" {
		sb.WriteString(fmt.Sprintf("%s{%s} \"%s\"\n", s.name, s.labelString, s.state))
	} else {
		sb.WriteString(fmt.Sprintf("%s \"%s\"\n", s.name, s.state))
	}
	return sb.String()
}
