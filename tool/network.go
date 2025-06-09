package tool

import (
	"context"
	"docker-mcp/cmd/logs"
	"encoding/json"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"strings"
)

func RegisterNetworkTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	logs.Info("RegisterNetworkTool called")
	RegisterNetworkListTool(ctx, srv, cli)
	RegisterNetworkCreateTool(ctx, srv, cli)
	RegisterNetworkRemoveTool(ctx, srv, cli)
	RegisterNetworkInspectTool(ctx, srv, cli)
	RegisterNetworkConnectTool(ctx, srv, cli)
	RegisterNetworkDisconnectTool(ctx, srv, cli)
	RegisterNetworkPruneTool(ctx, srv, cli)
}

func RegisterNetworkListTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_network_list",
		mcp.WithDescription("List Docker networks - equivalent to 'docker network ls' - Shows all networks on the system"),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		logs.Info("mcp_docker_network_list called")
		networks, err := cli.NetworkList(ctx, network.ListOptions{})
		if err != nil {
			logs.Error("NetworkList failed: %s", err.Error())
			return nil, err
		}
		logs.Info("NetworkList success, found %d networks", len(networks))
		result, _ := json.Marshal(networks)
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

func RegisterNetworkCreateTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_network_create",
		mcp.WithDescription("Create a Docker network - equivalent to 'docker network create' - Creates a new network for container communication"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Network name")),
		mcp.WithString("driver",
			mcp.DefaultString("bridge"),
			mcp.Description("Network driver (bridge, overlay, host, none, macvlan)")),
		mcp.WithString("subnet",
			mcp.Description("Subnet in CIDR format (e.g., 172.20.0.0/16)")),
		mcp.WithString("gateway",
			mcp.Description("Gateway IP address")),
		mcp.WithString("labels",
			mcp.Description("Labels in key=value format, separated by commas")),
		mcp.WithBoolean("internal",
			mcp.DefaultBool(false),
			mcp.Description("Create an internal network (no external connectivity)")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.Params.Arguments["name"].(string)
		driver := "bridge"
		if driverVal, ok := request.Params.Arguments["driver"]; ok {
			driver = driverVal.(string)
		}
		internal := false
		if internalVal, ok := request.Params.Arguments["internal"]; ok {
			internal = internalVal.(bool)
		}

		logs.InfoWithFields("mcp_docker_network_create called", map[string]interface{}{
			"name": name, "driver": driver, "internal": internal,
		})

		// 处理标签
		labels := make(map[string]string)
		if labelsVal, ok := request.Params.Arguments["labels"]; ok && labelsVal.(string) != "" {
			labelStr := labelsVal.(string)
			for _, label := range strings.Split(labelStr, ",") {
				parts := strings.SplitN(strings.TrimSpace(label), "=", 2)
				if len(parts) == 2 {
					labels[parts[0]] = parts[1]
				}
			}
		}

		// 构建IPAM配置
		ipamConfig := &network.IPAM{}
		if subnetVal, ok := request.Params.Arguments["subnet"]; ok && subnetVal.(string) != "" {
			subnet := subnetVal.(string)
			gateway := ""
			if gatewayVal, ok := request.Params.Arguments["gateway"]; ok {
				gateway = gatewayVal.(string)
			}

			ipamConfig.Config = []network.IPAMConfig{
				{
					Subnet:  subnet,
					Gateway: gateway,
				},
			}
		}

		createResp, err := cli.NetworkCreate(ctx, name, network.CreateOptions{
			Driver:   driver,
			Internal: internal,
			Labels:   labels,
			IPAM:     ipamConfig,
		})
		if err != nil {
			logs.ErrorWithFields("NetworkCreate failed", map[string]interface{}{"name": name, "error": err})
			return nil, err
		}
		logs.InfoWithFields("NetworkCreate success", map[string]interface{}{"name": name, "id": createResp.ID})
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

func RegisterNetworkRemoveTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_network_remove",
		mcp.WithDescription("Remove a Docker network - equivalent to 'docker network rm' - Deletes a network (must not be in use)"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Network name or ID to remove")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.Params.Arguments["name"].(string)
		logs.InfoWithFields("mcp_docker_network_remove called", map[string]interface{}{"name": name})

		err := cli.NetworkRemove(ctx, name)
		if err != nil {
			logs.ErrorWithFields("NetworkRemove failed", map[string]interface{}{"name": name, "error": err})
			return nil, err
		}
		logs.InfoWithFields("NetworkRemove success", map[string]interface{}{"name": name})
		result, _ := json.Marshal(map[string]interface{}{
			"status":  "success",
			"message": "Network removed successfully",
			"network": name,
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

func RegisterNetworkInspectTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_network_inspect",
		mcp.WithDescription("Inspect a Docker network - equivalent to 'docker network inspect' - Shows detailed network information"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Network name or ID to inspect")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name := request.Params.Arguments["name"].(string)
		logs.InfoWithFields("mcp_docker_network_inspect called", map[string]interface{}{"name": name})

		inspectResp, err := cli.NetworkInspect(ctx, name, network.InspectOptions{})
		if err != nil {
			logs.ErrorWithFields("NetworkInspect failed", map[string]interface{}{"name": name, "error": err})
			return nil, err
		}
		logs.InfoWithFields("NetworkInspect success", map[string]interface{}{"name": name})
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

func RegisterNetworkConnectTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_network_connect",
		mcp.WithDescription("Connect a container to a network - equivalent to 'docker network connect' - Attaches a container to a network"),
		mcp.WithString("network",
			mcp.Required(),
			mcp.Description("Network name or ID")),
		mcp.WithString("container",
			mcp.Required(),
			mcp.Description("Container name or ID")),
		mcp.WithString("ip",
			mcp.Description("Static IP address to assign to the container")),
		mcp.WithString("aliases",
			mcp.Description("Network aliases for the container, separated by commas")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		networkName := request.Params.Arguments["network"].(string)
		containerName := request.Params.Arguments["container"].(string)

		logs.InfoWithFields("mcp_docker_network_connect called", map[string]interface{}{
			"network": networkName, "container": containerName,
		})

		// 构建端点配置
		endpointConfig := &network.EndpointSettings{}
		if ipVal, ok := request.Params.Arguments["ip"]; ok && ipVal.(string) != "" {
			endpointConfig.IPAMConfig = &network.EndpointIPAMConfig{
				IPv4Address: ipVal.(string),
			}
		}
		if aliasesVal, ok := request.Params.Arguments["aliases"]; ok && aliasesVal.(string) != "" {
			aliases := strings.Split(aliasesVal.(string), ",")
			for i, alias := range aliases {
				aliases[i] = strings.TrimSpace(alias)
			}
			endpointConfig.Aliases = aliases
		}

		err := cli.NetworkConnect(ctx, networkName, containerName, endpointConfig)
		if err != nil {
			logs.ErrorWithFields("NetworkConnect failed", map[string]interface{}{
				"network": networkName, "container": containerName, "error": err,
			})
			return nil, err
		}
		logs.InfoWithFields("NetworkConnect success", map[string]interface{}{
			"network": networkName, "container": containerName,
		})
		result, _ := json.Marshal(map[string]interface{}{
			"status":    "success",
			"message":   "Container connected to network successfully",
			"network":   networkName,
			"container": containerName,
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

func RegisterNetworkDisconnectTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_network_disconnect",
		mcp.WithDescription("Disconnect a container from a network - equivalent to 'docker network disconnect' - Detaches a container from a network"),
		mcp.WithString("network",
			mcp.Required(),
			mcp.Description("Network name or ID")),
		mcp.WithString("container",
			mcp.Required(),
			mcp.Description("Container name or ID")),
		mcp.WithBoolean("force",
			mcp.DefaultBool(false),
			mcp.Description("Force disconnect the container")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		networkName := request.Params.Arguments["network"].(string)
		containerName := request.Params.Arguments["container"].(string)
		force := false
		if forceVal, ok := request.Params.Arguments["force"]; ok {
			force = forceVal.(bool)
		}

		logs.InfoWithFields("mcp_docker_network_disconnect called", map[string]interface{}{
			"network": networkName, "container": containerName, "force": force,
		})

		err := cli.NetworkDisconnect(ctx, networkName, containerName, force)
		if err != nil {
			logs.ErrorWithFields("NetworkDisconnect failed", map[string]interface{}{
				"network": networkName, "container": containerName, "error": err,
			})
			return nil, err
		}
		logs.InfoWithFields("NetworkDisconnect success", map[string]interface{}{
			"network": networkName, "container": containerName,
		})
		result, _ := json.Marshal(map[string]interface{}{
			"status":    "success",
			"message":   "Container disconnected from network successfully",
			"network":   networkName,
			"container": containerName,
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

func RegisterNetworkPruneTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_network_prune",
		mcp.WithDescription("Remove unused Docker networks - equivalent to 'docker network prune' - Cleans up networks not used by any container"),
		mcp.WithBoolean("force",
			mcp.DefaultBool(false),
			mcp.Description("Do not prompt for confirmation")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		force := false
		if forceVal, ok := request.Params.Arguments["force"]; ok {
			force = forceVal.(bool)
		}

		logs.InfoWithFields("mcp_docker_network_prune called", map[string]interface{}{"force": force})

		pruneResp, err := cli.NetworksPrune(ctx, filters.Args{})
		if err != nil {
			logs.ErrorWithFields("NetworksPrune failed", map[string]interface{}{"error": err})
			return nil, err
		}
		logs.InfoWithFields("NetworksPrune success", map[string]interface{}{
			"networks_deleted": len(pruneResp.NetworksDeleted),
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
