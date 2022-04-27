package elasticmetrics

import (
	"fmt"
	"net/http"
)

func MinimalElasticQuery(protocol string, service string, port int) error {
	url := fmt.Sprintf("%s://%s:%d/", protocol, service, port)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf(`Can't create minimal GET request: %v`, err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf(`Can't process minimal GET request: %v`, err)
	}

	defer resp.Body.Close()

	return nil
}
