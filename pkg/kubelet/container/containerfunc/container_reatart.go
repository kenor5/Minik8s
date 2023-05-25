package containerfunc

func ReStartContainer(containerName string) {

	StopContainer(containerName)
	StartContainer(containerName)
}
