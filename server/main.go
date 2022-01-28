package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nobonobo/rpicon/server/backend"
)

func main() {
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	script := "procon.py"
	flag.StringVar(&script, "script", script, "script path")
	flag.Parse()
	l, err := net.Listen("tcp", ListenAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	ctx, stop := signal.NotifyContext(context.Background(),
		os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	server := &http.Server{
		Addr:    ListenAddr,
		Handler: nil,
	}
	http.Handle("/api/", backend.New(script))
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()
	log.Println("http: Server started at", l.Addr())
	log.Fatal(server.Serve(l))
}
