package main

import "fmt"

var globalSessionStore SessionStore

func init() {
	sessionStore, err := NewFileSessionStore("./data/sessions.json")
	if err != nil {
		panic(fmt.Errorf("Error creating session store: %s", err))
	}

	globalSessionStore = sessionStore
}

type SessionStore interface {
	Find(string) (*Session, error)
	Save(*Session) error
	Delete(*Session) error
}

type FileSessionStore struct {
	Sessions map[string]Session
	*FileStore
}

func NewFileSessionStore(filename string) (*FileSessionStore, error) {
	filestore, err := NewFileStore(filename)
	if err != nil {
		return nil, err
	}

	store := &FileSessionStore{
		Sessions:  map[string]Session{},
		FileStore: filestore,
	}

	err = store.ReadJSONInto(store)
	if err != nil {
		return nil, err
	}

	return store, nil
}

func (store *FileSessionStore) Find(sessionID string) (*Session, error) {
	session, itemExists := store.Sessions[sessionID]
	if !itemExists {
		return nil, nil
	}

	return &session, nil
}

func (store *FileSessionStore) Save(session *Session) error {
	store.Sessions[session.ID] = *session

	return store.WriteJSONFrom(store)
}

func (store *FileSessionStore) Delete(session *Session) error {
	delete(store.Sessions, session.ID)

	return store.WriteJSONFrom(store)
}
