package main

import "net/http"

func main() {

	port := "3000"
	println("Server is running on port " + port)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	http.HandleFunc("/teachers", func(w http.ResponseWriter, r *http.Request) {	
		if(r.Method==http.MethodGet){
			w.Write([]byte("GET Teachers endpoint"))			
			return
		}	
		if(r.Method==http.MethodPost){
			w.Write([]byte("POST Teachers endpoint"))			
			return
		}
	})

	http.HandleFunc("/students", func(w http.ResponseWriter, r *http.Request) {		
		w.Write([]byte("Students endpoint"))			
	})
	err := http.ListenAndServe(port, nil)
	if err != nil {
		panic(err)
		
	}


}