package main

import (
	"crypto/tls"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gocityengineering/elastic-metrics/elasticmetrics"
	"github.com/gocityengineering/elastic-metrics/encodingutils"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// run tests defined in datadir
func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage: %s`, filepath.Base(os.Args[0]))
		flag.PrintDefaults()
		os.Exit(0)
	}

	configPath := flag.String("c", "/etc/elastic-metrics/config.yaml", "configuration path")
	schemaPath := "/etc/elastic-metrics/schema.yaml"

	flag.Parse()

	// don't add own timestamp - we're interested in the original timestamps only
	log.SetFlags(0)

	// validate configuration file
	err := encodingutils.ValidateFile(*configPath, "/etc/elastic-metrics/schema.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Config file %s is not valid against schema %s: %v\n", *configPath, schemaPath, err)
		os.Exit(1)
	}

	// parse configuration file
	var config = elasticmetrics.Config{
		ElasticProtocol:             "https",
		ElasticService:              "",
		ElasticPort:                 9200,
		KibanaUrl:                   "",
		PrometheusUrl:               "",
		QueryIntervalMinutes:        5,
		QueryFromMinutes:            6,
		MaxForwardedResultsPerQuery: 4,
		LogMatches:                  false,
		Verbose:                     false,
		Queries:                     []elasticmetrics.Query{},
		ElasticMetricsCounter:       prometheus.CounterVec{},
		DeduplicationHashMap:        map[uint32]bool{},
		EmailService:                "",
		EmailServicePort:            0,
		LeadingDatestampRegex:       `^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}\.\d{3}[\t ]*`,
	}
	err = elasticmetrics.ParseConfig(*configPath, &config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Can't load config file %s: %v\n", *configPath, err)
		os.Exit(1)
	}

	// elastic API requests will fail unless we set InsecureSkipVerify to true
	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	err = elasticmetrics.MinimalElasticQuery(config.ElasticProtocol, config.ElasticService, config.ElasticPort)
	if err != nil {
		log.Printf("Minimal Elasticsearch query failed: %v\n", err)
	}

	go func() {
		for {
			err = elasticmetrics.UpdateMetrics(&config)
			if err != nil {
				log.Printf("Can't update metrics: %v\n", err)
			}
			time.Sleep(time.Duration(config.QueryIntervalMinutes) * time.Minute)
		}
	}()

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
