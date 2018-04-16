package server

import (
	"io/ioutil"
	"os"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
)

func (s *Server) setView() *Server {
	v := os.Getenv("ME_VIEW")
	if v == "" {
		panic(`server: environment variable "ME_VIEW" not set`)
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
