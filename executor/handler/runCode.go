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
	Address   string `json:"address"`
	TxHash    string `json:"hash"` // Tx that will verify the user sent a tx to the blockchain
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

// Check if the request is valid. Otherwise will return an http.error
func validateRequest(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
	}
}

// Get the Request Body and checks if it is valid. Otherwise will return an http.error
func decodeRequestBody(w http.ResponseWriter, req *http.Request) codeRequest {
	var cReq codeRequest
	if err := json.NewDecoder(req.Body).Decode(&cReq); err != nil {
		badReq := fmt.Sprintf("Bad Req %s", req.Body)
		http.Error(w, badReq, http.StatusBadRequest)
	}
	return cReq
}

// Generate the environment for the code to be executed
func generateCodeEnvironment(req codeRequest) codeEnvironment {
	return codeEnvironment{
		TmpFolder:       tmpTestingFolder,
		ProblemFile:     fmt.Sprintf("p%s/p%s.gno", req.ProblemId, req.ProblemId),
		TestProblemFile: fmt.Sprintf("p%s/p%s_test.gno", req.ProblemId, req.ProblemId),
		Code:            req.Code,
		ProblemId:       req.ProblemId,
	}
}

func sendResponseBack(w http.ResponseWriter, resp EvalResponse) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func RunHandler(w http.ResponseWriter, r *http.Request) {
	validateRequest(w, r)
	w.Header().Set("Content-Type", "application/json")
	req := decodeRequestBody(w, r)
	// Environment to use for the code
	env := generateCodeEnvironment(req)

	// Execute the code and get the result
	execution, err := env.executeCode()

	// The execution ended and now you have to send it
	sendResultOnchain(req, execution)

	// Send the result to the user
	resp := EvalResponse{Output: execution, Error: err}
	sendResponseBack(w, resp)
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

// Setup the environment filers and folders for the code to be executed
func (env codeEnvironment) setupEnvironment() (string, error) {

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

// Copy the test file and the gno.mod file to the tmp folder
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

// Copies files into the tmp folder
func (env codeEnvironment) copyEnvProblem(file string) error {
	src := filepath.Join(headPath, problemsPath, file)
	dst := filepath.Join(headPath, env.TmpFolder, file)
	err := CopyFile(src, dst)
	if err != nil {
		return err
	}
	return nil
}

// Copies a file from src to dst
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
	path := filepath.Join(headPath, env.TmpFolder, env.TestProblemFile)
	dirName := filepath.Dir(path)
	file := filepath.Base(path)
	gnoTestFile := fmt.Sprintf("gno test %s", file)

	userCodeChannel := make(chan string)
	userErrorChannel := make(chan error)
	go runUserCode(dirName, gnoTestFile, userCodeChannel, userErrorChannel)

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
		go runSecurityTask(gnoTestFile)
	}

	return string(codeUserResult), nil
}

func runUserCode(testFileDir, testUserCode string, resultChannel chan<- string, errorChannel chan<- error) {
	cmdTestUserCode := fmt.Sprintf("cd %s; gno mod tidy; %s 2>&1", testFileDir, testUserCode)
	execCommand := exec.Command("sh", "-c", cmdTestUserCode)
	output, err := execCommand.Output()
	if err != nil {
		errorChannel <- err
	}
	resultChannel <- string(output)
}

func createFilterCommand(testUserCode string) string {
	return fmt.Sprintf("ps aux | grep '%s' | grep -v 'grep'", testUserCode)
}

func getPIDFromOutput(commandExecuted *exec.Cmd) ([]string, error) {
	bytes, err := commandExecuted.Output()
	if err != nil {
		return nil, err
	}
	output := string(bytes)
	pids := strings.Split(output, "\n")
	return pids, nil
}

func killProcessByPID(pid string) error {
	fmt.Printf("\nKilling pid %s...", pid)
	killCommand := fmt.Sprintf("kill %s", pid)
	killing := exec.Command("sh", "-c", killCommand)
	_, err := killing.Output()
	if err != nil {
		fmt.Printf("Kill Process crashed: %s: %s", pid, err)
		return err
	}
	return nil
}

func runSecurityTask(testUserCode string) error {
	SecurityTaskCommand := createFilterCommand(testUserCode)
	execSecurityTask := exec.Command("sh", "-c", SecurityTaskCommand)

	pids, err := getPIDFromOutput(execSecurityTask)
	if err != nil {
		return err
	}

	fmt.Printf("\nUser code timeout. Starting clean up")
	for _, process := range pids {
		if len(strings.TrimSpace(process)) == 0 {
			continue
		}
		fields := strings.Fields(process)
		pid := fields[1]
		err := killProcessByPID(pid)
		if err != nil {
			return err
		}
		fmt.Printf("done")

	}
	return nil
}

func sendResultOnchain(cReq codeRequest, result string) {
	account, client := SetupRegisterEnvironment(
		"/Users/iam-agf/Library/Application Support/gno",
		"g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5",
		"dev",
		"tcp://127.0.0.1:26657",
	)
	err := MakeTx(
		"gno.land/r/dev/shiken",
		"1000000ugnot",
		"",
		"AddNewScore",
		2000000,
		[]string{
			cReq.Address,
			cReq.ProblemId,
			result,
			"",
		},
		account,
		client,
	)
	if err != nil {
		panic(err)
	}
}
