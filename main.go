package main

import (
	"context"
	"docker-mcp/cmd"
	"docker-mcp/tool"
	"github.com/docker/docker/client"
	"github.com/mark3labs/mcp-go/server"
	"log"
)

func main() {
	//创建mcp server
	srv := server.NewMCPServer("docker-mcp-support", "1.0.0")

	// 获取配置
	var err error
	var cfg *cmd.Config
	ctx := context.Background()
	cfg, err = cmd.GetConfigFromArgs()
	cli, err := initDocker(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer cli.Close()
	tool.RegisterTool(ctx, srv, cli)
	//启动
	if err := server.ServeStdio(srv); err != nil {
		log.Fatalf("服务器错误: %v", err)
	}
}

func initDocker(ctx context.Context, cfg *cmd.Config) (*client.Client, error) {
	cli, err := client.NewClientWithOpts(client.WithHost(cfg.Path),
		client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	// 检查连接
	ping, err := cli.Ping(ctx)
	if err != nil {
		return nil, err
	}
	log.Printf("<UNK>: %v", ping)
	return cli, nil
}
