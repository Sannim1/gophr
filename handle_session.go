package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func HandleSessionNew(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	next := request.URL.Query().Get("next")

	templateData := map[string]interface{}{
		"Next": next,
	}
	RenderTemplate(responseWriter, request, "sessions/new", templateData)
}

func HandleSessionCreate(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	username := request.FormValue("username")
	password := request.FormValue("password")
	next := request.FormValue("next")

	templateData := map[string]interface{}{
		"Next": next,
	}

	user, err := FindUser(username, password)
	if err != nil {
		if IsValidationError(err) {
			templateData["Error"] = err
			templateData["User"] = user

			RenderTemplate(responseWriter, request, "sessions/new", templateData)

			return
		}
		panic(err)
	}

	session := FindOrCreateSession(responseWriter, request)
	session.UserID = user.ID
	err = globalSessionStore.Save(session)
	if err != nil {
		panic(err)
	}

	if next == "" {
		next = "/"
	}

	redirectURL := next + "?flash=Signed+in"

	http.Redirect(responseWriter, request, redirectURL, http.StatusFound)
}

func HandleSessionDestroy(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	session := RequestSession(request)

	if session != nil {
		err := globalSessionStore.Delete(session)
		if err != nil {
			panic(err)
		}
	}
	RenderTemplate(responseWriter, request, "sessions/destroy", nil)
}
