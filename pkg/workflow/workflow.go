package workflow

import (
	"sort"

	"dp/pkg/helpers"

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
	Jobs     Jobs
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

	var jobsMapSlice yaml.MapSlice

	if err := unmarshal(&jobsMapSlice); err != nil {
		return errors.Wrap(err, "fail to unmarshal jobs")
	}

	var sortedJobNames []string
	for _, item := range jobsMapSlice {
		sortedJobNames = append(sortedJobNames, item.Key.(string))
	}

	var indexOf = helpers.IndexOf(sortedJobNames)

	sort.SliceStable(*jobs, func(i, j int) bool {
		return indexOf((*jobs)[i].Name) < indexOf((*jobs)[j].Name)
	})

	return nil
}

func (tags *Tags) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var tagsMap yaml.MapSlice

	if err := unmarshal(&tagsMap); err == nil {
		for _, value := range tagsMap {
			*tags = append(*tags, Tag{Name: value.Key.(string), Value: value.Value.(string)})
		}
	} else {
		return errors.Wrap(err, "fail to unmarshal Tags")
	}

	return nil
}

type stepDefinition struct {
	Job          string
	AllowFailure bool `yaml:"allow_failure"`
}

func (step *Step) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var jobName string

	if err := unmarshal(&jobName); err == nil {
		*step = Step{Job: jobName, AllowFailure: false}
		return nil
	}

	var stepDefinition stepDefinition

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

	var argumentsMapSlice yaml.MapSlice
	if err := unmarshal(&argumentsMapSlice); err != nil {
		return errors.Wrap(err, "fail to unmarshal arguments")
	}
	var argumentNames []string
	for _, item := range argumentsMapSlice {
		argumentNames = append(argumentNames, item.Key.(string))
	}

	indexOf := helpers.IndexOf(argumentNames)

	sort.SliceStable(*arguments, func(i, j int) bool {
		return indexOf((*arguments)[i].Name) < indexOf((*arguments)[j].Name)
	})

	return nil
}

func (jobType *JobType) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var jobTypeString string
	if err := unmarshal(&jobTypeString); err != nil {
		return errors.Wrap(err, "invalid job type")
	}
	*jobType = asJobType(jobTypeString)
	return nil
}

func asJobType(jobTypeInterface interface{}) JobType {
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
