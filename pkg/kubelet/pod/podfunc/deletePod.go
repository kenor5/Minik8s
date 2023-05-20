package podfunc

import (
	docker "minik8s/pkg/kubelet/container/containerfunc"
)

func DeletePod(containerIds []string) error {
	for _, containerId := range containerIds {
		docker.StopContainer(containerId)

		_, err := docker.RemoveContainer(containerId)
		if err != nil {
			return err
		}
	}
	return nil
}
