package sqlconnect

import "database/sql"

func ConnectDb() *sql.DB{

	connectionString:="root:root@tcp(127.0.0.1:3306)/gograd"
	db,err:=sql.Open("mysql",connectionString)

	if err!=nil{
		panic(err)
	}
	return db;
}