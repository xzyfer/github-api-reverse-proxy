package main

import (
	"fmt"
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
	
	token := os.Getenv("AUTH_TOKEN")

	if token == "" {
		log.Fatal("$AUTH_TOKEN must be set")
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

	r.Header.Set("Authorization", fmt.Printf("%s OAUTH-TOKEN", os.Getenv("AUTH_TOKEN")))

	proxy := httputil.NewSingleHostReverseProxy(u)
	proxy.ServeHTTP(w, r)
}
