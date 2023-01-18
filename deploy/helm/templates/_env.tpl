{{- define "_env" }}
- name: DB_HOST
  value: host
- name: DB_PORT
  value: {{ pluck .Values.global.env .Values.app.db.port | first | default .Values.app.db.port._default | quote }}
- name: DB_DATABASE
  valueFrom:
    secretKeyRef:
      name: {{ .Chart.Name }}
      key: db_database
      optional: false
- name: DB_USERNAME
  valueFrom:
    secretKeyRef:
      name: {{ .Chart.Name }}
      key: db_username
      optional: false
- name: DB_PASSWORD
  valueFrom:
    secretKeyRef:
      name: {{ .Chart.Name }}
      key: db_password
      optional: false

- name: DEBUG
  value: {{ pluck .Values.global.env .Values.app.debug_mode | first | default .Values.app.debug_mode._default | quote }}
- name: LOGGER_LEVEL
  value: {{ pluck .Values.global.env .Values.app.logger.level | first | default .Values.app.logger.level._default | quote }}
- name: METRIC_HTTP_PORT
  value: {{ pluck .Values.global.env .Values.app.metric.http_port | first | default .Values.app.metric.http_port._default | quote }}
- name: METRIC_PATH
  value: {{ pluck .Values.global.env .Values.app.metric.path | first | default .Values.app.metric.path._default | quote }}

- name: GRPC_PORT
  value: {{ pluck .Values.global.env .Values.app.api.grpc_port | first | default .Values.app.api.grpc_port._default | quote }}

{{- if .Values.app.proxy.gateway }}
- name: PROXY_GATEWAY
  value: {{ .Values.app.proxy.gateway | quote }}
{{- end }}
- name: PROXY_LIST_PATH
  value: {{ pluck .Values.global.env .Values.app.proxy.list_path | first | default .Values.app.proxy.list_path._default | quote }}

- name: TELEGRAM_TOKEN
  valueFrom:
    secretKeyRef:
      name: {{ .Chart.Name }}
      key: telegram_token
      optional: false
- name: TELEGRAM_CHAT_ID
  valueFrom:
    secretKeyRef:
      name: {{ .Chart.Name }}
      key: telegram_chat_id
      optional: false

- name: PROVIDER_POSITION_REFRESH_DURATION
  value: {{ pluck .Values.global.env .Values.app.provider.position.refresh.duration | first | default .Values.app.provider.position.refresh.duration._default | quote }}
- name: PROVIDER_POSITION_REFRESH_DURATION_FLOATING
  value: {{ pluck .Values.global.env .Values.app.provider.position.refresh.duration_floating | first | default .Values.app.provider.position.refresh.duration_floating._default | quote }}
- name: PROVIDER_TRADERS_REFRESH_DURATION
  value: {{ pluck .Values.global.env .Values.app.provider.traders.refresh.duration | first | default .Values.app.provider.traders.refresh.duration._default | quote }}
{{- end }}