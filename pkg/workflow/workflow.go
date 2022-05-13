package workflow

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Workflow interface {
	Render() (string, error)
	Dag() error
}

type Tag struct {
	Name  string
	Value string
}

func Parse(projectDirectory string, content []byte) (Workflow, error) {
	var rawWorkflow map[string]interface{}
	err := yaml.Unmarshal(content, &rawWorkflow)
	if err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal yaml content")
	}

	if rawWorkflow["vendor"] == "glue" {
		return parseGlueWorkflow(projectDirectory, rawWorkflow)
	}

	return nil, errors.New("invalid vendor.")
}
