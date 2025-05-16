# Docker MCP

[English](README_EN.md) | 中文

Docker MCP 是一个通过模型-命令-插件（Model-Command-Plugin，MCP）框架提供 Docker 管理能力的工具。它允许您通过标准化接口与 Docker 容器和镜像进行交互。

## 功能特点

- **容器管理**
  - 列出容器
  - 运行容器
  - 启动/停止/重启容器
  - 删除容器
  - 查看容器日志和详情

- **镜像管理**
  - 列出镜像
  - 拉取镜像
  - 删除镜像（单个或批量）
  - 查看镜像详情

- **系统管理**
  - 检查 Docker 守护进程状态
  - 获取 Docker 版本信息
  - 查看系统信息
  - 检查磁盘使用情况

## 系统要求

- Go 1.24 或更高版本
- Docker（本地或远程）
- Docker API 访问权限

## 从源码构建

1. 克隆仓库：
   ```bash
   git clone https://github.com/yourusername/docker-mcp.git
   cd docker-mcp
   ```

2. 安装依赖：
   ```bash
   go mod download
   ```

3. 构建可执行文件：
   ```bash
   go build -o docker-mcp.exe
   ```

## 使用方法

您可以直接运行 Docker MCP：

```bash
./docker-mcp.exe
```

### 环境变量

- `DOCKER_PATH`：Docker 守护进程套接字路径或 TCP 端点（例如：`tcp://your-docker-server:2375`）

### 命令行参数

- `--path`：Docker 守护进程套接字路径或 TCP 端点（覆盖环境变量）

### 重要注意事项

为了使用远程 Docker API，您需要在 Docker 主机上启用 API 访问。有以下几种方式：

#### 方式一：修改 Docker 守护进程配置文件

1. 修改 Docker 守护进程配置文件（例如 `/etc/docker/daemon.json`），添加以下内容：
   ```json
   {
     "hosts": ["tcp://0.0.0.0:2375", "unix:///var/run/docker.sock"]
   }
   ```

2. 重启 Docker 服务：
   ```bash
   sudo systemctl restart docker
   ```

#### 方式二：修改 Docker 服务启动参数

1. 对于使用 systemd 的系统，编辑 Docker 服务文件：
   ```bash
   sudo systemctl edit docker.service
   ```

2. 添加以下内容：
   ```ini
   [Service]
   ExecStart=
   ExecStart=/usr/bin/dockerd -H fd:// -H tcp://0.0.0.0:2375
   ```

3. 重载 systemd 配置并重启 Docker：
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl restart docker
   ```

#### 验证 API 访问

确认 Docker API 已开启：
```bash
curl http://localhost:2375/version
```

**安全警告**：开放 2375 端口允许未经身份验证的 Docker API 访问。在生产环境中，建议使用 TLS 证书（2376 端口）或设置网络安全组/防火墙规则限制访问。仅在受信任的网络环境中使用此配置。

## Cursor 集成

Docker MCP 可以与 Cursor IDE 集成，直接在编辑器中提供 Docker 管理功能。

### 配置步骤

1. 打开 Cursor 设置
2. 导航到 MCP 配置部分
3. 在 Cursor 设置中添加以下配置：

```json
{
  "mcpServers": {
    "docker-mcp": {
      "command": "{your-build-path}/docker-mcp.exe",
      "args": [],
      "env": {
        "DOCKER_PATH": "tcp://your-docker-server:2375"
      }
    }
  }
}
```

4. 保存设置并重启 Cursor

### 配置选项

- `command`：docker-mcp 可执行文件的路径
- `args`：附加的命令行参数
- `env`：传递给可执行文件的环境变量
  - `DOCKER_PATH`：Docker 守护进程套接字路径或 TCP 端点

## 可用工具

### 容器工具

- `mcp_docker_container_list`：列出所有容器
- `mcp_docker_container_run`：运行 Docker 镜像
- `mcp_docker_container_start`：启动已停止的容器
- `mcp_docker_container_stop`：停止运行中的容器
- `mcp_docker_container_restart`：重启容器
- `mcp_docker_container_remove`：删除容器
- `mcp_docker_container_details`：获取容器详细信息
- `mcp_docker_container_log`：获取容器日志

### 镜像工具

- `mcp_docker_image_list`：列出所有 Docker 镜像
- `mcp_docker_image_pull`：拉取 Docker 镜像
- `mcp_docker_image_remove`：删除 Docker 镜像
- `mcp_docker_image_remove_batch`：批量删除多个 Docker 镜像
- `mcp_docker_image_details`：获取镜像详细信息

### 系统工具

- `mcp_docker_system_info`：测试 Docker 守护进程连接
- `mcp_docker_system_ping`：获取 Docker 详细系统信息
- `mcp_docker_system_server_version`：获取 Docker 版本信息
- `mcp_docker_system_disk_usage`：显示 Docker 磁盘使用情况

## 许可证

本项目采用 [MIT 许可证](LICENSE) 授权。 