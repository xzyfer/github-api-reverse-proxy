package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	
	http.HandleFunc("/", ProxyFunc)
	http.ListenAndServe(":" + port, nil)
}

func ProxyFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusNotFound)
		http.Error(w, "Only GET requests are supported", http.StatusForbidden)
		return
	}
	
	if repos, exists := os.LookupEnv("REPO_WHITELIST"); exists {
		found := false
		for _, repo := range repos {
			found = found || strings.HasPrefix(r.URL.Path, fmt.Sprintf("/repos/%s", repo))
		}
		if !found {
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, fmt.Sprintf("Operation not in whitelist: %s", r.URL.Path), http.StatusForbidden)
			return
		}
	}
	
	if token, exists := os.LookupEnv("AUTH_TOKEN"); exists {
		r.Header.Set("Authorization", fmt.Sprintf("token %s", token))
	}
	
	if ua, exists := os.LookupEnv("USER_AGENT"); exists {
		r.Header.Set("User-Agent", ua)
	} else {
		r.Header.Set("User-Agent", "github-api-reverse-proxy")
	}
	
	r.Header.Set("Pragma", "no-cache")
	r.Host = "api.github.com"

	proxy := httputil.NewSingleHostReverseProxy(&url.URL{
		Scheme: "https",
		Host: "api.github.com",
	})
	proxy.ServeHTTP(w, r)
}
