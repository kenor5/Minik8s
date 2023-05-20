package yamlParser

import (
	"fmt"
	"minik8s/entity"
	"testing"
)

var yamlPath = "../../test/pod2.yaml"
var yamlPath2 = "../../test/service_test.yaml"
var yamlPath3 = "../../test/nginx_deployment.yaml"

func TestParser(t *testing.T) {

	pod := &entity.Pod{}
	b, _ := ParseYaml(pod, yamlPath)

	service := &entity.Service{}
	s, _ := ParseYaml(service, yamlPath2)

	fmt.Println(pod)
	fmt.Println(service)

	if !b {
		fmt.Println("test error")
	}

	if !s {
		fmt.Println("parse service error")
	}

	deploy := &entity.Deployment{}
	d, _ := ParseYaml(deploy, yamlPath3)
	if !d {
		fmt.Println("parse deploy error")
	}
	fmt.Println(deploy)
}
