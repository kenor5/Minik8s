package entity

type Workflow struct {
	Kind string `json:"kind" yaml:"kind"`
	StartAt string `json:"startAt" yaml:"startAt"`
	WorkflowNodes []WorkflowNode `json:"WorkflowNodes,omitempty" yaml:"WorkflowNodes,omitempty"`
}

type WorkflowNode struct {
	Type string `json:"type" yaml:"type"`
    Name string `json:"Name" yaml:"Name"`

	Next string `json:"Next,omitempty" yaml:"Next,omitempty"`
	End string `json:"End,omitempty" yaml:"End,omitempty"`
	
	Choices []Choice `json:"Choices,omitempty" yaml:"Choices,omitempty"`
}

type Choice struct {
	Variable string `json:"variable" yaml:"variable"`
	Next string `json:"Next" yaml:"Next"`
	Condition string `json:"Condition" yaml:"Condition"`
	Number int64 `json:"Number" yaml:"Number"`
}

