package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	// "restapi/internal/api/middlewares"
	"sync"
)

 type  Teacher struct {
	ID   int
	FirstName string
	LastName string
	Class string
	Subject string
 }


var  teachers =make(map[int] Teacher)
var mutex = &sync.Mutex{}
var nextId=1

func init(){
	teachers[nextId]=Teacher{
		ID: nextId,
		FirstName: "john",
		LastName: "Doe",
		Class: "9",
		Subject: "Math",


	}
	nextId++;

	teachers[nextId]=Teacher{
		ID: nextId,
		FirstName: "Jane",
		LastName: "Smith",
		Class: "10",
		Subject: "Science",

	}
	nextId++;
	teachers[nextId]=Teacher{
		ID: nextId,
		FirstName: "Alice",
		LastName: "Johnson",
		Class: "11",
		Subject: "History",
}


}


func getTeacherHandler(w http.ResponseWriter, r *http.Request) {


	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")

	if idStr == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")


	 teacherList := make([]Teacher, 0, len(teachers))
	for _, teacher := range teachers {
		if (firstName=="" || teacher.FirstName==firstName)  && (lastName=="" || teacher.LastName==lastName) {
		teacherList = append(teacherList, teacher)
	}
 

	response := struct {
		Status string   `json:"status"`
		Count  int      `json:"count"`
		Data   []Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(teachers),
		Data:   teacherList,
	}
   
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
} else {

 
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}

	teacher, exists := teachers[id]
	if !exists {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}

	response := struct {
		Status string  `json:"status"`
		Data   Teacher `json:"data"`
	}{
		Status: "success",
		Data:   teacher,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
}




	 



func teachersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		getTeacherHandler(w,r)
	}
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Error parsing form", http.StatusBadRequest)
			return
		}
		firstName := r.FormValue("first_name")
		lastName := r.FormValue("last_name")
		class := r.FormValue("class")
		subject := r.FormValue("subject")
		if firstName == "" || lastName == "" || class == "" || subject == "" {
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}
		mutex.Lock()
		defer mutex.Unlock()
		teacher := Teacher{
			ID:        nextId,
			FirstName: firstName,
			LastName:  lastName,
			Class:     class,
			Subject:   subject,
		}
		teachers[nextId] = teacher
		nextId++
		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, "Teacher created with ID: %d\n", teacher.ID)
		return
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
