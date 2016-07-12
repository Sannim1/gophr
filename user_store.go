package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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
	filename       string
	usernameLookup map[string]User
	emailLookup    map[string]User
	Users          map[string]User
}

func NewFileUserStore(filename string) (*FileUserStore, error) {
	store := &FileUserStore{
		filename: filename,
		Users:    map[string]User{},
	}

	contents, err := ioutil.ReadFile(filename)
	if err != nil {
		// if file does not exist, means a new UserStore, so that's OK
		if os.IsNotExist(err) {
			return store, nil
		}

		return nil, err
	}

	err = json.Unmarshal(contents, store)
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

	contents, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(store.filename, contents, 0660)
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
