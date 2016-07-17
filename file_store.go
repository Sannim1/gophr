package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

// FileStore reads data objects from a file and writes data objects to a file
type FileStore struct {
	filename string
}

// NewFileStore creates a FileStore object that uses a specified file
func NewFileStore(filename string) (*FileStore, error) {
	store := &FileStore{
		filename: filename,
	}

	return store, nil
}

// ReadJSONInto reads the contents of the store's file as JSON into the specified structure
func (store *FileStore) ReadJSONInto(structure interface{}) error {
	contents, err := ioutil.ReadFile(store.filename)
	if err != nil {
		if os.IsNotExist(err) {
			// ignore error if file does not exist
			return nil
		}
		return err
	}

	err = json.Unmarshal(contents, structure)
	if err != nil {
		return err
	}

	return nil
}

// WriteJSONFrom writes the specified structure as JSON into the store
func (store *FileStore) WriteJSONFrom(structure interface{}) error {
	contents, err := json.MarshalIndent(structure, "", "  ")

	if err != nil {
		return err
	}

	return ioutil.WriteFile(store.filename, contents, 0660)
}
