package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type ProblemRequest struct {
	ProblemId string `json:"id"`
}

type ProblemContent struct {
	Statement string   `json:"statement"`
	Title     string   `json:"title"`
	Examples  []string `json:"examples"`
	Error     error    `json:"error"`
}

func GetProblem(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	path := r.URL.Path

	pathParts := strings.Split(path, "/")
	problemId := pathParts[len(pathParts)-1]

	problemContent, err := getProblem(problemId)
	if err != nil {
		problemContent = ProblemContent{Error: err}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(problemContent); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func getProblem(problemId string) (ProblemContent, error) {
	// Main Path + Problems Folder + Specific Problem folder + json file
	pathToProblemContent := filepath.Join(headPath, problemsFolder, fmt.Sprintf("p%s", problemId), fmt.Sprintf("p%s.json", problemId))
	jsonFile, err := os.Open(pathToProblemContent)
	if err != nil {
		return ProblemContent{}, err
	}
	defer jsonFile.Close()
	byteValue, _ := io.ReadAll(jsonFile)
	var result ProblemContent
	json.Unmarshal([]byte(byteValue), &result)

	return result, nil
}
