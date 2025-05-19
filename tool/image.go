package tool

import (
	"context"
	"docker-mcp/api"
	"docker-mcp/cmd/logs"
	"docker-mcp/resp"
	"encoding/json"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"strings"
)

func RegisterImageTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	logs.Info("RegisterImageTool called")
	RegisterImageListTool(ctx, srv, cli)
	RegisterImagePullTool(ctx, srv, cli)
	RegisterImageRemoveTool(ctx, srv, cli)
	RegisterImageRemoveBatchTool(ctx, srv, cli)
	RegisterImageDetailsTool(ctx, srv, cli)
}

func RegisterImageRemoveBatchTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_image_remove_batch",
		mcp.WithDescription("Remove multiple Docker images in batch - equivalent to 'docker rmi <image1> <image2>' - Deletes specified images from the system"),
		mcp.WithString("ids",
			mcp.Description("Comma-separated list of image names or IDs to remove, e.g., redis:v1.0.0,hello-world:latest")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ids := request.Params.Arguments["ids"].(string)
		logs.Info("mcp_docker_image_remove_batch called, ids: %s", ids)
		split := strings.Split(ids, ",")
		responses := make([]image.DeleteResponse, 0)
		for _, val := range split {
			logs.Info("Removing image: %s", val)
			res, err := cli.ImageRemove(ctx, val, image.RemoveOptions{
				Force: true,
			})
			if err != nil {
				logs.Error("Remove image failed: %s, error: %s", val, err.Error())
				return nil, err
			}
			logs.Info("Remove image success: %s", val)
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
		mcp.WithDescription("Remove a Docker image - equivalent to 'docker rmi <image>' - Deletes an image from the system"),
		mcp.WithString("id",
			mcp.Description("Image ID or image name with optional tag")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := request.Params.Arguments["id"].(string)
		logs.Info("mcp_docker_image_remove called, id: %s", id)
		res, err := cli.ImageRemove(ctx, id, image.RemoveOptions{
			Force: true,
		})
		if err != nil {
			logs.Error("Remove image failed: %s, error: %s", id, err.Error())
			return nil, err
		}
		logs.Info("Remove image success: %s", id)
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
		mcp.WithDescription("Pull a Docker image - equivalent to 'docker pull <image>' - Downloads an image from a registry"),
		mcp.WithString("image",
			mcp.Description("Image name to pull with optional tag")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.Params.Arguments["image"].(string)
		logs.Info("mcp_docker_image_pull called, image: %s", name)
		pullImage, err := api.PullImage(ctx, cli, name)
		if err != nil {
			logs.Error("Pull image failed: %s, error: %s", name, err.Error())
			return nil, err
		}
		logs.Info("Pull image success: %s", name)
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
		mcp.WithDescription("List all Docker images - equivalent to 'docker image ls' - Shows all images stored locally on the system"),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logs.Info("mcp_docker_image_list called")
		list, err := cli.ImageList(ctx, image.ListOptions{
			All: true,
		})
		if err != nil {
			logs.Error("ImageList failed: %s", err.Error())
			return nil, err
		}
		logs.Info("ImageList success, count: %d", len(list))
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
		mcp.WithDescription("Get detailed information about an image - equivalent to 'docker image inspect <image>' - Shows layers, configuration, and metadata"),
		mcp.WithString("id",
			mcp.Description("Image ID or image name with optional tag")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		id := request.Params.Arguments["id"].(string)
		logs.Info("mcp_docker_image_details called, id: %s", id)
		res, err := cli.ImageInspect(ctx, id)
		if err != nil {
			logs.Error("ImageInspect failed: %s, error: %s", id, err.Error())
			return nil, err
		}
		logs.Info("ImageInspect success: %s", id)
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
