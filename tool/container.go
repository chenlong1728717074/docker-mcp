package tool

import (
	"context"
	"docker-mcp/api"
	"docker-mcp/cmd/logs"
	"docker-mcp/resp"
	"encoding/json"
	"errors"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterContainerTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	RegisterContainerListTool(ctx, srv, cli)
	RegisterContainerRunTool(ctx, srv, cli)
	RegisterContainerStartTool(ctx, srv, cli)
	RegisterContainerStopTool(ctx, srv, cli)
	RegisterContainerRestartTool(ctx, srv, cli)
	RegisterContainerRemoveTool(ctx, srv, cli)
	RegisterContainerInspectTool(ctx, srv, cli)
	RegisterContainerLogsTool(ctx, srv, cli)
}

func RegisterContainerLogsTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_container_log",
		mcp.WithDescription("Get container logs - equivalent to 'docker logs <container-id>' - Shows output from the container application"),
		mcp.WithString("id",
			mcp.Description("Container ID or container name")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := request.Params.Arguments["id"].(string)
		logs.InfoWithFields("mcp_docker_container_log called", map[string]interface{}{"id": id})
		containerLogs, err := cli.ContainerLogs(ctx, id, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Tail:       "200",
		})
		if err != nil {
			logs.ErrorWithFields("ContainerLogs failed", map[string]interface{}{"id": id, "error": err})
			return nil, err
		}
		logs.InfoWithFields("ContainerLogs success", map[string]interface{}{"id": id})
		result, _ := json.Marshal(map[string]interface{}{
			"status": "success",
			"data":   containerLogs,
		})
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: string(result),
					Type: "text",
				},
			},
		}, nil
	})
}

func RegisterContainerInspectTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_container_details",
		mcp.WithDescription("Get detailed information about a container - equivalent to 'docker inspect <container-id>' - Shows configuration, volumes, networks, etc."),
		mcp.WithString("id",
			mcp.Description("Container ID or container name")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := request.Params.Arguments["id"].(string)
		logs.InfoWithFields("mcp_docker_container_details called", map[string]interface{}{"id": id})
		inspect, err := cli.ContainerInspect(ctx, id)
		if err != nil {
			logs.ErrorWithFields("ContainerInspect failed", map[string]interface{}{"id": id, "error": err})
			return nil, err
		}
		logs.InfoWithFields("ContainerInspect success", map[string]interface{}{"id": id})
		result, _ := json.Marshal(map[string]interface{}{
			"status": "success",
			"data":   inspect,
		})
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: string(result),
					Type: "text",
				},
			},
		}, nil
	})
}

func RegisterContainerRestartTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_container_restart",
		mcp.WithDescription("Restart a container - equivalent to 'docker restart <container-id>' - Gracefully stops and starts a container"),
		mcp.WithString("id",
			mcp.Description("Container ID or container name")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := request.Params.Arguments["id"].(string)
		timeout := 5
		logs.InfoWithFields("mcp_docker_container_restart called", map[string]interface{}{"id": id, "timeout": timeout})
		if err := cli.ContainerRestart(ctx, id, container.StopOptions{Timeout: &timeout}); err != nil {
			logs.ErrorWithFields("ContainerRestart failed", map[string]interface{}{"id": id, "error": err})
			return nil, err
		}
		logs.InfoWithFields("ContainerRestart success", map[string]interface{}{"id": id})
		result, _ := json.Marshal(map[string]string{
			"status": "success",
		})
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: string(result),
					Type: "text",
				},
			},
		}, nil
	})
}

func RegisterContainerStopTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_container_stop",
		mcp.WithDescription("Stop a running container - equivalent to 'docker stop <container-id>' - Sends SIGTERM signal to the main process"),
		mcp.WithString("id",
			mcp.Description("Container ID or container name")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := request.Params.Arguments["id"].(string)
		time := 5
		if err := cli.ContainerStop(ctx, id, container.StopOptions{Timeout: &time}); err != nil {
			return nil, err
		}
		result, _ := json.Marshal(map[string]string{
			"status": "success",
		})
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: string(result),
					Type: "text",
				},
			},
		}, nil
	})
}

func RegisterContainerStartTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_container_start",
		mcp.WithDescription("Start a stopped container - equivalent to 'docker start <container-id>' - Starts a previously created container"),
		mcp.WithString("id",
			mcp.Description("Container ID or container name")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := request.Params.Arguments["id"].(string)
		if err := cli.ContainerStart(ctx, id, container.StartOptions{}); err != nil {
			return nil, err
		}
		result, _ := json.Marshal(map[string]string{
			"status": "success",
		})
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: string(result),
					Type: "text",
				},
			},
		}, nil
	})
}

func RegisterContainerRemoveTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_container_remove",
		mcp.WithDescription("Remove a container - equivalent to 'docker rm <container-id>' - Automatically stops and removes the specified container"),
		mcp.WithString("id",
			mcp.Description("Container ID or container name")),
		mcp.WithBoolean("removeVolumes",
			mcp.Description("Whether to remove volumes associated with the container")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := request.Params.Arguments["id"].(string)
		//先关闭后删除
		if err := cli.ContainerStop(ctx, id, container.StopOptions{}); err != nil {
			return nil, err
		}
		removeVolumes := request.Params.Arguments["removeVolumes"].(bool)
		if err := cli.ContainerRemove(ctx, id, container.RemoveOptions{
			Force:         true,
			RemoveVolumes: removeVolumes,
		}); err != nil {
			return nil, err
		}
		result, _ := json.Marshal(map[string]string{
			"status": "success",
		})
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: string(result),
					Type: "text",
				},
			},
		}, nil
	})
}

func RegisterContainerRunTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_container_run",
		mcp.WithDescription("Run a Docker image - equivalent to 'docker run <image>' - Pulls the image (if not present locally), then creates and starts a container"),
		mcp.WithString("image",
			mcp.Required(),
			mcp.Description("Image name in format: [registry/][username/]name[:tag], e.g., redis or docker.io/library/redis:latest")),
		mcp.WithString("env",
			mcp.DefaultString(""),
			mcp.Description("Environment variables in format: VAR1=value1,VAR2=value2. For example: MYSQL_ROOT_PASSWORD=password,MYSQL_DATABASE=mydb")),
		mcp.WithString("containerName",
			mcp.DefaultString(""),
			mcp.Description("Assign a name to the container. If not specified, Docker will generate a random name")),
		mcp.WithString("ports",
			mcp.DefaultString(""),
			mcp.Description("Port mappings in format: [hostIP:]hostPort:containerPort[/protocol]. Multiple mappings separated by commas. Examples: 8080:80/tcp,127.0.0.1:5432:5432,3306:3306")),
		mcp.WithString("volumes",
			mcp.DefaultString(""),
			mcp.Description("Volume mappings in format: hostPath:containerPath[:mode]. Multiple volumes separated by commas. Examples: /data:/var/lib/mysql,/config:/etc/mysql/conf.d:ro")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		images, ok := request.Params.Arguments["image"].(string)
		if !ok || images == "" {
			return nil, errors.New("image parameter is required and must be a string")
		}
		env, containerName, ports, volumes := "", "", "", ""
		if val, ok := request.Params.Arguments["env"]; ok {
			env = val.(string)
		}
		if val, ok := request.Params.Arguments["containerName"]; ok {
			containerName = val.(string)
		}
		if val, ok := request.Params.Arguments["ports"]; ok {
			ports = val.(string)
		}
		if val, ok := request.Params.Arguments["volumes"]; ok {
			volumes = val.(string)
		}
		logs.Info("mcp_docker_container_run tool being visited: %s %s %s %s", env, containerName, ports, volumes)
		//拉取镜像
		pullMsg, err := api.PullImage(ctx, cli, images)
		logs.Info("mcp_docker_container_run tool image pull.....")
		if err != nil {
			logs.Error("mcp_docker_container_run tool image pull fail:", err.Error())
			return nil, err
		}
		logs.Info("mcp_docker_container_run tool container create.....")
		create, err := api.ContainerCreate(ctx, cli, images, env, containerName, ports, volumes)
		if err != nil {
			logs.Error("mcp_docker_container_run tool container create fail:", err.Error())
			return nil, err
		}
		logs.Info("mcp_docker_container_run tool container start.....")
		if err := api.ContainerStart(ctx, cli, create.ID); err != nil {
			logs.Error("mcp_docker_container_run tool container start fail:", err.Error())
			return nil, err
		}
		containers := resp.ContainerRun{
			PullMsg: pullMsg,
			Create: resp.ContainerCreate{
				ID:       create.ID,
				Warnings: create.Warnings,
			},
		}

		result, _ := json.Marshal(containers)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: string(result),
					Type: "text",
				},
			},
		}, nil
	})
}

func RegisterContainerListTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_container_list",
		mcp.WithDescription("List all containers - equivalent to 'docker ps -a' - Shows all containers (running and stopped) in the system"),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		list, err := cli.ContainerList(ctx, container.ListOptions{
			All: true,
		})
		if err != nil {
			return nil, err
		}
		containers := make([]resp.Container, 0, len(list))
		for _, ctr := range list {

			containers = append(containers, resp.Container{
				ID:      ctr.ID,
				Names:   ctr.Names,
				Command: ctr.Command,
				Created: ctr.Created,
				Ports:   ctr.Ports,
				Image:   ctr.Image,
			})
		}
		result, _ := json.Marshal(containers)
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: string(result),
					Type: "text",
				},
			},
		}, nil
	})
}
