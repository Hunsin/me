package fb

import (
	"context"
	"errors"
	"html/template"

	"cloud.google.com/go/firestore"
)

// An Image specifies the image URL and the base64 encoded placeholder
// while downloading the file.
type Image struct {
	Base64 template.URL `firestore:"data"`
	URL    string       `firestore:"url"`
}

type Content struct {
	Title string        `firestore:"title"`
	Name  string        `firestore:"name"`
	Info  template.HTML `firestore:"desc"`
}

type Value struct {
	Avatar, Background Image
	Content            Content
}

func (m *Model) update() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ctx := context.Background()

	c, err := m.app.Firestore(ctx)
	if err != nil {
		return err
	}

	d, err := c.Collection("render").Documents(ctx).GetAll()
	if err != nil {
		return err
	}

	m.src = make(map[string]*firestore.DocumentSnapshot)
	for _, snap := range d {
		m.src[snap.Ref.ID] = snap
	}

	return nil
}

func (m *Model) Values(lang string) (v Value, err error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if _, ok := m.src[lang]; !ok {
		return v, errors.New("fb: request language doesn't exist")
	}

	if err = m.src["avatar"].DataTo(&v.Avatar); err != nil {
		return
	}

	if err = m.src["background"].DataTo(&v.Background); err != nil {
		return
	}

	if err = m.src[lang].DataTo(&v.Content); err != nil {
		return
	}

	return
}
