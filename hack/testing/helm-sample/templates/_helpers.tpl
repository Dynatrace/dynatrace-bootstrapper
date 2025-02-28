{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "testing-sample.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "testing-sample.labels" -}}
helm.sh/chart: {{ include "testing-sample.chart" . }}
{{ include "testing-sample.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "testing-sample.selectorLabels" -}}
app.kubernetes.io/name: {{ include "testing-sample.chart" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

