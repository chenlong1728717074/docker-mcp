package resp

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/pkg/jsonmessage"
)

type Container struct {
	ID      string
	Names   []string
	Command string
	Created int64
	Ports   []container.Port
	Image   string
}

type ContainerRun struct {
	Create  ContainerCreate
	PullMsg []jsonmessage.JSONMessage
}

type ContainerCreate struct {
	ID       string   `json:"Id"`
	Warnings []string `json:"Warnings"`
}
