{{/*
Chart name, truncated to 63 chars.
*/}}
{{- define "sentinel.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Fully qualified app name, truncated to 63 chars.
*/}}
{{- define "sentinel.fullname" -}}
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
Common labels.
*/}}
{{- define "sentinel.labels" -}}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/part-of: sentinel
{{- end }}

{{/*
API server labels.
*/}}
{{- define "sentinel.api.labels" -}}
{{ include "sentinel.labels" . }}
app.kubernetes.io/name: sentinel-api
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: api
{{- end }}

{{/*
API server selector labels.
*/}}
{{- define "sentinel.api.selectorLabels" -}}
app.kubernetes.io/name: sentinel-api
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: api
{{- end }}

{{/*
Model server labels.
*/}}
{{- define "sentinel.modelServer.labels" -}}
{{ include "sentinel.labels" . }}
app.kubernetes.io/name: sentinel-model-server
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: model-server
{{- end }}

{{/*
Model server selector labels.
*/}}
{{- define "sentinel.modelServer.selectorLabels" -}}
app.kubernetes.io/name: sentinel-model-server
app.kubernetes.io/instance: {{ .Release.Name }}
app.kubernetes.io/component: model-server
{{- end }}
