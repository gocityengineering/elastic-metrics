package elasticmetrics

import (
	"testing"
)

func TestGetPrometheusUrl(t *testing.T) {
	var tests = []struct {
		description string
		url         string
		label       string
		expected    string
	}{
		{"with_special_chars", "https://prometheus.dev-euw2.common.lpge.co", "Some ${label}'", "https://prometheus.dev-euw2.common.lpge.co/graph?g0.expr=sum%28increase%28elastic_metrics_total%7Bquery%3D%22Some+%24%7Blabel%7D%27%22%7D%5B5m%5D%29%29&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=12h"},
		{"with_spaces", "https://prometheus.dev-euw2.common.lpge.co", "Some other label", "https://prometheus.dev-euw2.common.lpge.co/graph?g0.expr=sum%28increase%28elastic_metrics_total%7Bquery%3D%22Some+other+label%22%7D%5B5m%5D%29%29&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=12h"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := getPrometheusUrl(test.url, test.label)
			if result != test.expected {
				t.Errorf("Unexpected Prometheus URL for inputs %s and %s: expected %s, got %s", test.url, test.label, test.expected, result)
			}
		})
	}
}
