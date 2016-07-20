package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HandleImageNew displays a new image form
func HandleImageNew(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	RenderTemplate(responseWriter, request, "images/new", nil)
}

// HandleImageCreate creates a new image
func HandleImageCreate(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	if request.FormValue("url") != "" {
		HandleImageCreateFromURL(responseWriter, request)
		return
	}

	HandleImageCreateFromFile(responseWriter, request)
}

// HandleImageCreateFromURL creates a new image from a specified image URL
func HandleImageCreateFromURL(responseWriter http.ResponseWriter, request *http.Request) {

	templateData := make(map[string]interface{})

	user := RequestUser(request)

	image := NewImage(user)
	image.Description = request.FormValue("description")

	err := image.CreateFromURL(request.FormValue("url"))

	if err != nil {
		if IsValidationError(err) {
			templateData["Error"] = err
			templateData["ImageURL"] = request.FormValue("url")
			templateData["Image"] = image

			RenderTemplate(responseWriter, request, "images/new", templateData)

			return
		}

		panic(err)
	}

	redirectURL := "/?flash=Image+Uploaded+Successfully"
	http.Redirect(responseWriter, request, redirectURL, http.StatusFound)
}

// HandleImageCreateFromFile creates a new image from an uploaded image file
func HandleImageCreateFromFile(responseWriter http.ResponseWriter, request *http.Request) {

	templateData := make(map[string]interface{})
	user := RequestUser(request)

	image := NewImage(user)
	image.Description = request.FormValue("description")
	templateData["Image"] = image

	file, headers, err := request.FormFile("file")

	// Check if a file was actually uploaded
	if file == nil {
		templateData["Error"] = errNoImage
		RenderTemplate(responseWriter, request, "images/new", templateData)

		return
	}

	// A file was uploaded but an error occurred
	if err != nil {
		panic(err)
	}

	defer file.Close()

	err = image.CreateFromFile(file, headers)
	if err != nil {
		templateData["Error"] = err
		RenderTemplate(responseWriter, request, "images/new", templateData)

		return
	}

	redirectURL := "/?flash=Image+Uploaded+Successfully"
	http.Redirect(responseWriter, request, redirectURL, http.StatusFound)
}
