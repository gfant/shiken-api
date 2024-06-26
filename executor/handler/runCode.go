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

type codeEnvironment struct {
	TmpFolder       string // Tmp file name for storing all testings
	ProblemFile     string // Name created for Problem file
	TestProblemFile string // Name created for Testing file
	Result          string // Result of code execution
	Code            string // Data Related to code of user
	ProblemId       string // Data Related to problem requested by user
}

func Run(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")

	var req codeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		badReq := fmt.Sprintf("Bad Req %s", r.Body)
		http.Error(w, badReq, http.StatusBadRequest)
		return
	}

	env := codeEnvironment{
		TmpFolder:       tmpTestingFolder,
		ProblemFile:     fmt.Sprintf("p%s/p%s.gno", req.ProblemId, req.ProblemId),
		TestProblemFile: fmt.Sprintf("p%s/p%s_test.gno", req.ProblemId, req.ProblemId),
		Code:            req.Code,
		ProblemId:       req.ProblemId,
	}

	execution, err := env.executeCode()

	resp := EvalResponse{Output: execution, Error: err}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func (env codeEnvironment) executeCode() (string, error) {
	var output string

	_, err := env.setupEnvironment()
	if err != nil {
		return "", err
	}

	// Executes a bash command to run this code
	output, err = env.evalCode()
	if err != nil {
		return "", err
	}
	return output, nil
}

func (env codeEnvironment) setupEnvironment() (string, error) {
	// Files required to test and their structure

	path := filepath.Join(headPath, env.TmpFolder, env.ProblemFile)
	if err := generateFolderPath(path); err != nil {
		return "", err
	}
	if err := generateFileToPath(env.Code, path); err != nil {
		return "", err
	}
	if err := env.copyTestFileAndGnoMod(env.ProblemId); err != nil {
		return "", err
	}
	return path, nil
}

// Creates a tmp folder where everything will be tested
func generateFolderPath(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return nil
}

// Creates the gno file to test the code of the user
func generateFileToPath(code, path string) error {
	codebytes := []byte(code)
	if err := os.WriteFile(path, codebytes, 0744); err != nil {
		return err
	}
	return nil
}

func (env codeEnvironment) copyTestFileAndGnoMod(problemId string) error {
	// testfile
	if err := env.copyEnvProblem(env.TestProblemFile); err != nil {
		return err
	}
	// gnomod
	gnoMod := fmt.Sprintf("p%s/gno.mod", problemId)
	if err := env.copyEnvProblem(gnoMod); err != nil {
		return err
	}
	return nil
}

func (env codeEnvironment) copyEnvProblem(file string) error {
	src := filepath.Join(headPath, problemsPath, file)
	dst := filepath.Join(headPath, env.TmpFolder, file)
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

func (env codeEnvironment) evalCode() (string, error) {
	// Getting data for file and dirs
	testFilePath := filepath.Join(headPath, env.TmpFolder, env.TestProblemFile)
	testFileDir := filepath.Dir(testFilePath)
	testFile := filepath.Base(testFilePath)
	testUserCodeCommand := fmt.Sprintf("gno test %s", testFile)

	userCodeChannel := make(chan string)
	userErrorChannel := make(chan error)
	go runUserCode(testFileDir, testUserCodeCommand, userCodeChannel, userErrorChannel)

	time.Sleep(3 * time.Second)

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
		go runSecurityTask(testUserCodeCommand)
	}

	return string(codeUserResult), nil
}

func runUserCode(testFileDir, testUserCode string, resultChannel chan<- string, errorChannel chan<- error) {
	TestUserCodeCommand := fmt.Sprintf("cd %s; gno mod tidy; %s 2>&1", testFileDir, testUserCode)
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
