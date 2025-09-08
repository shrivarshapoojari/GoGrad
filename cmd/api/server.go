package main

import (
	"fmt"
	"net/http"
	"strings"
	"restapi/internal/api/middlewares"
)

func main() {

	port := "3000"
	println("Server is running on port " + port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	http.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {

			w.Write([]byte("GET Teachers endpoint"))

			return
		}
		if r.Method == http.MethodPost {
			w.Write([]byte("POST Teachers endpoint"))
			return
		}
	})

	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Students endpoint"))
		
	})

	//path parameters
	http.HandleFunc("/teachers/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL.Path)
		path := strings.TrimPrefix(r.URL.Path, "/teachers/")
		userId := strings.TrimSuffix(path, "/")
		fmt.Println("User ID:", userId)
		w.Write([]byte("Teacher ID: " + userId))
	})

	//query parameters

	http.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query().Get("q")

		fmt.Println("Query Parameters:", query)
	})

    server := &http.Server{
		Addr:   ":" + port,
		Handler: middlewares.ResponseTimeMiddleWare(http.DefaultServeMux),
	}


	err :=server.ListenAndServe()
	if err != nil {
		panic(err)

	}

}
