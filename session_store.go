package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

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
	filename string
	Sessions map[string]Session
}

func NewFileSessionStore(filename string) (*FileSessionStore, error) {
	store := &FileSessionStore{
		Sessions: map[string]Session{},
		filename: filename,
	}

	contents, err := ioutil.ReadFile(filename)

	if err != nil {
		if os.IsNotExist(err) {
			// ignore error if file does not exist, means its a new session
			return store, nil
		}

		return nil, err
	}

	err = json.Unmarshal(contents, store)
	if err != nil {
		return nil, err
	}

	return store, err
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

	contents, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(store.filename, contents, 0660)
}

func (store *FileSessionStore) Delete(session *Session) error {
	delete(store.Sessions, session.ID)

	contents, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(store.filename, contents, 0660)
}
