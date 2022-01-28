//go:build develop
// +build develop

package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

const ListenAddr = ":8080"

func init() {
	u, err := url.Parse("http://localhost:3000/")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", httputil.NewSingleHostReverseProxy(u))
}
