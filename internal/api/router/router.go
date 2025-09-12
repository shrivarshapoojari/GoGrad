package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func Router() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.RootHandler)

	mux.HandleFunc("/teachers/", handlers.TeachersHandler) // Handles both /teachers and /teachers/{id}
	mux.HandleFunc("/students/", handlers.StudentsHandler) // Handles both /students and /students/{id}
	return mux

}
