package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
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
	absolutePathToProblems := filepath.Join(headPath, problemsFolder)
	lsCommand := fmt.Sprintf(`cd %s; ls -la | grep "p\d"`, absolutePathToProblems)
	listExec := exec.Command("sh", "-c", lsCommand)
	output, err := listExec.Output()
	if err != nil {
		return []string{}, err
	}
	listArr := strings.Split(string(output), "\n")
	listArr = listArr[:len(listArr)-1]
	response := []string{}
	for idx, line := range listArr {
		parts := strings.Split(line, " ")
		last := parts[len(parts)-1]
		problem := fmt.Sprintf("Problem %s : %s", strconv.Itoa(idx+1), last)
		response = append(response, problem)
	}
	return response, nil
}
