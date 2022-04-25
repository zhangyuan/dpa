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
		Tags: Tags{
			Tag{Name: "lob", Value: "sales"},
		},
		Jobs: Jobs{Job{
			Name:        "ingestion",
			Description: "extract log from excel to s3",
			Entrypoint:  "raw/ingestion.py",
			Type:        Python,
			Arguments: Arguments{
				{Name: "source_path", Value: "s3://sourceBucket/source/"},
				{Name: "raw_path", Value: "s3://rawStorageBucket/raw/"},
				{Name: "schema", Value: map[interface{}]interface{}{
					"id":          "int",
					"description": "string",
				}},
			},
			Tags: Tags{
				{Name: "team", Value: "fantastic-team"},
				{Name: "region", Value: "us-west-1"},
			},
		}, Job{
			Name:        "transform",
			Description: "transform and load",
			Entrypoint:  "transformations/transform.sql",
			Type:        GlueSQL,
			Arguments:   Arguments{},
			Tags:        Tags{},
		}},
	}

	assert.Equal(t, expected, workflow)
}
