package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// HandleHome displays the home page
func HandleHome(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	images, err := globalImageStore.FindAll(0)
	if err != nil {
		panic(err)
	}

	templateData := map[string]interface{}{
		"Images": images,
	}

	RenderTemplate(responseWriter, request, "index/home", templateData)
}
