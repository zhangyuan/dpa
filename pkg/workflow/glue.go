package workflow

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

type GlueWorkflow struct {
	Workflow
	Name        string
	Description string
	Jobs        []GlueJob
	Schedule    Schedule
}

type GlueJob struct {
	Name        string
	Description string
	Type        string
	Entrypoint  string
	Args        interface{}
	Requires    []RequiredJob
	Tags        []Tag
}

type RequiredJob struct {
	JobName string
}

type Schedule struct {
	Cron string
}

const (
	PythonJob  string = "python"
	GlueSQLJob string = "glue-sql"
	DummyJob   string = "dummy"
)

func ParseJobType(rawType string) (string, error) {
	switch rawType {
	case PythonJob:
		return PythonJob, nil
	case GlueSQLJob:
		return GlueSQLJob, nil
	case DummyJob:
		return DummyJob, nil
	}

	return "", errors.Errorf("invalid job type %v", rawType)
}

func parseGlueWorkflow(rawWorkflow map[string]interface{}) (*GlueWorkflow, error) {
	rawJobs := rawWorkflow["jobs"].(map[interface{}]interface{})

	jobs := []GlueJob{}

	for key, value := range rawJobs {
		properties := value.(map[interface{}]interface{})

		jobType, err := ParseJobType(properties["type"].(string))
		if err != nil {
			return nil, err
		}

		var entrypoint string

		if properties["entrypoint"] == nil {
			entrypoint = ""
		} else {
			entrypoint = properties["entrypoint"].(string)
		}

		requiredJobs := []RequiredJob{}

		if properties["requires"] != nil {
			rawRequires := properties["requires"].([]interface{})
			for _, rawRequiredJob := range rawRequires {
				rawRequiredJob := rawRequiredJob.(map[interface{}]interface{})

				requiredJobs = append(requiredJobs, RequiredJob{
					JobName: rawRequiredJob["job_name"].(string),
				})
			}
		}

		tags := []Tag{}
		if properties["tags"] != nil {
			rawTags := properties["tags"].(map[interface{}]interface{})
			for name, value := range rawTags {
				tags = append(tags, Tag{
					Name:  name.(string),
					Value: value.(string),
				})
			}
		}

		job := GlueJob{
			Name:        key.(string),
			Description: properties["description"].(string),
			Type:        jobType,
			Entrypoint:  entrypoint,
			Args:        properties["args"],
			Requires:    requiredJobs,
			Tags:        tags,
		}

		jobs = append(jobs, job)
	}

	var schedule Schedule
	if rawWorkflow["schedule"] != nil {
		rawSchedule := rawWorkflow["schedule"].(map[interface{}]interface{})
		if rawSchedule["cron"] != nil {
			schedule = Schedule{
				Cron: rawSchedule["cron"].(string),
			}
		}
	}

	return &GlueWorkflow{
		Name:        rawWorkflow["name"].(string),
		Description: rawWorkflow["description"].(string),
		Jobs:        jobs,
		Schedule:    schedule,
	}, nil
}

func (workflow *GlueWorkflow) Render() (string, error) {
	stack := map[string]interface{}{}
	stack["AWSTemplateFormatVersion"] = "2010-09-09"
	stack["Description"] = workflow.Description

	awsGlueWorkflow := map[string]interface{}{
		"Type": "AWS::Glue::Workflow",
		"Properties": map[string]string{
			"Description": workflow.Description,
			"Name":        workflow.Name,
		},
	}

	resources := map[string]interface{}{
		workflow.Name: awsGlueWorkflow,
	}
	stack["Resources"] = resources

	template, err := yaml.Marshal(&stack)
	if err != nil {
		return "", errors.New("fail to marshal to yaml")
	}

	return string(template), nil
}
