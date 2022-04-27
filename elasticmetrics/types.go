package elasticmetrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

// configuration file
// see config/config.yaml

type Config struct {
	ElasticProtocol             string  `json:"elasticProtocol"`
	ElasticService              string  `json:"elasticService"`
	ElasticPort                 int     `json:"elasticPort"`
	KibanaUrl                   string  `json:"kibanaUrl"`
	PrometheusUrl               string  `json:"prometheusUrl"`
	QueryIntervalMinutes        int     `json:"queryIntervalMinutes"`
	QueryFromMinutes            int     `json:"queryFromMinutes"`
	MaxForwardedResultsPerQuery int     `json:"maxForwardedResultsPerQuery"`
	LogMatches                  bool    `json:"logMatches"`
	Verbose                     bool    `json:"verbose"`
	Queries                     []Query `json:"queries"`
	ElasticMetricsCounter       prometheus.CounterVec
	DeduplicationHashMap        map[uint32]bool
	EmailService                string `json:"emailService"`
	EmailServicePort            int    `json:"emailServicePort"`
	LeadingDatestampRegex       string `json:"leadingDatestampRegex"`
}

type Query struct {
	// required properties
	Label        string `json:"label"`
	LuceneQuery  string `json:"luceneQuery"`
	DefaultField string `json:"defaultField"`
	Team         string `json:"team"`
	// optional properties
	AddLabels      bool     `json:"addLabels"`
	AlertThreshold int      `json:"alertThreshold"`
	Compact        bool     `json:"compact"`
	Ignore         bool     `json:"ignore"`
	Recipients     []string `json:"recipients"`
}

// POST data structure for Lucene query
// with timestamp range
//
// query:
//   bool:
//     must:
//       query_string:
//         query: err*r
//         default_field: message
//       filter:
//         range:
//           @timestamp:
//             gte: "now-10m"

type Data struct {
	Query QueryMap `json:"query"`
}

type QueryMap struct {
	Bool BoolMap `json:"bool"`
}

type BoolMap struct {
	Must   MustMap   `json:"must"`
	Filter FilterMap `json:"filter"`
}

type MustMap struct {
	QueryString QueryStringMap `json:"query_string"`
}

type QueryStringMap struct {
	Query        string `json:"query"`
	DefaultField string `json:"default_field"`
}

type FilterMap struct {
	Range RangeMap `json:"range"`
}

type RangeMap struct {
	Timestamp TimestampMap `json:"@timestamp"`
}

type TimestampMap struct {
	GreaterThanOrEqual string `json:"gte"`
}

// Result structs
// access to:
// .Hits.Total.Value
// .Hits.Hits[*].Source.Timestamp
// .Hits.Hits[*].Source.Message
// .Hits.Hits[*].Source.Kubernetes.Namespace
// .Hits.Hits[*].Source.Kubernetes.Container.Name
// .Hits.Hits[*].Source.Kubernetes.Pod.Name
// .Hits.Hits[*].Source.Kubernetes.Node.Labels.Region
type Result struct {
	Hits HitsMap `json:"hits"`
}

type HitsMap struct {
	Total HitsMapTotal `json:"total"`
	Hits  []Hit        `json:"hits"`
}

type HitsMapTotal struct {
	Value int `json:"value"`
}

type Hit struct {
	Source HitSource `json:"_source"`
}

type HitSource struct {
	Timestamp  string        `json:"@timestamp"`
	Message    string        `json:"message"`
	Kubernetes KubernetesMap `json:"kubernetes"`
}

type KubernetesMap struct {
	Container ContainerMap `json:"container"`
	Namespace string       `json:"namespace"`
	Node      NodeMap      `json:"node"`
	Pod       PodMap       `json:"pod"`
}

type ContainerMap struct {
	Name string `json:"name"`
}

type PodMap struct {
	Name string `json:"name"`
}

type NodeMap struct {
	Labels NodeLabelsMap `json:"labels"`
}

type NodeLabelsMap struct {
	Region string `json:"topology_kubernetes_io/region"`
}

// Slack payload

// {
// 	"blocks": [
// 		{
// 			"type": "section",
// 			"text": {
// 				"type": "mrkdwn",
// 				"text": "This is a mrkdwn section block :ghost: *this is bold*, and ~this is crossed out~, and <https://google.com|this is a link>"
// 			}
// 		}
// 	]
// }

type SlackPayload struct {
	Blocks []Block `json:"blocks"`
}

type Block struct {
	Type string  `json:"type"` // section
	Text TextMap `json:"text"`
}

type TextMap struct {
	Type string `json:"type"` // mrkdown
	Text string `json:"text"`
}

type EmailServicePayload struct {
	Recipients   []string `json:"recipients"`
	Subject      string   `json:"subject"`
	MarkdownBody string   `json:"markdownBody"`
}
