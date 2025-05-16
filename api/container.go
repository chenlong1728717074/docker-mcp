package api

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

func ContainerCreate(ctx context.Context, cli *client.Client, name string) (container.CreateResponse, error) {
	return cli.ContainerCreate(ctx,
		&container.Config{
			Image: name,
		},
		&container.HostConfig{}, nil, nil, "")

}
func ContainerStart(ctx context.Context, cli *client.Client, containerID string) error {
	return cli.ContainerStart(ctx, containerID, container.StartOptions{})
}
