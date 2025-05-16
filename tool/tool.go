package tool

import (
	"context"
	"github.com/docker/docker/client"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	RegisterSystemTool(ctx, srv, cli)
	RegisterContainerTool(ctx, srv, cli)
	RegisterImageTool(ctx, srv, cli)
}
