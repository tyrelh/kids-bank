package renderer

import (
	"net/http"
	"text/template"
)

func RenderIndex(w http.ResponseWriter, r *http.Request) {
	templ := template.Must(template.ParseFiles("index.html"))
	err := templ.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
