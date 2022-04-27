# Elasticsearch metrics

This service watches Elasticsearch logs for certain patterns defined as code in YAML by users. When a match is made:

* the Prometheus counter `elastic_metrics_total{query="QUERYNAME"}` increases for each reported match
* if `logMatches` is true, the lines are written to stdout (up to a configurable number of matches)
* if the query's `team` property is set, the lines are written to Slack
* if the query's `recipients` array is not empty, the lines are also emailed to each email address listed

Users will mostly edit objects found in the `queries` array; each corresponds to one alert.

Individual queries carry three fields `label`, `luceneQuery` and `defaultField`. These define the Lucene query that is sent to Elasticsearch. All other query properties are optional.

See the code example below for an explanation of what each property is used for. As you add queries, be sure to see what other alerts are already found in the ConfigMap you're editing.

Many queries use a combination of regular expressions, must includes, must not includes and boolean logic. The label is used in several places: you will see it as label `query` for the Prometheus counter `elastic_metrics_total`; you will see it at the top of your Slack notification; and it is the subject of any emails sent.

```yaml
# let's get DNS names out of the way
# elastic-metrics queries Elasticsearch and matches link to
# Kibana for logs and Prometheus for metrics
elasticService: elasticsearch-master.elastic-system.svc.cluster.local
elasticPort: 9200
kibanaUrl: https://kibana.dev.acme.com
prometheusUrl: https://prometheus.dev.acme.com
# lower is better for timely alerts
# interested in hourly rate of increase? Prometheus is your friend
queryIntervalMinutes: 5
# this value should be > queryInterval and  < 2 * queryInterval
# a given log entry should appear in no more than two consecutive queries
# as the deduplication hashmap resets after guarding against repeat entries
queryFromMinutes: 7
# don't overwhelm your Slack channel
# this is the cap on results forwarded for a single query
maxForwardedResultsPerQuery: 3
# it's recommended to set logMatches to false; some queries may pick up the output of
# the elastic-metrics service by mistake, causing a feedback loop
logMatches: false
queries:
  - label: "EventFailedScheduling"
    luceneQuery: 'kubernetes.container.name:/ev.*exporter/ AND message:(+"FailedScheduling" -"namespace=kube\-system")'
    defaultField: "*"
    team: "platform" # one of platform, ecom, passManagement (note the upper-case 'M')
    ignore: false # optional: defaults to false
    compact: false # optional: defaults to false; don't send log lines to Slack
    addLabels: false # optional: add region, namespace and pod name to log lines
    recipients: [] # optional: array of email addresses
```

## Notes on Lucene queries
The query DSL is documented at [www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html](https://www.elastic.co/guide/en/elasticsearch/reference/current/query-dsl-query-string-query.html).

You can set the `defaultField` property to `message`, but this can cause feedback loops as the engine might pick up the output of `elastic-metrics` itself reporting on previous matchews.

Regular expressions have to match an entire field. It is more useful for short fields such as `kubernetes.container.name` than it is for the core `message` field.

The `+` and `-` operators force Lucene to retrieve only matches that satisfy these requirements. `message:(+"foobar")` indicates that matches MUST contain the string `foobar`. Lucene applies a complex system of weights and priorities that can get in the way when trying to pin down specific error messages we look for in the logs. (Arguably a query language  designed for logs would start from a case-sensitive, exact match. Lucene's approach is diametrically opposed to this. It revels in fuzziness and proximity matches.)

## How do I add new alerts?
Clone this repo and edit the configmap manifests in folder `chart/templates/`.

Add your own queries to the `queries` list in the nested YAML document. You can use `*` for the defaultField to search across all fields.

## One of my alerts is too noisy; how do I silence it without deleting it altogether?
Add property `ignore` and set it to `true`:
```yaml
queries:
  - label: "NoisyQuery"
    luceneQuery: "error"
    defaultField: "*"
    team: "platform"
    ignore: true
```

## What about Prometheus, PagerDuty and Grafana?
Prometheus will increment counter `elastic_metrics_total` for each match that is made. The `label` property can be found in the label `query`. Follow the metrics link in any given Slack or email match to see an example of how we can use these time series for alerting.

## How do Lucene queries work?
Lucene queries are the native Elasticsearch query engine. Kibana Query Language is a close relative.

For the most part you specify `FIELD:WORD` or `FIELD:"WORD1 WORD2"`. You can then combine multiple expressions as follows:

```
message:error* AND kubernetes.container.name:nginx
```

This will retrieve log lines containing "Error", "errors" and "erroring" among many others from nginx containers.

The `defaultField` property can be pushed to one side by setting it to `*`.
