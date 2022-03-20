package asciiweb

import (
	"html/template"
	"net/http"
)

type Error struct {
	Str  string
	Type int
}

// Function for processsing errors from the site and internal problems.
func Err(Str string, Status int, w http.ResponseWriter) {
	Info := Error{Str, Status}
	val, err := template.ParseFiles("templates/error.html")
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(Status)
	err = val.Execute(w, Info)
	if err != nil {
		http.Error(w, http.StatusText(500), http.StatusInternalServerError)
		return
	}
}
