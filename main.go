package main

import (
	"flag"
	"html/template"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/vit0rr/go-chat/trace"
)

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the application.")
	flag.Parse()
	r := newRoom()
	r.tracer = trace.New(os.Stdout)

	http.Handle("/chat", MustAuth(&templateHandler{filename: "templates/chat.html"}))
	http.Handle("/login", &templateHandler{filename: "templates/login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)

	go r.run()

	log.Printf("Starting web server on %s%s", "http://localhost", *addr)

	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(t.filename))
	})
	t.templ.Execute(w, r)
}
