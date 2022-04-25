package workflow

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type JobType int64

const (
	Unknown JobType = iota
	Python
	GlueSQL
)

type Argument struct {
	Name  string
	Value interface{}
}

type Arguments []Argument

type Tag struct {
	Name  string
	Value string
}

type Tags []Tag

type Job struct {
	Name        string
	Description string
	Type        JobType
	Entrypoint  string
	Arguments   Arguments
	Tags        Tags
}

type Jobs []Job

type Workflow struct {
	Version string
	Name    string
	Jobs    Jobs `yaml:"jobs"`
	Tags    Tags
}

func (e *Jobs) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var jobsMap map[string]map[string]interface{}

	if err := unmarshal(&jobsMap); err == nil {
		for jobName, properties := range jobsMap {
			arguments := Arguments{}
			if properties["args"] != nil {
				arguments = *AsArguments(properties["args"])
			}

			tags := Tags{}
			if properties["tags"] != nil {
				tags = *AsTags(properties["tags"])
			}

			job := Job{
				Name:        jobName,
				Description: AsDescription(properties["description"]),
				Type:        AsJobType(properties["type"].(string)),
				Entrypoint:  properties["entrypoint"].(string),
				Arguments:   arguments,
				Tags:        tags,
			}
			*e = append(*e, job)
		}
	} else {
		return errors.Wrap(err, "fail to unmarshal Jobs")
	}

	return nil
}

func (tags *Tags) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tagsMap map[string]string

	if err := unmarshal(&tagsMap); err == nil {
		for name, value := range tagsMap {
			*tags = append(*tags, Tag{Name: name, Value: value})
		}
	} else {
		return errors.Wrap(err, "fail to unmarshal Tags")
	}

	return nil
}

func AsDescription(descriptionInterface interface{}) string {
	return descriptionInterface.(string)
}

func AsJobType(jobTypeInterface interface{}) JobType {
	switch jobTypeInterface {
	case "python":
		return Python
	case "glue-sql":
		return GlueSQL
	default:
		return Unknown
	}
}

func AsArguments(argumentsInterface interface{}) *Arguments {
	arguments := Arguments{}
	argumementsList := argumentsInterface.(map[interface{}]interface{})
	for nameInterface, value := range argumementsList {
		name := nameInterface.(string)
		arguments = append(arguments, Argument{Name: name, Value: value})
	}
	return &arguments
}

func AsTags(tagsInterface interface{}) *Tags {
	tags := Tags{}
	argumementsList := tagsInterface.(map[interface{}]interface{})
	for nameInterface, value := range argumementsList {
		tags = append(tags, Tag{Name: nameInterface.(string), Value: value.(string)})
	}
	return &tags
}

func Parse(content []byte) (*Workflow, error) {
	workflow := Workflow{}
	err := yaml.Unmarshal(content, &workflow)
	if err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal yaml content")
	}

	return &workflow, nil
}
