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

var teachers = make(map[int]models.Teacher)
var mutex = &sync.Mutex{}
var nextId = 1

func init() {
	teachers[nextId] = models.Teacher{
		ID:        nextId,
		FirstName: "john",
		LastName:  "Doe",
		Class:     "9",
		Subject:   "Math",
	}
	nextId++

	teachers[nextId] = models.Teacher{
		ID:        nextId,
		FirstName: "Jane",
		LastName:  "Smith",
		Class:     "10",
		Subject:   "Science",
	}
	nextId++
	teachers[nextId] = models.Teacher{
		ID:        nextId,
		FirstName: "Alice",
		LastName:  "Johnson",
		Class:     "11",
		Subject:   "History",
	}

}


func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
	 GetTeacherHandler(w, r)
	}
	if r.Method == http.MethodPost {
	 AddTeacherHandler(w, r)
	}
}


func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("Adding new teachers\n")
	mutex.Lock()
	defer mutex.Unlock()

	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		fmt.Print("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	addedTeachers := make([]models.Teacher, 0, len(newTeachers))
	for _, t := range newTeachers {
		if t.FirstName == "" || t.LastName == "" || t.Class == "" || t.Subject == "" {
			http.Error(w, "Missing required fields in one of the teachers", http.StatusBadRequest)
			return
		}
		t.ID = nextId
		teachers[nextId] = t
		nextId++
		addedTeachers = append(addedTeachers, t)
	}
	response := struct {
		Status string           `json:"status"`
		Count  int              `json:"count"`
		Data   []models.Teacher `json:"data"`
	}{
		Status: "success",
		Count:  len(addedTeachers),
		Data:   addedTeachers,
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func GetTeacherHandler(w http.ResponseWriter, r *http.Request) {

	path := strings.TrimPrefix(r.URL.Path, "/teachers/")
	idStr := strings.TrimSuffix(path, "/")

	if idStr == "" {
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")

		teacherList := make([]models.Teacher, 0, len(teachers))
		for _, teacher := range teachers {
			if (firstName == "" || teacher.FirstName == firstName) && (lastName == "" || teacher.LastName == lastName) {
				teacherList = append(teacherList, teacher)
			}
		}

		// ✅ MOVED OUTSIDE THE LOOP - SEND RESPONSE ONLY ONCE
		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList), // ✅ Fixed: should be len(teacherList), not len(teachers)
			Data:   teacherList,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
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
			Status string         `json:"status"`
			Data   models.Teacher `json:"data"`
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
