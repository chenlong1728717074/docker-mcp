package tool

import (
	"context"
	"docker-mcp/cmd/logs"
	"encoding/json"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
	"github.com/docker/docker/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"strings"
)

// RegisterVolumeTool volume tool
func RegisterVolumeTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	logs.Info("RegisterVolumeTool called")
	RegisterVolumeListTool(ctx, srv, cli)
	RegisterVolumeCreateTool(ctx, srv, cli)
	RegisterVolumeRemoveTool(ctx, srv, cli)
	RegisterVolumeInspectTool(ctx, srv, cli)
	RegisterVolumePruneTool(ctx, srv, cli)
}

func RegisterVolumeListTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_volume_list",
		mcp.WithDescription("List Docker volumes - equivalent to 'docker volume ls' - Shows all volumes on the system"),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logs.Info("mcp_docker_volume_list called")
		listResp, err := cli.VolumeList(ctx, volume.ListOptions{})
		if err != nil {
			logs.Error("VolumeList failed: %s", err.Error())
			return nil, err
		}
		logs.Info("VolumeList success, found %d volumes", len(listResp.Volumes))
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

func RegisterVolumeCreateTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_volume_create",
		mcp.WithDescription("Create a Docker volume - equivalent to 'docker volume create' - Creates a new volume for data persistence"),
		mcp.WithString("name",
			mcp.Description("Volume name (optional, Docker will generate one if not provided)")),
		mcp.WithString("driver",
			mcp.DefaultString("local"),
			mcp.Description("Volume driver (default: local)")),
		mcp.WithString("labels",
			mcp.Description("Labels in key=value format, separated by commas (e.g., env=prod,app=web)")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := ""
		if nameVal, ok := request.Params.Arguments["name"]; ok {
			name = nameVal.(string)
		}
		driver := "local"
		if driverVal, ok := request.Params.Arguments["driver"]; ok {
			driver = driverVal.(string)
		}

		logs.InfoWithFields("mcp_docker_volume_create called", map[string]interface{}{"name": name, "driver": driver})

		// 处理标签
		labels := make(map[string]string)
		if labelsVal, ok := request.Params.Arguments["labels"]; ok && labelsVal.(string) != "" {
			labelStr := labelsVal.(string)
			// 解析 key=value,key2=value2 格式
			for _, label := range strings.Split(labelStr, ",") {
				parts := strings.SplitN(strings.TrimSpace(label), "=", 2)
				if len(parts) == 2 {
					labels[parts[0]] = parts[1]
				}
			}
		}

		createResp, err := cli.VolumeCreate(ctx, volume.CreateOptions{
			Name:   name,
			Driver: driver,
			Labels: labels,
		})
		if err != nil {
			logs.ErrorWithFields("VolumeCreate failed", map[string]interface{}{"name": name, "error": err})
			return nil, err
		}
		logs.InfoWithFields("VolumeCreate success", map[string]interface{}{"name": createResp.Name})
		result, _ := json.Marshal(createResp)
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

func RegisterVolumeRemoveTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_volume_remove",
		mcp.WithDescription("Remove a Docker volume - equivalent to 'docker volume rm' - Deletes a volume (must not be in use)"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Volume name to remove")),
		mcp.WithBoolean("force",
			mcp.DefaultBool(false),
			mcp.Description("Force removal of the volume")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.Params.Arguments["name"].(string)
		force := false
		if forceVal, ok := request.Params.Arguments["force"]; ok {
			force = forceVal.(bool)
		}

		logs.InfoWithFields("mcp_docker_volume_remove called", map[string]interface{}{"name": name, "force": force})

		err := cli.VolumeRemove(ctx, name, force)
		if err != nil {
			logs.ErrorWithFields("VolumeRemove failed", map[string]interface{}{"name": name, "error": err})
			return nil, err
		}
		logs.InfoWithFields("VolumeRemove success", map[string]interface{}{"name": name})
		result, _ := json.Marshal(map[string]interface{}{
			"status":  "success",
			"message": "Volume removed successfully",
			"volume":  name,
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

func RegisterVolumeInspectTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_volume_inspect",
		mcp.WithDescription("Inspect a Docker volume - equivalent to 'docker volume inspect' - Shows detailed volume information"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Volume name to inspect")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.Params.Arguments["name"].(string)
		logs.InfoWithFields("mcp_docker_volume_inspect called", map[string]interface{}{"name": name})

		inspectResp, err := cli.VolumeInspect(ctx, name)
		if err != nil {
			logs.ErrorWithFields("VolumeInspect failed", map[string]interface{}{"name": name, "error": err})
			return nil, err
		}
		logs.InfoWithFields("VolumeInspect success", map[string]interface{}{"name": name})
		result, _ := json.Marshal(inspectResp)
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

func RegisterVolumePruneTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_volume_prune",
		mcp.WithDescription("Remove unused Docker volumes - equivalent to 'docker volume prune' - Cleans up volumes not used by any container"),
		mcp.WithBoolean("force",
			mcp.DefaultBool(false),
			mcp.Description("Do not prompt for confirmation")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		force := false
		if forceVal, ok := request.Params.Arguments["force"]; ok {
			force = forceVal.(bool)
		}

		logs.InfoWithFields("mcp_docker_volume_prune called", map[string]interface{}{"force": force})

		pruneResp, err := cli.VolumesPrune(ctx, filters.Args{})
		if err != nil {
			logs.ErrorWithFields("VolumesPrune failed", map[string]interface{}{"error": err})
			return nil, err
		}
		logs.InfoWithFields("VolumesPrune success", map[string]interface{}{
			"volumes_deleted": len(pruneResp.VolumesDeleted),
			"space_reclaimed": pruneResp.SpaceReclaimed,
		})
		result, _ := json.Marshal(pruneResp)
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
