package main

import (
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
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
	unauthenticatedRouter.GET("/login", HandleSessionNew)
	unauthenticatedRouter.POST("/login", HandleSessionCreate)
	unauthenticatedRouter.ServeFiles("/assets/*filepath", http.Dir("assets/"))

	authenticatedRouter := NewRouter()
	authenticatedRouter.GET("/images/new", HandleNewImage)
	authenticatedRouter.GET("/sign-out", HandleSessionDestroy)
	authenticatedRouter.GET("/account", HandleUserEdit)
	authenticatedRouter.POST("/account", HandleUserUpdate)

	middleware := Middleware{}
	middleware.Add(unauthenticatedRouter)
	// middleware.Add(http.HandlerFunc(AuthenticateRequest))
	middleware.Add(http.HandlerFunc(RequireLogin))
	middleware.Add(authenticatedRouter)

	log.Fatal(http.ListenAndServe(":3000", middleware))
}

type NotFound struct{}

func (notFound NotFound) ServeHTTP(responseWriter http.ResponseWriter, request *http.Request) {

}

func NewRouter() *httprouter.Router {
	router := httprouter.New()
	router.NotFound = new(NotFound)

	return router
}
