package prometheusgin

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

type Histogram struct {
	name        string
	help        string
	buckets     []float64
	counts      []int
	sum         float64
	labels      map[string]string
	mu          sync.Mutex
	labelString string
}

func NewHistogram(name, help string, buckets []float64, labels map[string]string) *Histogram {
	labelStr := formatLabels(labels)
	sort.Float64s(buckets)
	return &Histogram{
		name:        name,
		help:        help,
		buckets:     append([]float64{}, buckets...),
		counts:      make([]int, len(buckets)+1),
		labels:      labels,
		labelString: labelStr,
	}
}

func (h *Histogram) Observe(v float64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.sum += v
	for i, b := range h.buckets {
		if v <= b {
			h.counts[i]++
		}
	}
	h.counts[len(h.counts)-1]++
}

func (h *Histogram) Export() string {
	h.mu.Lock()
	defer h.mu.Unlock()
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("# HELP %s %s\n", h.name, h.help))
	sb.WriteString(fmt.Sprintf("# TYPE %s histogram\n", h.name))
	cumulativeCount := 0
	for i, b := range h.buckets {
		cumulativeCount += h.counts[i]
		if h.labelString != "" {
			sb.WriteString(fmt.Sprintf("%s_bucket{%s,le=\"%s\"} %d\n", h.name, h.labelString, formatFloat(b), cumulativeCount))
		} else {
			sb.WriteString(fmt.Sprintf("%s_bucket{le=\"%s\"} %d\n", h.name, formatFloat(b), cumulativeCount))
		}
	}
	cumulativeCount += h.counts[len(h.counts)-1]
	if h.labelString != "" {
		sb.WriteString(fmt.Sprintf("%s_bucket{%s,le=\"+Inf\"} %d\n", h.name, h.labelString, cumulativeCount))
	} else {
		sb.WriteString(fmt.Sprintf("%s_bucket{le=\"+Inf\"} %d\n", h.name, cumulativeCount))
	}
	if h.labelString != "" {
		sb.WriteString(fmt.Sprintf("%s_sum{%s} %s\n", h.name, h.labelString, formatFloat(h.sum)))
		sb.WriteString(fmt.Sprintf("%s_count{%s} %d\n", h.name, h.labelString, cumulativeCount))
	} else {
		sb.WriteString(fmt.Sprintf("%s_sum %s\n", h.name, formatFloat(h.sum)))
		sb.WriteString(fmt.Sprintf("%s_count %d\n", h.name, cumulativeCount))
	}
	return sb.String()
}
