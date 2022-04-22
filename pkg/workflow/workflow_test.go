package workflow

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	yamlFileContent, _ := ioutil.ReadFile("example-0.1.yaml")
	workflow, err := Parse(yamlFileContent)

	assert.Nil(t, err)

	assert.Equal(t, &Workflow{
		Version: "0.1",
		Name:    "my-workflow",
		Jobs: Jobs{Job{
			Name:        "ingestion",
			Description: "extract log from excel to s3",
		}, Job{
			Name:        "transform",
			Description: "load log to ods",
		}},
	}, workflow)
}
