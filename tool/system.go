package tool

import (
	"context"
	"docker-mcp/resp"
	"encoding/json"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"strings"
)

func RegisterSystemTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	RegisterPingTool(ctx, srv, cli)
	RegisterInfoTool(ctx, srv, cli)
	RegisterServiceVersionTool(ctx, srv, cli)
	RegisterDiskUsageTool(ctx, srv, cli)
}

func RegisterInfoTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_system_info",
		mcp.WithDescription("Test Docker daemon connectivity - equivalent to 'docker info' (simplified) - Verifies if Docker daemon is running and returns basic information"),
	)

	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		ping, err := cli.Ping(ctx)
		if err != nil {
			return nil, err
		}
		system := resp.System{
			APIVersion:       ping.APIVersion,
			OSType:           ping.OSType,
			Experimental:     ping.Experimental,
			BuilderVersion:   ping.APIVersion,
			NodeState:        string(ping.SwarmStatus.NodeState),
			ControlAvailable: ping.SwarmStatus.ControlAvailable,
		}
		result, _ := json.Marshal(system)
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
func RegisterPingTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_system_ping",
		mcp.WithDescription("Get detailed Docker system information - equivalent to 'docker info' - Shows containers, images, drivers, storage, and other system details"),
	)

	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		info, err := cli.Info(ctx)
		if err != nil {
			return nil, err
		}
		result, _ := json.Marshal(info)
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

func RegisterServiceVersionTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_system_server_version",
		mcp.WithDescription("Get Docker version information - equivalent to 'docker version' - Shows version numbers and API version for compatibility assessment"),
	)

	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		svi, err := cli.ServerVersion(ctx)
		if err != nil {
			return nil, err
		}
		result, _ := json.Marshal(svi)
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

func RegisterDiskUsageTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_system_disk_usage",
		mcp.WithDescription("Show Docker disk usage - equivalent to 'docker system df' - Displays space used by containers, images, volumes, and build cache"),
		mcp.WithString("options",
			mcp.Description("Optional comma-separated list of resource types to include: container, image, volume, build-cache (e.g., 'image,volume,container')")),
	)

	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		params := make([]types.DiskUsageObject, 0)
		if val, ok := request.Params.Arguments["options"].(string); ok && val != "" {
			for _, v := range strings.Split(val, ",") {
				params = append(params, types.DiskUsageObject(v))
			}
		}
		svi, err := cli.DiskUsage(ctx, types.DiskUsageOptions{
			Types: params,
		})
		if err != nil {
			return nil, err
		}
		result, _ := json.Marshal(svi)
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
