package main

import (
	"html/template"
	"log"
	"net/http"
	"sync"
)

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(t.filename))
	})
	t.templ.Execute(w, nil)
}

func main() {
	r := newRoom()
	http.Handle("/", &templateHandler{filename: "templates/chat.html"})
	http.Handle("/room", r)

	go r.run()

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
