package main

import (
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func getURLsFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	urls := []string{}

	doc.Find("a[href]").Each(func(_ int, s *goquery.Selection) {
		href, ok := s.Attr("href")
		if !ok || href == "" {
			return
		}

		u, err := url.Parse(href)
		if err != nil {
			return
		}

		urls = append(urls, baseURL.ResolveReference(u).String())
	})

	return urls, nil
}

func getImagesFromHTML(htmlBody string, baseURL *url.URL) ([]string, error) {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlBody))
	if err != nil {
		return nil, err
	}

	urls := []string{}

	doc.Find("img[src]").Each(func(_ int, s *goquery.Selection) {
		src, ok := s.Attr("src")
		if !ok || src == "" {
			return
		}

		u, err := url.Parse(src)
		if err != nil {
			return
		}

		urls = append(urls, baseURL.ResolveReference(u).String())
	})

	return urls, nil
}
