package scale

import (
	"minik8s/entity"
	"testing"
)

func TestGeneratePrometheusTargets(t *testing.T) {
	type args struct {
		nodes []*entity.Node
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Test modify config file",
			args: args{nodes: []*entity.Node{
				{
					Ip: "192.168.1.0",
				},
				{
					Ip: "192.168.2.0",
				},
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := GeneratePrometheusTargets(tt.args.nodes); (err != nil) != tt.wantErr {
				t.Errorf("GeneratePrometheusTargets() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
