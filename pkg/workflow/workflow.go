package workflow

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Workflow interface{}

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
