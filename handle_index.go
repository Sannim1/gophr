package main

import(
    "net/http"
    "github.com/julienschmidt/httprouter"
)

func HandleHome(responseWriter http.ResponseWriter, request * http.Request, params httprouter.Params) {
    // Display home page
    RenderTemplate(responseWriter, request, "index/home", nil)
}
