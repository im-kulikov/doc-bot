package query

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Command struct {
	Name        string
	Link        string
	Description string
}

type SearchResults struct {
	Commands []Command
}

var urlRegexp = regexp.MustCompile(`http[s]?://(?:[a-zA-Z]|[0-9]|[$-_@.&+]|[!*\(\),]|(?:%[0-9a-fA-F][0-9a-fA-F]))+`)

func GetLink(query string, s *goquery.Selection) string {
	if val, ok := s.Find(query).Attr("href"); ok {
		return "https://godoc.org" + strings.TrimSpace(val)
	}

	return ""
}

func trim(str string, fn func(r rune) bool) string {
	var result []rune

	for _, i := range str {
		if !fn(i) {
			result = append(result, i)
		}
	}

	return string(result)
}

func replaceLinks(str string) string {
	items := urlRegexp.FindAllString(str, -1)
	for _, item := range items {
		nstr := fmt.Sprintf("[link](%s)", item)
		str = strings.Replace(str, item, nstr, -1)
	}
	return str
}

func SearchInDocuments(query string) (*SearchResults, error) {
	var (
		results  SearchResults
		url      = fmt.Sprintf("https://godoc.org/?q=%s", query)
		doc, err = goquery.NewDocument(url)
	)

	if err != nil {
		return nil, err
	}

	results.Commands = make([]Command, 0)

	// Find the review items
	doc.Find(".container > .table.table-condensed tr").Each(func(i int, s *goquery.Selection) {
		if len(results.Commands) >= maxResults {
			return
		}

		name := s.Find("td:first-child a").Text()
		name = strings.TrimSpace(name)

		link := GetLink("td:first-child a", s)

		description := s.Find("td:last-child").Text()
		description = strings.TrimSpace(description)
		description = trim(description, func(r rune) bool {
			return r == 8203
		})

		description = replaceLinks(description)

		if len(name) < 5 && len(description) < 5 {
			return
		}

		// For each item found, get the band and title
		results.Commands = append(results.Commands, Command{
			Name:        name,
			Link:        link,
			Description: description,
		})
	})

	return &results, nil
}
