package tool

import (
	"context"
	"docker-mcp/cmd/logs"
	"encoding/json"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// RegisterVolumeTool volume tool
func RegisterVolumeTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	RegisterVolumeListTool(ctx, srv, cli)
}

func RegisterVolumeListTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_volume_list",
		mcp.WithDescription("List Docker volumes - equivalent to 'docker volume ls'"+
			" - Shows all volumes on the system"),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logs.Info("mcp_docker_volume_list called")
		listResp, err := cli.VolumeList(ctx, volume.ListOptions{})
		if err != nil {
			return nil, err
		}
		result, _ := json.Marshal(listResp)
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
