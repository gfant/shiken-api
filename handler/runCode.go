package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type codeRequest struct {
	Code      string `json:"code"`
	ProblemId string `json:"id"`
}

type EvalResponse struct {
	Output string `json:"output"`
	Error  error  `json:"error"`
}

func RunCode(w http.ResponseWriter, r *http.Request) {
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
	var output string

	err := setupEnvironment(code, problemId)
	if err != nil {
		return "", err
	}
	/*
		// Executes a bash command to run this code
		cmd := exec.Command("gno", "test", filepath)
		bytes, err := cmd.Output()
		if err != nil {
			return "0", nil
		}
		output = string(bytes[:])
	*/
	return output, nil
}

func setupEnvironment(code, problemId string) error {
	fileroute := fmt.Sprintf("p%s/p%s.gno", problemId, problemId)

	path := filepath.Join(tmpTestingFolder, fileroute)
	if err := generateCodeFolder(path); err != nil {
		return err
	}
	if err := generateCodeFile(code, path); err != nil {
		return err
	}
	if err := copyTestFile(problemId); err != nil {
		return err
	}
	return nil
}

// Creates a tmp folder where everything will be tested
func generateCodeFolder(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return nil
}

// Creates the gno file to test the code of the user
func generateCodeFile(code, path string) error {
	codebytes := []byte(code)
	if err := os.WriteFile(path, codebytes, 0744); err != nil {
		return err
	}
	return nil
}

func copyTestFile(problemId string) error {
	testFile := fmt.Sprintf("p%s/p%s_test.gno", problemId, problemId)
	src := filepath.Join(problemsPath, testFile)
	dst := filepath.Join(tmpTestingFolder, testFile)
	fmt.Println(src, dst)
	err := CopyFile(src, dst)
	if err != nil {
		return err
	}
	return nil
}

func CopyFile(src, dst string) error {
	// Open the source file
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	// Create the destination file
	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	// Copy the contents from source to destination
	_, err = io.Copy(destinationFile, sourceFile)
	if err != nil {
		return err
	}

	// Ensure all data is written to the destination file
	err = destinationFile.Sync()
	if err != nil {
		return err
	}

	return nil
}
