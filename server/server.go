package server

import (
	"html/template"
	"net/http"

	"github.com/Hunsin/me/fb"
)

// A Server is a http Handler which carries a path multiplexor
// and Firebase client.
type Server struct {
	mdl *fb.Model
	mux *http.ServeMux
	tmp *template.Template
}

// ServeHTTP implements the http Handler interface.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// New returns a pointer to an initialized Server. Static files under
// the dir is served under URL path "/public". The template will parse
// the view file.
func New(m *fb.Model, dir, view string) *Server {
	if m == nil {
		panic("server: A nil fb.Model is applied")
	}

	s := &Server{
		mux: http.NewServeMux(),
		mdl: m,
		tmp: template.New("view"),
	}
	return s.setMux(dir).setView(view)
}
