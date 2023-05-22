package scale

import (
	"minik8s/entity"
	"minik8s/tools/yamlParser"
	"testing"
)

func TestResourceController_GetPodCPUUsage(t *testing.T) {
	type fields struct {
		PodsName       map[string]*entity.Pod
		metricsManager MetricsManager
	}
	type args struct {
		pod *entity.Pod
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "Test Get Pod CPUUsage",
			args:    args{pod: &entity.Pod{}},
			wantErr: false,
		},
	}

	yamlParser.ParseYaml(tests[0].args.pod, "/home/zhaoxi/go/src/minik8s/test/pod2.yaml")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := &ResourceController{
				PodsName:       tt.fields.PodsName,
				metricsManager: tt.fields.metricsManager,
			}
			_, err := rc.GetPodCPUUsage(tt.args.pod)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPodCPUUsage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			//if got != tt.want {
			//	t.Errorf("GetPodCPUUsage() got = %v, want %v", got, tt.want)
			//}
		})
	}
}
