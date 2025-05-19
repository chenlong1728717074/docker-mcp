package main

import (
	"context"
	"docker-mcp/cmd"
	"docker-mcp/cmd/logs"
	"docker-mcp/tool"
	"github.com/docker/docker/client"
	"github.com/mark3labs/mcp-go/server"
	"path/filepath"
)

func main() {
	//创建mcp server
	srv := server.NewMCPServer("docker-mcp-support", "1.0.0")
	logs.Info("Starting Docker MCP service")
	// 获取配置
	ctx := context.Background()

	cfg, err := cmd.GetConfigFromArgs()
	if err != nil {
		logs.Fatal("Docker MCP service configuration failed to load: %v", err)
	}
	logs.Info("Docker MCP initialization configuration %v", cfg)
	cli, err := initDocker(ctx, cfg)
	if err != nil {
		logs.Fatal("Docker connection failed: %v", err)
	}
	defer cli.Close()

	tool.RegisterTool(ctx, srv, cli)

	//启动
	if err := server.ServeStdio(srv); err != nil {
		logs.Fatal("Docker MCP service failed to start: %v", err)
	}
}

func initDocker(ctx context.Context, cfg *cmd.Config) (*client.Client, error) {
	opts := []client.Opt{
		client.WithHost(cfg.Path),
	}
	if cfg.CertPath != "" {
		caFile := filepath.Join(cfg.CertPath, "ca.pem")
		certFile := filepath.Join(cfg.CertPath, "cert.pem")
		keyFile := filepath.Join(cfg.CertPath, "key.pem")
		opts = append(opts, client.WithTLSClientConfig(caFile, certFile, keyFile))
	}

	cli, err := client.NewClientWithOpts(opts...)
	if err != nil {
		return nil, err
	}

	// 检查连接
	ping, err := cli.Ping(ctx)
	if err != nil {
		cli.Close() // 如果Ping失败，确保关闭客户端
		return nil, err
	}

	logs.Info("Connected to Docker SUCCESS, API version: %v", ping.APIVersion)
	return cli, nil
}
