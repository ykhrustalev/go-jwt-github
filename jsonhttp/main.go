package jsonhttp

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Message struct {
	Message string
	Error   string
}

func ErrorResponse(w http.ResponseWriter, message string, err error, httpCode int) {
	fmt.Printf("%s, error: %v\n", message, err)
	obj := Message{message, ""}
	if err != nil {
		obj.Error = err.Error()
	}
	Response(w, obj, httpCode)
}

func Response200(w http.ResponseWriter, obj interface{}) {
	Response(w, obj, http.StatusOK)
}
func Response(w http.ResponseWriter, obj interface{}, httpCode int) {
	body, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpCode)
	w.Write(body)
}

func ErrorResponse500(w http.ResponseWriter, message string, err error) {
	ErrorResponse(w, message, err, http.StatusInternalServerError)
}
