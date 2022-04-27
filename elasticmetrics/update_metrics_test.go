package elasticmetrics

import (
	"testing"
)

func TestModifyTimestamp(t *testing.T) {
	var tests = []struct {
		description      string
		elasticTimestamp string
		expected         string
	}{
		{"elasticTimestamp01", "2022-01-10T14:39:44.284Z", "2022-01-10 14:39:44.284"},
		{"elasticTimestamp02", "1980-01-02T15:16:17.000Z", "1980-01-02 15:16:17.000"},
	}

	for _, test := range tests {
		t.Run(test.description, func(t *testing.T) {
			got := modifyTimestamp(test.elasticTimestamp)

			if got != test.expected {
				t.Errorf("Unexpected conversion result (got=%s expected=%s)", got, test.expected)
			}
		})
	}
}
