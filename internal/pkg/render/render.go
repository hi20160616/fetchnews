package render

import (
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/hi20160616/fetchnews/config"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Page struct {
	Title string
	Data  interface{}
}

var templates = template.New("")

func init() {
	templates.Funcs(template.FuncMap{
		"summary":       Summary,
		"smartTime":     SmartTime,
		"smartLongTime": SmartLongTime,
	})
	// tmplPath := filepath.Join("../../../templates", "default") // for TestValidReq
	tmplPath := filepath.Join(config.Data.WebServer.Tmpl, "default")
	pattern := filepath.Join(tmplPath, "*.html")
	templates = template.Must(templates.ParseGlob(pattern))
}

func Derive(w http.ResponseWriter, tmpl string, p *Page) {
	if err := templates.ExecuteTemplate(w, tmpl+".html", p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("err template: %s.html\n\terror: %#v", tmpl, err)
	}
}

func Summary(des string) string {
	dRune := []rune(des)
	if len(dRune) <= 300 {
		return des
	}
	return string(dRune[:300])
}

func parseWithZone(t time.Time) time.Time {
	loc := time.FixedZone("UTC", 8*60*60)
	return t.In(loc)

}

func SmartTime(t *timestamppb.Timestamp) string {
	return parseWithZone(t.AsTime()).Format("[15:04][01.02]")
}

func SmartLongTime(t time.Time) string {
	return parseWithZone(t).String()
}
