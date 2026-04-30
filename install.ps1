# WatchUp Agent Installation Script for Windows
# Usage: iwr -useb https://raw.githubusercontent.com/tomurashigaraki22/watchup-agent-v2/main/install.ps1 | iex

$ErrorActionPreference = "Stop"

# Configuration
$REPO = "tomurashigaraki22/watchup-agent-v2"
$INSTALL_DIR = "$env:ProgramFiles\WatchUp\Agent"
$CONFIG_DIR = "$env:ProgramData\WatchUp\Agent"
$BINARY_NAME = "watchup-agent.exe"

# Colors for output
function Write-ColorOutput($ForegroundColor) {
    $fc = $host.UI.RawUI.ForegroundColor
    $host.UI.RawUI.ForegroundColor = $ForegroundColor
    if ($args) {
        Write-Output $args
    }
    $host.UI.RawUI.ForegroundColor = $fc
}

function Write-Info($message) {
    Write-ColorOutput Cyan $message
}

function Write-Success($message) {
    Write-ColorOutput Green $message
}

function Write-Warning($message) {
    Write-ColorOutput Yellow $message
}

function Write-Error($message) {
    Write-ColorOutput Red $message
}

# Check if running as Administrator
function Test-Administrator {
    $currentUser = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    return $currentUser.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
}

# Detect architecture
function Get-Architecture {
    $arch = $env:PROCESSOR_ARCHITECTURE
    switch ($arch) {
        "AMD64" { return "amd64" }
        "ARM64" { return "arm64" }
        default {
            Write-Error "Unsupported architecture: $arch"
            exit 1
        }
    }
}

# Get latest release version
function Get-LatestVersion {
    Write-Info "Fetching latest version..."
    try {
        $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$REPO/releases/latest"
        $version = $response.tag_name
        Write-Success "Latest version: $version"
        return $version
    } catch {
        Write-Warning "Could not fetch latest version, using 'latest'"
        return "latest"
    }
}

# Download and install binary
function Install-Binary {
    param($version, $arch)
    
    Write-Info "Downloading WatchUp Agent..."
    
    $downloadUrl = "https://github.com/$REPO/releases/download/$version/watchup-agent-windows-$arch.exe"
    $tempFile = "$env:TEMP\$BINARY_NAME"
    
    try {
        Invoke-WebRequest -Uri $downloadUrl -OutFile $tempFile
    } catch {
        Write-Error "Failed to download binary from $downloadUrl"
        Write-Warning "Please build from source or download manually"
        exit 1
    }
    
    # Create installation directory
    Write-Info "Installing binary to $INSTALL_DIR..."
    New-Item -ItemType Directory -Force -Path $INSTALL_DIR | Out-Null
    
    # Copy binary
    Copy-Item -Path $tempFile -Destination "$INSTALL_DIR\$BINARY_NAME" -Force
    
    # Cleanup
    Remove-Item -Path $tempFile -Force
    
    Write-Success "Binary installed successfully!"
}

# Setup configuration
function Setup-Config {
    Write-Info "Setting up configuration..."
    
    # Create config directory
    New-Item -ItemType Directory -Force -Path $CONFIG_DIR | Out-Null
    
    # Download example config
    $configUrl = "https://raw.githubusercontent.com/$REPO/main/config.example.yaml"
    $configPath = "$CONFIG_DIR\config.yaml"
    
    try {
        Invoke-WebRequest -Uri $configUrl -OutFile $configPath
    } catch {
        Write-Warning "Could not download config, creating default..."
        
        $defaultConfig = @"
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
  token_file: "$($CONFIG_DIR -replace '\\', '\\')\agent_token"
"@
        Set-Content -Path $configPath -Value $defaultConfig
    }
    
    Write-Success "Configuration created at $configPath"
}

