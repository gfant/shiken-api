package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type ProblemRequest struct {
	ProblemId string `json:"id"`
}

type ProblemContent struct {
	Content string `json:"content"`
	Error   error  `json:"error"`
}

func GetProblem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ProblemRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		badReq := fmt.Sprintf("Bad Req %s", r.Body)
		http.Error(w, badReq, http.StatusBadRequest)
		return
	}

	problemId := req.ProblemId
	execution, err := getProblem(problemId)
	resp := ProblemContent{Content: execution, Error: err}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func getProblem(problemId string) (string, error) {
	return problemId, nil
}
