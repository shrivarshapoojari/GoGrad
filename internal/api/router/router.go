package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func Router() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.RootHandler)

	mux.HandleFunc("/teachers", handlers.TeachersHandler)  // Handles /teachers
	mux.HandleFunc("/teachers/", handlers.TeachersHandler) // Handles /teachers/{id}
	mux.HandleFunc("/students", handlers.StudentsHandler)  // Handles /students
	mux.HandleFunc("/students/", handlers.StudentsHandler) // Handles /students/{id}
	return mux

}
