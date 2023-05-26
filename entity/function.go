package entity

type Function struct {
	Kind string `json:"kind" yaml:"kind"`
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	FunctionStatus FunctionStatus `josn:"functionStatus,omitempty" yaml:"functionStatus,omitempty"`
	FunctionPath string `json:"functionPath" yaml:"functionPath"`
	RequirementPath string `json:"requirementPath" yaml:"requirementPath"`	
}

type FunctionStatus struct {
	Status string `json:"status,omitempty" yaml:"status,omitempty"`
	PodTemplate Pod `json:"podTemplate,omitempty" yaml:"podTemplate,omitempty"`
	FunctionPods []FunctionPod  `json:"functionPod,omitempty" yaml:"functionPod,omitempty"`
	AccessTimes int32 `json:"accessTimes,omitempty" yaml:"accessTimes,omitempty"`
}

type FunctionPod struct {
	PodName string `json:"podName,omitempty" yaml:"podName,omitempty"`
	PodIp string `json:"podIp,omitempty" yaml:"podIp,omitempty"`
}
