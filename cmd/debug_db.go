package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Connect to database
	connectionString := "root:root@tcp(127.0.0.1:3306)/gograd"
	db, err := sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Error pinging database:", err)
	}
	fmt.Println("✅ Connected to database successfully")

	// Check if teachers table exists
	rows, err := db.Query("SHOW TABLES LIKE 'teachers'")
	if err != nil {
		log.Fatal("Error checking for teachers table:", err)
	}
	defer rows.Close()

	if rows.Next() {
		fmt.Println("✅ Teachers table exists")

		// Check table structure
		rows2, err := db.Query("DESCRIBE teachers")
		if err != nil {
			log.Fatal("Error describing teachers table:", err)
		}
		defer rows2.Close()

		fmt.Println("\n📋 Teachers table structure:")
		for rows2.Next() {
			var field, fieldType, null, key, defaultVal, extra sql.NullString
			err := rows2.Scan(&field, &fieldType, &null, &key, &defaultVal, &extra)
			if err != nil {
				log.Fatal("Error scanning table structure:", err)
			}
			fmt.Printf("Column: %s, Type: %s\n", field.String, fieldType.String)
		}
	} else {
		fmt.Println("❌ Teachers table does NOT exist")
		fmt.Println("You need to create the table first!")
	}

	// Try to count existing records
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM teachers").Scan(&count)
	if err != nil {
		fmt.Printf("❌ Error counting teachers: %v\n", err)
	} else {
		fmt.Printf("✅ Teachers table has %d records\n", count)
	}
}
