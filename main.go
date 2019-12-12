package main

import (
	"flag"
	"log"
	"net/http"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/gomniauth/providers/github"
	"github.com/stretchr/gomniauth/providers/google"
)

type templateHnadler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHnadler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	t.templ.Execute(w, r)
}

func main() {
	var addr = flag.String("addr", ":8080", "Adres aplikacji internetowej")
	flag.Parse()
	gomniauth.SetSecurityKey("123456789")
	gomniauth.WithProviders(
		github.New("id", "key", "http://localhost:8080/auth/callback/github"),
		google.New("461495886139-h9jrtuo5eipstjnf84uruf1q0u3cfkrt.apps.googleusercontent.com", "QlF6B9ENdZ5CNO_y_e3pZmUo", "http://localhost:8080/auth/callback/google"),
		facebook.New("id", "key", "http://localhost:8080/auth/callback/facebook"),
	)
	r := newRoom()
	//r.tracer = trace.New(os.Stdout)
	http.Handle("/chat", MustAuth(&templateHnadler{filename: "chat.html"}))
	http.Handle("/login", &templateHnadler{filename: "login.html"})
	http.HandleFunc("/auth/", loginHandler)
	http.Handle("/room", r)
	go r.run()
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
