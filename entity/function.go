package entity

type Function struct {
	Kind string `json:"kind" yaml:"kind"`
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	FunctionPath string `json:"functionPath" yaml:"functionPath"`
	RequirementPath string `json:"requirementPath" yaml:"requirementPath"`	
}