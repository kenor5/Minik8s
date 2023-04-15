package yamlParser

import (
	"minik8s/entity"
	"fmt"
	"testing"
)

var yamlPath = "../../test/pod1.yaml"

func TestParser(t *testing.T) {

	pod := &entity.Pod{}
	b, _ := parseYaml(pod, yamlPath)

	fmt.Println(pod)

	if !b {
		fmt.Println("test error")
	}
}
