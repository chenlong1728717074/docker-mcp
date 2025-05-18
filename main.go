package main

import (
	"context"
	"docker-mcp/cmd"
	"docker-mcp/tool"
	"github.com/docker/docker/client"
	"github.com/mark3labs/mcp-go/server"
	"log"
	"path/filepath"
)

func main() {
	//创建mcp server
	srv := server.NewMCPServer("docker-mcp-support",
		"1.0.0")

	// 获取配置
	ctx := context.Background()
	cfg, err := cmd.GetConfigFromArgs()
	if err != nil {
		log.Fatalf("配置加载失败: %v", err)
	}

	cli, err := initDocker(ctx, cfg)
	if err != nil {
		log.Fatalf("Docker连接失败: %v", err)
	}
	defer cli.Close()

	tool.RegisterTool(ctx, srv, cli)

	//启动
	if err := server.ServeStdio(srv); err != nil {
		log.Fatalf("服务器错误: %v", err)
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

	log.Printf("已连接到Docker API版本: %v", ping.APIVersion)
	return cli, nil
}
