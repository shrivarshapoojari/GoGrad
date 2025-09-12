package sqlconnect

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql" 
)

func ConnectDb() *sql.DB {

	connectionString := "root:root@tcp(127.0.0.1:3306)/gograd"
	db, err := sql.Open("mysql", connectionString)

	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to database successfully")
	return db
}
