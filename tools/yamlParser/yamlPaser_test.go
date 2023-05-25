package yamlParser

import (
	"fmt"
	"minik8s/entity"
	"minik8s/tools/log"
	"testing"
)

var yamlPath = "../../test/pod2.yaml"
var yamlPath2 = "../../test/service_test.yaml"
var yamlPath3 = "../../test/nginx_deployment.yaml"
var yamlPath4 = "../../test/hpa_test.yaml"
var yamlPath5 = "../../test/pod4_zx.yaml"

func TestParser(t *testing.T) {

	pod := &entity.Pod{}
	b, _ := ParseYaml(pod, yamlPath5)

	service := &entity.Service{}
	s, _ := ParseYaml(service, yamlPath2)

	hpa := &entity.HorizontalPodAutoscaler{}
	h, _ := ParseYaml(hpa, yamlPath4)

	fmt.Println(pod)
	fmt.Println(service)
	fmt.Println(hpa)

	if !b {
		log.PrintE("test error")
	}

	if !s {
		log.PrintE("parse service error")
	}
	if !h {
		log.PrintE("parse hpa error")
	}

	//deploy := &entity.Deployment{}
	//d, _ := ParseYaml(deploy, yamlPath3)
	//if !d {
	//	fmt.Println("parse deploy error")
	//}
	//fmt.Println(deploy)
}
