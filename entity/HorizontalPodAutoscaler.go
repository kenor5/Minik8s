package entity

import "time"

type HorizontalPodAutoscaler struct {
	Kind     string     `json:"kind" yaml:"kind"`
	Metadata ObjectMeta `json:"metadata" yaml:"metadata"`
	Spec     HpaSpec    `json:"spec" yaml:"spec"`
	Status   HpaStatus  `json:"status" yaml:"status"`
}

type HpaSpec struct {
	ScaleTargetRef ScaleTargetRef `json:"scaleTargetRef" yaml:"scaleTargetRef"`
	ScaleInterval  int32          `json:"scaleInterval,omitempty" yaml:"scaleInterval,omitempty"`
	MinReplicas    int32          `json:"minReplicas" yaml:"minReplicas"`
	MaxReplicas    int32          `json:"maxReplicas" yaml:"maxReplicas"`
	Metrics        []MetricSpec   `json:"metrics" yaml:"metrics"`
}
type ScaleTargetRef struct {
	Kind string `json:"kind" yaml:"kind"`
	Name string `json:"name" yaml:"name"`
}

type MetricSpec struct {
	Type     string             `json:"type" yaml:"type"`
	Resource ResourceMetricSpec `json:"resource,omitempty" yaml:"resource,omitempty"`
}

type ResourceMetricSpec struct {
	Name   string       `json:"name" yaml:"name"`
	Target MetricTarget `json:"target" yaml:"target"`
}

type MetricTarget struct {
	Type               string `json:"type" yaml:"type"`
	AverageUtilization string `json:"averageUtilization,omitempty" yaml:"averageUtilization,omitempty"`
	AverageValue       string `json:"averageValue,omitempty" yaml:"averageValue,omitempty"`
}
type HpaStatus struct {
	ObservedGeneration int64          `json:"observedGeneration" yaml:"observedGeneration"`
	LastScaleTime      time.Time      `json:"lastScaleTime" yaml:"lastScaleTime"`
	CurrentReplicas    int32          `json:"currentReplicas" yaml:"currentReplicas"`
	DesiredReplicas    int32          `json:"desiredReplicas" yaml:"desiredReplicas"`
	CurrentMetrics     []MetricStatus `json:"currentMetrics,omitempty" yaml:"currentMetrics,omitempty"`
}
type MetricStatus struct {
	Type           string               `json:"type" yaml:"type"`
	ResourceStatus ResourceMetricStatus `json:"resource,omitempty" yaml:"resource,omitempty"`
}

type ResourceMetricStatus struct {
	Name    string            `json:"name" yaml:"name"`
	Current MetricValueStatus `json:"current,omitempty" yaml:"current,omitempty"`
	//DescribedValue DescribedMetricStatus `json:"describedValue,omitempty" yaml:"describedValue,omitempty"`
}

type MetricValueStatus struct {
	AverageValue       string `json:"averageValue,omitempty" yaml:"averageValue,omitempty"`
	AverageUtilization string `json:"averageUtilization,omitempty" yaml:"averageUtilization,omitempty"`
}

//type DescribedMetricStatus struct {
//	CurrentAverageValue       string            `json:"currentAverageValue,omitempty" yaml:"currentAverageValue,omitempty"`
//	CurrentAverageUtilization int32             `json:"currentAverageUtilization,omitempty" yaml:"currentAverageUtilization,omitempty"`
//	currentValue              MetricValueStatus `json:"currentValue,omitempty" yaml:"currentValue,omitempty"`
//}
