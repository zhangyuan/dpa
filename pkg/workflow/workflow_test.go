package workflow

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestParseGlueWorkflow(t *testing.T) {
	yamlFileContent, _ := ioutil.ReadFile("example-0.1.yaml")
	workflow, err := Parse("fixtures/sampleproject", yamlFileContent)

	assert.Nil(t, err)

	expected := &GlueWorkflow{
		ProjectDirectory: "fixtures/sampleproject",
		Name:             "my-workflow",
		Description:      "my workflow",
		Schedule: Schedule{
			Cron: "00 20 * * ? *",
		},
		IamRole: "iam-role-arn",
		PythonModules: []string{
			"sampleproject",
		},
		Tags: []Tag{
			{Name: "lob", Value: "sales"},
		},
		Jobs: []GlueJob{
			{
				Name:        "ingestion",
				Description: "extract log from excel to s3",
				Type:        PythonJob,
				Entrypoint:  "sampleproject/jobs/ingestion.py",
				Args: map[interface{}]interface{}{
					"source_path": "s3://sourceBucket/source/",
					"raw_path":    "s3://rawStorageBucket/raw/",
					"schema": map[interface{}]interface{}{
						"id":          "int",
						"description": "string",
					},
				},
				Role:     "iam-role-arn",
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
				Type:        PythonJob,
				Entrypoint:  "sampleproject/jobs/transform.py",
				Args: map[interface{}]interface{}{
					"years": []interface{}{2021, 2022},
				},
				Role: "iam-role-arn",
				Requires: []RequiredJob{
					{
						JobName: "ingestion",
					},
				},
				Tags: []Tag{},
			},
		},
	}

	assert.Equal(t, expected, workflow)
}

func TestRenderGlueWorkflow(t *testing.T) {
	workflow := glueWorkflowFixture()

	rendered, err := workflow.Render()
	assert.Nil(t, err)

	expectedYaml, _ := ioutil.ReadFile("example-glue.yaml")
	expectedData, expectedJsonErr := yaml2Json(string(expectedYaml))
	assert.Nil(t, expectedJsonErr)

	renderedJson, renderedJsonErr := yaml2Json(rendered)
	assert.Nil(t, renderedJsonErr)

	assert.JSONEq(t, expectedData, renderedJson)
}

func yaml2Json(yamlContent string) (string, error) {
	var yamlData interface{}
	if err := yaml.Unmarshal([]byte(yamlContent), &yamlData); err != nil {
		return "", err
	}

	if jsonData, err := json.Marshal(yamlData); err != nil {
		return "", err
	} else {
		fmt.Println(string(jsonData))
		return string(jsonData), nil
	}
}

func TestGetDag(t *testing.T) {
	workflow := glueWorkflowFixture()
	err := workflow.Dag()
	assert.Nil(t, err)
}

func glueWorkflowFixture() Workflow {
	workflow := GlueWorkflow{
		ProjectDirectory: "fixtures/lakefoundation-project",
		Name:             "my-workflow",
		Description:      "my workflow",
		IamRole:          "iam-role-arn",
		Schedule: Schedule{
			Cron: "00 20 * * ? *",
		},
		Tags: []Tag{
			{Name: "lob", Value: "sales"},
		},
		Jobs: []GlueJob{
			{
				Name:        "ingestion",
				Description: "extract log from excel to s3",
				Type:        PythonJob,
				Entrypoint:  "ingestion/ingestion.py",
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
				Type:        PythonJob,
				Entrypoint:  "transformations/transform.py",
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
		},
	}

	return &workflow
}
