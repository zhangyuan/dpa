package build

import (
	"dpa/pkg/workflow"
	"io/ioutil"
	"path"
	"strings"

	"github.com/pkg/errors"
)

func Build(projectDirectoryOrDpPath string) error {
	var workflowPath string
	var workingDirectory string
	if strings.HasSuffix(projectDirectoryOrDpPath, ".yaml") {
		workflowPath = projectDirectoryOrDpPath
		workingDirectory = path.Dir(workflowPath)
	} else {
		workingDirectory = strings.TrimSuffix(projectDirectoryOrDpPath, "/")
		workflowPath = strings.TrimSuffix(projectDirectoryOrDpPath, "/") + "/dp.yaml"

	}
	yamlFileContent, err := ioutil.ReadFile(workflowPath)

	if err != nil {
		return errors.Wrap(err, "fail to read dp.yaml")
	}

	workflow, err := workflow.Parse(workingDirectory, yamlFileContent)
	if err != nil {
		return errors.Wrap(err, "fail to parse workflow")
	}

	return workflow.Build()
}
