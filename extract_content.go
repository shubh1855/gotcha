package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getHeadingFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	if heading := doc.Find("h1").First(); heading.Length() > 0 {
		return strings.TrimSpace(heading.Text())
	}

	if heading := doc.Find("h2").First(); heading.Length() > 0 {
		return strings.TrimSpace(heading.Text())
	}

	return ""
}

func getFirstParagraphFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return ""
	}

	main := doc.Find("main").First()
	if main.Length() > 0 {
		if p := main.Find("p").First(); p.Length() > 0 {
			return strings.TrimSpace(p.Text())
		}
	}

	if p := doc.Find("p").First(); p.Length() > 0 {
		return strings.TrimSpace(p.Text())
	}

	return ""
}
