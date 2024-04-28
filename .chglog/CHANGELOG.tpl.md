{{ range .Versions }}
{{ $tag := .Tag.Name }}
<a name="{{ .Tag.Name }}"></a>
## {{ if .Tag.Previous }}[{{ .Tag.Name }}]({{ $.Info.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }}){{ else }}{{ .Tag.Name }}{{ end }} <kbd>{{ datetime "2006/01/02" .Tag.Date }}</kbd>

{{ range .CommitGroups -}}
{{ $type := .Title }}
{{ range .Commits -}}{{$type}} {{ .Subject }}<br>{{ end }}
{{ end -}}

{{ end -}}
