package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HandleUserShow displays a page containing the user's images
func HandleUserShow(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {

	userID := params.ByName("userID")

	user, err := globalUserStore.Find(userID)
	if err != nil {
		panic(err)
	}

	// 404, user does not exists
	if user == nil {
		http.NotFound(responseWriter, request)
		return
	}

	images, err := globalImageStore.FindAllByUser(user, 0)
	if err != nil {
		panic(err)
	}

	templateData := map[string]interface{}{
		"Images": images,
		"User":   user,
	}

	RenderTemplate(responseWriter, request, "users/show", templateData)
}

// HandleNewUser displays a page to register new users
func HandleNewUser(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	// Display new user page
	RenderTemplate(responseWriter, request, "users/new", nil)
}

// HandleUserCreate creates a new user from submitted parameters
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

	session := NewSession(responseWriter)
	session.UserID = user.ID

	err = globalSessionStore.Save(session)
	if err != nil {
		panic(err)
	}

	redirectURL := "/?flash_message=User+created&msg_type=success"
	http.Redirect(responseWriter, request, redirectURL, http.StatusFound)
}

// HandleUserEdit displays a page from which a user can update their details
func HandleUserEdit(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	user := RequestUser(request)
	templateData := map[string]interface{}{
		"User": user,
	}
	RenderTemplate(responseWriter, request, "users/edit", templateData)
}

// HandleUserUpdate persists a user's updated details
func HandleUserUpdate(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	currentUser := RequestUser(request)
	email := request.FormValue("email")
	currentPassword := request.FormValue("currentPassword")
	newPassword := request.FormValue("newPassword")

	user, err := UpdateUser(currentUser, email, currentPassword, newPassword)
	templateData := map[string]interface{}{
		"User": user,
	}
	if err != nil {
		if IsValidationError(err) {
			templateData["Error"] = err.Error()
			RenderTemplate(responseWriter, request, "users/edit", templateData)

			return
		}
		panic(err)
	}

	err = globalUserStore.Save(*currentUser)
	if err != nil {
		panic(err)
	}

	redirectURL := "/account?flash_message=User+updated&msg_type=success"
	http.Redirect(responseWriter, request, redirectURL, http.StatusFound)
}
