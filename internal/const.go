package internal

import "regexp"

const (
	maxResults      = 15
	InfoTpl         = `Version: %s BuildTime %s`
	HelloTpl        = `Привет, %s`
	SearchAnswerTpl = "Поиск по запросу %q"
	HelpTpl         = `
Доступные команды:

/start - приветствие
/help - помощь

search <keyword> - поиск ключевого слова
`
	SearchTpl = `
{{ if .Commands }}
Вот что мне удалось найти:

{{ range .Commands }}- [{{ .Name }}]({{ .Link }}) - {{ .Description }}
{{ end }}

{{ else }}
Извините, поиск не дал результата
{{ end }}
`
)

var (
	SearchRe = regexp.MustCompile(`(?i)search\s(.*)`)
	QueryRe  = regexp.MustCompile(`\s`)
)
