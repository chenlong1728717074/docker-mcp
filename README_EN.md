# Docker MCP

English | [中文](README.md)

Docker MCP is a tool that provides Docker management capabilities through the Model-Command-Plugin (MCP) framework. It allows you to interact with Docker containers and images through a standardized interface.

## Features

- **Container Management**
  - List containers
  - Run containers
  - Start/stop/restart containers
  - Remove containers
  - View container logs and details

- **Image Management**
  - List images
  - Pull images
  - Remove images (single or batch)
  - View image details

- **System Management**
  - Check Docker daemon status
  - Get Docker version information
  - View system information
  - Check disk usage

## Requirements

- Go 1.24 or higher
- Docker (local or remote)
- Docker API access

## Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/docker-mcp.git
   cd docker-mcp
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Build the executable:
   ```bash
   go build -o docker-mcp.exe
   ```

## Usage

You can run Docker MCP directly with:

```bash
./docker-mcp.exe
```

### Environment Variables

- `DOCKER_PATH`: Docker daemon socket path or TCP endpoint (e.g., `tcp://your-docker-server:2375`)

### Command-line Arguments

- `--path`: Docker daemon socket path or TCP endpoint (overrides environment variable)

### Important Notes

To use the remote Docker API, you need to enable API access on your Docker host. There are several ways to do this:

#### Method 1: Modify Docker Daemon Configuration File

1. Edit the Docker daemon configuration file (e.g., `/etc/docker/daemon.json`), and add the following content:
   ```json
   {
     "hosts": ["tcp://0.0.0.0:2375", "unix:///var/run/docker.sock"]
   }
   ```

2. Restart the Docker service:
   ```bash
   sudo systemctl restart docker
   ```

#### Method 2: Modify Docker Service Startup Parameters

1. For systems using systemd, edit the Docker service file:
   ```bash
   sudo systemctl edit docker.service
   ```

2. Add the following content:
   ```ini
   [Service]
   ExecStart=
   ExecStart=/usr/bin/dockerd -H fd:// -H tcp://0.0.0.0:2375
   ```

3. Reload systemd configuration and restart Docker:
   ```bash
   sudo systemctl daemon-reload
   sudo systemctl restart docker
   ```

#### Verify API Access

Confirm that the Docker API is enabled:
```bash
curl http://localhost:2375/version
```

**Security Warning**: Opening port 2375 allows unauthenticated access to the Docker API. In production environments, it is recommended to use TLS certificates (port 2376) or set up network security groups/firewall rules to restrict access. Only use this configuration in trusted network environments.

## Cursor Integration

Docker MCP can be integrated with Cursor IDE to provide Docker management capabilities directly within the editor.

### Configuration Steps

1. Open Cursor settings
2. Navigate to the MCP configuration section
3. Add the following configuration to your Cursor settings:

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

4. Save the settings and restart Cursor

### Configuration Options

- `command`: Path to the docker-mcp executable
- `args`: Additional command-line arguments
- `env`: Environment variables to pass to the executable
  - `DOCKER_PATH`: Docker daemon socket path or TCP endpoint

## Available Tools

### Container Tools

- `mcp_docker_container_list`: List all containers
- `mcp_docker_container_run`: Run a Docker image
- `mcp_docker_container_start`: Start a stopped container
- `mcp_docker_container_stop`: Stop a running container
- `mcp_docker_container_restart`: Restart a container
- `mcp_docker_container_remove`: Remove a container
- `mcp_docker_container_details`: Get detailed information about a container
- `mcp_docker_container_log`: Get container logs

### Image Tools

- `mcp_docker_image_list`: List all Docker images
- `mcp_docker_image_pull`: Pull a Docker image
- `mcp_docker_image_remove`: Remove a Docker image
- `mcp_docker_image_remove_batch`: Remove multiple Docker images in batch
- `mcp_docker_image_details`: Get detailed information about an image

### System Tools

- `mcp_docker_system_info`: Test Docker daemon connectivity
- `mcp_docker_system_ping`: Get detailed Docker system information
- `mcp_docker_system_server_version`: Get Docker version information
- `mcp_docker_system_disk_usage`: Show Docker disk usage

## License

This project is licensed under the [MIT License](LICENSE). 