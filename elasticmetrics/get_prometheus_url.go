package elasticmetrics

import (
	"fmt"
	"net/url"
)

func getPrometheusUrl(prometheusUrl, label string) string {
	expr := fmt.Sprintf("sum(increase(elastic_metrics_total{query=\"%s\"}[5m]))", label)
	return fmt.Sprintf("%s/graph?g0.expr=%s&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=12h", prometheusUrl, url.QueryEscape(expr))
}
