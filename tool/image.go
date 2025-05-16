package tool

import (
	"context"
	"docker-mcp/api"
	"docker-mcp/resp"
	"encoding/json"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"strings"
)

func RegisterImageTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	RegisterImageListTool(ctx, srv, cli)
	RegisterImagePullTool(ctx, srv, cli)
	RegisterImageRemoveTool(ctx, srv, cli)
	RegisterImageRemoveBatchTool(ctx, srv, cli)
	RegisterImageDetailsTool(ctx, srv, cli)
}

func RegisterImageRemoveBatchTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_image_remove_batch",
		mcp.WithDescription("Batch delete Docker images, the tool will execute the command: Docker RMI in batches "),
		mcp.WithString("ids",
			mcp.Description("docker name or ID of the image。"+
				"param use, splice together。"+
				"example: redis:v1.0.0,hello-world:latest")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ids := request.Params.Arguments["ids"].(string)
		split := strings.Split(ids, ",")
		responses := make([]image.DeleteResponse, 0)
		for _, val := range split {
			res, err := cli.ImageRemove(ctx, val, image.RemoveOptions{
				Force: true,
			})
			if err != nil {
				return nil, err
			}
			responses = append(responses, res...)
		}

		result, _ := json.Marshal(responses)
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

func RegisterImageRemoveTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_image_remove",
		mcp.WithDescription("remove Docker image , command:docker rmi "),
		mcp.WithString("id",
			mcp.Description("docker image id")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := request.Params.Arguments["id"].(string)
		res, err := cli.ImageRemove(ctx, id, image.RemoveOptions{
			Force: true,
		})
		if err != nil {
			return nil, err
		}
		result, _ := json.Marshal(res)
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

func RegisterImagePullTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_image_pull",
		mcp.WithDescription("pull Docker image , command:docker pull "),
		mcp.WithString("image",
			mcp.Description("docker image")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.Params.Arguments["image"].(string)
		pullImage, err := api.PullImage(ctx, cli, name)
		if err != nil {
			return nil, err
		}
		result, _ := json.Marshal(pullImage)
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

func RegisterImageListTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_image_list",
		mcp.WithDescription("get image list,equivalent to a command:docker image ls"),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		list, err := cli.ImageList(ctx, image.ListOptions{
			All: true,
		})
		if err != nil {
			return nil, err
		}
		images := make([]resp.Image, 0)
		for _, summary := range list {
			tag := ""
			if len(summary.RepoTags) > 0 {
				parts := strings.SplitN(summary.RepoTags[0], ":", 2)
				if len(parts) == 2 {
					tag = parts[1]
				}
			}
			images = append(images, resp.Image{
				Repository: summary.RepoTags,
				Tag:        tag,
				ImageID:    summary.ID,
				Created:    summary.Created,
				Size:       summary.Size,
			})
		}
		result, _ := json.Marshal(images)
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

func RegisterImageDetailsTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_image_details",
		mcp.WithDescription("get Docker details , command:docker image inspect <image-name-or-id>"),
		mcp.WithString("id",
			mcp.Description("docker (image id or image id)")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := request.Params.Arguments["id"].(string)
		res, err := cli.ImageInspect(ctx, id)
		if err != nil {
			return nil, err
		}
		result, _ := json.Marshal(res)
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
