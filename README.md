# WatchUp Agent

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.19+-blue.svg)](https://golang.org)
[![Platform](https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20macOS-lightgrey.svg)](https://github.com/tomurashigaraki22/watchup-agent-v2)

A production-grade, open source Go agent that collects comprehensive system metrics and sends them to the WatchUp monitoring platform. The agent uses secure device linking authentication and runs as a lightweight background service.

## 🌟 Open Source

**WatchUp Agent is fully open source** under the MIT License. You can:
- ✅ **Audit the code** for security and functionality
- ✅ **Customize and modify** for your specific needs  
- ✅ **Contribute features** and improvements
- ✅ **Self-host** against custom backends
- ✅ **Use in enterprise** environments with confidence

The agent integrates with the **WatchUp platform** (proprietary) but can be adapted for other monitoring backends.

## Features

- **Device Linking Authentication**: Secure one-time setup similar to GitHub CLI
- **Comprehensive System Metrics**: CPU, Memory, Disk, Network monitoring
- **Extended Network Monitoring**: Active connections, port checks, latency monitoring
- **Real-time Collection**: Configurable intervals (1s to hours)
- **Robust Communication**: Retry logic with exponential backoff
- **Production Ready**: Graceful shutdown, error handling, comprehensive logging
- **Lightweight**: Minimal resource usage (<5% CPU, <50MB RAM)
- **Cross-platform**: Works on Windows, Linux, macOS

## Quick Start

### One-Line Installation

#### Linux & macOS:
```bash
curl -fsSL https://raw.githubusercontent.com/tomurashigaraki22/watchup-agent-v2/main/install.sh | bash
```

#### Windows (PowerShell as Administrator):
```powershell
iwr -useb https://raw.githubusercontent.com/watchup/watchup-agent/main/install.ps1 | iex
```

### Manual Installation

#### Option 1: Download Pre-built Binary

1. **Download the latest release** for your platform:
   - Visit [Releases](https://github.com/tomurashigaraki22/watchup-agent-v2/releases/latest)
   - Download the appropriate binary:
     - Linux: `watchup-agent-linux-amd64`
     - macOS: `watchup-agent-darwin-amd64` or `watchup-agent-darwin-arm64`
     - Windows: `watchup-agent-windows-amd64.exe`

2. **Install the binary**:
   
   **Linux/macOS**:
   ```bash
   # Make executable
   chmod +x watchup-agent-*
   
   # Move to system path
   sudo mv watchup-agent-* /usr/local/bin/watchup-agent
   ```
   
   **Windows**:
   ```powershell
   # Move to Program Files
   Move-Item watchup-agent-windows-amd64.exe "C:\Program Files\WatchUp\watchup-agent.exe"
   ```

3. **Create configuration**:
   ```bash
   # Download example config
   curl -fsSL https://raw.githubusercontent.com/watchup/watchup-agent/main/config.example.yaml -o config.yaml
   
   # Edit config.yaml and set your server_id
   nano config.yaml
   ```

4. **Run the agent**:
   ```bash
   ./watchup-agent
   ```

#### Option 2: Build from Source

1. **Prerequisites**:
   - Go 1.19 or later
   - Git

2. **Clone and build**:
   ```bash
   git clone https://github.com/tomurashigaraki22/watchup-agent-v2.git
   cd watchup-agent
   go build -o watchup-agent cmd/agent/main.go cmd/agent/setup.go
   ```

3. **Run**:
   ```bash
   ./watchup-agent
   ```

### First Run Setup

1. **Configure server ID**:
   - Agent will prompt for a **server ID** if not configured
   - Enter a unique identifier for your server (e.g., "web-prod-01")
   - **Important**: Server ID must be unique within your WatchUp account

2. **Link to your account**:
   - Visit the provided URL in your browser
   - Enter the displayed code to approve the agent
   - Agent will start collecting metrics automatically

> **Note**: Each agent is linked to your WatchUp user account. The server_id you choose must be unique within your account but can be reused by other users.

### Running as a Service

#### Linux (systemd):
```bash
# Create service file
sudo tee /etc/systemd/system/watchup-agent.service > /dev/null <<EOF
[Unit]
Description=WatchUp Monitoring Agent
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
ExecStart=/usr/local/bin/watchup-agent
WorkingDirectory=/etc/watchup-agent
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF

# Start and enable service
sudo systemctl daemon-reload
sudo systemctl start watchup-agent
sudo systemctl enable watchup-agent

# Check status
sudo systemctl status watchup-agent

# View logs
sudo journalctl -u watchup-agent -f
```

#### macOS (launchd):
```bash
# Create plist file
tee ~/Library/LaunchAgents/com.watchup.agent.plist > /dev/null <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.watchup.agent</string>
    <key>ProgramArguments</key>
    <array>
        <string>/usr/local/bin/watchup-agent</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
</dict>
</plist>
EOF

# Load service
launchctl load ~/Library/LaunchAgents/com.watchup.agent.plist

# Check status
launchctl list | grep watchup
```

#### Windows (Service):
```powershell
# Install as Windows Service
sc.exe create WatchUpAgent binPath= "C:\Program Files\WatchUp\watchup-agent.exe" start= auto

# Start service
Start-Service WatchUpAgent

# Check status
Get-Service WatchUpAgent
```

## Configuration

See `config.example.yaml` for all available options.

### Core Settings:
- `server_id`: Unique identifier for this server
- `endpoint`: Your backend API URL
- `interval`: How often to send metrics (e.g., "5s", "1m")
- `metrics`: Which metrics to collect (cpu, memory, disk, network, connections)

### Network Monitoring (Phase 4):
- `ports`: List of ports to monitor for availability
- `latency_checks`: List of hosts/services to check latency

Example:
```yaml
# Enable connection monitoring
metrics:
  connections: true

# Monitor specific ports
ports:
  - port: 80
    name: "HTTP"
    host: "localhost"
  - port: 443
    name: "HTTPS"
    host: "mysite.com"

# Check latency to external services
latency_checks:
  - host: "8.8.8.8"
    name: "Google DNS"
    type: "tcp"
    port: 53
  - host: "https://api.myservice.com"
    name: "My API"
    type: "http"
```

## Development Status

This project is being built in phases:

- [x] **Phase 0**: Project Setup & Foundation ✅
- [x] **Phase 1**: Authentication System (Device Linking) ✅
- [x] **Phase 2**: Core Metrics Collection ✅
- [x] **Phase 3**: Main Agent Loop & Communication ✅
- [x] **Phase 4**: Network Monitoring (Extended) ✅
- [ ] **Phase 5**: Process Monitoring
- [ ] **Phase 6**: Service Monitoring
- [ ] **Phase 7**: Deployment & Production Readiness
- [ ] **Phase 8**: Testing & Refinement

## Architecture

```
watchup-agent/
├── cmd/agent/          # Main application entry point
├── internal/
│   ├── auth/          # Authentication & device linking
│   ├── config/        # Configuration management
│   ├── metrics/       # System metrics collection
│   └── client/        # HTTP client & communication
├── config.yaml        # Agent configuration
└── README.md
```

## Requirements

- Go 1.19 or later
- Network access to WatchUp backend API
- Standard user permissions (no root required for basic metrics)

## 🤝 Contributing

We welcome contributions! This is an open source project under the MIT License.

### Ways to Contribute:
- **Bug Reports**: Use GitHub Issues
- **Feature Requests**: Describe your use case in an issue
- **Code Contributions**: Fork, create a feature branch, and submit a PR
- **Documentation**: Help improve our docs and examples

### Development Setup

1. **Fork and clone**:
   ```bash
   git clone https://github.com/YOUR_USERNAME/watchup-agent.git
   cd watchup-agent
   ```

2. **Install dependencies**:
   ```bash
   go mod download
   ```

3. **Make your changes**:
   ```bash
   # Create a feature branch
   git checkout -b feature/amazing-feature
   
   # Make your changes
   # ...
   
   # Test your changes
   go test ./...
   ```

4. **Build and test**:
   ```bash
   # Build for your platform
   go build -o watchup-agent cmd/agent/main.go cmd/agent/setup.go
   
   # Test the binary
   ./watchup-agent
   ```

5. **Build for all platforms** (optional):
   ```bash
   # Linux/macOS
   ./build.sh
   
   # Windows
   .\build.ps1
   ```

6. **Submit a pull request**:
   ```bash
   git add .
   git commit -m "Add amazing feature"
   git push origin feature/amazing-feature
   ```
   Then open a PR on GitHub!

### Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/metrics/

# Run with verbose output
go test -v ./...
```

### Code Style

- Follow standard Go conventions
- Use `gofmt` to format code
- Add comments for exported functions
- Write tests for new functionality

See [CONTRIBUTING.md](CONTRIBUTING.md) for detailed guidelines.

---

## 🚀 Deployment & Releases

### Creating a Release

1. **Tag a new version**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **GitHub Actions will automatically**:
   - Build binaries for all platforms
   - Create a GitHub release
   - Upload all binaries as release assets

### Manual Release Build

```bash
# Build all platforms
./build.sh v1.0.0

# Binaries will be in dist/
ls dist/
```

### Supported Platforms

| OS | Architecture | Binary Name |
|----|--------------|-------------|
| Linux | amd64 | `watchup-agent-linux-amd64` |
| Linux | arm64 | `watchup-agent-linux-arm64` |
| Linux | armv7 | `watchup-agent-linux-armv7` |
| macOS | amd64 (Intel) | `watchup-agent-darwin-amd64` |
| macOS | arm64 (M1/M2) | `watchup-agent-darwin-arm64` |
| Windows | amd64 | `watchup-agent-windows-amd64.exe` |
| Windows | arm64 | `watchup-agent-windows-arm64.exe` |

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🔗 Links

- **WatchUp Platform**: [https://watchup.site](https://watchup.site)
- **Documentation**: [MONITORING_CAPABILITIES.md](MONITORING_CAPABILITIES.md)
- **Deployment Guide**: [DEPLOYMENT.md](DEPLOYMENT.md)
- **Architecture**: [OPEN_SOURCE_ARCHITECTURE.md](OPEN_SOURCE_ARCHITECTURE.md)
- **Contributing**: [CONTRIBUTING.md](CONTRIBUTING.md)
- **Issues**: [GitHub Issues](https://github.com/tomurashigaraki22/watchup-agent-v2/issues)

## 🙋 Support

- **Documentation**: Check the docs in this repository
- **Community**: GitHub Discussions
- **Issues**: GitHub Issues for bugs and feature requests
- **Security**: Email security@watchup.com for security issues

---

**Made with ❤️ by the WatchUp team and contributors**