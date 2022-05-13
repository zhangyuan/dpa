package python

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePythonRequirements(t *testing.T) {
	modules, err := GetPythonRequirements("fixtures")
	assert.Nil(t, err)

	assert.Equal(t, PythonModules{
		{Name: "openpyxl", Version: "3.0.9"},
		{Name: "Office365-REST-Python-Client", Version: ""},
	}, modules)
}
