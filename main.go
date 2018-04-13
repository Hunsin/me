package main

import (
	"flag"
	"log"
	"net/http"
	"strings"

	"github.com/Hunsin/me/fb"
	"github.com/Hunsin/me/server"
	"golang.org/x/crypto/acme/autocert"
)

type domains []string

func (d *domains) Set(v string) error {
	*d = append(*d, v)
	return nil
}

func (d *domains) String() string {
	return strings.Join(*d, " ")
}

func main() {
	var ds domains
	flag.Var(&ds, "d", "Server domains")
	flag.Parse()
	if ds.String() == "" {
		log.Fatal("me: at least a domain flag must provided")
	}

	c, err := fb.New()
	if err != nil {
		log.Fatal(err)
	}

	srv := server.New(c)

	log.Fatal(http.Serve(autocert.NewListener(ds...), srv))
}
