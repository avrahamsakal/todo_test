{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "todo_test.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" | lower -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "todo_test.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" | lower -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "todo_test.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" | lower -}}
{{- end -}}


{{/*
Returns the hostname.
If the hostname is set in `global.hosts.todo_test.name`, that will be returned,
otherwise the hostname will be assembed using `todo_test` as the prefix, and the `malaasot.todo_test` function.
*/}}
{{- define "todo_test.hostname" -}}
{{- coalesce .Values.global.hosts.todo_test.name (include "malaasot.assembleHost"  (dict "name" "todo_test" "context" . )) -}}
{{- end -}}
