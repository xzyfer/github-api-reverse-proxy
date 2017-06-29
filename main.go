package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

const proxyHost = "https://api.github.com/"

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	
	http.HandleFunc("/", ProxyFunc)
	http.ListenAndServe(":" + port, nil)
}

func ProxyFunc(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(proxyHost)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	token := os.Getenv("AUTH_TOKEN")
	if token != "" {
		r.Header.Set("Authorization", fmt.Sprintf("token %s", os.Getenv("AUTH_TOKEN")))
	}
	
	ua := os.Getenv("USER_AGENT")
	if ua == "" {
		ua = "github-api-reverse-proxy"
	}
	r.Header.Set("User-Agent", os.Getenv("USER_AGENT"))
	
	r.Header.Set("Pragma", "no-cache")

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}
