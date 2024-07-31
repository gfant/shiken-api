package main

import (
	"fmt"
	"os"
	"text/template"
)

type Test struct {
	Test   string
	Result string
}

type ProblemConfig struct {
	Id      int
	Problem string
}

func main() {
	problem := 1
	problemConfig := map[int]ProblemConfig{
		1: ProblemConfig{
			Id:      1,
			Problem: TemplateP1Tests,
		},
	}
	chosenProblemConfig := problemConfig[problem]
	produceTest(chosenProblemConfig)
}

func produceTest(config ProblemConfig) {
	testTemplate := config.Problem
	problemId := config.Id

	// Crear un nuevo objeto de plantilla
	tmpl, err := template.New("test").Parse(testTemplate)
	if err != nil {
		panic(err)
	}

	// Start the test generator
	tests := []Test{}
	for i := 0; i < 10; i++ {
		test, testResult := TestGeneratorP1()
		tests = append(tests, Test{Test: test, Result: testResult})
	}

	filename := fmt.Sprintf("p%d.gno", problemId)
	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Rellenar la plantilla y escribir a stdout
	err = tmpl.Execute(file, tests)
	if err != nil {
		panic(err)
	}

}
