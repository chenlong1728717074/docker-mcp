package tool

import (
	"context"
	"docker-mcp/api"
	"docker-mcp/resp"
	"encoding/json"
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
		mcp.WithDescription("get Docker logs  , command:docker logs "),
		mcp.WithString("id",
			mcp.Description("docker container id")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := request.Params.Arguments["id"].(string)

		logs, err := cli.ContainerLogs(ctx, id, container.LogsOptions{
			ShowStdout: true,
			ShowStderr: true,
			Tail:       "200",
		})
		if err != nil {
			return nil, err
		}

		result, _ := json.Marshal(map[string]interface{}{
			"status": "success",
			"data":   logs,
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
		mcp.WithDescription("get Docker details  , command:docker inspect "),
		mcp.WithString("id",
			mcp.Description("docker container id")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := request.Params.Arguments["id"].(string)

		inspect, err := cli.ContainerInspect(ctx, id)
		if err != nil {
			return nil, err
		}

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
		mcp.WithDescription("restart Docker container , command:docker restart "),
		mcp.WithString("id",
			mcp.Description("docker container id")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := request.Params.Arguments["id"].(string)
		time := 5
		if err := cli.ContainerRestart(ctx, id, container.StopOptions{Timeout: &time}); err != nil {
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

func RegisterContainerStopTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_container_stop",
		mcp.WithDescription("stop Docker container, command:docker stop "),
		mcp.WithString("id",
			mcp.Description("docker container id")),
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
		mcp.WithDescription("start Docker container , command:docker start "),
		mcp.WithString("id",
			mcp.Description("docker container id")),
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
		mcp.WithDescription("remove Docker container,automatically stop before deletion  , command:docker remove "),
		mcp.WithString("id",
			mcp.Description("docker container id")),
		mcp.WithBoolean("removeVolumes",
			mcp.Description("want to delete the mounted volume?")),
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
		mcp.WithDescription("Run existing/non-existent images, similar to Docker run。"+
			"Example: Docker RMI Hello World。"+
			"This tool will first pull the image, then create and start the container。"),
		mcp.WithString("image", mcp.Description("image name "+
			"Example: redis  or docker.io/library/redis")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		images := request.Params.Arguments["image"].(string)

		//拉取镜像
		pullMsg, err := api.PullImage(ctx, cli, images)
		if err != nil {
			return nil, err
		}
		create, err := api.ContainerCreate(ctx, cli, images)
		if err != nil {
			return nil, err
		}
		if err := api.ContainerStart(ctx, cli, create.ID); err != nil {
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
		mcp.WithDescription("get container list,equivalent to a command:docker ps -a"),
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
