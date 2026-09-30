package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"sync/atomic"
	"testing"
)

func newIntegrationConfig(serverURL string, maxPages, maxDepth int) *config {
	baseURL, err := url.Parse(serverURL)
	if err != nil {
		panic(err)
	}

	cfg := &config{
		pages:              make(map[string]PageData),
		baseURL:            baseURL,
		mu:                 &sync.Mutex{},
		concurrencyControl: make(chan struct{}, 5),
		wg:                 &sync.WaitGroup{},
		maxPages:           maxPages,
		maxDepth:           maxDepth,
		userAgent:          defaultUserAgent,
		maxRedirects:       defaultMaxRedirects,
	}

	cfg.client = cfg.newHTTPClient()

	return cfg
}

func runIntegrationCrawl(cfg *config, rawURL string) {
	cfg.wg.Add(1)

	go func() {
		cfg.concurrencyControl <- struct{}{}
		cfg.crawlPage(rawURL, 0)
	}()

	cfg.wg.Wait()
}

func TestCrawlerIntegration_MultiPageCrawl(t *testing.T) {
	var server *httptest.Server

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, err := fmt.Fprint(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>Home</h1>
				<p>Home page</p>
				<a href="/about">About</a>
				<a href="/contact">Contact</a>
			</body>
			</html>
		`)
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, err := fmt.Fprint(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>About</h1>
				<p>About page</p>
			</body>
			</html>
		`)
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	mux.HandleFunc("/contact", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		_, err := fmt.Fprint(w, `
			<!DOCTYPE html>
			<html>
			<body>
				<h1>Contact</h1>
				<p>Contact page</p>
			</body>
			</html>
		`)
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	server = httptest.NewServer(mux)
	defer server.Close()

	cfg := newIntegrationConfig(server.URL, 10, -1)
	runIntegrationCrawl(cfg, server.URL)

	if cfg.stats.PagesCrawled != 3 {
		t.Fatalf(
			"PagesCrawled = %d, want 3",
			cfg.stats.PagesCrawled,
		)
	}

	if len(cfg.pages) != 3 {
		t.Fatalf(
			"pages count = %d, want 3",
			len(cfg.pages),
		)
	}

	homeURL := server.URL + "/"
	aboutURL := server.URL + "/about"
	contactURL := server.URL + "/contact"

	for _, pageURL := range []string{
		homeURL,
		aboutURL,
		contactURL,
	} {
		if _, ok := cfg.pages[pageURL]; !ok {
			t.Errorf("expected crawled page %q", pageURL)
		}
	}
}

func TestCrawlerIntegration_SameDomainRestriction(t *testing.T) {
	var externalRequests atomic.Int32

	externalServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			externalRequests.Add(1)
			w.Header().Set("Content-Type", "text/html")

			_, err := fmt.Fprint(w, `
				<html>
				<body>
					<h1>External</h1>
				</body>
				</html>
			`)
			if err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		},
	))
	defer externalServer.Close()

	mainServer := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")

			_, err := fmt.Fprintf(w, `
				<html>
				<body>
					<h1>Home</h1>
					<a href="%s">External</a>
				</body>
				</html>
			`, externalServer.URL)
			if err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		},
	))
	defer mainServer.Close()

	cfg := newIntegrationConfig(mainServer.URL, 10, -1)
	runIntegrationCrawl(cfg, mainServer.URL)

	if cfg.stats.PagesCrawled != 1 {
		t.Fatalf(
			"PagesCrawled = %d, want 1",
			cfg.stats.PagesCrawled,
		)
	}

	if got := externalRequests.Load(); got != 0 {
		t.Fatalf(
			"external requests = %d, want 0",
			got,
		)
	}
}

func TestCrawlerIntegration_MaxDepth(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		_, err := fmt.Fprint(w, `
			<html>
			<body>
				<h1>Root</h1>
				<a href="/level1">Level 1</a>
			</body>
			</html>
		`)
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	mux.HandleFunc("/level1", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		_, err := fmt.Fprint(w, `
			<html>
			<body>
				<h1>Level 1</h1>
				<a href="/level2">Level 2</a>
			</body>
			</html>
		`)
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	mux.HandleFunc("/level2", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		_, err := fmt.Fprint(w, `
			<html>
			<body>
				<h1>Level 2</h1>
			</body>
			</html>
		`)
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cfg := newIntegrationConfig(server.URL, 10, 1)
	runIntegrationCrawl(cfg, server.URL)

	if cfg.stats.PagesCrawled != 2 {
		t.Fatalf(
			"PagesCrawled = %d, want 2",
			cfg.stats.PagesCrawled,
		)
	}

	if _, ok := cfg.pages[server.URL+"/"]; !ok {
		t.Error("root page was not crawled")
	}

	if _, ok := cfg.pages[server.URL+"/level1"]; !ok {
		t.Error("level 1 page was not crawled")
	}

	if _, ok := cfg.pages[server.URL+"/level2"]; ok {
		t.Error("level 2 page should not have been crawled")
	}
}

func TestCrawlerIntegration_MaxPages(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")

		_, err := fmt.Fprint(w, `
			<html>
			<body>
				<h1>Root</h1>
				<a href="/one">One</a>
				<a href="/two">Two</a>
				<a href="/three">Three</a>
			</body>
			</html>
		`)
		if err != nil {
			t.Errorf("failed to write response: %v", err)
		}
	})

	for _, path := range []string{
		"/one",
		"/two",
		"/three",
	} {
		mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Type", "text/html")

			_, err := fmt.Fprintf(w, `
				<html>
				<body>
					<h1>%s</h1>
				</body>
				</html>
			`, path)
			if err != nil {
				t.Errorf("failed to write response: %v", err)
			}
		})
	}

	server := httptest.NewServer(mux)
	defer server.Close()

	cfg := newIntegrationConfig(server.URL, 2, -1)
	runIntegrationCrawl(cfg, server.URL)

	if len(cfg.pages) != 2 {
		t.Fatalf(
			"pages count = %d, want 2",
			len(cfg.pages),
		)
	}

	if cfg.stats.PagesCrawled != 2 {
		t.Fatalf(
			"PagesCrawled = %d, want 2",
			cfg.stats.PagesCrawled,
		)
	}
}
