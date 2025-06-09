package tool

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterNetworkTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	RegisterNetworkListTool(ctx, srv, cli)
}

func RegisterNetworkListTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {

}
