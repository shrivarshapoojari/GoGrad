package models

type Student struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Class     string `json:"class"`
	Age       int    `json:"age"`
	Grade     string `json:"grade"`
}
