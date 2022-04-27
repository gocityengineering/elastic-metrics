package elasticmetrics

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

func UpdateMetrics(config *Config) error {
	newDeduplicationHashMap := make(map[uint32]bool)
	emailService := config.EmailService
	emailServicePort := config.EmailServicePort
	// verbose := config.Verbose

	// if log lines begin with an identical datestamp, the stamp should go
	// users are welcome to add other formats they'd like to omit in config
	leadingDatestampRegex := regexp.MustCompile(config.LeadingDatestampRegex)

	start := time.Now()

	// process queries
	slackBuffer := ""
	emailBuffer := ""
  hasBothOptionalUrls := len(config.KibanaUrl) > 0 && len(config.PrometheusUrl) > 0
  hasNeitherOptionalUrl := len(config.KibanaUrl) == 0 && len(config.PrometheusUrl) == 0

	for _, query := range config.Queries {

		// skip if set to be ignored in YAML
		if query.Ignore == true {
			continue
		}

		result, err := luceneQuery(config, query.LuceneQuery, query.DefaultField)

		if err != nil {
			return fmt.Errorf("can't query Elasticsearch: %v", err)
		}

		alertThreshold := 0
		if query.AlertThreshold > 0 {
			alertThreshold = query.AlertThreshold
		}

		compact := query.Compact
		recipients := query.Recipients

		// skip if no matches
		total := result.Hits.Total.Value
		if total == 0 {
			continue
		}

		// increase counter even if below threshold
		// the count may include duplicates
		config.ElasticMetricsCounter.With(prometheus.Labels{"query": query.Label}).Add(float64(total))

		// skip if below threshold
		if total < alertThreshold {
			continue
		}

		// don't log anything yet - there may be hits in other containers/namespaces
		counter := 0

    // messy output code should be handled by a function?
    kibanaUrl := getKibanaUrl(config.KibanaUrl, query.LuceneQuery)
    prometheusUrl := getPrometheusUrl(config.PrometheusUrl, query.Label)

    // build Slack buffer
		slackBuffer = fmt.Sprintf(":warning: *%s*\n", query.Label)
    emailBuffer = fmt.Sprintf("**%s**\n\n", query.Label)

    if len(config.KibanaUrl) > 0 {
      slackBuffer += fmt.Sprintf("<%s|*Logs*>", kibanaUrl)
      emailBuffer += fmt.Sprintf("[Logs](%s)", kibanaUrl)
    }
    if hasBothOptionalUrls {
      slackBuffer += " | "
      emailBuffer += " | "
    }
    if len(config.PrometheusUrl) > 0 {
      slackBuffer += fmt.Sprintf("<%s|*Metrics*>", prometheusUrl)
      emailBuffer += fmt.Sprintf("[Metrics](%s)", prometheusUrl)
    }
    if !hasNeitherOptionalUrl {
      slackBuffer += "\n"
      emailBuffer += "\n\n"
    }
    slackBuffer += fmt.Sprintf("*Query*: %s\n*Default field*: %s\n*Matches*: %d\n", query.LuceneQuery, query.DefaultField, total)
    emailBuffer += fmt.Sprintf("**Query**: %s\n\n**Default field**: %s\n\n**Matches**: %d\n\n", query.LuceneQuery, query.DefaultField, total)

		for _, hit := range result.Hits.Hits {
			timestamp := hit.Source.Timestamp
			timestamp = modifyTimestamp(timestamp)

			message := hit.Source.Message
			message = leadingDatestampRegex.ReplaceAllString(message, "")
			labels := ""
			if query.AddLabels == true {
				labels = fmt.Sprintf(" region=%s namespace=%s pod=%s",
					hit.Source.Kubernetes.Node.Labels.Region,
					hit.Source.Kubernetes.Namespace,
					hit.Source.Kubernetes.Pod.Name)
			}

			logLine := fmt.Sprintf("%s%s%s", timestamp, labels, message)
			hash := hash(logLine)

			// consider both previous and current run for deduplication
			if config.DeduplicationHashMap[hash] || newDeduplicationHashMap[hash] {
				log.Printf("skipped duplicate\n")
				continue
			}

			counter = counter + 1

			// remember hash for next pass
			newDeduplicationHashMap[hash] = true

			// update Slack if compact flag not set
			if !compact {
				slackBuffer += fmt.Sprintf("> *%s*%s %s\n", timestamp, labels, message)
			}

			// update email if recipients listed
			if len(recipients) > 0 {
				emailBuffer += fmt.Sprintf(">**%s**%s\n<pre>%s</pre>\n\n", timestamp, labels, message)
			}

			// log only if requested
			if config.LogMatches == true {
				log.Printf("%s%s %s\n", timestamp, labels, logLine)
			}

			// only forward up to the limit set
			if counter >= config.MaxForwardedResultsPerQuery {
				break
			}
		}

		config.DeduplicationHashMap = newDeduplicationHashMap
		if counter == 0 {
			continue
		}

		// Slack if requested
		if query.Team != "" {
			err = slackQuery(query.Team, slackBuffer)
			if err != nil {
				return fmt.Errorf("can't notify Slack: %v", err)
			}
		}

		// Email if requested
		if len(recipients) > 0 {
			err = emailServiceQuery(recipients, query.Label, emailBuffer, emailService, emailServicePort)
			if err != nil {
				return fmt.Errorf("can't mail recipients: %v", err)
			}
		}
	}

	elapsed := time.Since(start)
	if elapsed > time.Duration(config.QueryIntervalMinutes)*time.Minute {
		log.Printf("Queries took longer (%s) than interval set (%dm)\n", elapsed, config.QueryIntervalMinutes)
	}

	return nil
}

func modifyTimestamp(s string) string {
	// example timestamp: 2022-01-10T14:39:44.284Z
	// remove trailing "Z"
	s = s[:len(s)-1]
	// replace T
	s = strings.ReplaceAll(s, "T", " ")
	return s
}
