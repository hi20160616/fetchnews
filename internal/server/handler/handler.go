package handler

import (
	"context"
	"log"
	"net/http"
	"regexp"

	pb "github.com/hi20160616/fetchnews-api/proto/v1"
	"github.com/hi20160616/fetchnews/configs"
	"github.com/hi20160616/fetchnews/internal/server/render"
	"github.com/hi20160616/fetchnews/internal/service"
	tmpl "github.com/hi20160616/fetchnews/templates"
)

var validPath = regexp.MustCompile("^/(list|article|search)/(.*?)$")

// makeHandler invoke fn after path valided, and arrange args from url to object: `&render.Page{}`
func makeHandler(fn func(http.ResponseWriter, *http.Request, *render.Page)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
		}
		fn(w, r, &render.Page{})
	}
}

// GetHandler is a handler merger and a router for mutipl handler
func GetHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// The "/" pattern matches everything, so we need to check
		// that we're at the root here.
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			return
		}
		homeHandler(w, req)
	})
	mux.Handle("/s/", http.StripPrefix("/s/", http.FileServer(http.FS(tmpl.FS))))
	mux.HandleFunc("/list/", makeHandler(listArticlesHandler))
	mux.HandleFunc("/article/", makeHandler(getArticleHandler))
	mux.HandleFunc("/search/", makeHandler(searchArticlesHandler))
	return mux
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	render.Derive(w, "home", &render.Page{Title: "Home", Data: configs.Data.MS})
}

func listArticlesHandler(w http.ResponseWriter, r *http.Request, p *render.Page) {
	msTitle := r.URL.Query().Get("v")
	ds, err := service.ListArticles(context.Background(), &pb.ListArticlesRequest{}, msTitle)
	if err != nil {
		log.Println(err)
	}
	p.Title = msTitle
	p.Data = ds.Articles
	render.Derive(w, "list", p)
}

func getArticleHandler(w http.ResponseWriter, r *http.Request, p *render.Page) {
	msTitle := r.URL.Query().Get("website")
	id := r.URL.Query().Get("id")
	a, err := service.GetArticle(context.Background(), &pb.GetArticleRequest{Id: id}, msTitle)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	p.Data = a
	render.Derive(w, "article", p)
}

func searchArticlesHandler(w http.ResponseWriter, r *http.Request, p *render.Page) {
	kw := r.URL.Query().Get("v")
	as, err := service.SearchArticles(context.Background(), &pb.SearchArticlesRequest{Keyword: kw})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	p.Data = as
	render.Derive(w, "search", p)
}
