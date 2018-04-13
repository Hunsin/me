package fb

import (
	"context"
	"log"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go"
)

type Model struct {
	app *firebase.App
	src map[string]*firestore.DocumentSnapshot
	mu  sync.RWMutex
}

func New() (*Model, error) {
	a, err := firebase.NewApp(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	m := &Model{a, make(map[string]*firestore.DocumentSnapshot), sync.RWMutex{}}
	if err := m.update(); err != nil {
		return nil, err
	}

	go func() {
		for {
			time.Sleep(time.Minute)
			if err := m.update(); err != nil {
				log.Panicln(err)
			}
		}
	}()

	return m, nil
}
