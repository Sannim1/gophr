package main

import (
    "net/http"
    "github.com/julienschmidt/httprouter"
)

func HandleNewUser(responseWriter http.ResponseWriter, request *http.Request, params httprouter.Params)  {
    // Display new user page
    RenderTemplate(responseWriter, request, "users/new", nil)
}
