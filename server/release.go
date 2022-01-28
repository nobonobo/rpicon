//go:build !develop
// +build !develop

package main

import (
	"embed"
	"io/fs"
	"log"
	"net/http"
)

const ListenAddr = ":80"

//go:generate sh -c "cd frontend; npm run build"
//go:embed frontend/build/*
///go:embed frontend/build/_app/assets/pages/__layout.svelte-*.css
///go:embed frontend/build/_app/pages/__layout.svelte-*.js
var content embed.FS

func init() {
	pub, err := fs.Sub(content, "frontend/build")
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", http.FileServer(http.FS(pub)))
}
