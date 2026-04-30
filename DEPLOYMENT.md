# WatchUp Agent - Deployment Guide

This guide covers deploying the WatchUp Agent in various environments and scenarios.

## Table of Contents
- [Quick Deployment](#quick-deployment)
- [Production Deployment](#production-deployment)
- [Docker Deployment](#docker-deployment)
- [Cloud Platforms](#cloud-platforms)
- [Configuration Management](#configuration-management)
- [Monitoring & Troubleshooting](#monitoring--troubleshooting)

---

## Quick Deployment

### Single Server

**Linux/macOS**:
```bash
curl -fsSL https://raw.githubusercontent.com/tomurashigaraki22/watchup-agent-v2/main/install.sh | bash
```

**Windows**:
```powershell
iwr -useb https://raw.githubusercontent.com/watchup/watchup-agent/main/install.ps1 | iex
```

### Multiple Servers

Create a deployment script:

```bash
#!/bin/bash
# deploy-agents.sh

SERVERS=(
    "user@server1.example.com"
    "user@server2.example.com"
    "user@server3.example.com"
)

for SERVER in "${SERVERS[@]}"; do
    echo "Deploying to $SERVER..."
    ssh "$SERVER" 'curl -fsSL https://raw.githubusercontent.com/tomurashigaraki22/watchup-agent-v2/main/install.sh | bash'
done
```

---

## Production Deployment

### Prerequisites

1. **System Requirements**:
   - Go 1.19+ (if building from source)
   - Network access to `v2-server.watchup.site`
   - Standard user permissions (no root required for basic metrics)

2. **Firewall Rules**:
   - Outbound HTTPS (443) to `v2-server.watchup.site`
   - No inbound ports required

### Step-by-Step Production Deployment

#### 1. Download and Install

```bash
# Create installation directory
sudo mkdir -p /opt/watchup-agent
cd /opt/watchup-agent

# Download latest release
LATEST_VERSION=$(curl -s https://api.github.com/repos/watchup/watchup-agent/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
curl -fsSL "https://github.com/watchup/watchup-agent/releases/download/${LATEST_VERSION}/watchup-agent-linux-amd64" -o watchup-agent

# Make executable
chmod +x watchup-agent

# Create symlink
sudo ln -sf /opt/watchup-agent/watchup-agent /usr/local/bin/watchup-agent
```

#### 2. Configure

```bash
# Create config directory
sudo mkdir -p /etc/watchup-agent

# Download example config
sudo curl -fsSL https://raw.githubusercontent.com/watchup/watchup-agent/main/config.example.yaml -o /etc/watchup-agent/config.yaml

# Edit configuration
sudo nano /etc/watchup-agent/config.yaml
```

**Production Configuration Example**:
```yaml
server_id: "prod-web-01"
endpoint: "https://v2-server.watchup.site"
interval: 30s  # Production: 30s-60s recommended

metrics:
  cpu: true
  memory: true
  disk: true
  network: true
  connections: true  # Enable for detailed monitoring

auth:
  token_file: "/etc/watchup-agent/agent_token"

# Production port monitoring
ports:
  - port: 80
    name: "HTTP"
    host: "localhost"
    timeout: "3s"
  - port: 443
    name: "HTTPS"
    host: "localhost"
    timeout: "3s"
  - port: 22
    name: "SSH"
    host: "localhost"
    timeout: "2s"

# Production latency checks
latency_checks:
  - host: "8.8.8.8"
    name: "Google DNS"
    type: "tcp"
    port: 53
    timeout: "3s"
  - host: "v2-server.watchup.site"
    name: "WatchUp API"
    type: "tcp"
    port: 443
    timeout: "5s"
```

#### 3. Install as System Service

**Linux (systemd)**:
```bash
sudo tee /etc/systemd/system/watchup-agent.service > /dev/null <<EOF
[Unit]
Description=WatchUp Monitoring Agent
Documentation=https://github.com/watchup/watchup-agent
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/watchup-agent
WorkingDirectory=/etc/watchup-agent
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

# Security hardening
NoNewPrivileges=true
PrivateTmp=true
ProtectSystem=strict
ProtectHome=true
ReadWritePaths=/etc/watchup-agent

[Install]
WantedBy=multi-user.target
EOF

# Reload systemd
sudo systemctl daemon-reload

# Enable and start
sudo systemctl enable watchup-agent
sudo systemctl start watchup-agent

# Check status
sudo systemctl status watchup-agent
```

#### 4. Link Agent to Account

```bash
# View logs to get linking code
sudo journalctl -u watchup-agent -f

# You'll see output like:
# 🔗 Link this agent to your Watchup account:
#    Visit: https://v2-server.watchup.site/agent-link
#    Enter code: XK92-PQ
```

Visit the URL and enter the code to link the agent to your account.

#### 5. Verify Operation

```bash
# Check service status
sudo systemctl status watchup-agent

# View recent logs
sudo journalctl -u watchup-agent -n 50

# Check if metrics are being sent
sudo journalctl -u watchup-agent | grep "Metrics sent successfully"
```

---

## Docker Deployment

### Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o watchup-agent cmd/agent/main.go cmd/agent/setup.go

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy binary from builder
COPY --from=builder /build/watchup-agent .

# Copy config
COPY config.yaml .

# Run as non-root user
RUN adduser -D -u 1000 watchup
USER watchup

CMD ["./watchup-agent"]
```

### Docker Compose

```yaml
version: '3.8'

services:
  watchup-agent:
    build: .
    container_name: watchup-agent
    restart: unless-stopped
    volumes:
      - ./config.yaml:/app/config.yaml:ro
      - ./agent_token:/app/agent_token
    environment:
      - TZ=UTC
    network_mode: host  # Required for accurate network metrics
```

### Build and Run

```bash
# Build image
docker build -t watchup-agent .

# Run container
docker run -d \
  --name watchup-agent \
  --restart unless-stopped \
  --network host \
  -v $(pwd)/config.yaml:/app/config.yaml:ro \
  -v $(pwd)/agent_token:/app/agent_token \
  watchup-agent

# View logs
docker logs -f watchup-agent
```

---

## Cloud Platforms

### AWS EC2

```bash
#!/bin/bash
# User Data script for EC2 instance

# Install agent
curl -fsSL https://raw.githubusercontent.com/tomurashigaraki22/watchup-agent-v2/main/install.sh | bash

# Configure with instance metadata
INSTANCE_ID=$(ec2-metadata --instance-id | cut -d " " -f 2)
sudo sed -i "s/server_id: \"\"/server_id: \"aws-ec2-${INSTANCE_ID}\"/" /etc/watchup-agent/config.yaml

# Start service
sudo systemctl start watchup-agent
sudo systemctl enable watchup-agent
```

### Google Cloud Platform

```bash
#!/bin/bash
# Startup script for GCE instance

# Install agent
curl -fsSL https://raw.githubusercontent.com/tomurashigaraki22/watchup-agent-v2/main/install.sh | bash

# Configure with instance metadata
INSTANCE_NAME=$(curl -H "Metadata-Flavor: Google" http://metadata.google.internal/computeMetadata/v1/instance/name)
sudo sed -i "s/server_id: \"\"/server_id: \"gcp-${INSTANCE_NAME}\"/" /etc/watchup-agent/config.yaml

# Start service
sudo systemctl start watchup-agent
sudo systemctl enable watchup-agent
```

### Azure VM

```bash
#!/bin/bash
# Custom script extension for Azure VM

# Install agent
curl -fsSL https://raw.githubusercontent.com/tomurashigaraki22/watchup-agent-v2/main/install.sh | bash

# Configure with instance metadata
VM_NAME=$(curl -H Metadata:true "http://169.254.169.254/metadata/instance/compute/name?api-version=2021-02-01&format=text")
sudo sed -i "s/server_id: \"\"/server_id: \"azure-${VM_NAME}\"/" /etc/watchup-agent/config.yaml

# Start service
sudo systemctl start watchup-agent
sudo systemctl enable watchup-agent
```

### DigitalOcean Droplet

```bash
#!/bin/bash
# Cloud-init script for DigitalOcean

# Install agent
curl -fsSL https://raw.githubusercontent.com/tomurashigaraki22/watchup-agent-v2/main/install.sh | bash

# Configure
DROPLET_NAME=$(curl -s http://169.254.169.254/metadata/v1/hostname)
sudo sed -i "s/server_id: \"\"/server_id: \"do-${DROPLET_NAME}\"/" /etc/watchup-agent/config.yaml

# Start service
sudo systemctl start watchup-agent
sudo systemctl enable watchup-agent
```

---

## Configuration Management

### Ansible

```yaml
# playbook.yml
---
- name: Deploy WatchUp Agent
  hosts: all
  become: yes
  
  vars:
    watchup_version: "latest"
    watchup_endpoint: "https://v2-server.watchup.site"
    watchup_interval: "30s"
  
  tasks:
    - name: Download installation script
      get_url:
        url: https://raw.githubusercontent.com/watchup/watchup-agent/main/install.sh
        dest: /tmp/install-watchup.sh
        mode: '0755'
    
    - name: Run installation script
      shell: /tmp/install-watchup.sh
      args:
        creates: /usr/local/bin/watchup-agent
    
    - name: Configure agent
      template:
        src: config.yaml.j2
        dest: /etc/watchup-agent/config.yaml
        mode: '0644'
    
    - name: Start and enable service
      systemd:
        name: watchup-agent
        state: started
        enabled: yes
```

### Terraform

```hcl
# main.tf
resource "null_resource" "install_watchup_agent" {
  count = length(var.server_ips)
  
  connection {
    type        = "ssh"
    host        = var.server_ips[count.index]
    user        = "root"
    private_key = file(var.ssh_private_key)
  }
  
  provisioner "remote-exec" {
    inline = [
      "curl -fsSL https://raw.githubusercontent.com/tomurashigaraki22/watchup-agent-v2/main/install.sh | bash",
      "sed -i 's/server_id: \"\"/server_id: \"${var.server_names[count.index]}\"/' /etc/watchup-agent/config.yaml",
      "systemctl start watchup-agent",
      "systemctl enable watchup-agent"
    ]
  }
}
```

### Chef

```ruby
# recipes/default.rb
remote_file '/tmp/install-watchup.sh' do
  source 'https://raw.githubusercontent.com/watchup/watchup-agent/main/install.sh'
  mode '0755'
  action :create
end

execute 'install_watchup_agent' do
  command '/tmp/install-watchup.sh'
  creates '/usr/local/bin/watchup-agent'
end

template '/etc/watchup-agent/config.yaml' do
  source 'config.yaml.erb'
  mode '0644'
  notifies :restart, 'service[watchup-agent]'
end

service 'watchup-agent' do
  action [:enable, :start]
end
```

---

## Monitoring & Troubleshooting

### Health Checks

```bash
# Check if agent is running
systemctl is-active watchup-agent

# Check if metrics are being sent
journalctl -u watchup-agent --since "5 minutes ago" | grep "Metrics sent successfully"

# Check for errors
journalctl -u watchup-agent -p err --since today

# Check network connectivity
curl -I https://v2-server.watchup.site
```

### Common Issues

#### Agent Not Starting

```bash
# Check service status
sudo systemctl status watchup-agent

# View detailed logs
sudo journalctl -u watchup-agent -n 100 --no-pager

# Check configuration
sudo watchup-agent --validate-config  # If supported
```

#### Authentication Issues

```bash
# Check token file
ls -la /etc/watchup-agent/agent_token

# Re-authenticate
sudo rm /etc/watchup-agent/agent_token
sudo systemctl restart watchup-agent
# Follow linking instructions in logs
```

#### Network Issues

```bash
# Test connectivity
curl -v https://v2-server.watchup.site/health

# Check firewall
sudo iptables -L -n | grep 443

# Check DNS resolution
nslookup v2-server.watchup.site
```

### Performance Tuning

```yaml
# For high-frequency monitoring (development/testing)
interval: 5s

# For production (recommended)
interval: 30s

# For low-resource systems
interval: 60s
metrics:
  connections: false  # Disable expensive metrics
```

---

## Security Best Practices

1. **Token Security**:
   - Store tokens with restricted permissions (0600)
   - Never commit tokens to version control
   - Rotate tokens periodically

2. **Network Security**:
   - Use HTTPS only
   - Restrict outbound connections to WatchUp API
   - Monitor for unauthorized access attempts

3. **System Security**:
   - Run with minimal required permissions
   - Keep agent updated
   - Monitor agent logs for anomalies

4. **Configuration Security**:
   - Protect config files (0644 permissions)
   - Use environment variables for sensitive data
   - Audit configuration changes

---

## Backup & Recovery

### Backup Configuration

```bash
# Backup config and token
sudo tar -czf watchup-agent-backup.tar.gz \
  /etc/watchup-agent/config.yaml \
  /etc/watchup-agent/agent_token

# Store securely
scp watchup-agent-backup.tar.gz backup-server:/backups/
```

### Restore Configuration

```bash
# Restore from backup
sudo tar -xzf watchup-agent-backup.tar.gz -C /

# Restart service
sudo systemctl restart watchup-agent
```

---

## Scaling

### Load Balancer Monitoring

Deploy agents on each backend server with unique server_ids:

```yaml
# lb-backend-01
server_id: "lb-backend-01"

# lb-backend-02
server_id: "lb-backend-02"
```

### Container Orchestration

For Kubernetes, use DaemonSet to deploy on all nodes:

```yaml
apiVersion: apps/v1
kind: DaemonSet
metadata:
  name: watchup-agent
spec:
  selector:
    matchLabels:
      name: watchup-agent
  template:
    metadata:
      labels:
        name: watchup-agent
    spec:
      hostNetwork: true
      containers:
      - name: watchup-agent
        image: watchup/agent:latest
        volumeMounts:
        - name: config
          mountPath: /app/config.yaml
          subPath: config.yaml
      volumes:
      - name: config
        configMap:
          name: watchup-agent-config
```

---

## Support

For deployment issues:
- **Documentation**: Check this guide and README.md
- **Issues**: https://github.com/watchup/watchup-agent/issues
- **Community**: GitHub Discussions
- **Security**: security@watchup.com