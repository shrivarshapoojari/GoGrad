package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
	"strconv"
	"strings"
)

// Database connection (initialized once)
var db *sql.DB

func init() {
	// Initialize database connection
	db = sqlconnect.ConnectDb()
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

	var newTeachers []models.Teacher
	err := json.NewDecoder(r.Body).Decode(&newTeachers)
	if err != nil {
		fmt.Print("Error decoding request body:", err)
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	addedTeachers := make([]models.Teacher, 0, len(newTeachers))

	// Prepare SQL statement for inserting teachers
	insertSQL := "INSERT INTO teachers (first_name, last_name, class, subject) VALUES (?, ?, ?, ?)"

	for _, t := range newTeachers {
		if t.FirstName == "" || t.LastName == "" || t.Class == "" || t.Subject == "" {
			http.Error(w, "Missing required fields in one of the teachers", http.StatusBadRequest)
			return
		}

		// Execute INSERT query
		result, err := db.Exec(insertSQL, t.FirstName, t.LastName, t.Class, t.Subject)
		if err != nil {
			fmt.Printf("Error inserting teacher: %v\n", err)
			http.Error(w, "Error adding teacher to database", http.StatusInternalServerError)
			return
		}

		// Get the auto-generated ID
		id, err := result.LastInsertId()
		if err != nil {
			fmt.Printf("Error getting last insert ID: %v\n", err)
			http.Error(w, "Error retrieving teacher ID", http.StatusInternalServerError)
			return
		}

		// Set the ID and add to response
		t.ID = int(id)
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
		// Get all teachers with optional filtering
		firstName := r.URL.Query().Get("first_name")
		lastName := r.URL.Query().Get("last_name")

		var teacherList []models.Teacher
		var rows *sql.Rows
		var err error

		// Build query based on filters
		if firstName != "" || lastName != "" {
			query := "SELECT id, first_name, last_name, class, subject FROM teachers WHERE 1=1"
			args := []interface{}{}

			if firstName != "" {
				query += " AND first_name = ?"
				args = append(args, firstName)
			}
			if lastName != "" {
				query += " AND last_name = ?"
				args = append(args, lastName)
			}

			rows, err = db.Query(query, args...)
		} else {
			// Get all teachers
			rows, err = db.Query("SELECT id, first_name, last_name, class, subject FROM teachers")
		}

		if err != nil {
			fmt.Printf("Error querying teachers: %v\n", err)
			http.Error(w, "Error retrieving teachers from database", http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Scan results
		for rows.Next() {
			var teacher models.Teacher
			err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Class, &teacher.Subject)
			if err != nil {
				fmt.Printf("Error scanning teacher row: %v\n", err)
				http.Error(w, "Error processing teacher data", http.StatusInternalServerError)
				return
			}
			teacherList = append(teacherList, teacher)
		}

		// Check for any error during iteration
		if err = rows.Err(); err != nil {
			fmt.Printf("Error iterating teacher rows: %v\n", err)
			http.Error(w, "Error processing teacher data", http.StatusInternalServerError)
			return
		}

		response := struct {
			Status string           `json:"status"`
			Count  int              `json:"count"`
			Data   []models.Teacher `json:"data"`
		}{
			Status: "success",
			Count:  len(teacherList),
			Data:   teacherList,
		}

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Error encoding response", http.StatusInternalServerError)
			return
		}
	} else {
		// Get specific teacher by ID
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
			return
		}

		var teacher models.Teacher
		query := "SELECT id, first_name, last_name, class, subject FROM teachers WHERE id = ?"
		row := db.QueryRow(query, id)

		err = row.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Class, &teacher.Subject)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Teacher not found", http.StatusNotFound)
				return
			}
			fmt.Printf("Error querying teacher by ID: %v\n", err)
			http.Error(w, "Error retrieving teacher from database", http.StatusInternalServerError)
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
