package workflow

import (
	"sort"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type JobType int64

const (
	Unknown JobType = iota
	Python
	GlueSQL
	Dummy
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

type Schedule struct {
	Cron string
}

type Job struct {
	Name        string
	Description string
	Type        JobType
	Entrypoint  string
	Arguments   Arguments `yaml:"args"`
	Tags        Tags
}

type Jobs []Job

type Step struct {
	Job          string
	AllowFailure bool
}

type Steps []Step
type Workflow struct {
	Version  string
	Name     string
	Tags     Tags
	Schedule Schedule
	Jobs     Jobs `yaml:"jobs"`
	Steps    Steps
}

func (jobs *Jobs) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var jobsMap map[string]Job

	if err := unmarshal(&jobsMap); err == nil {
		for name, job := range jobsMap {
			job.Name = name

			*jobs = append(*jobs, job)
		}
	} else {
		return errors.Wrap(err, "fail to unmarshal Jobs")
	}

	sort.SliceStable(*jobs, func(i, j int) bool {
		return (*jobs)[i].Name < (*jobs)[j].Name
	})

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

type StepDefinition struct {
	Job          string
	AllowFailure bool `yaml:"allow_failure"`
}

func (step *Step) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var jobName string

	if err := unmarshal(&jobName); err == nil {
		*step = Step{Job: jobName, AllowFailure: false}
		return nil
	}

	var stepDefinition StepDefinition

	if err := unmarshal(&stepDefinition); err == nil {
		*step = Step{Job: stepDefinition.Job, AllowFailure: stepDefinition.AllowFailure}
		return nil
	}

	return errors.New("could not unmarshal step")
}

func (arguments *Arguments) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var argumentsMap map[string]interface{}
	if err := unmarshal(&argumentsMap); err == nil {
		if len(argumentsMap) == 0 {
			*arguments = Arguments{}
		} else {
			for name, value := range argumentsMap {
				*arguments = append(*arguments, Argument{
					Name:  name,
					Value: value,
				})
			}
		}
	} else {
		return errors.Wrap(err, "fail to unmarshal Arguments")
	}

	sort.SliceStable(*arguments, func(i, j int) bool {
		return (*arguments)[i].Name < (*arguments)[j].Name
	})

	return nil
}

func (jobType *JobType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var jobTypeString string
	if err := unmarshal(&jobTypeString); err != nil {
		return errors.Wrap(err, "invalid job type")
	}
	*jobType = AsJobType(jobTypeString)
	return nil
}

func AsJobType(jobTypeInterface interface{}) JobType {
	switch jobTypeInterface {
	case "python":
		return Python
	case "glue-sql":
		return GlueSQL
	case "dummy":
		return Dummy
	default:
		return Unknown
	}
}

func Parse(content []byte) (*Workflow, error) {
	workflow := Workflow{}
	err := yaml.Unmarshal(content, &workflow)
	if err != nil {
		return nil, errors.Wrap(err, "fail to unmarshal yaml content")
	}

	return &workflow, nil
}
