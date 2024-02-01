# Libraries
| Name | Version | License |
| ---- | ------- | ------- |
{{- range . }}
| {{ .Name }} | {{ .Version }} | [{{ .LicenseName }}]({{ .LicenseURL }}) |
{{- end }}