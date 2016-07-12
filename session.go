package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	// keep users logged in for 3 days
	sessionLength     = 3 * 24 * time.Hour
	sessionCookieName = "GophrSession"
	sessionIDLength   = 20
)

type Session struct {
	ID     string
	UserID string
	Expiry time.Time
}

func NewSession(responseWriter http.ResponseWriter) *Session {
	expiry := time.Now().Add(sessionLength)

	session := &Session{
		ID:     GenerateID("sess", sessionIDLength),
		Expiry: expiry,
	}

	cookie := http.Cookie{
		Name:    sessionCookieName,
		Value:   session.ID,
		Expires: expiry,
	}

	http.SetCookie(responseWriter, &cookie)

	return session
}

func RequestSession(request *http.Request) *Session {
	cookie, err := request.Cookie(sessionCookieName)
	if err != nil {
		return nil
	}

	session, err := globalSessionStore.Find(cookie.Value)
	if err != nil {
		panic(err)
	}

	if session == nil {
		return nil
	}

	if session.Expired() {
		globalSessionStore.Delete(session)
		return nil
	}

	return session
}

func RequestUser(request *http.Request) *User {
	session := RequestSession(request)

	if session == nil {
		return nil
	}

	user, err := globalUserStore.Find(session.UserID)
	if err != nil {
		panic(err)
	}

	return user
}

func RequireLogin(responseWriter http.ResponseWriter, request *http.Request) {
	// let the request pass if we've got a user
	user := RequestUser(request)

	if user != nil {
		return
	}

	query := url.Values{}
	query.Add("next", url.QueryEscape(request.URL.String()))

	redirectPath := fmt.Sprintf("/login?%s", query.Encode())

	http.Redirect(responseWriter, request, redirectPath, http.StatusFound)
}

func (session *Session) Expired() bool {
	return session.Expiry.Before(time.Now())
}
