package containerfunc

import (
	"fmt"
	"minik8s/tools/log"
	"os/exec"
	"strings"
)

func CheckContainerRunning(containerName string) (bool, error) {
	cmd := exec.Command("docker", "ps", "-q", "--filter", fmt.Sprintf("name=%s", containerName))
	output, err := cmd.Output()
	if err != nil {
		log.PrintW(err)
		return false, err
	}

	containers := strings.TrimSpace(string(output))
	if containers == "" {
		// 没有找到容器，返回false
		return false, nil
	}

	return true, nil
}