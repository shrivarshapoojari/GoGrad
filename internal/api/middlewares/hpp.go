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




