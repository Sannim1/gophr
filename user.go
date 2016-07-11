package main

import (
    "fmt"
    "encoding/json"
    "golang.org/x/crypto/bcrypt"
)

type User struct {
    ID              string
    Email           string
    Username        string
    HashedPassword  string
}

const (
    passwordLength  = 8
    hashCost        = 10
    userIDLength    = 16
)

func NewUser(username, email, password string) (User, []error) {
    errors := make([]error, 0)

    user := User {
        Email: email,
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

func (user User) String () string {
    jsonUser, _ := json.MarshalIndent(user, "", "  ")
    return fmt.Sprintf(string(jsonUser))
}