# Install Windows Service
function Install-Service {
    Write-Info "Installing Windows Service..."
    
    $serviceName = "WatchUpAgent"
    $displayName = "WatchUp Monitoring Agent"
    $description = "Collects system metrics and sends them to WatchUp platform"
    $binaryPath = "$INSTALL_DIR\$BINARY_NAME"
    
    # Check if service already exists
    $existingService = Get-Service -Name $serviceName -ErrorAction SilentlyContinue
    if ($existingService) {
        Write-Warning "Service already exists. Stopping and removing..."
        Stop-Service -Name $serviceName -Force -ErrorAction SilentlyContinue
        sc.exe delete $serviceName | Out-Null
        Start-Sleep -Seconds 2
    }
    
    # Create service using sc.exe (more reliable than New-Service for this use case)
    $createResult = sc.exe create $serviceName binPath= "`"$binaryPath`"" start= auto DisplayName= $displayName
    
    if ($LASTEXITCODE -eq 0) {
        # Set description
        sc.exe description $serviceName $description | Out-Null
        
        # Set recovery options (restart on failure)
        sc.exe failure $serviceName reset= 86400 actions= restart/60000/restart/60000/restart/60000 | Out-Null
        
        Write-Success "Windows Service installed!"
    } else {
        Write-Warning "Could not install as Windows Service. You can run the agent manually."
    }
}

# Add to PATH
function Add-ToPath {
    Write-Info "Adding to system PATH..."
    
    $currentPath = [Environment]::GetEnvironmentVariable("Path", "Machine")
    
    if ($currentPath -notlike "*$INSTALL_DIR*") {
        [Environment]::SetEnvironmentVariable(
            "Path",
            "$currentPath;$INSTALL_DIR",
            "Machine"
        )
        Write-Success "Added to system PATH"
    } else {
        Write-Info "Already in system PATH"
    }
}

# Print next steps
function Show-NextSteps {
    Write-Host ""
    Write-Success "╔════════════════════════════════════════════════════════════╗"
    Write-Success "║  WatchUp Agent Installation Complete! 🎉                  ║"
    Write-Success "╚════════════════════════════════════════════════════════════╝"
    Write-Host ""
    Write-Info "📋 Next Steps:"
    Write-Host ""
    Write-Host "1. " -NoNewline
    Write-Warning "Configure your server ID:"
    Write-Host "   notepad $CONFIG_DIR\config.yaml"
    Write-Host ""
    Write-Host "2. " -NoNewline
    Write-Warning "Start the agent:"
    Write-Host "   Start-Service WatchUpAgent"
    Write-Host ""
    Write-Host "3. " -NoNewline
    Write-Warning "Check status:"
    Write-Host "   Get-Service WatchUpAgent"
    Write-Host ""
    Write-Host "4. " -NoNewline
    Write-Warning "View logs:"
    Write-Host "   Get-EventLog -LogName Application -Source WatchUpAgent -Newest 50"
    Write-Host ""
    Write-Host "Or run manually:"
    Write-Host "   cd $CONFIG_DIR"
    Write-Host "   $INSTALL_DIR\$BINARY_NAME"
    Write-Host ""
    Write-Info "📚 Documentation:"
    Write-Host "   https://github.com/$REPO"
    Write-Host ""
    Write-Info "🔗 Link your agent:"
    Write-Host "   The agent will display a link and code on first run"
    Write-Host ""
}

# Main installation flow
function Main {
    Write-Host ""
    Write-Success "╔════════════════════════════════════════════════════════════╗"
    Write-Success "║         WatchUp Agent Installer (Windows)                 ║"
    Write-Success "║         https://watchup.site                              ║"
    Write-Success "╚════════════════════════════════════════════════════════════╝"
    Write-Host ""
    
    # Check administrator privileges
    if (-not (Test-Administrator)) {
        Write-Error "This script must be run as Administrator!"
        Write-Warning "Right-click PowerShell and select 'Run as Administrator'"
        exit 1
    }
    
    $arch = Get-Architecture
    Write-Info "Detected architecture: windows-$arch"
    
    $version = Get-LatestVersion
    Install-Binary -version $version -arch $arch
    Setup-Config
    Install-Service
    Add-ToPath
    Show-NextSteps
}

# Run main function
Main