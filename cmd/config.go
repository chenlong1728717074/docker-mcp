package cmd

import (
	"errors"
	"flag"
	"os"
)

type Config struct {
	Path     string
	CertPath string
}

// 从命令行参数获取数据库配置
func GetConfigFromArgs() (*Config, error) {
	var config Config
	// 定义命令行参数 - 优先使用环境变量作为默认值
	//"tcp://101.126.149.147:2375"
	flag.StringVar(&config.Path, "path", os.Getenv("DOCKER_PATH"), "docker addr")
	flag.StringVar(&config.CertPath, "cert", os.Getenv("DOCKER_CERT"), "docker addr")

	// 解析命令行参数
	flag.Parse()
	if config.Path == "" {
		return nil, errors.New("Failed to obtain initialization configuration")
	}
	return &config, nil
}
