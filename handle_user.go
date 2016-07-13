package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func HandleNewUser(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	// Display new user page
	RenderTemplate(responseWriter, request, "users/new", nil)
}

func HandleUserCreate(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	// Process creating a new user
	user, errors := NewUser(
		request.FormValue("username"),
		request.FormValue("email"),
		request.FormValue("password"),
	)

	templateData := map[string]interface{}{
		"User": user,
	}

	if len(errors) > 0 {
		templateErrors := make([]string, len(errors))
		// loop through the error slice to ensure that all errors are validation errors
		// panic, if otherwise
		for _, err := range errors {
			if !IsValidationError(err) {
				panic(err)
			}
			templateErrors = append(templateErrors, err.Error())
		}

		templateData["Errors"] = templateErrors
		RenderTemplate(responseWriter, request, "users/new", templateData)
		return
	}

	err := globalUserStore.Save(user)
	if err != nil {
		panic(err)
	}

	http.Redirect(responseWriter, request, "/?flash=User+created", http.StatusFound)
}
