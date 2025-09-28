package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"restapi/internal/models"
	"restapi/internal/repository/sqlconnect"
	"strconv"
	"strings"
)

// Email validation regex pattern
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

// isValidEmail validates email format
func isValidEmail(email string) bool {
	return emailRegex.MatchString(email)
}

func TeachersHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Print("TeachersHandler invoked\n")
	fmt.Print(r.Method + "\n")
	if r.Method == http.MethodGet {
		GetTeacherHandler(w, r)
	}
	if r.Method == http.MethodPost {
		fmt.Print("Adding new teacher\n")
		AddTeacherHandler(w, r)
	}
	if r.Method == http.MethodPut {
		fmt.Print("Updating teacher\n")
		updateTeacherHandler(w, r)
	}
	if r.Method == http.MethodPatch {
		fmt.Print("Patching teacher\n")
		patchTeacherHandler(w, r)
	}
	if r.Method == http.MethodDelete {
		fmt.Print("Deleting teacher\n")
		deleteTeacherHandler(w, r)
	}
}

func AddTeacherHandler(w http.ResponseWriter, r *http.Request) {
	db := sqlconnect.ConnectDb()
	defer db.Close()
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
	insertSQL := "INSERT INTO teachers (first_name, last_name, email, class, subject) VALUES (?, ?, ?, ?, ?)"

	for _, t := range newTeachers {
		if t.FirstName == "" || t.LastName == "" || t.Email == "" || t.Class == "" || t.Subject == "" {
			http.Error(w, "Missing required fields in one of the teachers", http.StatusBadRequest)
			return
		}

		// Validate email format
		if !isValidEmail(t.Email) {
			http.Error(w, fmt.Sprintf("Invalid email format: %s", t.Email), http.StatusBadRequest)
			return
		}

		// Execute INSERT query
		result, err := db.Exec(insertSQL, t.FirstName, t.LastName, t.Email, t.Class, t.Subject)
		if err != nil {
			fmt.Printf("Error inserting teacher: %v\n", err)
			http.Error(w, "Error adding teacher to database", http.StatusInternalServerError)
			return
		}
		fmt.Println(result)
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
	db := sqlconnect.ConnectDb()
	defer db.Close()
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
			query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE 1=1"
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
			rows, err = db.Query("SELECT id, first_name, last_name, email, class, subject FROM teachers")
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
			err := rows.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
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
		query := "SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?"
		row := db.QueryRow(query, id)

		err = row.Scan(&teacher.ID, &teacher.FirstName, &teacher.LastName, &teacher.Email, &teacher.Class, &teacher.Subject)
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

func updateTeacherHandler(w http.ResponseWriter, r *http.Request) {
	db := sqlconnect.ConnectDb()
	defer db.Close()
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}
	var updatedTeacher models.Teacher
	err = json.NewDecoder(r.Body).Decode(&updatedTeacher)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update teacher in the database
	query := "UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?"
	_, err = db.Exec(query, updatedTeacher.FirstName, updatedTeacher.LastName, updatedTeacher.Email, updatedTeacher.Class, updatedTeacher.Subject, id)
	if err != nil {
		fmt.Printf("Error updating teacher: %v\n", err)
		http.Error(w, "Error updating teacher in database", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string         `json:"status"`
		Data   models.Teacher `json:"data"`
	}{
		Status: "success",
		Data:   updatedTeacher,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

// PATCH /teachers/{id}

func patchTeacherHandler(w http.ResponseWriter, r *http.Request) {
	db := sqlconnect.ConnectDb()
	defer db.Close()
	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}
	var updatedFields map[string]interface{}
	err = json.NewDecoder(r.Body).Decode(&updatedFields)

	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var existingTeacher models.Teacher
	err = db.QueryRow("SELECT id, first_name, last_name, email, class, subject FROM teachers WHERE id = ?", id).Scan(
		&existingTeacher.ID,
		&existingTeacher.FirstName,
		&existingTeacher.LastName,
		&existingTeacher.Email,
		&existingTeacher.Class,
		&existingTeacher.Subject,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Teacher not found", http.StatusNotFound)
			return
		}
		fmt.Printf("Error retrieving existing teacher: %v\n", err)
		http.Error(w, "Error retrieving teacher from database", http.StatusInternalServerError)
		return
	}
	// for field, value := range updatedFields {
	// 	switch field {
	// 	case "first_name":
	// 		existingTeacher.FirstName = value.(string)
	// 	case "last_name":
	// 		existingTeacher.LastName = value.(string)
	// 	case "email":
	// 		existingTeacher.Email = value.(string)
	// 	case "class":
	// 		existingTeacher.Class = value.(string)
	// 	case "subject":
	// 		existingTeacher.Subject = value.(string)
	// 	}
	// }

	// apply update using reflect
	teacherVal := reflect.ValueOf(&existingTeacher).Elem()
	teacherType := teacherVal.Type()

	for k, v := range updatedFields {
		for i := 0; i < teacherVal.NumField(); i++ {
			field := teacherType.Field(i)
			if field.Tag.Get("json") == k {

				if teacherVal.Field(i).CanSet() {
					teacherVal.Field(i).Set(reflect.ValueOf(v).Convert(teacherVal.Field(i).Type()))
				}

			}
		}
	}

	// Save the updated teacher record back to the database
	_, err = db.Exec("UPDATE teachers SET first_name = ?, last_name = ?, email = ?, class = ?, subject = ? WHERE id = ?",
		existingTeacher.FirstName,
		existingTeacher.LastName,
		existingTeacher.Email,
		existingTeacher.Class,
		existingTeacher.Subject,
		existingTeacher.ID,
	)
	if err != nil {
		fmt.Printf("Error updating teacher: %v\n", err)
		http.Error(w, "Error updating teacher in database", http.StatusInternalServerError)
		return
	}

	response := struct {
		Status string         `json:"status"`
		Data   models.Teacher `json:"data"`
	}{
		Status: "success",
		Data:   existingTeacher,
	}
	fmt.Print("Patching teacher\n")

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}

}

func deleteTeacherHandler(w http.ResponseWriter, r *http.Request) {
	db := sqlconnect.ConnectDb()
	defer db.Close()
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, "/teachers/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid teacher ID", http.StatusBadRequest)
		return
	}
	result, err := db.Exec("DELETE FROM teachers WHERE id = ?", id)
	if err != nil {
		fmt.Printf("Error deleting teacher: %v\n", err)
		http.Error(w, "Error deleting teacher from database", http.StatusInternalServerError)
		return
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error getting rows affected: %v\n", err)
		http.Error(w, "Error checking deletion result", http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "Teacher not found", http.StatusNotFound)
		return
	}
	response := struct {
		Status string `json:"status"`
	}{
		Status: "success",
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
