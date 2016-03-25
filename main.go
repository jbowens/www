package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	listenAddr = ":8080"
)

func main() {
	initRoutes(http.DefaultServeMux)
	s := &http.Server{
		Addr:           listenAddr,
		Handler:        http.DefaultServeMux,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err := s.ListenAndServe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
