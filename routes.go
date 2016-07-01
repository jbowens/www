package www

import (
	"html/template"
	"log"
	"net/http"
	"path"
	"time"

	"github.com/jbowens/assets"
	"github.com/jbowens/www/blog"
)

var (
	htmlTemplates = map[string]*template.Template{
		"index.html": nil,
		"about.html": nil,
		"post.html":  nil,
	}
)

var (
	css    assets.Bundle
	fonts  assets.Bundle
	images assets.Bundle
)

func Serve(listenAddr string) error {
	for f := range htmlTemplates {
		htmlTemplates[f] = template.Must(template.ParseFiles("static/html/"+f, "static/html/layout.html"))
	}

	var err error
	css, err = assets.Development("static/css")
	if err != nil {
		return err
	}
	fonts, err = assets.Development("static/fonts")
	if err != nil {
		return err
	}
	images, err = assets.Development("static/images")
	if err != nil {
		return err
	}

	initRoutes(http.DefaultServeMux, css)
	s := &http.Server{
		Addr:           listenAddr,
		Handler:        http.DefaultServeMux,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return s.ListenAndServe()
}

func initRoutes(mux *http.ServeMux, css assets.Bundle) {
	innerMux := http.NewServeMux()
	innerMux.HandleFunc("/about", handlerAbout)
	innerMux.HandleFunc("/p/", handlerBlogPost)
	innerMux.Handle("/static/css/", http.StripPrefix("/static/css/", css))
	innerMux.Handle("/static/fonts/", http.StripPrefix("/static/fonts/", fonts))
	innerMux.Handle("/static/images/", http.StripPrefix("/static/images/", images))
	innerMux.HandleFunc("/", handlerCatchall)

	mux.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
		log.Printf("%s %s\n", req.Method, req.URL)
		innerMux.ServeHTTP(rw, req)
	})
}

type sharedTemplateParams struct {
	IncludeCSS []string
}

func handlerIndex(rw http.ResponseWriter, req *http.Request) {
	params := struct {
		sharedTemplateParams
		Posts []blog.Post
	}{Posts: blog.Posts()}
	params.IncludeCSS = css.RelativePaths()

	err := htmlTemplates["index.html"].ExecuteTemplate(rw, "base", params)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func handlerAbout(rw http.ResponseWriter, req *http.Request) {
	params := struct{ sharedTemplateParams }{}
	params.IncludeCSS = css.RelativePaths()

	err := htmlTemplates["about.html"].ExecuteTemplate(rw, "base", params)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func handlerBlogPost(rw http.ResponseWriter, req *http.Request) {
	p, ok := blog.PostByID(path.Base(req.URL.Path))
	if !ok {
		http.NotFound(rw, req)
		return
	}

	params := struct {
		sharedTemplateParams
		Post blog.Post
	}{Post: p}
	params.IncludeCSS = css.RelativePaths()

	err := htmlTemplates["post.html"].ExecuteTemplate(rw, "base", params)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func handlerCatchall(rw http.ResponseWriter, req *http.Request) {
	switch req.URL.Path {
	case "/":
		handlerIndex(rw, req)
	default:
		http.NotFound(rw, req)
	}
}
