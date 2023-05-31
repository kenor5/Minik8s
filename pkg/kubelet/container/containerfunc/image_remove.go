package containerfunc

import (
	"context"
	"fmt"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"log"
)

func DeleteImage(imageName string) error {
	cli, _ := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	defer cli.Close()

	imageExists, err := ImageExist(imageName)
	if err != nil {
		return fmt.Errorf("failed to check if image exists: %v", err)
	}

	if !imageExists {
		log.Printf("Image '%s' does not exist.\n", imageName)
		return nil
	}

	_, err = cli.ImageRemove(context.Background(), imageName, types.ImageRemoveOptions{
		Force:         true,
		PruneChildren: true,
	})
	if err != nil {
		return fmt.Errorf("failed to remove image: %v", err)
	}

	log.Printf("Image '%s' has been successfully removed.\n", imageName)

	return nil
}