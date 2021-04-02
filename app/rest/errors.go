package rest

import (
	"github.com/go-chi/render"
	"log"
	"net/http"
)

const (
	ErrServerInternal = 0 // server internal error
	ErrJsonDecode     = 1 // failed unmarshalling incoming request
)

func SendErrorJSON(w http.ResponseWriter, r *http.Request, httpStatusCode int, err error, errCode int, details string) {
	log.Printf("[DEBUG] %d, %+v, %d ", httpStatusCode, err, errCode)
	render.Status(r, httpStatusCode)
	render.JSON(w, r, map[string]interface{}{"error": err.Error(), "code": errCode, "details": details})
}
