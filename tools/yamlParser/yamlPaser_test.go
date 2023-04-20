package yamlParser

import (
	"fmt"
	"minik8s/entity"
	"testing"
)

var yamlPath = "../../test/pod1.yaml"

func TestParser(t *testing.T) {

	pod := &entity.Pod{}
	b, _ := ParseYaml(pod, yamlPath)

	fmt.Println(pod)

	if !b {
		fmt.Println("test error")
	}
}
