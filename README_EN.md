# Docker MCP

Docker MCP is a Docker management tool based on the Model Context Protocol (MCP), providing comprehensive Docker operation capabilities for AI assistants. Through the MCP protocol, AI assistants can directly manage Docker containers, images, networks, volumes, and system resources.

## Features

### 🐳 Container Management
- List, create, start, stop, and restart containers
- Remove containers and view detailed container information
- Real-time container log viewing
- Support for advanced configurations like environment variables, port mapping, and volume mounting

### 🖼️ Image Management
- List local images and pull remote images
- Remove single or batch remove images
- View detailed image information and layer structure

### 🌐 Network Management
- Create, remove, and inspect Docker networks
- Connect and disconnect containers from networks
- List all networks and prune unused networks
- Support for custom network configurations (driver, subnet, gateway, etc.)

### 💾 Volume Management
- Create, remove, and inspect Docker volumes
- List all volumes and prune unused volumes
- Support for custom volume drivers and options

### ⚙️ System Management
- Check Docker daemon connectivity status
- Get system information and version information
- Monitor disk usage

### 🔐 Authentication Support
- Docker Registry login authentication
- Support for private registry access

## System Requirements

- Go 1.21 or higher
- Docker Engine (local or remote)
- Supported operating systems: Linux, macOS, Windows

## Building from Source

1. Clone the repository:
   ```bash
   git clone https://github.com/chenlong1728717074/docker-mcp.git
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

- `DOCKER_PATH`: Docker daemon socket path or TCP endpoint (e.g., `tcp://your-docker-server:2375` or TLS-enabled `tcp://your-docker-server:2376`)
- `DOCKER_CERT`: Path to TLS certificate directory (required when using port 2376 with TLS authentication). This directory must contain the following three files:
  - `ca.pem`: CA certificate file
  - `cert.pem`: Client certificate file
  - `key.pem`: Client private key file

### Command-line Arguments

- `--path`: Docker daemon socket path or TCP endpoint (overrides environment variable)
- `--cert`: Path to TLS certificate directory (overrides environment variable). The directory structure is the same as required for `DOCKER_CERT`

### Important Notes

To use the remote Docker API, you need to enable API access on your Docker host. There are several ways to do this:

#### Method 1: Modify Docker Daemon Configuration File

1. Edit the Docker daemon configuration file (e.g., `/etc/docker/daemon.json`), and add the following content:
   ```json
   {
     "hosts": ["tcp://0.0.0.0:2375", "unix:///var/run/docker.sock"]
   }
   ```

   Or enable TLS (recommended for production environments):
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

   Note: The certificate file paths in the above configuration need to match the actual certificate file paths on the server. Additionally, the client needs to connect using client certificates issued by the same CA.

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
   
   Or enable TLS (recommended for production environments):
   ```ini
   [Service]
   ExecStart=
   ExecStart=/usr/bin/dockerd -H fd:// -H tcp://0.0.0.0:2376 --tlsverify --tlscacert=/path/to/ca.pem --tlscert=/path/to/cert.pem --tlskey=/path/to/key.pem
   ```

   Note: The certificate files in the above configuration must be issued by the same CA as the client's certificates to ensure mutual authentication security.

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

For TLS-enabled connections:
```bash
curl --cacert /path/to/ca.pem --cert /path/to/cert.pem --key /path/to/key.pem https://localhost:2376/version
```

**Security Warning**: Opening port 2375 allows unauthenticated access to the Docker API. In production environments, it is recommended to use TLS certificates (port 2376) or set up network security groups/firewall rules to restrict access. Only use port 2375 configuration in trusted network environments.

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
        "DOCKER_PATH": "tcp://your-docker-server:2375",//tls:2376
        "DOCKER_CERT": "{your-cert-path}" // Directory containing ca.pem, cert.pem, and key.pem
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
  - `DOCKER_CERT`: Path to TLS certificate directory (required when using TLS-enabled connections)

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

#### Generating TLS Certificates

To use TLS secure connections, you need to generate three certificate files: ca.pem, cert.pem, and key.pem. You can generate them using the following steps:

1. Install OpenSSL tool

2. Generate CA private key and certificate:
   ```bash
   openssl genrsa -out ca-key.pem 4096
   openssl req -new -x509 -days 365 -key ca-key.pem -out ca.pem
   ```

3. Generate server key and certificate signing request:
   ```bash
   openssl genrsa -out server-key.pem 4096
   openssl req -subj "/CN=your-docker-server" -new -key server-key.pem -out server.csr
   ```

4. Create server certificate:
   ```bash
   openssl x509 -req -days 365 -in server.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out server-cert.pem
   ```

5. Generate client key and certificate signing request:
   ```bash
   openssl genrsa -out key.pem 4096
   openssl req -subj "/CN=client" -new -key key.pem -out client.csr
   ```

6. Create client certificate:
   ```bash
   openssl x509 -req -days 365 -in client.csr -CA ca.pem -CAkey ca-key.pem -CAcreateserial -out cert.pem
   ```

7. Set correct file permissions:
   ```bash
   chmod 0400 ca-key.pem key.pem server-key.pem
   chmod 0444 ca.pem server-cert.pem cert.pem
   ```

8. On the server side, configure:
   - ca.pem (CA certificate)
   - server-cert.pem (renamed to cert.pem)
   - server-key.pem (renamed to key.pem)

9. On the client side, use:
   - ca.pem (CA certificate)
   - cert.pem (client certificate)
   - key.pem (client private key)