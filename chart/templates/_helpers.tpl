{{/*
Expand the name of the chart.
*/}}
{{- define "elastic-metrics.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "elastic-metrics.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "elastic-metrics.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "elastic-metrics.labels" -}}
helm.sh/chart: {{ include "elastic-metrics.chart" . }}
{{ include "elastic-metrics.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
dora-metrics/enabled: "true"
{{- end }}

{{/*
Selector labels
*/}}
{{- define "elastic-metrics.selectorLabels" -}}
app.kubernetes.io/name: {{ include "elastic-metrics.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app: {{ include "elastic-metrics.name" . }}
prometheus-enabled: "true"
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "elastic-metrics.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "elastic-metrics.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}
