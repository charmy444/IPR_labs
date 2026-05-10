{{/*
Имя чарта
*/}}
{{- define "postgres-infra.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Полное имя ресурса
*/}}
{{- define "postgres-infra.fullname" -}}
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
Метки стандарта Helm
*/}}
{{- define "postgres-infra.labels" -}}
helm.sh/chart: {{ include "postgres-infra.chart" . }}
{{ include "postgres-infra.selectorLabels" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}

{{- define "postgres-infra.selectorLabels" -}}
app.kubernetes.io/name: {{ include "postgres-infra.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
app: postgres
{{- end }}

{{- define "postgres-infra.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" }}
{{- end }}
