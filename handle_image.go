package main

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HandleImageNew displays a new image form
func HandleImageNew(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params) {
	RenderTemplate(responseWriter, request, "images/new", nil)
}

// HandleImageShow displays a single image on it's own page
func HandleImageShow(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	imageID := params.ByName("imageID")
	image, err := globalImageStore.Find(imageID)
	if err != nil {
		panic(err)
	}

	// check image exists
	if image == nil {
		http.NotFound(responseWriter, request)
		return
	}

	// Get image owner
	user, err := globalUserStore.Find(image.UserID)
	if err != nil {
		panic(err)
	}

	if user == nil {
		panic(fmt.Errorf("Could not find user %s for image:%s", image.UserID, image.ID))
	}

	templateData := map[string]interface{}{
		"User":  user,
		"Image": image,
	}

	RenderTemplate(responseWriter, request, "images/show", templateData)
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
