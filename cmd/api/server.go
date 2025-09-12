package main

import (
	"net/http"
	"restapi/internal/api/router"
)

func main() {

	port := "3000"
	println("Server is running on port " + port)

	// Use the router to get configured mux
	mux := router.Router()

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	err := server.ListenAndServe()
	if err != nil {
		panic(err)

	}

}

// type Middleware func(http.Handler) http.Handler

// func applyMiddlewares(handler http.Handler, middlewares ...Middleware) http.Handler {

// 	for _, middleware := range middlewares {
// 		handler = middleware(handler)
// 	}
// 	return handler

// }
