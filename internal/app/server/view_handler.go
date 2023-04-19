package server

import (
	"html/template"
	"net/http"
)

func viewHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../../template/sse.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
	tmpl.Execute(w, nil)
}
