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

	"github.com/catatsuy/bento/server"
)

var (
	appVersion string
)

func main() {
	var (
		port int
	)

	flag.IntVar(&port, "port", 8080, "port to listen")
	flag.Parse()

	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGTERM)
	signal.Notify(sigchan, syscall.SIGINT)

	var l net.Listener
	var err error

	l, err = net.ListenTCP("tcp", &net.TCPAddr{Port: port})
	if err != nil {
		panic(err.Error())
	}

	s := &http.Server{
		Handler: server.New(appVersion),
	}

	go func() {
		log.Println(s.Serve(l))
	}()

	<-sigchan

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err = s.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
