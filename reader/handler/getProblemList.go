package handler

import (
	"encoding/json"
	"net/http"
)

type problemListRequest struct {
	List []string `json:"problems"`
}

func GetProblemList(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	list, err := getProblemList()
	if err != nil {
		http.Error(w, "A problem ocurred", http.StatusInternalServerError)
	}

	resp := problemListRequest{List: list}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func getProblemList() ([]string, error) {
	return []string{"1", "2", "3"}, nil
}
