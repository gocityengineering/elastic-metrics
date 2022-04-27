package elasticmetrics

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func luceneQuery(config *Config, query string, defaultField string) (Result, error) {
	var result = Result{}
	var data = Data{}

	// fetch response
	// data struct:
	// query:
	//   bool:
	//     must:
	//       query_string:
	//         query: QUERY
	//         default_field: DEFAULTFIELD
	//       filter:
	//         range:
	//           @timestamp:
	//             gte: "now-CONFIG.QUERYFROMMINUTESm"
	data.Query.Bool.Must.QueryString.Query = query
	data.Query.Bool.Must.QueryString.DefaultField = defaultField
	data.Query.Bool.Filter.Range.Timestamp.GreaterThanOrEqual = fmt.Sprintf("now-%dm", config.QueryFromMinutes)
	bytes, err := json.Marshal(data)
	if err != nil {
		return result, fmt.Errorf("failed to marshal data object")
	}
	body := strings.NewReader(string(bytes))

	// curl -X POST 'elasticsearch-master.elastic-system.svc.cluster.local:9200/_search?pretty' -H 'Content-Type: application/json -d'{...}'
	req, err := http.NewRequest("POST", fmt.Sprintf("%s://%s:%d/_search", config.ElasticProtocol, config.ElasticService, config.ElasticPort), body)
	if err != nil {
		return result, fmt.Errorf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("failed to send POST request: %v", err)
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(&result)
	if err != nil {
		return result, fmt.Errorf("failed to parse response body: %v", err)
	}

	return result, nil
}
