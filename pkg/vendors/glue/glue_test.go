package glue

import (
	"dp/pkg/workflow"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRender(t *testing.T) {
	workflow := &workflow.Workflow{
		Version: "0.1",
		Name:    "my-workflow",
		Jobs: workflow.Jobs{
			workflow.Job{
				Name:        "ingestion",
				Description: "ingest from raw bucket",
				Entrypoint:  "raw/ingestion.py",
				Type:        workflow.Python,
				Arguments: workflow.Arguments{
					{Name: "source_path", Value: "s3://sourceBucket/source/"},
					{Name: "raw_path", Value: "s3://rawStorageBucket/raw/"},
					{Name: "schema", Value: map[interface{}]interface{}{
						"id":          "int",
						"description": "string",
					}},
				},
				Tags: workflow.Tags{
					{Name: "team", Value: "fantastic-team"},
					{Name: "region", Value: "us-west-1"},
				},
			},
		},
	}

	writer := NewwGlueWriter(workflow)

	output, err := writer.Render()

	assert.Nil(t, err)

	expectedOutput := `{
		"ingestion": {
			"Type" : "AWS::Glue::Job",
			"Properties": {
				"Name": "ingestion",
				"Description": "ingest from raw bucket",
				"Command": {
					"Name": "pythonshell",
					"PythonVersion": "3",
					"ScriptLocation": "raw/ingestion.py"
				},
				"MaxCapacity": 2,
				"ExecutionProperty": {
					"MaxConcurrentRuns": 2
				},
				"DefaultArguments": {
					"--enable-metrics": "true",
					"--enable-continuous-cloudwatch-log": "true"
				},
				"GlueVersion": "3.0"
			}
		}
	}`

	fmt.Println(expectedOutput)
	fmt.Println(output)
	assert.JSONEq(t, expectedOutput, output)
}
