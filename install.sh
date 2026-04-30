#!/bin/bash
# WatchUp Agent Installation Script
# Supports: Linux and macOS
# Usage: curl -fsSL https://raw.githubusercontent.com/tomurashigaraki22/watchup-agent-v2/main/install.sh | bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="tomurashigaraki22/watchup-agent-v2"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="/etc/watchup-agent"
SERVICE_DIR="/etc/systemd/system"
BINARY_NAME="watchup-agent"

# Detect OS and Architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$ARCH" in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        armv7l)
            ARCH="armv7"
            ;;
        *)
            echo -e "${RED}Unsupported architecture: $ARCH${NC}"
            exit 1
            ;;
    esac
    
    case "$OS" in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        *)
            echo -e "${RED}Unsupported operating system: $OS${NC}"
            exit 1
            ;;
    esac
    
    echo -e "${BLUE}Detected platform: ${OS}-${ARCH}${NC}"
}

# Get latest release version
get_latest_version() {
    echo -e "${BLUE}Fetching latest version...${NC}"
    VERSION=$(curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
    
    if [ -z "$VERSION" ]; then
        echo -e "${YELLOW}Could not fetch latest version, using 'latest'${NC}"
        VERSION="latest"
    else
        echo -e "${GREEN}Latest version: ${VERSION}${NC}"
    fi
}

# Download and install binary
install_binary() {
    echo -e "${BLUE}Downloading WatchUp Agent...${NC}"
    
    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/watchup-agent-${OS}-${ARCH}"
    
    # Create temporary directory
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"
    
    # Download binary
    if ! curl -fsSL "$DOWNLOAD_URL" -o "$BINARY_NAME"; then
        echo -e "${RED}Failed to download binary from ${DOWNLOAD_URL}${NC}"
        echo -e "${YELLOW}Attempting to build from source...${NC}"
        build_from_source
        return
    fi
    
    # Make executable
    chmod +x "$BINARY_NAME"
    
    # Install binary
    echo -e "${BLUE}Installing binary to ${INSTALL_DIR}...${NC}"
    sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
    
    # Cleanup
    cd - > /dev/null
    rm -rf "$TMP_DIR"
    
    echo -e "${GREEN}Binary installed successfully!${NC}"
}

# Build from source if binary not available
build_from_source() {
    echo -e "${BLUE}Building from source...${NC}"
    
    # Check if Go is installed
    if ! command -v go &> /dev/null; then
        echo -e "${RED}Go is not installed. Please install Go 1.19+ and try again.${NC}"
        echo -e "${YELLOW}Visit: https://golang.org/doc/install${NC}"
        exit 1
    fi
    
    # Clone repository
    TMP_DIR=$(mktemp -d)
    cd "$TMP_DIR"
    
    echo -e "${BLUE}Cloning repository...${NC}"
    git clone "https://github.com/${REPO}.git"
    cd watchup-agent-v2
    
    # Build
    echo -e "${BLUE}Building binary...${NC}"
    go build -o "$BINARY_NAME" cmd/agent/main.go cmd/agent/setup.go
    
    # Install
    echo -e "${BLUE}Installing binary to ${INSTALL_DIR}...${NC}"
    sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
    
    # Cleanup
    cd - > /dev/null
    rm -rf "$TMP_DIR"
    
    echo -e "${GREEN}Built and installed from source!${NC}"
}

# Create configuration directory and file
setup_config() {
    echo -e "${BLUE}Setting up configuration...${NC}"
    
    # Create config directory
    sudo mkdir -p "$CONFIG_DIR"
    
    # Download example config
    if ! sudo curl -fsSL "https://raw.githubusercontent.com/${REPO}/main/config.example.yaml" -o "${CONFIG_DIR}/config.yaml"; then
        echo -e "${YELLOW}Could not download config, creating default...${NC}"
        sudo tee "${CONFIG_DIR}/config.yaml" > /dev/null <<EOF
# WatchUp Agent Configuration
server_id: ""
endpoint: "https://v2-server.watchup.site"
interval: 5s

metrics:
  cpu: true
  memory: true
  disk: true
  network: true
  connections: false

auth:
  token_file: "${CONFIG_DIR}/agent_token"
EOF
    fi
    
    echo -e "${GREEN}Configuration created at ${CONFIG_DIR}/config.yaml${NC}"
}

# Install systemd service (Linux only)
install_service() {
    if [ "$OS" != "linux" ]; then
        echo -e "${YELLOW}Systemd service installation skipped (not Linux)${NC}"
        return
    fi
    
    echo -e "${BLUE}Installing systemd service...${NC}"
    
    sudo tee "${SERVICE_DIR}/watchup-agent.service" > /dev/null <<EOF
[Unit]
Description=WatchUp Monitoring Agent
Documentation=https://github.com/${REPO}
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=root
ExecStart=${INSTALL_DIR}/${BINARY_NAME}
WorkingDirectory=${CONFIG_DIR}
Restart=always
RestartSec=10
StandardOutput=journal
StandardError=journal

[Install]
WantedBy=multi-user.target
EOF
    
    # Reload systemd
    sudo systemctl daemon-reload
    
    echo -e "${GREEN}Systemd service installed!${NC}"
}

# Install launchd service (macOS only)
install_launchd() {
    if [ "$OS" != "darwin" ]; then
        return
    fi
    
    echo -e "${BLUE}Installing launchd service...${NC}"
    
    PLIST_PATH="$HOME/Library/LaunchAgents/com.watchup.agent.plist"
    
    mkdir -p "$HOME/Library/LaunchAgents"
    
    tee "$PLIST_PATH" > /dev/null <<EOF
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>Label</key>
    <string>com.watchup.agent</string>
    <key>ProgramArguments</key>
    <array>
        <string>${INSTALL_DIR}/${BINARY_NAME}</string>
    </array>
    <key>WorkingDirectory</key>
    <string>${CONFIG_DIR}</string>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardOutPath</key>
    <string>/tmp/watchup-agent.log</string>
    <key>StandardErrorPath</key>
    <string>/tmp/watchup-agent.error.log</string>
</dict>
</plist>
EOF
    
    echo -e "${GREEN}Launchd service installed!${NC}"
}

# Print next steps
print_next_steps() {
    echo ""
    echo -e "${GREEN}╔════════════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║  WatchUp Agent Installation Complete! 🎉                  ║${NC}"
    echo -e "${GREEN}╚════════════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${BLUE}📋 Next Steps:${NC}"
    echo ""
    echo -e "1. ${YELLOW}Configure your server ID:${NC}"
    echo -e "   sudo nano ${CONFIG_DIR}/config.yaml"
    echo ""
    
    if [ "$OS" = "linux" ]; then
        echo -e "2. ${YELLOW}Start the agent:${NC}"
        echo -e "   sudo systemctl start watchup-agent"
        echo ""
        echo -e "3. ${YELLOW}Enable auto-start on boot:${NC}"
        echo -e "   sudo systemctl enable watchup-agent"
        echo ""
        echo -e "4. ${YELLOW}Check status:${NC}"
        echo -e "   sudo systemctl status watchup-agent"
        echo ""
        echo -e "5. ${YELLOW}View logs:${NC}"
        echo -e "   sudo journalctl -u watchup-agent -f"
    elif [ "$OS" = "darwin" ]; then
        echo -e "2. ${YELLOW}Start the agent:${NC}"
        echo -e "   launchctl load ~/Library/LaunchAgents/com.watchup.agent.plist"
        echo ""
        echo -e "3. ${YELLOW}Check status:${NC}"
        echo -e "   launchctl list | grep watchup"
        echo ""
        echo -e "4. ${YELLOW}View logs:${NC}"
        echo -e "   tail -f /tmp/watchup-agent.log"
    else
        echo -e "2. ${YELLOW}Run the agent manually:${NC}"
        echo -e "   cd ${CONFIG_DIR} && sudo ${INSTALL_DIR}/${BINARY_NAME}"
    fi
    
    echo ""
    echo -e "${BLUE}📚 Documentation:${NC}"
    echo -e "   https://github.com/${REPO}"
    echo ""
    echo -e "${BLUE}🔗 Link your agent:${NC}"
    echo -e "   The agent will display a link and code on first run"
    echo ""
}

# Main installation flow
main() {
    echo -e "${GREEN}"
    echo "╔════════════════════════════════════════════════════════════╗"
    echo "║         WatchUp Agent Installer                           ║"
    echo "║         https://watchup.site                              ║"
    echo "╚════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    
    # Check if running as root for system-wide installation
    if [ "$EUID" -eq 0 ]; then
        echo -e "${YELLOW}Running as root. Installing system-wide.${NC}"
    else
        echo -e "${YELLOW}Not running as root. Will use sudo for system operations.${NC}"
    fi
    
    detect_platform
    get_latest_version
    install_binary
    setup_config
    install_service
    install_launchd
    print_next_steps
}

# Run main function
main