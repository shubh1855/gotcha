package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func newTestConfig(maxRedirects int) *config {
	cfg := &config{
		userAgent:    defaultUserAgent,
		maxRedirects: maxRedirects,
	}

	cfg.client = cfg.newHTTPClient()

	return cfg
}

func TestFetchHTML_NoRedirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, "<html><body>Hello</body></html>"); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	}))
	defer server.Close()

	cfg := newTestConfig(10)

	html, err := cfg.fetchHTML(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(html, "Hello") {
		t.Fatalf("expected HTML response")
	}
}

func TestFetchHTML_SingleRedirect(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/final", http.StatusFound)
	})

	mux.HandleFunc("/final", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, "<html><body>Redirect Success</body></html>"); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cfg := newTestConfig(10)

	html, err := cfg.fetchHTML(server.URL)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !strings.Contains(html, "Redirect Success") {
		t.Fatalf("expected redirected page")
	}
}

func TestFetchHTML_TooManyRedirects(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/a", http.StatusFound)
	})

	mux.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/b", http.StatusFound)
	})

	mux.HandleFunc("/b", func(w http.ResponseWriter, r *http.Request) {
		if _, err := fmt.Fprint(w, "<html><body>Done</body></html>"); err != nil {
			t.Fatalf("failed to write response: %v", err)
		}
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cfg := newTestConfig(1)

	_, err := cfg.fetchHTML(server.URL)
	if err == nil {
		t.Fatal("expected redirect limit error")
	}
}

func TestFetchHTML_RedirectLoop(t *testing.T) {
	mux := http.NewServeMux()

	mux.HandleFunc("/a", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/b", http.StatusFound)
	})

	mux.HandleFunc("/b", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/a", http.StatusFound)
	})

	server := httptest.NewServer(mux)
	defer server.Close()

	cfg := newTestConfig(5)

	_, err := cfg.fetchHTML(server.URL + "/a")
	if err == nil {
		t.Fatal("expected redirect loop error")
	}
}
