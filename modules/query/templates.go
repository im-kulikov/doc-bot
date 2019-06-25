package query

import "text/template"

func SearchTemplate() (*template.Template, error) {
	return template.New("searchResult").Parse(SearchTpl)
}
