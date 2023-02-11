package apic

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type errorResponse struct {
	Message string
}

func respondWithError(w http.ResponseWriter, code int, err error) {
	errorResponse := errorResponse{fmt.Sprintf("%v", err)}
	response, _ := json.Marshal(errorResponse)
	respondWithJSON(w, code, response)
}

func respondWithJSON(w http.ResponseWriter, code int, response []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_, err := w.Write(response)
	if err != nil {
		handleAppError(w, err)
	}
}

func handleAppError(w http.ResponseWriter, err error) {
	switch err.(type) {
		
	default:
		respondWithError(w, http.StatusInternalServerError, err)
		return
	}
}
