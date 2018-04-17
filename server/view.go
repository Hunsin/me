package server

import (
	"io/ioutil"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

// setView minifies the file v and parses the template.
func (s *Server) setView(v string) *Server {
	if v == "" {
		panic(`server: view file must provided`)
	}

	b, err := ioutil.ReadFile(v)
	if err != nil {
		panic(err)
	}

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFunc("text/javascript", js.Minify)

	n, err := m.Bytes("text/html", b)
	if err != nil {
		panic(err)
	}

	_, err = s.tmp.Parse(string(n))
	if err != nil {
		panic(err)
	}

	return s
}
