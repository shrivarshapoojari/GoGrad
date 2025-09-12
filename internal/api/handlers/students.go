package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"restapi/internal/models"
	"strconv"
	"strings"
	"sync"
)

var students = make(map[int]models.Student)
var studentMutex = &sync.Mutex{}
var nextStudentId = 1

func init() {
	students[nextStudentId] = models.Student{
		ID:        nextStudentId,
		FirstName: "Emma",
		LastName:  "Watson",
		Class:     "10A",
		Age:       16,
		Grade:     "A",
	}
	nextStudentId++

	students[nextStudentId] = models.Student{
		ID:        nextStudentId,
		FirstName: "Liam",
		LastName:  "Johnson",
		Class:     "9B",
		Age:       15,
		Grade:     "B+",
	}
	nextStudentId++

	students[nextStudentId] = models.Student{
		ID:        nextStudentId,
		FirstName: "Sophia",
		LastName:  "Davis",
		Class:     "11C",
		Age:       17,
		Grade:     "A-",
	}
	nextStudentId++
}

func AddStudentHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Adding new student\n")
	studentMutex.Lock()
	defer studentMutex.Unlock()

	var student models.Student
	if err := json.NewDecoder(r.Body).Decode(&student); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if student.FirstName == "" || student.LastName == "" || student.Class == "" {
		http.Error(w, "Missing required fields: first_name, last_name, class", http.StatusBadRequest)
		return
	}

	// Set ID and add to map
	student.ID = nextStudentId
	students[nextStudentId] = student
	nextStudentId++

	response := struct {
		Status  string         `json:"status"`
		Message string         `json:"message"`
		Data    models.Student `json:"data"`
	}{
		Status:  "success",
		Message: "Student added successfully",
		Data:    student,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func GetStudentHandler(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimPrefix(r.URL.Path, "/students/")
	idStr := strings.TrimSuffix(path, "/")

	if idStr == "" {
		// Get all students or filter by query parameters
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")
		class := r.URL.Query().Get("class")
		grade := r.URL.Query().Get("grade")

		studentList := make([]models.Student, 0, len(students))
		for _, student := range students {
			// Apply filters if provided
			if (firstName == "" || student.FirstName == firstName) &&
				(lastName == "" || student.LastName == lastName) &&
				(class == "" || student.Class == class) &&
				(grade == "" || student.Grade == grade) {
				studentList = append(studentList, student)
			}
		}

		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Student `json:"data"`
		}{
			Status: "success",
			Count:  len(studentList),
			Data:   studentList,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	} else {
		// Get specific student by ID
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid student ID", http.StatusBadRequest)
			return
		}

		student, exists := students[id]
		if !exists {
			http.Error(w, "Student not found", http.StatusNotFound)
			return
		}

		response := struct {
			Status string         `json:"status"`
			Data   models.Student `json:"data"`
		}{
			Status: "success",
			Data:   student,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	}
}

// StudentsHandler routes requests based on HTTP method
func StudentsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		GetStudentHandler(w, r)
	case http.MethodPost:
		AddStudentHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
