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

	expected := &Workflow{
		Version: "0.1",
		Name:    "my-workflow",
		Jobs: Jobs{Job{
			Name:        "ingestion",
			Description: "extract log from excel to s3",
			Entrypoint:  "raw/ingestion.py",
			Type:        Python,
		}, Job{
			Name:        "transform",
			Description: "transform and load",
			Entrypoint:  "transformations/transform.sql",
			Type:        GlueSQL,
		}},
	}

	assert.Equal(t, expected, workflow)
}
