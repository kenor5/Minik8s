package ControllerManager

import (
	"minik8s/entity"
	"minik8s/tools/yamlParser"
	"testing"
)

// 测试ApiServer 对deployment的解析和写入
func TestApplyDeployment(t *testing.T) {
	type args struct {
		deployment *entity.Deployment
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Valid deployment should not return error",
			args: args{
				deployment: &entity.Deployment{},
			},
			wantErr: false,
		},
	}
	yamlParser.ParseYaml(tests[0].args.deployment, "/home/zhaoxi/go/src/minik8s/test/nginx_deployment.yaml")
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err := ApplyDeployment(tt.args.deployment); (err != nil) != tt.wantErr {
				t.Errorf("ApplyDeployment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// 测试ApiServer 对deployment相关内容的删除
func TestDeleteDeployment(t *testing.T) {
	type args struct {
		DeploymentName string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "Valid deployment should not return error",
			args: args{
				DeploymentName: "",
			},
			wantErr: false,
		},
	}
	deployment := &entity.Deployment{}

	yamlParser.ParseYaml(deployment, "/home/zhaoxi/go/src/minik8s/test/nginx_deployment.yaml")
	tests[0].args.DeploymentName = deployment.Metadata.Name
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := DeleteDeployment(tt.args.DeploymentName); (err != nil) != tt.wantErr {
				t.Errorf("DeleteDeployment() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
