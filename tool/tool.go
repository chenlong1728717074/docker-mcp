package tool

import (
	"context"
	"docker-mcp/cmd/logs"
	"github.com/docker/docker/client"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	logs.Info("RegisterTool called")
	RegisterSystemTool(ctx, srv, cli)
	RegisterContainerTool(ctx, srv, cli)
	RegisterImageTool(ctx, srv, cli)
}
