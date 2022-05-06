package workflow

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Workflow interface{}

type Tag struct {
	Name  string
	Value string
}

func Parse(content []byte) (Workflow, error) {
	var rawWorkflow map[string]interface{}
	err := yaml.Unmarshal(content, &rawWorkflow)
	if err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal yaml content")
	}

	if rawWorkflow["vendor"] == "glue" {
		return parseGlueWorkflow(rawWorkflow)
	}

	return nil, errors.New("invalid vendor.")
}
