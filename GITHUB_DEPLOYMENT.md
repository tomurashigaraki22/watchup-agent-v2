# GitHub Deployment Guide

This guide explains how to deploy the WatchUp Agent using GitHub releases and automated builds.

## Table of Contents
- [Initial Setup](#initial-setup)
- [Creating Releases](#creating-releases)
- [Automated Builds](#automated-builds)
- [Distribution](#distribution)

---

## Initial Setup

### 1. Create GitHub Repository

```bash
# Initialize git repository
git init

# Add all files
git add .

# Create initial commit
git commit -m "Initial commit: WatchUp Agent v1.0.0"

# Add remote repository
git remote add origin https://github.com/YOUR_USERNAME/watchup-agent.git

# Push to GitHub
git push -u origin main
```

### 2. Configure GitHub Actions

The repository includes `.github/workflows/release.yml` which automatically:
- Builds binaries for all platforms
- Creates GitHub releases
- Uploads release assets

**No additional configuration needed!** The workflow triggers automatically on version tags.

---

## Creating Releases

### Method 1: Command Line (Recommended)

```bash
# 1. Ensure all changes are committed
git add .
git commit -m "Release v1.0.0"
git push

# 2. Create and push a version tag
git tag v1.0.0
git push origin v1.0.0

# 3. GitHub Actions will automatically:
#    - Build binaries for all platforms
#    - Create a release
#    - Upload all binaries
```

### Method 2: GitHub Web Interface

1. Go to your repository on GitHub
2. Click "Releases" → "Create a new release"
3. Click "Choose a tag" → Type `v1.0.0` → "Create new tag"
4. Fill in release title and description
5. Click "Publish release"
6. GitHub Actions will build and upload binaries

### Version Numbering

Follow [Semantic Versioning](https://semver.org/):
- **v1.0.0** - Major release (breaking changes)
- **v1.1.0** - Minor release (new features, backward compatible)
- **v1.0.1** - Patch release (bug fixes)

---

## Automated Builds

### Build Matrix

The GitHub Actions workflow builds for:

| Platform | Architecture | Output File |
|----------|--------------|-------------|
| Linux | amd64 | `watchup-agent-linux-amd64` |
| Linux | arm64 | `watchup-agent-linux-arm64` |
| Linux | armv7 | `watchup-agent-linux-armv7` |
| macOS | amd64 (Intel) | `watchup-agent-darwin-amd64` |
| macOS | arm64 (M1/M2) | `watchup-agent-darwin-arm64` |
| Windows | amd64 | `watchup-agent-windows-amd64.exe` |
| Windows | arm64 | `watchup-agent-windows-arm64.exe` |

### Build Process

1. **Trigger**: Push a tag starting with `v` (e.g., `v1.0.0`)
2. **Build**: GitHub Actions compiles binaries for all platforms
3. **Release**: Creates a GitHub release with all binaries
4. **Assets**: Binaries are uploaded as release assets

### Monitoring Builds

1. Go to your repository on GitHub
2. Click "Actions" tab
3. View the "Build and Release" workflow
4. Check build status and logs

---

## Distribution

### Installation Scripts

The repository includes installation scripts that automatically download the latest release:

**Linux/macOS** (`install.sh`):
```bash
curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/watchup-agent/main/install.sh | bash
```

**Windows** (`install.ps1`):
```powershell
iwr -useb https://raw.githubusercontent.com/YOUR_USERNAME/watchup-agent/main/install.ps1 | iex
```

### Customizing Installation Scripts

Update the `REPO` variable in both scripts:

**install.sh**:
```bash
REPO="YOUR_USERNAME/watchup-agent"
```

**install.ps1**:
```powershell
$REPO = "YOUR_USERNAME/watchup-agent"
```

### Direct Download Links

Users can download binaries directly:

```
https://github.com/YOUR_USERNAME/watchup-agent/releases/download/v1.0.0/watchup-agent-linux-amd64
https://github.com/YOUR_USERNAME/watchup-agent/releases/download/v1.0.0/watchup-agent-darwin-amd64
https://github.com/YOUR_USERNAME/watchup-agent/releases/download/v1.0.0/watchup-agent-windows-amd64.exe
```

### Latest Release API

Get the latest version programmatically:

```bash
# Get latest version tag
curl -s https://api.github.com/repos/YOUR_USERNAME/watchup-agent/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'

# Download latest Linux binary
LATEST=$(curl -s https://api.github.com/repos/YOUR_USERNAME/watchup-agent/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')
curl -fsSL "https://github.com/YOUR_USERNAME/watchup-agent/releases/download/${LATEST}/watchup-agent-linux-amd64" -o watchup-agent
```

---

## Release Checklist

Before creating a release:

- [ ] All tests pass (`go test ./...`)
- [ ] Version number updated in code (if applicable)
- [ ] CHANGELOG.md updated with changes
- [ ] Documentation updated
- [ ] All changes committed and pushed
- [ ] Tag created with correct version number
- [ ] GitHub Actions build succeeds
- [ ] Release notes written
- [ ] Installation scripts tested

---

## Troubleshooting

### Build Fails

**Check GitHub Actions logs**:
1. Go to "Actions" tab
2. Click on the failed workflow
3. Review error messages

**Common issues**:
- Missing dependencies: Run `go mod tidy`
- Syntax errors: Run `go build` locally
- Test failures: Run `go test ./...` locally

### Release Not Created

**Verify**:
- Tag starts with `v` (e.g., `v1.0.0`)
- Tag was pushed to GitHub: `git push origin v1.0.0`
- GitHub Actions has permissions to create releases

**Fix permissions**:
1. Go to repository Settings
2. Actions → General
3. Workflow permissions → "Read and write permissions"

### Binaries Not Uploaded

**Check**:
- Build completed successfully
- All platform builds succeeded
- Release was created (not just a tag)

**Manual upload**:
1. Build locally: `./build.sh v1.0.0`
2. Go to GitHub release
3. Edit release
4. Upload binaries from `dist/` folder

---

## Advanced Configuration

### Custom Build Flags

Edit `.github/workflows/release.yml`:

```yaml
- name: Build binary
  env:
    GOOS: ${{ matrix.os }}
    GOARCH: ${{ matrix.arch }}
    CGO_ENABLED: 0
  run: |
    go build \
      -ldflags="-s -w -X main.Version=${{ github.ref_name }}" \
      -o watchup-agent-${{ matrix.os }}-${{ matrix.arch }} \
      cmd/agent/main.go cmd/agent/setup.go
```

### Pre-release Builds

Create pre-release versions:

```bash
# Tag as pre-release
git tag v1.0.0-beta.1
git push origin v1.0.0-beta.1

# Mark as pre-release in GitHub
# Edit release → Check "This is a pre-release"
```

### Build Notifications

Add Slack/Discord notifications to GitHub Actions:

```yaml
- name: Notify on success
  if: success()
  uses: 8398a7/action-slack@v3
  with:
    status: ${{ job.status }}
    text: 'Release ${{ github.ref_name }} built successfully!'
    webhook_url: ${{ secrets.SLACK_WEBHOOK }}
```

---

## Continuous Deployment

### Automatic Updates

Users can set up automatic updates:

```bash
#!/bin/bash
# update-agent.sh

CURRENT_VERSION=$(watchup-agent --version 2>/dev/null | grep -oP 'v\d+\.\d+\.\d+')
LATEST_VERSION=$(curl -s https://api.github.com/repos/YOUR_USERNAME/watchup-agent/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ "$CURRENT_VERSION" != "$LATEST_VERSION" ]; then
    echo "Updating from $CURRENT_VERSION to $LATEST_VERSION"
    curl -fsSL https://raw.githubusercontent.com/YOUR_USERNAME/watchup-agent/main/install.sh | bash
    systemctl restart watchup-agent
fi
```

### Version Pinning

Pin to specific version in installation:

```bash
# Install specific version
VERSION="v1.0.0"
curl -fsSL "https://github.com/YOUR_USERNAME/watchup-agent/releases/download/${VERSION}/watchup-agent-linux-amd64" -o watchup-agent
```

---

## Best Practices

1. **Semantic Versioning**: Follow semver for version numbers
2. **Release Notes**: Write clear, detailed release notes
3. **Testing**: Test releases before marking as latest
4. **Changelog**: Maintain CHANGELOG.md with all changes
5. **Security**: Sign releases with GPG (optional but recommended)
6. **Documentation**: Update docs with each release
7. **Deprecation**: Announce breaking changes in advance

---

## Example Release Workflow

```bash
# 1. Finish development
git checkout main
git pull

# 2. Update version and changelog
echo "v1.1.0" > VERSION
nano CHANGELOG.md

# 3. Commit changes
git add .
git commit -m "Prepare release v1.1.0"
git push

# 4. Create and push tag
git tag -a v1.1.0 -m "Release v1.1.0: Add network latency monitoring"
git push origin v1.1.0

# 5. Wait for GitHub Actions to complete
# 6. Edit release notes on GitHub
# 7. Announce release to users
```

---

## Support

For deployment issues:
- **GitHub Actions**: Check workflow logs
- **Issues**: https://github.com/YOUR_USERNAME/watchup-agent/issues
- **Documentation**: README.md and DEPLOYMENT.md