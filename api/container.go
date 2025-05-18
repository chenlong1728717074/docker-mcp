package api

import (
	"context"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"strings"
)

func ContainerCreate(ctx context.Context, cli *client.Client, name, env, containerName, ports, volumes string) (container.CreateResponse, error) {
	var envs []string
	if env != "" {
		envs = strings.Split(env, ",")
	}
	// 处理端口映射
	exposedPorts, portBindings := buildPort(ports)
	// 处理卷挂载
	containerVolumes, binds := buildVolumes(volumes)
	return cli.ContainerCreate(ctx,
		&container.Config{
			Image:        name,
			Env:          envs,
			ExposedPorts: exposedPorts,
			Volumes:      containerVolumes,
		},
		&container.HostConfig{
			PortBindings: portBindings,
			Binds:        binds,
		}, nil, nil, containerName)

}

// 修改函数返回两个值：暴露的端口和端口映射
func buildPort(ports string) (nat.PortSet, nat.PortMap) {
	exposedPorts := nat.PortSet{}
	portBindings := nat.PortMap{}

	if ports == "" {
		return exposedPorts, portBindings
	}

	// 按逗号分割端口映射列表
	portList := strings.Split(ports, ",")

	for _, portMapping := range portList {
		// 跳过空条目
		if portMapping == "" {
			continue
		}

		// 解析映射关系
		parts := strings.Split(strings.TrimSpace(portMapping), ":")

		// 处理不同的格式
		if len(parts) == 1 {
			// 格式: containerPort - 暴露容器端口但不映射
			containerPort := nat.Port(parts[0])
			// 检查是否已包含协议，如果没有则默认为TCP
			if !strings.Contains(string(containerPort), "/") {
				containerPort = nat.Port(string(containerPort) + "/tcp")
			}
			// 添加到暴露端口
			exposedPorts[containerPort] = struct{}{}
			// 空的绑定列表表示端口已暴露但未映射
			portBindings[containerPort] = []nat.PortBinding{}

		} else if len(parts) == 2 {
			// 格式: hostPort:containerPort - 将主机端口映射到容器端口
			hostPort := parts[0]
			containerPort := nat.Port(parts[1])
			// 检查是否已包含协议，如果没有则默认为TCP
			if !strings.Contains(string(containerPort), "/") {
				containerPort = nat.Port(string(containerPort) + "/tcp")
			}

			// 添加到暴露端口
			exposedPorts[containerPort] = struct{}{}

			// 添加端口绑定
			portBindings[containerPort] = []nat.PortBinding{
				{
					HostIP:   "0.0.0.0", // 默认绑定到所有IPv4地址
					HostPort: hostPort,
				},
			}

		} else if len(parts) == 3 {
			// 格式: hostIP:hostPort:containerPort - 指定主机IP、主机端口和容器端口
			hostIP := parts[0]
			hostPort := parts[1]
			containerPort := nat.Port(parts[2])
			// 检查是否已包含协议，如果没有则默认为TCP
			if !strings.Contains(string(containerPort), "/") {
				containerPort = nat.Port(string(containerPort) + "/tcp")
			}

			// 添加到暴露端口
			exposedPorts[containerPort] = struct{}{}

			// 添加端口绑定
			portBindings[containerPort] = []nat.PortBinding{
				{
					HostIP:   hostIP,
					HostPort: hostPort,
				},
			}
		}
	}

	return exposedPorts, portBindings
}

func buildVolumes(volumes string) (map[string]struct{}, []string) {
	// containerVolumes用于container.Config.Volumes（匿名卷）
	containerVolumes := make(map[string]struct{})

	// binds用于container.HostConfig.Binds（绑定挂载）
	var binds []string

	if volumes == "" {
		return containerVolumes, binds
	}

	// 按逗号分割卷挂载列表
	volumeList := strings.Split(volumes, ",")

	for _, vol := range volumeList {
		// 跳过空条目
		if vol == "" {
			continue
		}

		vol = strings.TrimSpace(vol)

		// 解析挂载规范
		parts := strings.Split(vol, ":")

		// 处理不同的格式
		if len(parts) == 1 {
			// 格式: /container/path - 匿名卷
			containerPath := parts[0]
			containerVolumes[containerPath] = struct{}{}

		} else if len(parts) >= 2 {
			// 格式: /host/path:/container/path[:mode] - 绑定挂载
			// 或者: volume_name:/container/path[:mode] - 命名卷
			// 直接添加到binds列表，让Docker处理详细解析
			// Docker会正确处理两部分(path:path)或三部分(path:path:mode)的格式
			binds = append(binds, vol)
		}
	}

	return containerVolumes, binds
}

func ContainerStart(ctx context.Context, cli *client.Client, containerID string) error {
	return cli.ContainerStart(ctx, containerID, container.StartOptions{})
}
