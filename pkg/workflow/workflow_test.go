package workflow

import (
	"dp/pkg/helpers"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderGlueWorkflow(t *testing.T) {
	yamlFileContent, _ := ioutil.ReadFile("fixtures/sampleproject/sampleproject.yaml")
	workflow, err := Parse("fixtures/sampleproject", yamlFileContent)

	rendered, err := workflow.Render()
	assert.Nil(t, err)

	expectedYaml, _ := ioutil.ReadFile("fixtures/sampleproject/infra/stack.yaml")
	expectedData, expectedJsonErr := helpers.YAML2Json(string(expectedYaml))
	assert.Nil(t, expectedJsonErr)

	renderedJson, renderedJsonErr := helpers.YAML2Json(rendered)
	assert.Nil(t, renderedJsonErr)

	assert.JSONEq(t, expectedData, renderedJson)
}
