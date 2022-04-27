package elasticmetrics

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/ghodss/yaml"
	"github.com/prometheus/client_golang/prometheus"
)

func ParseConfig(configPath string, config *Config) error {
	byteArray, err := ioutil.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("can't read configuration file %s: %v", configPath, err)
	}

	jsonArray, err := yaml.YAMLToJSON(byteArray)
	if err != nil {
		return fmt.Errorf("can't convert configuration file %s to JSON: %v", configPath, err)
	}

	err = json.Unmarshal(jsonArray, config)
	if err != nil {
		return fmt.Errorf("can't parse JSON configuration: %v", err)
	}

	// register collector
	config.ElasticMetricsCounter = *prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "elastic_metrics_total",
		Help: "counter for log entries matching rules",
	},
		[]string{
			"query",
		})
	prometheus.MustRegister(config.ElasticMetricsCounter)

	// initialise deduplication map
	config.DeduplicationHashMap = make(map[uint32]bool)

	return nil
}
