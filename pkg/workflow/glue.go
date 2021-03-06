package workflow

import (
	"dpa/pkg/helpers"
	"fmt"
	"os"
	"strings"

	"github.com/heimdalr/dag"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"gopkg.in/yaml.v2"
)

type GlueWorkflow struct {
	Workflow
	ProjectDirectory string
	Name             string
	Description      string
	Jobs             []GlueJob
	Schedule         Schedule
	IamRole          string
	PythonModules    []string
	ArtifactsPath    string
	Tags             []Tag
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
	PySparkJob string = "pyspark"
	DummyJob   string = "dummy"
)

type PythonModule struct {
	Name    string
	Version string
}

func ParseJobType(rawType string) (string, error) {
	switch rawType {
	case PythonJob:
		return PythonJob, nil
	case PySparkJob:
		return PySparkJob, nil
	case DummyJob:
		return DummyJob, nil
	}

	return "", errors.Errorf("invalid job type %v", rawType)
}

func parseGlueWorkflow(projectDirectory string, rawWorkflow map[string]interface{}) (*GlueWorkflow, error) {
	var iamRole string

	if rawWorkflow["iam_role"] != nil {
		iamRole = rawWorkflow["iam_role"].(string)
	}

	var artifactsPath string
	if rawWorkflow["artifacts_path"] != nil {
		artifactsPath = strings.TrimSuffix(rawWorkflow["artifacts_path"].(string), "/")
	}

	jobs := []GlueJob{}
	rawJobs := rawWorkflow["jobs"].(map[interface{}]interface{})

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

	tags := []Tag{}
	if rawWorkflow["tags"] != nil {
		rawTags := rawWorkflow["tags"].(map[interface{}]interface{})

		for k, v := range rawTags {
			tags = append(tags, Tag{Name: k.(string), Value: v.(string)})
		}
	}

	var pythonModules []string
	if rawWorkflow["python_modules"] != nil {
		rawPythonModules := rawWorkflow["python_modules"].([]interface{})
		pythonModules = lo.Map(rawPythonModules, func(moduleName interface{}, i int) string {
			return moduleName.(string)
		})
	} else {
		pythonModules = []string{}
	}

	return &GlueWorkflow{
		ProjectDirectory: projectDirectory,
		Name:             rawWorkflow["name"].(string),
		Description:      rawWorkflow["description"].(string),
		IamRole:          iamRole,
		Jobs:             jobs,
		Schedule:         schedule,
		Tags:             tags,
		PythonModules:    pythonModules,
		ArtifactsPath:    artifactsPath,
	}, nil
}

