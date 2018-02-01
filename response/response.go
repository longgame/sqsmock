package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Error *map[string]interface{} `json:"error"`
	Data  *interface{}            `json:"data"`
}

func (r *Response) Respond(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(r)
}
