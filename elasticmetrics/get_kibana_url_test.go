package elasticmetrics

import (
	"testing"
)

func TestGetKibanaUrl(t *testing.T) {
	var tests = []struct {
		description string
		url         string
		query       string
		expected    string
	}{
		{"with_space", "https://mykibana.com", "message:\"some error\" AND kubernetes.container.name:\"mycontainer\"", "https://mykibana.com/app/logs/stream?logFilter=(language:lucene,query:%27message%3A%22some+error%22+AND+kubernetes.container.name%3A%22mycontainer%22%27)"},
		{"with_wildcard", "https://mykibana.com", "message:err*r", "https://mykibana.com/app/logs/stream?logFilter=(language:lucene,query:%27message%3Aerr%2Ar%27)"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			result := getKibanaUrl(test.url, test.query)
			if result != test.expected {
				t.Errorf("Unexpected Kibana URL for inputs %s and %s: expected %s, got %s", test.url, test.query, test.expected, result)
			}
		})
	}
}