func (workflow *GlueWorkflow) Render() (string, error) {
	stack := map[string]interface{}{}
	stack["AWSTemplateFormatVersion"] = "2010-09-09"
	stack["Description"] = workflow.Description

	awsGlueWorkflowProperties := map[string]interface{}{
		"Description": workflow.Description,
		"Name":        workflow.Name,
	}

	if workflow.Tags != nil && len(workflow.Tags) > 0 {
		tags := make(map[string]interface{})
		lo.ForEach(workflow.Tags, func(t Tag, i int) {
			tags[t.Name] = t.Value
		})
		awsGlueWorkflowProperties["Tags"] = tags
	}

	awsGlueWorkflow := map[string]interface{}{
		"Type":       "AWS::Glue::Workflow",
		"Properties": awsGlueWorkflowProperties,
	}

	resources := map[string]interface{}{
		workflow.ResourceId(): awsGlueWorkflow,
	}

	rootJobsNames, err := workflow.GetRootJobs()
	if err != nil {
		return "", errors.Wrap(err, "fail to get root jobs")
	}

	workflowStartTrigger := map[string]interface{}{
		"Type": "AWS::Glue::Trigger",
		"Properties": map[string]interface{}{
			"Name":            fmt.Sprintf("trigger-start-%s", workflow.Name),
			"Type":            "SCHEDULED",
			"Schedule":        fmt.Sprintf("cron(%s)", workflow.Schedule.Cron),
			"StartOnCreation": true,
			"WorkflowName": map[string]interface{}{
				"Ref": workflow.ResourceId(),
			},
			"Actions": lo.Map(rootJobsNames, func(jobName string, i int) map[string]interface{} {
				return map[string]interface{}{
					"JobName": map[string]interface{}{
						"Ref": ToJobResourceId(jobName),
					},
				}
			}),
		},
	}

	workflowStartResourceName := fmt.Sprintf("TriggerStart%s", workflow.ResourceId())

	resources[workflowStartResourceName] = workflowStartTrigger

	// render job
	for _, job := range workflow.Jobs {
		resourceName := job.ResourceId()

		var jobRole string
		var commandName string
		var defaultArguments map[string]interface{}

		if len(job.Role) > 0 {
			jobRole = job.Role
		} else {
			jobRole = workflow.IamRole
		}

		defaultArguments = map[string]interface{}{}

		if len(workflow.PythonModules) > 0 {
			extraPyFiles := lo.Map(workflow.PythonModules, func(moduleName string, i int) string {
				if strings.HasPrefix(moduleName, "s3://") {
					return moduleName
				}
				return fmt.Sprintf("%s/%s", workflow.ArtifactsPath, moduleName)
			})
			defaultArguments["--extra-py-files"] = strings.Join(extraPyFiles, ",")
		}

		if job.Args != nil {
			arguments, err := json.Marshal(job.Args)
			if err != nil {
				message := fmt.Sprintf("invalid arguments: %v", job.Args)
				return "", errors.Wrap(err, message)
			}

			defaultArguments["--arguments"] = string(arguments)
		}

		var scriptLocation string
		if !strings.HasPrefix(job.Entrypoint, "s3://") && workflow.ArtifactsPath != "" {
			scriptLocation = workflow.ArtifactsPath + "/" + job.Entrypoint
		} else {
			scriptLocation = job.Entrypoint
		}

		var jobTags map[string]interface{}

		if job.Tags != nil && len(job.Tags) > 0 {
			jobTags = make(map[string]interface{})
			lo.ForEach(job.Tags, func(t Tag, i int) {
				jobTags[t.Name] = t.Value
			})
		}

		if job.Type == PythonJob {
			commandName = "pythonshell"

			properties := map[string]interface{}{
				"Command": map[string]interface{}{
					"Name":           commandName,
					"PythonVersion":  "3",
					"ScriptLocation": scriptLocation,
				},
				"DefaultArguments": defaultArguments,
				"Role":             jobRole,
				"Name":             job.Name,
			}

			if jobRole != "" {
				properties["Role"] = jobRole
			}

			if defaultArguments != nil {
				properties["DefaultArguments"] = defaultArguments
			}

			if len(jobTags) > 0 {
				properties["Tags"] = jobTags
			}

			resources[resourceName] = map[string]interface{}{
				"Type":       "AWS::Glue::Job",
				"Properties": properties,
			}
		} else if job.Type == PySparkJob {
			commandName = "glueetl"

			properties := map[string]interface{}{
				"Command": map[string]interface{}{
					"Name":           commandName,
					"PythonVersion":  "3",
					"ScriptLocation": scriptLocation,
				},
				"DefaultArguments": defaultArguments,
				"Role":             jobRole,
				"Name":             job.Name,
				"GlueVersion":      "3.0",
				"MaxCapacity":      2,
				"ExecutionProperty": map[string]interface{}{
					"MaxConcurrentRuns": 1,
				},
			}

			if jobRole != "" {
				properties["Role"] = jobRole
			}

			if defaultArguments != nil {
				properties["DefaultArguments"] = defaultArguments
			}

			if len(jobTags) > 0 {
				properties["Tags"] = jobTags
			}

			resources[resourceName] = map[string]interface{}{
				"Type":       "AWS::Glue::Job",
				"Properties": properties,
			}
		}
	}

	// render trigger
	for _, job := range workflow.Jobs {
		if lo.Contains(rootJobsNames, job.Name) {
			continue
		}

		conditions := lo.Map(job.Requires, func(rj RequiredJob, i int) map[string]interface{} {
			return map[string]interface{}{
				"JobName": map[string]interface{}{
					"Ref": ToJobResourceId(rj.JobName),
				},
				"LogicalOperator": "EQUALS",
				"State":           "SUCCEEDED",
			}
		})

		properties := map[string]interface{}{
			"Description":     fmt.Sprintf("trigger %s", job.Name),
			"StartOnCreation": true,
		}

		if len(conditions) > 0 {
			properties["Predicate"] = map[string]interface{}{
				"Conditions": conditions,
				"Logical":    "AND",
			}
		}

		actions := []map[string]interface{}{
			{
				"JobName": map[string]interface{}{
					"Ref": job.ResourceId(),
				},
			},
		}
		properties["Actions"] = actions
		properties["Type"] = "CONDITIONAL"
		properties["WorkflowName"] = map[string]interface{}{
			"Ref": workflow.ResourceId(),
		}

		trigger := map[string]interface{}{
			"Type":       "AWS::Glue::Trigger",
			"Properties": properties,
		}

		triggerName := ToTriggerResourceId(job.Name)

		resources[triggerName] = trigger
	}

	stack["Resources"] = resources

	template, err := yaml.Marshal(&stack)
	if err != nil {
		return "", errors.New("fail to marshal to yaml")
	}

	return string(template), nil
}

