package handler

import (
	"encoding/json"
	"net/http"
)

type pingRequest struct {
	Ping string `json:"ping"`
}

func GetPing(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	resp := pingRequest{Ping: "Is working"}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
