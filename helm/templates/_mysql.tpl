{{/* ######### mysql templates */}}
{{/* vim: set filetype=mustache: */}}

{{/*
Returns the mysql hostname.
If the hostname is set in `global.hosts.mysql.name`, that will be returned,
otherwise the hostname will be assembled using `mysql` as the prefix, and the `malaasot.assembleHost` function.
*/}}
{{- define "malaasot.mysql.hostname" -}}
{{- coalesce .Values.global.hosts.mysql.name (include "malaasot.assembleHost"  (dict "name" "mysql" "context" . )) -}}
{{- end -}}

{{- define "malaasot.mysql.tcp.hostname" -}}
{{- coalesce .Values.global.hosts.mysql.name (include "malaasot.assembleTCPHost"  (dict "name" "mysql" "context" . )) -}}
{{- end -}}

{{/*
Return the db hostname
If an external mysql host is provided, it will use that, otherwise it will fallback
to the service name
This overrides the upstream mysql chart so that we can deterministically
use the name of the service the upstream chart creates
*/}}
{{- define "malaasot.mysql.host" -}}
{{- coalesce .Values.global.mysql.host (printf "%s-%s" .Release.Name "mysql") -}}
{{- end -}}

{{/*
Return the replicaDb hostname
If an external mysql host is provided, it will use that, otherwise it will fallback
to the service name
This overrides the upstream mysql chart so that we can deterministically
use the name of the service the upstream chart creates
*/}}
{{- define "malaasot.mysql.replicaHost" -}}
{{- coalesce .Values.global.mysql.replication.host (printf "%s-%s" .Release.Name "mysql-replica") -}}
{{- end -}}

{{/*
Alias of malaasot.mysql.host
*/}}
{{- define "mysql.fullname" -}}
{{- template "malaasot.mysql.host" . -}}
{{- end -}}


{{/*
Return the db database name
*/}}
{{- define "malaasot.mysql.mysqlDatabase" -}}
{{- coalesce .Values.global.mysql.database "malaasot" -}}
{{- end -}}

{{/*
Return the db database name device_ordering
*/}}
{{- define "malaasot.mysql.mysqlDeviceDb" -}}
{{- coalesce .Values.global.mysql.deviceDatabase "device_ordering" -}}
{{- end -}}

{{/*
Return the db database name planpusher
*/}}
{{- define "malaasot.mysql.mysqlPlanpusherDb" -}}
{{- coalesce .Values.global.mysql.planpusherDatabase "planpusher" -}}
{{- end -}}

{{/*
Return the db database name editor
*/}}
{{- define "malaasot.mysql.mysqlEditorDb" -}}
{{- coalesce .Values.global.mysql.editorDatabase "editor" -}}
{{- end -}}

{{/*
Return the db database name food
*/}}
{{- define "malaasot.mysql.mysqlFoodDb" -}}
{{- coalesce .Values.global.mysql.foodDatabase "food_dev" -}}
{{- end -}}

{{/*
Return the db database name robin2
*/}}
{{- define "malaasot.mysql.mysqlRobin2Db" -}}
{{- coalesce .Values.global.mysql.robinDatabase "robin2" -}}
{{- end -}}

{{/*

Return the db database name reporting
*/}}
{{- define "malaasot.mysql.mysqlReportingDb" -}}
{{- coalesce .Values.global.mysql.reportingDatabase "reporting" -}}
{{- end -}}

{{/*
Return the db database name solera
*/}}
{{- define "malaasot.mysql.mysqlSoleraDb" -}}
{{- coalesce .Values.global.mysql.soleraDatabase "solerahack" -}}
{{- end -}}


{{/*
Return the db username
If the mysql username is provided, it will use that, otherwise it will fallback
to "malaasot" default
*/}}
{{- define "malaasot.mysql.mysqlUser" -}}
{{- coalesce .Values.global.mysql.username "malaasotUser" -}}
{{- end -}}

{{/*
Return the db port
If the mysql port is provided, it will use that, otherwise it will fallback
to 3306 default
*/}}
{{- define "malaasot.mysql.port" -}}
{{- coalesce .Values.global.mysql.port 3306 -}}
{{- end -}}

{{/*
Return the db conn string
If the mysql port is provided, it will use that, otherwise it will fallback
to 3306 default
*/}}
{{- define "malaasot.mysql.connString" -}}
{{- $user := include "malaasot.mysql.mysqlUser" . }}
{{- $pass := .Values.global.mysql.mysqlPassword }}
{{- $host := include "malaasot.mysql.host" . }}
{{- $port := include "malaasot.mysql.port" . }}
{{- $db := include "malaasot.mysql.mysqlDatabase" . }}
{{- printf "mysql+mysqldb://%s:%s@%s:%s/%s?charset=utf8&use_unicode=0" $user $pass $host $port $db }}
{{- end -}}

{{/*
Return the secret name
Defaults to a release-based name and falls back to .Values.global.mysql.secretName
  when using an external mysql
*/}}
{{- define "malaasot.mysql.password.secret" -}}
{{- default (printf "%s-%s" .Release.Name "mysql") .Values.global.mysql.password.secret | quote -}}
{{- end -}}

{{/*
Alias of malaasot.mysql.password.secret to override upstream mysql chart naming
*/}}
{{- define "mysql.secretName" -}}
{{- template "malaasot.mysql.password.secret" . -}}
{{- end -}}

{{/*
Return the name of the key in a secret that contains the mysql password
Uses `mysql` to match upstream mysql chart when not using an
  external mysql
*/}}
{{- define "malaasot.mysql.password.key" -}}
{{- default "mysql-password" .Values.global.mysql.password.key | quote -}}
{{- end -}}

{{- define "malaasot.mysql.rootPassword.key" -}}
{{- default "mysql-root-password" .Values.global.mysql.rootPassword.key | quote -}}
{{- end -}}
