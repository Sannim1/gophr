package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func HandleHome(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params) {
	// Display home page
	RenderTemplate(responseWriter, request, "index/home", nil)
}
