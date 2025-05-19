package tool

import (
	"context"
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
		mcp.WithDescription("Login to Docker Registry,Equivalent to a command: docker login "),
		//mcp.WithString("args",
		//	mcp.Description("Docker registry username")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		//args := request.Params.Arguments["args"].(string)

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
