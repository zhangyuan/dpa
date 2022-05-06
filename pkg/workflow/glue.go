package workflow

import "github.com/pkg/errors"

type GlueWorkflow struct {
	Workflow
	Name string
	Jobs []GlueJob
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
	return &GlueWorkflow{
		Name: rawWorkflow["name"].(string),
		Jobs: jobs,
	}, nil
}
