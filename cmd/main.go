package main

import (
	"fmt"
	"os"

	"github.com/jbowens/www"
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

	err = www.Serve(listenAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}
