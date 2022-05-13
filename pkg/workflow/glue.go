package workflow

import (
	"fmt"

	"github.com/heimdalr/dag"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"gopkg.in/yaml.v2"
)

type GlueWorkflow struct {
	Workflow
	Name        string
	Description string
	Jobs        []GlueJob
	Schedule    Schedule
	IamRole     string
}

type GlueJob struct {
	Name        string
	Description string
	Type        string
	Entrypoint  string
	Args        interface{}
	Requires    []RequiredJob
	Tags        []Tag
	Role        string
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
	var iamRole string
	if rawWorkflow["iam_role"] != nil {
		iamRole = rawWorkflow["iam_role"].(string)
	}

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

		var jobRole string
		if properties["roe"] != nil {
			jobRole = properties["roe"].(string)
		} else {
			jobRole = iamRole
		}

		job := GlueJob{
			Name:        key.(string),
			Description: properties["description"].(string),
			Type:        jobType,
			Entrypoint:  entrypoint,
			Args:        properties["args"],
			Requires:    requiredJobs,
			Role:        jobRole,
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
		IamRole:     iamRole,
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
			"Name":        fmt.Sprintf("Workflow_%s", workflow.Name),
		},
	}

	resources := map[string]interface{}{
		workflow.Name: awsGlueWorkflow,
	}

	for _, job := range workflow.Jobs {
		if job.Type != PythonJob {
			continue
		}

		resourceName := fmt.Sprintf("Job_%s", job.Name)
		var commandName string
		if job.Type == PythonJob {
			commandName = "pythonshell"
			var jobRole string
			if len(job.Role) > 0 {
				jobRole = job.Role
			} else {
				jobRole = workflow.IamRole
			}
			arguments, err := json.Marshal(job.Args)
			if err != nil {
				message := fmt.Sprintf("invalid arguments: %v", job.Args)
				return "", errors.Wrap(err, message)
			}

			resources[resourceName] = map[string]interface{}{
				"Type": "AWS::Glue::Job",
				"Properties": map[string]interface{}{
					"Command": map[string]interface{}{
						"Name":           commandName,
						"PythonVersion":  "3",
						"ScriptLocation": job.Entrypoint,
					},
					"DefaultArguments": map[string]interface{}{
						"--arguments": string(arguments),
					},
					"Role": jobRole,
				},
			}
		}
	}

	for _, job := range workflow.Jobs {
		conditions := lo.Map(job.Requires, func(rj RequiredJob, i int) map[string]interface{} {
			return map[string]interface{}{
				"JobName":         rj.JobName,
				"LogicalOperator": "EQUALS",
				"State":           "SUCCEEDED",
			}
		})

		properties := map[string]interface{}{
			"Description": fmt.Sprintf("trigger %s", job.Name),
		}

		if len(conditions) > 0 {
			properties["Predicate"] = map[string]interface{}{
				"Conditions": conditions,
			}
		}
		trigger := map[string]interface{}{
			"Type":       "AWS::Glue::Trigger",
			"Properties": properties,
		}

		triggerName := fmt.Sprintf("Trigger_%s", job.Name)
		resources[triggerName] = trigger
	}

	stack["Resources"] = resources

	template, err := yaml.Marshal(&stack)
	if err != nil {
		return "", errors.New("fail to marshal to yaml")
	}

	return string(template), nil
}

// WIP
func (workflow *GlueWorkflow) Dag() error {
	d := dag.NewDAG()

	jobs := &workflow.Jobs

	jobNameToVertex := make(map[string]string)

	for _, job := range *jobs {
		v1, _ := d.AddVertex(job.Name)
		jobNameToVertex[job.Name] = v1
	}

	jobNames := lo.Map(workflow.Jobs, func(t GlueJob, i int) string {
		return t.Name
	})

	for _, job := range *jobs {
		for _, requiredJob := range job.Requires {
			requiredJobName, ok := lo.Find(jobNames, func(jobName string) bool {
				return requiredJob.JobName == jobName
			})
			if !ok {
				return errors.Errorf("invalid job name %s", requiredJob.JobName)
			}

			_ = d.AddEdge(jobNameToVertex[requiredJobName], jobNameToVertex[job.Name])
		}
	}

	// var walk func(d *dag.DAG, id string) error

	// walk = func(d *dag.DAG, id string) error {
	// 	jobName, err := d.GetVertex(id)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if children, err := d.GetChildren(id); err != nil {
	// 		return err
	// 	} else {
	// 		for childId, childJobName := range children {
	// 			fmt.Sprintln("%s => %s", jobName, childJobName)
	// 			if err := walk(d, childId); err != nil {
	// 				return err
	// 			}
	// 		}
	// 	}
	// 	return nil
	// }

	// for id := range d.GetRoots() {
	// 	walk(d, id)
	// }
	// fmt.Println("")

	for id, jobName := range d.GetVertices() {
		fmt.Println("jobName: ", jobName)
		children, _ := d.GetChildren(id)

		for _, child := range children {
			fmt.Println("child:", child)
		}

		ancestors, _ := d.GetAncestors(id)
		for _, a := range ancestors {
			fmt.Println("ancestor:", a)
		}

		fmt.Println("")
		fmt.Println("")
	}

	return nil
}
