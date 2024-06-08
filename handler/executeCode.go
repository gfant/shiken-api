package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
)

type codeRequest struct {
	Code      string `json:"code"`
	ProblemId string `json:"id"`
}

type EvalResponse struct {
	Output string `json:"output"`
	Error  error  `json:"error"`
}

func GetCode(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req codeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		badReq := fmt.Sprintf("Bad Req %s", r.Body)
		http.Error(w, badReq, http.StatusBadRequest)
		return
	}

	execution, err := executeCode(req.Code, req.ProblemId)

	resp := EvalResponse{Output: execution, Error: err}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func executeCode(code, problemId string) (string, error) {

	if err := generateCodeFolder(code, problemId); err != nil {
		return "", err
	}
	if err := generateCodeFile(code, problemId); err != nil {
		return "", err
	}

	testFilepath := fmt.Sprintf("p%s.gno", problemId)
	// Executes a bash command to run this code
	cmd := exec.Command("gno", "test", testFilepath)
	bytes, err := cmd.Output()
	if err != nil {
		return "0", nil
	}
	output := string(bytes[:])
	return output, nil
}

func generateCodeFolder(code, problemId string) error {
	// Creates a tmp folder where everything will be tested
	path := fmt.Sprintf("tmp/%s", problemId)
	if err := os.MkdirAll(path, 0755); err != nil {
		return err
	}
}

func generateCodeFile(code, problemId string) error {
	// Creates the gno file to test the code of the user
	codebytes := []byte(code)
	filepath := fmt.Sprintf("/tmp/%s/p%s.gno", problemId, problemId)

	if err := os.WriteFile(filepath, codebytes, 0644); err != nil {
		return "", err
	}
}