func (workflow *GlueWorkflow) Build() error {
	workingDirectory := strings.TrimSuffix(workflow.ProjectDirectory, "/")

	buildDirectory := workingDirectory + "/build"

	projectDirectory := buildDirectory + "/project"
	err := helpers.Mkdirp(projectDirectory)
	if err != nil {
		return errors.Wrap(err, "fail to create project directory")
	}

	cloudformationFilename := "dp-cloudformation.yaml"
	outputPath := projectDirectory + "/" + cloudformationFilename

	outputContent, err := workflow.Render()
	if err != nil {
		return errors.Wrap(err, "fail to render workflow")
	}

	err = os.WriteFile(outputPath, []byte(outputContent), 0644)
	if err != nil {
		return errors.Wrap(err, "fail to write to "+outputPath)
	}

	artifactsDirectory := buildDirectory + "/artifacts"
	err = helpers.Mkdirp(artifactsDirectory)
	if err != nil {
		return errors.Wrap(err, "fail to create artifacts directory")
	}

	for _, job := range workflow.Jobs {
		if job.Type == PythonJob || job.Type == PySparkJob {
			if strings.HasPrefix(job.Entrypoint, "s3://") {
				continue
			}
			fromPath := workingDirectory + "/" + job.Entrypoint
			toPath := artifactsDirectory + "/" + job.Entrypoint

			if err := helpers.Copy(fromPath, toPath); err != nil {
				return err
			}
		}
	}

	if err := helpers.Copy(outputPath, artifactsDirectory+"/"+cloudformationFilename); err != nil {
		return err
	}

	return nil
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

func (workflow *GlueWorkflow) ResourceId() string {
	return fmt.Sprintf("Workflow%s", NormalizeResourceId(workflow.Name))
}

func (workflow *GlueWorkflow) GetRootJobs() ([]string, error) {
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
				return nil, errors.Errorf("invalid job name %s", requiredJob.JobName)
			}

			_ = d.AddEdge(jobNameToVertex[requiredJobName], jobNameToVertex[job.Name])
		}
	}

	rootJobNames := []string{}

	for _, jobName := range d.GetRoots() {
		rootJobNames = append(rootJobNames, jobName.(string))
	}

	return rootJobNames, nil

}

func (glueJob *GlueJob) ResourceId() string {
	return ToJobResourceId(glueJob.Name)
}

func ToJobResourceId(jobName string) string {
	return fmt.Sprintf("Job%s", NormalizeResourceId(jobName))
}

func ToTriggerResourceId(jobName string) string {
	return fmt.Sprintf("Trigger%s", NormalizeResourceId(jobName))
}

func NormalizeResourceId(name string) string {
	name = strings.ReplaceAll(name, "_", " ")
	name = strings.ReplaceAll(name, "-", " ")
	name = strings.Title(name)
	return strings.ReplaceAll(name, " ", "")
}
