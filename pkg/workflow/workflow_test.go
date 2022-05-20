package workflow

import (
	"dpa/pkg/helpers"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseAndRenderCloudformationForGlueWorkflow(t *testing.T) {
	yamlFileContent, _ := ioutil.ReadFile("fixtures/sampleproject/sampleproject.yaml")
	workflow, err := Parse("fixtures/sampleproject", yamlFileContent)
	assert.Nil(t, err)

	rendered, err := workflow.Render()
	assert.Nil(t, err)

	expectedYaml, _ := ioutil.ReadFile("fixtures/sampleproject/infra/stack.yaml")
	expectedData, expectedJsonErr := helpers.YAML2Json(string(expectedYaml))
	assert.Nil(t, expectedJsonErr)

	renderedJson, renderedJsonErr := helpers.YAML2Json(rendered)
	assert.Nil(t, renderedJsonErr)

	assert.JSONEq(t, expectedData, renderedJson)
}
