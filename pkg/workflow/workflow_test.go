package workflow

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGlueWorkflow(t *testing.T) {
	yamlFileContent, _ := ioutil.ReadFile("example-0.1.yaml")
	workflow, err := Parse(yamlFileContent)

	assert.Nil(t, err)

	expected := &GlueWorkflow{
		Name: "my-workflow",
		Jobs: []GlueJob{
			{
				Name:        "ingestion",
				Description: "extract log from excel to s3",
				Type:        PythonJob,
				Entrypoint:  "raw/ingestion.py",
				Args: map[interface{}]interface{}{
					"source_path": "s3://sourceBucket/source/",
					"raw_path":    "s3://rawStorageBucket/raw/",
					"schema": map[interface{}]interface{}{
						"id":          "int",
						"description": "string",
					},
				},
				Requires: []RequiredJob{},
				Tags: []Tag{
					{
						Name: "team", Value: "fantastic-team",
					},
					{
						Name: "region", Value: "us-west-1",
					},
				},
			},
			{
				Name:        "transformation",
				Description: "transform and load",
				Type:        GlueSQLJob,
				Entrypoint:  "transformations/transform.sql",
				Args: map[interface{}]interface{}{
					"years": []interface{}{2021, 2022},
				},
				Requires: []RequiredJob{
					{
						JobName: "ingestion",
					},
				},
				Tags: []Tag{},
			},
			{
				Name:        "notification",
				Description: "dummy job",
				Type:        DummyJob,
				Entrypoint:  "",
				Args:        nil,
				Requires: []RequiredJob{
					{
						JobName: "transformation",
					},
				},
				Tags: []Tag{},
			},
		},
	}

	assert.Equal(t, expected, workflow)
}
