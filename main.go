package main

import (
	"crypto/tls"
	"flag"
	"net/http"

	bv "github.com/Hunsin/beaver"
	"github.com/Hunsin/me/fb"
	"github.com/Hunsin/me/server"
	"golang.org/x/crypto/acme/autocert"
)

type config struct {
	Domain []string `json:"domains"`
	Public string   `json:"public_dir"`
	View   string   `json:"view_file"`
	Cert   string   `json:"cert_dir"`
}

// valid checks each field of the config.
func (c config) valid() {
	if len(c.Domain) == 0 {
		bv.Fatal("me: at least a domain must provided")
	}
	if c.Public == "" {
		bv.Warn("me: public_dir not provided")
	}
	if c.View == "" {
		bv.Fatal("me: view_file must provided")
	}
	if c.Cert == "" {
		bv.Warn("me: cert_dir not provided")
	}
}

// httpToHTTPS redirects the connection to the security one.
func httpToHTTPS(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+r.RequestURI, http.StatusFound)
}

func main() {

	// get config filepath
	path := flag.String("c", "", "config filepath")
	flag.Parse()
	if *path == "" {
		bv.Fatal("me: config filepath needed")
	}

	// parse config
	var cfg config
	if err := bv.JSON(&cfg).Open(*path); err != nil {
		bv.Fatal(err)
	}
	cfg.valid()

	// initialize the fb.Model
	f, err := fb.New()
	if err != nil {
		bv.Fatal(err)
	}

	// create HTTP listener which handles Let's Encrypt challenge
	m := &autocert.Manager{
		Cache:      autocert.DirCache(cfg.Cert),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(cfg.Domain...),
	}
	go func() {
		h := m.HTTPHandler(http.HandlerFunc(httpToHTTPS))
		bv.Fatal(http.ListenAndServe(":http", h))
	}()

	// create main server and listen at HTTPS port
	s := &http.Server{
		Addr:      ":https",
		TLSConfig: &tls.Config{GetCertificate: m.GetCertificate},
		Handler:   server.New(f, cfg.Public, cfg.View),
	}
	bv.Fatal(s.ListenAndServeTLS("", ""))
}
