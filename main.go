package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/coreos/go-systemd/v22/activation"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	time.Sleep(1 * time.Second)
	now := time.Now()
	io.WriteString(w, now.Format(time.Stamp))
	io.WriteString(w, "\n")
}

func main() {
	listeners, err := activation.Listeners()
	if err != nil {
		log.Fatal(err)
	}

	if len(listeners) == 0 {
		log.Fatal("Unexpected number of socket activation fds")
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/", HelloServer)
	srv := &http.Server{
		Handler: mux,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGHUP)
		s := <-sigint
		log.Printf("HTTP server Reload: %v", s)

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Printf("HTTP server Shutdown: %v", err)
		}
		close(idleConnsClosed)
	}()

	if err := srv.Serve(listeners[0]); err != http.ErrServerClosed {
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}
