package main

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func teachersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		handlers.GetTeacherHandler(w, r)
	}
	if r.Method == http.MethodPost {
		handlers.AddTeacherHandler(w, r)
	}
}

func studentsHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Students endpoint"))
}

func main() {

	port := "3000"
	println("Server is running on port " + port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	http.HandleFunc("/teachers/", teachersHandler)
	http.HandleFunc("/students", studentsHandler)

	// handler := applyMiddlewares(
	// 	http.DefaultServeMux,
	// 	middlewares.CORS,
	// 	middlewares.SecurityHeaders,
	// 	middlewares.ResponseTimeMiddleWare,
	// )

	server := &http.Server{
		Addr:    ":" + port,
		Handler: http.DefaultServeMux,
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
