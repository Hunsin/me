package server

import (
	"net/http"
	"os"
	"strings"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

func direct(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	lang := strings.ToLower(r.Header.Get("Accept-Language"))
	if strings.Contains(lang, "tw") || strings.Contains(lang, "hk") {
		http.Redirect(w, r, "/tw", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/en", http.StatusFound)
}

func (s *Server) render(w http.ResponseWriter, r *http.Request) {
	if p, ok := w.(http.Pusher); ok {
		p.Push("/static/style.css", &http.PushOptions{Method: "GET"})
	}

	lang := r.URL.Path[1:]

	v, err := s.mdl.Values(lang)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	s.tmp.Execute(w, v)
}

func (s *Server) setMux() *Server {
	static := os.Getenv("ME_PUBLIC_DIR")
	if static == "" {
		panic(`server: environment variable "ME_PUBLIC_DIR" not set`)
	}

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/javascript", js.Minify)

	fs := http.FileServer(http.Dir(static))

	s.mux.Handle("/static/", http.StripPrefix("/static/", m.Middleware(fs)))

	s.mux.HandleFunc("/en", s.render)

	s.mux.HandleFunc("/tw", s.render)

	s.mux.HandleFunc("/", direct)

	return s
}
