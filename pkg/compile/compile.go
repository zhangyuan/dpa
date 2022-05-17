package compile

import (
	"dp/pkg/helpers"
	"dp/pkg/workflow"
	"io/ioutil"
	"os"
	"strings"

	"github.com/pkg/errors"
)

func Compile(projectDirectory string) error {

	workingDirectory := strings.TrimSuffix(projectDirectory, "/")

	buildDirectory := workingDirectory + "/build"
	err := helpers.Mkdirp(buildDirectory)
	if err != nil {
		return errors.Wrap(err, "fail to create build directory")
	}
	outputPath := buildDirectory + "/dp-cloudformation.yaml"

	workflowPath := strings.TrimSuffix(projectDirectory, "/") + "/dp.yaml"

	yamlFileContent, err := ioutil.ReadFile(workflowPath)
	if err != nil {
		return errors.Wrap(err, "fail to read dp.yaml")
	}

	workflow, err := workflow.Parse(workingDirectory, yamlFileContent)
	if err != nil {
		return errors.Wrap(err, "fail to parse workflow")
	}

	outputContent, err := workflow.Render()
	if err != nil {
		return errors.Wrap(err, "fail to render workflow")
	}

	err = os.WriteFile(outputPath, []byte(outputContent), 0644)
	if err != nil {
		return errors.Wrap(err, "fail to write to "+outputPath)
	}

	return nil
}
