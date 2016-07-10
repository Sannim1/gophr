package main

import(
    "net/http"
)

func AuthenticateRequest(responseWriter http.ResponseWriter, request *http.Request) {
    // Redirect the user to login if they're not authenticated
    authenticated := false

    if ! authenticated {
        http.Redirect(responseWriter, request, "/register", http.StatusFound)
    }
}
