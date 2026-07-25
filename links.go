package main

import "net/url"

func classifyLinks(baseURL *url.URL, links []string) (internal, external []string) {
	internal = make([]string, 0, len(links))
	external = make([]string, 0)

	for _, link := range links {
		parsedURL, err := url.Parse(link)
		if err != nil {
			continue
		}

		if parsedURL.Host == baseURL.Host {
			internal = append(internal, link)
		} else {
			external = append(external, link)
		}
	}

	return internal, external
}
