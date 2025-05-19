package tool

import (
	"context"
	"encoding/json"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterAuthTool(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	RegisterRegistry(ctx, srv, cli)
}

func RegisterRegistry(ctx context.Context, srv *server.MCPServer, cli *client.Client) {
	tool := mcp.NewTool("mcp_docker_auth_registry",
		mcp.WithDescription("Login to Docker Registry,Equivalent to a command: docker login "),
		mcp.WithString("username",
			mcp.Required(),
			mcp.Description("Docker registry username")),
		mcp.WithString("password",
			mcp.Required(),
			mcp.Description("Docker registry password")),
		mcp.WithString("serverAddress",
			mcp.DefaultString("https://index.docker.io/v1/"),
			mcp.Description("Docker registry address, default is Docker Hub")),
	)
	srv.AddTool(tool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		username := request.Params.Arguments["username"].(string)
		password := request.Params.Arguments["password"].(string)
		serverAddress := request.Params.Arguments["serverAddress"].(string)

		loginResp, err := cli.RegistryLogin(ctx, registry.AuthConfig{
			Username:      username,
			Password:      password,
			ServerAddress: serverAddress,
		})
		if err != nil {
			return nil, err
		}
		result, _ := json.Marshal(map[string]interface{}{
			"status":       "success",
			"login_status": loginResp.Status,
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
