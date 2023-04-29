package utils

import (
	"fmt"
	"testing"
)

func TestGetField(t *testing.T) {
	// 文件名不加.yaml
	res, err := GetField("../../../test", "pod1", "kind")
	if err != nil {
		t.Error("get field err")
	} else {
		fmt.Print(res)
	}

}
