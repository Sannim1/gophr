package main

import (
	"encoding/json"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID             string
	Email          string
	Username       string
	HashedPassword string
}

const (
	passwordLength = 8
	hashCost       = 10
	userIDLength   = 16
)

func NewUser(username, email, password string) (User, []error) {

	var errors []error

	user := User{
		Email:    email,
		Username: username,
	}

	if username == "" {
		errors = append(errors, errNoUsername)
	}

	if email == "" {
		errors = append(errors, errNoEmail)
	}

	if password == "" {
		errors = append(errors, errNoPassword)
	}

	if len(password) < passwordLength {
		errors = append(errors, errPasswordTooShort)
	}

	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		errors = append(errors, err)
	}
	if existingUser != nil {
		errors = append(errors, errUsernameExists)
	}

	existingUser, err = globalUserStore.FindByEmail(email)
	if err != nil {
		errors = append(errors, err)
	}
	if existingUser != nil {
		errors = append(errors, errEmailExists)
	}

	if len(errors) > 0 {
		return user, errors
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)

	user.HashedPassword = string(hashedPassword)
	user.ID = GenerateID("usr", userIDLength)

	if err != nil {
		errors = append(errors, err)
	}

	return user, errors
}

func FindUser(username, password string) (*User, error) {
	userToBeFound := &User{
		Username: username,
	}

	existingUser, err := globalUserStore.FindByUsername(username)
	if err != nil {
		return userToBeFound, err
	}

	if existingUser == nil {
		return userToBeFound, errCredentialsIncorrect
	}

	comparePasswords := bcrypt.CompareHashAndPassword(
		[]byte(existingUser.HashedPassword),
		[]byte(password),
	)

	if comparePasswords != nil {
		return userToBeFound, errCredentialsIncorrect
	}

	return existingUser, nil
}

func UpdateUser(user *User, email, currentPassword, newPassword string) (User, error) {
	userToBeUpdated := *user
	userToBeUpdated.Email = email

	// check if the email exists
	existingUser, err := globalUserStore.FindByEmail(email)
	if err != nil {
		return userToBeUpdated, err
	}

	if existingUser != nil && existingUser.ID != user.ID {
		return userToBeUpdated, errEmailExists
	}

	// At this point we can update the email address
	user.Email = email

	// No current password? Or current password equals new password? Dont try updating the password
	if currentPassword == "" || currentPassword == newPassword {
		return userToBeUpdated, nil
	}

	comparePasswords := bcrypt.CompareHashAndPassword(
		[]byte(existingUser.HashedPassword),
		[]byte(currentPassword),
	)
	if comparePasswords != nil {
		return userToBeUpdated, errPasswordIncorrect
	}

	if newPassword == "" {
		return userToBeUpdated, errNoPassword
	}

	if len(newPassword) < passwordLength {
		return userToBeUpdated, errPasswordTooShort
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), hashCost)
	user.HashedPassword = string(hashedPassword)

	return userToBeUpdated, err
}

func (user User) String() string {
	jsonUser, _ := json.MarshalIndent(user, "", "  ")
	return fmt.Sprintf(string(jsonUser))
}
