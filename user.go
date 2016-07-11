package main

import "golang.org/x/crypto/bcrypt"

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

func NewUser(username, email, password string) (User, error) {
    user := User {
        Email: email,
        Username: username,
    }

    if username == "" {
        return user, errNoUsername
    }

    if email == "" {
        return user, errNoEmail
    }

    if password == "" {
        return user, errNoPassword
    }

    if len(password) < passwordLength {
        return user, errPasswordTooShort
    }

    existingUser, err := globalUserStore.FindByUsername(username)
    if err != nil {
        return user, err
    }
    if existingUser != nil {
        return user, errUsernameExists
    }

    existingUser, err = globalUserStore.FindByEmail(email)
    if err != nil {
        return user, err
    }
    if existingUser != nil {
        return user, errEmailExists
    }

    hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), hashCost)

    user.HashedPassword = string(hashedPassword)
    user.ID = GenerateID("usr", userIDLength)

    return user, err
}
