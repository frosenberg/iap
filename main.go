package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func singleJoiningSlash(a, b string) string {
	aslash := strings.HasSuffix(a, "/")
	bslash := strings.HasPrefix(b, "/")
	switch {
	case aslash && bslash:
		return a + b[1:]
	case !aslash && !bslash:
		return a + "/" + b
	}
	return a + b
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}
	http.HandleFunc("/", handleHTTP)
	log.Printf("listening on port: %s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
func makeProxyHandler(target *url.URL, username string) http.HandlerFunc {
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = singleJoiningSlash(target.Path, req.URL.Path)
		if target.RawQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = target.RawQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = target.RawQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}

		// Clear authentication so target service can't see credentials
		req.Header.Set("Authentication", "")
		// Tell target service which user is authenticated
		req.Header.Set("X-IAP-User", username)
	}
	proxyHandler := &httputil.ReverseProxy{Director: director}
	return proxyHandler.ServeHTTP
}

func handleHTTP(w http.ResponseWriter, req *http.Request) {
	username := "f@gmail.com"
	targetUrl, err := url.Parse("http://localhost:9001/")
	if err != nil {
		http.Error(w, "Failed to parse target server url", 500)
		return
	}
	proxyHandler := makeProxyHandler(targetUrl, username)
	proxyHandler(w, req)
}
