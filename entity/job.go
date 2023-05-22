package entity

type Job struct {
	Kind string `json:"kind" yaml:"kind"`
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	JobStatus JobStatus `josn:"jobstatus,omitempty" yaml:"jobstatus,omitempty"`
	CudaPath string `json:"cudaPath" yaml:"cudaPath"`
	SlurmPath string `json:"slurmPath" yaml:"slurmPath"`
}

type JobStatus struct {
	Status string `json:"status,omitempty" yaml:"status,omitempty"`
	JobID string `json:"jobId,omitempty" yaml:"jobId,omitempty"`
    Result string `json:"result,omitempty" yaml:"result,omitempty"`
}

// type JobFile struct {
// 	CudaPath string `json:"cudaPath" yaml:"cudaPath"`
// 	SlurmPath string `json:"slurmPath" yaml:"slurmPath"`
// 	CudaData []byte `json:"cudaData" yaml:"cudaData"`
//     SlurmData []byte `json:"slurmData" yaml:"slurmData"`
// }