package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	listenFlag     = flag.String("listen", ":4211", "address and port to listen")
	maxTimeoutFlag = flag.String("max-timeout", "10m", "maximum timeout")
	maxTimeout     time.Duration
)

func main() {
	var err error
	// Flags and parse
	flag.Parse()
	maxTimeout, err = time.ParseDuration(*maxTimeoutFlag)
	if err != nil {
		fmt.Println("Invalid timeout format: ", err.Error())
		os.Exit(1)
	}
	if len(flag.Args()) > 0 {
		fmt.Println("Too many arguments")
		flag.Usage()
		os.Exit(1)
	}

	// Mux
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
		return
	})
	mux.HandleFunc("/timeout", handleTime)

	// Server
	server := &http.Server{
		Addr:    *listenFlag,
		Handler: mux,
	}
	serverCh := make(chan struct{})
	go func() {
		log.Printf("[INFO] server is listening on %s\n", *listenFlag)
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("[ERR] server exited with: %s", err)
		}
		close(serverCh)
	}()
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	<-signalCh

	log.Printf("[INFO] received interrupt, shutting down...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("[ERR] failed to shutdown server: %s", err)
	}
	os.Exit(2)
}

func handleTime(w http.ResponseWriter, r *http.Request) {
	var err error
	var duration time.Duration
	rawTimeout := r.URL.Query().Get("duration")
	if rawTimeout != "" {
		duration, err = time.ParseDuration(rawTimeout)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(string(err.Error())))
			return
		}
	} else {
		duration = maxTimeout
	}

	if duration > maxTimeout {
		duration = maxTimeout
	}
	time.Sleep(duration)
	w.WriteHeader(http.StatusOK)
	return
}
