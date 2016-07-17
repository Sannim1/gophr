package main

import (
	"fmt"
	"strings"
)

var globalUserStore UserStore

func init() {
	store, err := NewFileUserStore("./data/users.json")
	if err != nil {
		panic(fmt.Errorf("Error creating user store: %s", err))
	}

	globalUserStore = store
}

type UserStore interface {
	Find(string) (*User, error)
	FindByEmail(string) (*User, error)
	FindByUsername(string) (*User, error)
	Save(User) error
}

type FileUserStore struct {
	Users          map[string]User
	usernameLookup map[string]User
	emailLookup    map[string]User
	*FileStore
}

func NewFileUserStore(filename string) (*FileUserStore, error) {
	filestore, err := NewFileStore(filename)
	if err != nil {
		return nil, err
	}
	store := &FileUserStore{
		Users:     map[string]User{},
		FileStore: filestore,
	}

	err = store.ReadJSONInto(store)
	if err != nil {
		return nil, err
	}

	store.loadLookupMaps()

	return store, nil
}

func (store *FileUserStore) loadLookupMaps() {

	for _, user := range store.Users {
		store.addUserToLookupMaps(user)
	}
}

func (store *FileUserStore) addUserToLookupMaps(user User) {
	if len(store.usernameLookup) == 0 {
		store.usernameLookup = map[string]User{}
	}
	if len(store.emailLookup) == 0 {
		store.emailLookup = map[string]User{}
	}
	store.usernameLookup[strings.ToLower(user.Username)] = user
	store.emailLookup[strings.ToLower(user.Email)] = user
}

func (store FileUserStore) Save(user User) error {
	store.Users[user.ID] = user

	err := store.WriteJSONFrom(store)
	if err != nil {
		return err
	}

	store.addUserToLookupMaps(user)

	return nil
}

func (store FileUserStore) Find(userId string) (*User, error) {
	user, itemExists := store.Users[userId]

	if itemExists {
		return &user, nil
	}

	return nil, nil
}

func (store FileUserStore) FindByUsername(username string) (*User, error) {
	if username == "" {
		return nil, nil
	}

	user, itemExists := store.usernameLookup[strings.ToLower(username)]
	if itemExists {
		return &user, nil
	}

	return nil, nil
}

func (store FileUserStore) FindByEmail(email string) (*User, error) {
	if email == "" {
		return nil, nil
	}

	user, itemExists := store.emailLookup[strings.ToLower(email)]
	if itemExists {
		return &user, nil
	}

	return nil, nil
}
