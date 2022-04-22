package workflow

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type Job struct {
	Name        string
	Description string
	Type        string
	Entrypoint  string
	args        map[string]string
	tags        map[string]string
}

type Jobs []Job

type Workflow struct {
	Version string
	Name    string
	Jobs    Jobs `yaml:"jobs"`
}

func (e *Jobs) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var jobsMap map[string]map[string]interface{}

	if err := unmarshal(&jobsMap); err == nil {
		for jobName, properties := range jobsMap {
			job := Job{
				Name:        jobName,
				Description: properties["description"].(string),
			}
			*e = append(*e, job)
		}
	}

	return nil

}

func Parse(content []byte) (*Workflow, error) {
	workflow := Workflow{}
	err := yaml.Unmarshal(content, &workflow)
	if err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal yaml content")
	}

	return &workflow, nil
}
