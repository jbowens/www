package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jbowens/www/blog"
)

const (
	listenAddr = ":8080"
)

func main() {
	err := blog.Load("blog/markdown")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}

	initRoutes(http.DefaultServeMux)
	s := &http.Server{
		Addr:           listenAddr,
		Handler:        http.DefaultServeMux,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	err = s.ListenAndServe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
