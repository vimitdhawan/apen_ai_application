package handler

import (
	"context"
	"encoding/json"
	gogpt "github.com/sashabaranov/go-gpt3"
	"net/http"
	"strings"
)

func TextCompletionHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		headers := []string{"Content-Type", "Accept", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Access-Control-Allow-Methods", "Origin", "X-Requested-With"}
		methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "PATCH"}
		w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))
		w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))
		return
	}

	var cr gogpt.CompletionRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&cr); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	c := gogpt.NewClient(strings.Split(r.Header.Get("Authorization"), " ")[1])
	ctx := context.Background()
	resp, err := c.CreateCompletion(ctx, cr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, resp)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}
