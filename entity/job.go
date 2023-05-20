package entity

type Job struct {
	Kind string `json:"kind" yaml:"kind"`
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	CudaPath string `json:"cudaPath" yaml:"cudaPath"`
	SlurmPath string `json:"slurmPath" yaml:"slurmPath"`
}

type JobFile struct {
	CudaPath string `json:"cudaPath" yaml:"cudaPath"`
	SlurmPath string `json:"slurmPath" yaml:"slurmPath"`
	CudaData []byte `json:"cudaData" yaml:"cudaData"`
    SlurmData []byte `json:"slurmData" yaml:"slurmData"`
}