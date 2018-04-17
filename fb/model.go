package fb

import (
	"context"
	"errors"
	"os"
	"sync"
	"time"

	"cloud.google.com/go/firestore"
	"firebase.google.com/go"
	bv "github.com/Hunsin/beaver"
	"google.golang.org/api/option"
)

type Model struct {
	app *firebase.App
	src map[string]*firestore.DocumentSnapshot
	mu  sync.RWMutex
}

func New() (*Model, error) {
	cred := os.Getenv("FIREBASE_CREDENTIAL_FILE")
	if cred == "" {
		return nil, errors.New(`fb: environment variable "FIREBASE_CREDENTIAL_FILE" not set`)
	}

	a, err := firebase.NewApp(context.Background(), nil, option.WithCredentialsFile(cred))
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
				bv.Warn(err)
			}
		}
	}()

	return m, nil
}
