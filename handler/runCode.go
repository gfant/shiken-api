package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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

	_, testFile, err := setupEnvironment(code, problemId)
	if err != nil {
		return "", err
	}

	// Executes a bash command to run this code
	output, err = EvalCode(testFile)
	if err != nil {
		return "", err
	}
	return output, nil
}

func setupEnvironment(code, problemId string) (string, string, error) {
	fileroute := fmt.Sprintf("p%s/p%s.gno", problemId, problemId)
	testFileroute := fmt.Sprintf("p%s/p%s_test.gno", problemId, problemId)

	path := filepath.Join(tmpTestingFolder, fileroute)
	if err := generateCodeFolder(path); err != nil {
		return "", "", err
	}
	if err := generateCodeFile(code, path); err != nil {
		return "", "", err
	}
	if err := copyTestFileAndGnoMod(problemId); err != nil {
		return "", "", err
	}
	return path, testFileroute, nil
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

func copyTestFileAndGnoMod(problemId string) error {
	testFile := fmt.Sprintf("p%s/p%s_test.gno", problemId, problemId)
	// testfile
	if err := copyFromSource(testFile); err != nil {
		return err
	}
	// gnomod
	gnoMod := fmt.Sprintf("p%s/gno.mod", problemId)
	if err := copyFromSource(gnoMod); err != nil {
		return err
	}
	return nil
}

func copyFromSource(file string) error {
	src := filepath.Join(problemsPath, file)
	dst := filepath.Join(tmpTestingFolder, file)
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

func EvalCode(testFileRoute string) (string, error) {
	// Getting data for file and dirs
	testFileDir := filepath.Dir(testFileRoute)
	testFile := filepath.Base(testFileRoute)
	tmpTestFileDir := filepath.Join(tmpTestingFolder, testFileDir)

	testUserCode := fmt.Sprintf("gno test %s", testFile)

	userCodeChannel := make(chan string)
	userErrorChannel := make(chan error)
	//SecurityTaskChannel := make(chan string)

	go runUserCode(tmpTestFileDir, testUserCode, userCodeChannel, userErrorChannel)

	time.Sleep(1 * time.Second)

	var codeUserResult string
	select {
	case codeUserResult = <-userCodeChannel:
		fmt.Printf("Code submitted: %s", codeUserResult)
	default:
		codeUserResult = ""
	}
	var err error
	select {
	case err = <-userErrorChannel:
		if err != nil {
			return "", err
		}
	default:
		err = nil
	}
	if codeUserResult == "" && err == nil {
		go runSecurityTask(testUserCode)
	}

	return string(codeUserResult), nil
}

func runUserCode(tmpTestFileDir, testUserCode string, resultChannel chan<- string, errorChannel chan<- error) {
	TestUserCodeCommand := fmt.Sprintf("export GNOTESTPATH=$(pwd); cd $GNOTESTPATH/%s; gno mod tidy; %s 2>&1", tmpTestFileDir, testUserCode)
	execCommand := exec.Command("sh", "-c", TestUserCodeCommand)
	output, err := execCommand.Output()

	if err != nil {
		errorChannel <- err
	}
	resultChannel <- string(output)
}

func runSecurityTask(testUserCode string) error {
	SecurityTaskCommand := fmt.Sprintf("ps aux | grep '%s' | grep -v 'grep'", testUserCode)
	execSecurityTask := exec.Command("sh", "-c", SecurityTaskCommand)
	bytes, err := execSecurityTask.Output()
	if err != nil {
		return err
	}
	output := string(bytes)

	pids := strings.Split(output, "\n")
	fmt.Printf("\nUser code timeout. Starting clean up")
	for _, process := range pids {
		if len(strings.TrimSpace(process)) == 0 {
			continue
		}
		fields := strings.Fields(process)
		pid := fields[1]
		fmt.Printf("\nKilling pid %s...", pid)
		killCommand := fmt.Sprintf("kill %s", pid)
		killing := exec.Command("sh", "-c", killCommand)
		_, err := killing.Output()
		if err != nil {
			fmt.Printf("Kill Process crashed: %s: %s", pid, err)
			return err
		}
		fmt.Printf("done")

	}
	return nil
}
