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
		mcp.WithDescription("Test if Docker daemon is online and obtain minimal version information"),
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
		mcp.WithDescription("Get Docker global system information (number of containers, images, drivers, storage, etc.) monitoring, dashboard"),
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
		mcp.WithDescription("Obtain Docker version related information (version number, API version) compatibility assessment"),
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
		mcp.WithDescription("Equivalent to command: Docker system df,"+
			"You can use it to see which containers are taking up the most space. "+
			"The disk consumption of containers that have been stopped but not deleted. Decide whether to clear the space."),
		mcp.WithString("options",
			mcp.Description("Optional parameters:container(ContainerObject represents a container DiskUsageObject.),"+
				"image(ImageObject represents an image DiskUsageObject),"+
				"volume(VolumeObject represents a volume DiskUsageObject)"+
				",build-cache(BuildCacheObject represents a build-cache DiskUsageObject.)。"+
				"attention: parameters are used for ,"+
				"example:image,volume,container "),
		),
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
