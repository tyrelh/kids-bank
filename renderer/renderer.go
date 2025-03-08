package renderer

import (
	"embed"
	"net/http"
	"text/template"
)

//go:embed *.html
var templates embed.FS

func RenderIndex(w http.ResponseWriter, r *http.Request) {
	templ, err := template.ParseFS(templates, "index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = templ.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
