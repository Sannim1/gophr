package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func main() {
	// mux := http.NewServeMux()
	//
	// mux.HandleFunc("/", func(responseWriter http.ResponseWriter, request *http.Request) {
	//     RenderTemplate(responseWriter, request, "index/home", nil)
	// })
	//
	// assetServer := http.FileServer(http.Dir("assets/"))
	//
	// mux.Handle("/assets/", http.StripPrefix("/assets/", assetServer))
	//
	// http.ListenAndServe(":3000", mux)

	unauthenticatedRouter := NewRouter()
	unauthenticatedRouter.GET("/", HandleHome)
	unauthenticatedRouter.GET("/register", HandleNewUser)
	unauthenticatedRouter.POST("/register", HandleUserCreate)
	// unauthenticatedRouter.Handle("POST", "/register", HandleUserCreate)
	unauthenticatedRouter.ServeFiles("/assets/*filepath", http.Dir("assets/"))

	authenticatedRouter := NewRouter()
	authenticatedRouter.GET("/images/new", HandleNewImage)

	middleware := Middleware{}
	middleware.Add(unauthenticatedRouter)
	middleware.Add(http.HandlerFunc(AuthenticateRequest))
	middleware.Add(authenticatedRouter)

	http.ListenAndServe(":3000", middleware)
}

type NotFound struct{}

func (notFound NotFound) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {

}

func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = new(NotFound)

	return router
}
