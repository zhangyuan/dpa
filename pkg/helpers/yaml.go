package helpers

import (
	jsoniter "github.com/json-iterator/go"
	"gopkg.in/yaml.v2"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func YAML2Json(yamlContent string) (string, error) {
	var yamlData interface{}
	if err := yaml.Unmarshal([]byte(yamlContent), &yamlData); err != nil {
		return "", err
	}

	if jsonData, err := json.Marshal(yamlData); err != nil {
		return "", err
	} else {
		return string(jsonData), nil
	}
}
