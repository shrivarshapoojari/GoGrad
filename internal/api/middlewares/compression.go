package middlewares

import (
	"compress/gzip"
	"fmt"
	"net/http"
	"strings"
)




func Compression(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {		
      

		if !strings.Contains(r.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-Encoding", "gzip")
		gz := gzip.NewWriter(w)

		defer gz.Close()
        
		w=&gzipResponseWriter{ResponseWriter: w,Writer: gz}

		next.ServeHTTP(w,r)
fmt.Println("After next")

	})}



type gzipResponseWriter struct{
	http.ResponseWriter
	Writer *gzip.Writer
}


func(gw *gzipResponseWriter) Write(b []byte)(int,error){

	return gw.Writer.Write(b)
}
