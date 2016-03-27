package www

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"path/filepath"
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
	css = map[string][]byte{}
)

var (
	staticFileServer = http.StripPrefix("/static/", http.FileServer(http.Dir("static")))
)

func Serve(listenAddr string) error {
	for f := range htmlTemplates {
		htmlTemplates[f] = template.Must(template.ParseFiles("static/html/"+f, "static/html/layout.html"))
	}

	cssBundle := assets.Dir("static/css").MustAllFiles().MustFilter(
		assets.Concat(),
		assets.Fingerprint(),
	)
	for _, asset := range cssBundle.Assets() {
		b, err := ioutil.ReadAll(asset.Contents())
		if err != nil {
			return err
		}
		css[asset.FileName()] = b
	}

	initRoutes(http.DefaultServeMux)
	s := &http.Server{
		Addr:           listenAddr,
		Handler:        http.DefaultServeMux,
		ReadTimeout:    30 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	return s.ListenAndServe()
}

func initRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/about", handlerAbout)
	mux.HandleFunc("/static/", handlerStatic)
	mux.HandleFunc("/p/", handlerBlogPost)
	mux.HandleFunc("/", handlerCatchall)
}

type sharedTemplateParams struct {
	IncludeCSS []string
}

func handlerIndex(rw http.ResponseWriter, req *http.Request) {
	tmpl := htmlTemplates["index.html"]

	params := struct{ sharedTemplateParams }{}
	for f := range css {
		params.IncludeCSS = append(params.IncludeCSS, f)
	}

	err := tmpl.ExecuteTemplate(rw, "base", params)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}
}

func handlerAbout(rw http.ResponseWriter, req *http.Request) {
	params := struct{ sharedTemplateParams }{}
	for f := range css {
		params.IncludeCSS = append(params.IncludeCSS, f)
	}

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
	for f := range css {
		params.IncludeCSS = append(params.IncludeCSS, f)
	}

	err := htmlTemplates["post.html"].ExecuteTemplate(rw, "base", params)
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
