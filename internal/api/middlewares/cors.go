package middlewares

import( "net/http"
"fmt"
)

func cors (next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {

		origin:=r.Header.Get("Accept-Encoding")
		fmt.Println("Origin:", origin)
		w.Header().Set("Access-Control-Allow-Headers","Content-Type,Authorization" )
		w.Header().Set("Access-Control-Allow-Methods","GET,POST,PUT,DELETE,OPTIONS" )
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")


		next.ServeHTTP(w, r)
	})}
