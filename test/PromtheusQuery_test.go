package main

import "testing"

func TestPromtheusQuery(t *testing.T) {
	type args struct {
		query string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "查询CPU sum",
			args: args{query: "sum(container_cpu_usage_seconds_total{name=\"autoscale\"})"},
		},
		{
			name: "查询Memory sum",
			args: args{query: "container_memory_usage_bytes{name=\"autoscale\"}"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			PromtheusQuery(tt.args.query)
		})
	}
}
