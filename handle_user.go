package main

import (
    "net/http"
    "github.com/julienschmidt/httprouter"
)

func HandleNewUser(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params)  {
    // Display new user page
    RenderTemplate(responseWriter, request, "users/new", nil)
}

func HandleUserCreate(responseWriter http.ResponseWriter, request *http.Request, _ httprouter.Params)  {
    // Process creating a new user
    user, err := NewUser(
        request.FormValue("username"),
        request.FormValue("email"),
        request.FormValue("password"),
    )

    templateData := map[string]interface{}{
        "User": user,
    }

    if err != nil {
        if IsValidationError(err) {
            templateData["Error"] = err.Error()
            RenderTemplate(responseWriter, request, "users/new", templateData)

            return
        }
        panic(err)
    }

    err = globalUserStore.Save(user)
    if err != nil {
        panic(err)
    }

    http.Redirect(responseWriter, request, "/?flash=User+created", http.StatusFound)
}
