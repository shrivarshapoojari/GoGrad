package router

import (
	"net/http"
	"restapi/internal/api/handlers"
)

func Router() *http.ServeMux {

	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.RootHandler)

	 
	mux.HandleFunc("GET /teachers/", handlers.GetTeacherHandler)  
	mux.HandleFunc("GET /teachers/{id}", handlers.GetTeacherHandler)  
	mux.HandleFunc("POST /teachers/", handlers.AddTeacherHandler)  
	mux.HandleFunc("PUT /teachers/{id}", handlers.UpdateTeacherHandler) 
	mux.HandleFunc("PATCH /teachers/{id}", handlers.PatchTeacherHandler)  
	mux.HandleFunc("DELETE /teachers/{id}", handlers.DeleteTeacherHandler)  
	mux.HandleFunc("/students", handlers.StudentsHandler)   
	mux.HandleFunc("/students/", handlers.StudentsHandler)  
	return mux

}
