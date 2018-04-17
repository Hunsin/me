package server

import (
	"net/http"
	"os"
	"strings"

	bv "github.com/Hunsin/beaver"
)

// direct redirects the request to different languages depending on
// header "Accept-Language".
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

// render renders the HTML page depending on different language
// specified by request path.
func (s *Server) render(w http.ResponseWriter, r *http.Request) {
	if p, ok := w.(http.Pusher); ok {
		p.Push("/public/css/fontawesome-all.min.css", &http.PushOptions{Method: "GET"})
		p.Push("/public/webfonts/fa-brands-400.woff2", &http.PushOptions{Method: "GET"})
	}

	lang := r.URL.Path[1:]

	v, err := s.mdl.Values(lang)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	if err = s.tmp.Execute(w, v); err != nil {
		bv.Error(err)
	}
}

// serMux initialize the s.mux. If pub is specified, it serves static
// files in the directory under URL path "/public".
func (s *Server) setMux(pub string) *Server {
	s.mux.HandleFunc("/en", s.render)

	s.mux.HandleFunc("/tw", s.render)

	s.mux.HandleFunc("/", direct)

	if pub != "" {
		if info, err := os.Stat(pub); err != nil {
			panic(err)
		} else if !info.IsDir() {
			panic(`server: pub must be a directory`)
		}

		fs := http.FileServer(http.Dir(pub))
		s.mux.Handle("/public/", http.StripPrefix("/public/", fs))
	}

	return s
}
