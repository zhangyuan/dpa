package workflow

type GlueWorkflow struct {
	Workflow
}

func parseGlueWorkflow(rawWorkflow map[string]interface{}) (*GlueWorkflow, error) {
	return &GlueWorkflow{}, nil
}
