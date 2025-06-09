# Docker MCP

Docker MCP 是一个基于 Model Context Protocol (MCP) 的 Docker 管理工具，为 AI 助手提供了完整的 Docker 操作能力。通过 MCP 协议，AI 助手可以直接管理 Docker 容器、镜像、网络、卷和系统资源。

## 功能特点

### 🐳 容器管理
- 列出、创建、启动、停止、重启容器
- 删除容器和查看容器详细信息
- 实时查看容器日志
- 支持环境变量、端口映射、卷挂载等高级配置

### 🖼️ 镜像管理
- 列出本地镜像和拉取远程镜像
- 删除单个或批量删除镜像
- 查看镜像详细信息和层级结构

### 🌐 网络管理
- 创建、删除、检查 Docker 网络
- 连接和断开容器与网络
- 列出所有网络和清理未使用的网络
- 支持自定义网络配置（驱动、子网、网关等）

### 💾 卷管理
- 创建、删除、检查 Docker 卷
- 列出所有卷和清理未使用的卷
- 支持自定义卷驱动和选项

### ⚙️ 系统管理
- 检查 Docker 守护进程连接状态
- 获取系统信息和版本信息
- 监控磁盘使用情况

### 🔐 认证支持
- Docker Registry 登录认证
- 支持私有镜像仓库访问

## 系统要求

- Go 1.21 或更高版本
- Docker Engine（本地或远程）
- 支持的操作系统：Linux、macOS、Windows

## 从源码构建

1. 克隆仓库：
   ```bash
   git clone https://gitee.com/a-little-dragon/docker-mcp.git
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

- `DOCKER_PATH`：Docker 守护进程套接字路径或 TCP 端点（例如：`tcp://your-docker-server:2375` 或启用TLS的 `tcp://your-docker-server:2376`）
- `DOCKER_CERT`：TLS证书目录路径（当使用2376端口带TLS验证时需要）。该目录必须包含以下三个文件：
  - `ca.pem`：CA证书文件
  - `cert.pem`：客户端证书文件
  - `key.pem`：客户端私钥文件

### 命令行参数

- `--path`：Docker 守护进程套接字路径或 TCP 端点（覆盖环境变量）
- `--cert`：TLS证书目录路径（覆盖环境变量）。目录结构同上述`DOCKER_CERT`要求

### 重要注意事项

为了使用远程 Docker API，您需要在 Docker 主机上启用 API 访问。有以下几种方式：

#### 方式一：修改 Docker 守护进程配置文件

1. 修改 Docker 守护进程配置文件（例如 `/etc/docker/daemon.json`），添加以下内容：
   ```json
   {
     "hosts": ["tcp://0.0.0.0:2375", "unix:///var/run/docker.sock"]
   }
   ```

   或者启用TLS（推荐用于生产环境）：
   ```json
   {
     "hosts": ["tcp://0.0.0.0:2376", "unix:///var/run/docker.sock"],
     "tls": true,
     "tlsverify": true,
     "tlscacert": "/path/to/ca.pem",
     "tlscert": "/path/to/cert.pem",
     "tlskey": "/path/to/key.pem"
   }
   ```

   注意：以上配置中的证书文件路径需要与服务器上的实际证书文件路径一致。同时，客户端需要使用相同的CA签发的客户端证书进行连接。

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

   或者启用TLS（推荐用于生产环境）：
   ```ini
   [Service]
   ExecStart=
   ExecStart=/usr/bin/dockerd -H fd:// -H tcp://0.0.0.0:2376 --tlsverify --tlscacert=/path/to/ca.pem --tlscert=/path/to/cert.pem --tlskey=/path/to/key.pem
   ```

   注意：以上配置中的证书文件与客户端使用的证书必须由同一个CA签发，以确保相互认证的安全性。

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

对于启用TLS的连接：
```bash
curl --cacert /path/to/ca.pem --cert /path/to/cert.pem --key /path/to/key.pem https://localhost:2376/version
```

**安全警告**：开放 2375 端口允许未经身份验证的 Docker API 访问。在生产环境中，建议使用 TLS 证书（2376 端口）或设置网络安全组/防火墙规则限制访问。仅在受信任的网络环境中使用2375端口的配置。

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
        "DOCKER_PATH": "tcp://your-docker-server:2375", //tls:2376
        "DOCKER_CERT": "{your-cert-path}" // 包含ca.pem、cert.pem和key.pem的目录路径
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
  - `DOCKER_CERT`：TLS证书目录路径（使用启用TLS的连接时需要提供）

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

#### 生成TLS证书

为了使用TLS安全连接，您需要生成三个证书文件：ca.pem、cert.pem和key.pem。您可以使用以下步骤生成：

1. 安装OpenSSL工具

2. 生成CA私钥和证书：
   ```bash
   openssl genrsa -out ca-key.pem 4096
   openssl req -new -x509 -days 365 -key ca-key.pem -out ca.pem
   ```

3. 生成服务器密钥和证书签名请求：
   ```bash
   openssl genrsa -out server-key.pem 4096
   openssl req -subj "/CN=your-docker-server" -new -key server-key.pem -out server.csr
   ```

4. 创建服务器证书：
   ```bash
   openssl x509 -req -days 365 -in server.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem
   ```

5. 生成客户端密钥和证书签名请求：
   ```bash
   openssl genrsa -out key.pem 4096
   openssl req -subj "/CN=client" -new -key key.pem -out client.csr
   ```

6. 创建客户端证书：
   ```bash
   openssl x509 -req -days 365 -in client.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out cert.pem
   ```

7. 设置正确的文件权限：
   ```bash
   chmod 0400 ca-key.pem key.pem server-key.pem
   chmod 0444 ca.pem server-cert.pem cert.pem
   ```

8. 在服务器端配置：
   - ca.pem (CA证书)
   - server-cert.pem (重命名为cert.pem)
   - server-key.pem (重命名为key.pem)

9. 在客户端使用：
   - ca.pem (CA证书)
   - cert.pem (客户端证书)
   - key.pem (客户端私钥)