package ScaleController

import (
	"minik8s/entity"
	"minik8s/pkg/apiserver/scale"
	"minik8s/tools/yamlParser"
	"testing"
)

func TestAutoscalerManager_startAutoscalerMonitor(t *testing.T) {
	type fields struct {
		metricsManager *scale.MetricsManager
		autoscalers    map[string]*entity.HorizontalPodAutoscaler
	}
	type args struct {
		autoscaler *entity.HorizontalPodAutoscaler
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name:   "测试自动扩缩HPA",
			fields: fields{metricsManager: nil},
			args:   args{autoscaler: nil},
		},
	}
	tests[0].fields.metricsManager = scale.NewMetricsManager()
	hpa := &entity.HorizontalPodAutoscaler{}
	_, _ = yamlParser.ParseYaml(hpa, "/home/zhaoxi/go/src/minik8s/test/hpaAutoScaleTest.yaml")
	tests[0].args.autoscaler = hpa
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AM := &AutoscalerManager{
				metricsManager: tt.fields.metricsManager,
				autoscalers:    tt.fields.autoscalers,
			}
			AM.startAutoscalerMonitor(tt.args.autoscaler)
		})
	}
}
