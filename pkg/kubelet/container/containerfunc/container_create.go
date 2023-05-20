package containerfunc

import (
	"context"
	"fmt"
	"io"
	"minik8s/entity"
	"minik8s/tools/log"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

/**
* create a new container according to the given name and image
* if image dosen't exit, the function will try to pull it first
* return the new container's containerID
 */
func CreateContainer(Container entity.Container) string {

	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	defer cli.Close()
	image := Container.Image
	containerName := Container.Name
	error := EnsureImage(image)
	if error != nil {
		panic(error)
	}

	config := &container.Config{
		Image: image,
		Cmd:   Container.Command,
	}

	body, err := cli.ContainerCreate(context.Background(), config, &container.HostConfig{}, nil, nil, containerName)

	if err != nil {
		panic(err)
	}

	log.Print("ID: %s\n", body.ID)
	return body.ID
}

/**
* helper function
* used to ensure the image exists
 */
func EnsureImage(targetImage string) error {
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	defer cli.Close()

	exist, err := ImageExist(targetImage)
	if err != nil {
		return err
	}

	if exist {
		return nil
	}

	log.PrintE("image %s doesn't exist, automatically pulling\n", targetImage)

	reader, err := cli.ImagePull(context.Background(), targetImage, types.ImagePullOptions{})
	io.Copy(os.Stdout, reader)

	if err == nil {
		return nil
	}

	return fmt.Errorf("failed to ensure image %s", targetImage)
}

/**
* helper function
* used to check whether the image exist, return true or false
 */
func ImageExist(targetImage string) (bool, error) {
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	defer cli.Close()

	images, err := cli.ImageList(context.Background(), types.ImageListOptions{})
	if err != nil {
		return false, err
	}

	for _, image := range images {
		// fmt.Println(image)
		for _, tag := range image.RepoTags {
			// fmt.Printf("tag %s\n", tag)
			if tag == targetImage {
				log.Print("-----have found the image-----\n")
				return true, nil
			}
		}
	}

	return false, nil
}
