package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"

	"github.com/jbowens/assets"
)

var (
	htmlTemplates = map[string]*template.Template{
		"index.html": nil,
		"about.html": nil,
	}
)

var (
	css = map[string][]byte{}
)

var (
	staticFileServer = http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
)

func init() {
	for f := range htmlTemplates {
		htmlTemplates[f] = template.Must(template.ParseFiles("static/html/"+f, "static/html/layout.html"))
	}

	cssBundle := assets.Dir("static/css").MustAllFiles().MustFilter(
		assets.Concat(),
		assets.Fingerprint(),
		assets.WriteToDir("static/generated"),
	)
	for _, asset := range cssBundle.Assets() {
		b, err := ioutil.ReadAll(asset.Contents())
		if err != nil {
			panic(err)
		}
		css[asset.FileName()] = b
	}
}

func initRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/about", handlerAbout)
	mux.HandleFunc("/static/", handlerStatic)
	mux.HandleFunc("/", handlerCatchall)
}

func handlerIndex(rw http.ResponseWriter, req *http.Request) {
	tmpl := htmlTemplates["index.html"]

	params := struct {
		IncludeCSS []string
	}{}
	for f := range css {
		params.IncludeCSS = append(params.IncludeCSS, f)
	}

	err := tmpl.ExecuteTemplate(rw, "base", params)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func handlerAbout(rw http.ResponseWriter, req *http.Request) {
	tmpl := htmlTemplates["about.html"]

	params := struct {
		IncludeCSS []string
	}{}
	for f := range css {
		params.IncludeCSS = append(params.IncludeCSS, f)
	}

	err := tmpl.ExecuteTemplate(rw, "base", params)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func handlerStatic(rw http.ResponseWriter, req *http.Request) {
	b, ok := css[filepath.Base(req.URL.Path)]
	if ok {
		rw.Header().Set("Content-Type", "text/css")
		rw.Write(b)
		return
	}

	staticFileServer.ServeHTTP(rw, req)
}

func handlerCatchall(rw http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s\n", req.Method, req.URL)

	switch req.URL.Path {
	case "/":
		handlerIndex(rw, req)
	default:
		http.NotFound(rw, req)
	}
}
