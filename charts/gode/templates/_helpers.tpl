{{/*
Expand the name of the chart.
*/}}
{{- define "gode.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "gode.fullname" -}}
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
Create chart name and version for chart label.
*/}}
{{- define "gode.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels.
*/}}
{{- define "gode.labels" -}}
helm.sh/chart: {{ include "gode.chart" . }}
{{ include "gode.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels.
*/}}
{{- define "gode.selectorLabels" -}}
app.kubernetes.io/name: {{ include "gode.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Service account name.
*/}}
{{- define "gode.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "gode.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Container image.
*/}}
{{- define "gode.image" -}}
{{- $tag := default .Chart.AppVersion .Values.image.tag }}
{{- printf "%s:%s" .Values.image.repository $tag }}
{{- end }}

{{/*
PostgreSQL host.
*/}}
{{- define "gode.postgresql.host" -}}
{{- if .Values.postgresql.enabled }}
{{- printf "%s-postgresql" .Release.Name }}
{{- else }}
{{- .Values.externalPostgresql.host }}
{{- end }}
{{- end }}

{{/*
PostgreSQL port.
*/}}
{{- define "gode.postgresql.port" -}}
{{- if .Values.postgresql.enabled }}
{{- print "5432" }}
{{- else }}
{{- .Values.externalPostgresql.port | toString }}
{{- end }}
{{- end }}

{{/*
PostgreSQL database.
*/}}
{{- define "gode.postgresql.database" -}}
{{- if .Values.postgresql.enabled }}
{{- .Values.postgresql.auth.database }}
{{- else }}
{{- .Values.externalPostgresql.database }}
{{- end }}
{{- end }}

{{/*
PostgreSQL username.
*/}}
{{- define "gode.postgresql.username" -}}
{{- if .Values.postgresql.enabled }}
{{- .Values.postgresql.auth.username }}
{{- else }}
{{- .Values.externalPostgresql.username }}
{{- end }}
{{- end }}

{{/*
PostgreSQL password secret name.
*/}}
{{- define "gode.postgresql.secretName" -}}
{{- if .Values.postgresql.enabled }}
  {{- if .Values.postgresql.auth.existingSecret }}
    {{- .Values.postgresql.auth.existingSecret }}
  {{- else }}
    {{- printf "%s-postgresql" .Release.Name }}
  {{- end }}
{{- else }}
  {{- if .Values.externalPostgresql.existingSecret }}
    {{- .Values.externalPostgresql.existingSecret }}
  {{- else }}
    {{- printf "%s-secret" (include "gode.fullname" .) }}
  {{- end }}
{{- end }}
{{- end }}

{{/*
PostgreSQL password secret key.
*/}}
{{- define "gode.postgresql.secretKey" -}}
{{- if .Values.postgresql.enabled }}
{{- print "password" }}
{{- else if .Values.externalPostgresql.existingSecret }}
{{- .Values.externalPostgresql.existingSecretKey }}
{{- else }}
{{- print "DB_PASSWORD" }}
{{- end }}
{{- end }}

{{/*
Redis host.
*/}}
{{- define "gode.redis.host" -}}
{{- if .Values.redis.enabled }}
{{- printf "%s-redis-master" .Release.Name }}
{{- else }}
{{- .Values.externalRedis.host }}
{{- end }}
{{- end }}

{{/*
Redis port.
*/}}
{{- define "gode.redis.port" -}}
{{- if .Values.redis.enabled }}
{{- print "6379" }}
{{- else }}
{{- .Values.externalRedis.port | toString }}
{{- end }}
{{- end }}

{{/*
Redis password secret name. Returns empty string if no password is configured.
*/}}
{{- define "gode.redis.secretName" -}}
{{- if .Values.redis.enabled }}
  {{- if .Values.redis.auth.enabled }}
    {{- if .Values.redis.auth.existingSecret }}
      {{- .Values.redis.auth.existingSecret }}
    {{- else }}
      {{- printf "%s-redis" .Release.Name }}
    {{- end }}
  {{- end }}
{{- else }}
  {{- if .Values.externalRedis.existingSecret }}
    {{- .Values.externalRedis.existingSecret }}
  {{- else if .Values.externalRedis.password }}
    {{- printf "%s-secret" (include "gode.fullname" .) }}
  {{- end }}
{{- end }}
{{- end }}

{{/*
Redis password secret key.
*/}}
{{- define "gode.redis.secretKey" -}}
{{- if .Values.redis.enabled }}
{{- print "redis-password" }}
{{- else if .Values.externalRedis.existingSecret }}
{{- .Values.externalRedis.existingSecretKey }}
{{- else }}
{{- print "REDIS_PASSWORD" }}
{{- end }}
{{- end }}

{{/*
Redis DB number.
*/}}
{{- define "gode.redis.db" -}}
{{- if .Values.redis.enabled }}
{{- print "0" }}
{{- else }}
{{- .Values.externalRedis.db | toString }}
{{- end }}
{{- end }}

{{/*
JWT secret name.
*/}}
{{- define "gode.jwt.secretName" -}}
{{- if .Values.jwt.existingSecret }}
{{- .Values.jwt.existingSecret }}
{{- else }}
{{- printf "%s-secret" (include "gode.fullname" .) }}
{{- end }}
{{- end }}

{{/*
JWT secret key.
*/}}
{{- define "gode.jwt.secretKey" -}}
{{- if .Values.jwt.existingSecret }}
{{- .Values.jwt.existingSecretKey }}
{{- else }}
{{- print "JWT_SECRET" }}
{{- end }}
{{- end }}
