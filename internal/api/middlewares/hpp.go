package middlewares

import (
	"net/http"

)



type HPPOptions struct{
	CheckQuery bool
	CheckBody bool
	CheckBodyOnlyForContentType string
	Whitelist []string
}




func Hpp(options HPPOptions) func(http.Handler) http.Handler {

	return func(next http.Handler) http.Handler {


		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
 
			if options.CheckQuery && r.Method == http.MethodPost {
				filterBodyParams(r, options.Whitelist)
			}




			next.ServeHTTP(w,r)


		})
	}
}


func filterBodyParams(r *http.Request, whitelist []string)  {
   
	err :=r.ParseForm()

	if err != nil {
		fmt.Println("Error parsing form:", err)
		return
	}

	for k,v := range r.Form{

		if len(v) > 1{
			r.Form.Set(k,v[0])
		}
		if !isWhiteListed(k,whitelist){
			 delete (r.Form,k)
		}

	}


}


func isWhiteListed(param string, whitelist []string) bool {

	for _, p := range whitelist {
		if p == param {
			return true
		}
		return  false
}