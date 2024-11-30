package prometheusgin

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type Summary struct {
	name         string
	help         string
	quantiles    []float64
	count        int
	sum          float64
	observations []float64
	labels       map[string]string
	mu           sync.Mutex
	labelString  string
}

func NewSummary(name, help string, quantiles []float64, labels map[string]string) *Summary {
	labelStr := formatLabels(labels)
	return &Summary{
		name:         name,
		help:         help,
		quantiles:    quantiles,
		labels:       labels,
		labelString:  labelStr,
		observations: []float64{},
	}
}

func (s *Summary) Observe(v float64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.count++
	s.sum += v
	s.observations = append(s.observations, v)
}

func (s *Summary) Export() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# HELP %s %s\n", s.name, s.help))
	sb.WriteString(fmt.Sprintf("# TYPE %s summary\n", s.name))

	if len(s.observations) > 0 {
		sorted := append([]float64{}, s.observations...)
		sort.Float64s(sorted)
		for _, q := range s.quantiles {
			pos := q * float64(len(sorted)+1)
			var quantile float64
			if pos < 1 {
				quantile = sorted[0]
			} else if pos >= float64(len(sorted)) {
				quantile = sorted[len(sorted)-1]
			} else {
				lower := sorted[int(pos)-1]
				upper := sorted[int(pos)]
				quantile = lower + (upper-lower)*(pos-float64(int(pos)))
			}
			if s.labelString != "" {
				sb.WriteString(fmt.Sprintf("%s{quantile=\"%s\",%s} %s\n", s.name, formatFloat(q), s.labelString, formatFloat(quantile)))
			} else {
				sb.WriteString(fmt.Sprintf("%s{quantile=\"%s\"} %s\n", s.name, formatFloat(q), formatFloat(quantile)))
			}
		}
	}

	if s.labelString != "" {
		sb.WriteString(fmt.Sprintf("%s_sum{%s} %s\n", s.name, s.labelString, formatFloat(s.sum)))
		sb.WriteString(fmt.Sprintf("%s_count{%s} %d\n", s.name, s.labelString, s.count))
	} else {
		sb.WriteString(fmt.Sprintf("%s_sum %s\n", s.name, formatFloat(s.sum)))
		sb.WriteString(fmt.Sprintf("%s_count %d\n", s.name, s.count))
	}

	return sb.String()
}
