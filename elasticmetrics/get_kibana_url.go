package elasticmetrics

import (
	"fmt"
	"net/url"
)

func getKibanaUrl(kibanaUrl, query string) string {
	// the query string is wrapped in single quotes (escaped to %27)
	safeQuery := url.QueryEscape("'" + query + "'")
	return fmt.Sprintf("%s/app/logs/stream?logFilter=(language:lucene,query:%s)",
		kibanaUrl, safeQuery)
}
