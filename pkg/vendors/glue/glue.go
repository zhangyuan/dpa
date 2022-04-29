package glue

import (
	"dp/pkg/workflow"
	"encoding/json"
)

type GlueWriter struct {
	Workflow *workflow.Workflow
}

func NewwGlueWriter(workflow *workflow.Workflow) *GlueWriter {
	return &GlueWriter{
		Workflow: workflow,
	}
}

func (w *GlueWriter) Render() (string, error) {
	glueJobs := make(map[string]interface{})

	for _, job := range w.Workflow.Jobs {
		var jobCommand = make(map[string]interface{})
		if job.Type == workflow.Python {
			jobCommand = map[string]interface{}{
				"Name":           "pythonshell",
				"PythonVersion":  "3",
				"ScriptLocation": job.Entrypoint,
			}

		}
		glueJobs[job.Name] = map[string]interface{}{
			"Type": "AWS::Glue::Job",
			"Properties": map[string]interface{}{
				"Name":        job.Name,
				"Description": job.Description,
				"Command":     jobCommand,
				"MaxCapacity": 2,
				"ExecutionProperty": map[string]interface{}{
					"MaxConcurrentRuns": 2,
				},
				"DefaultArguments": map[string]interface{}{
					"--enable-metrics":                   "true",
					"--enable-continuous-cloudwatch-log": "true",
				},
				"GlueVersion": "3.0",
			},
		}
	}

	if jsonContent, err := json.Marshal(glueJobs); err != nil {
		return "", err
	} else {
		return string(jsonContent), nil
	}
}
