#!/bin/bash
# Build script for WatchUp Agent
# Builds binaries for all supported platforms

set -e

VERSION=${1:-"dev"}
OUTPUT_DIR="dist"

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Building WatchUp Agent v${VERSION}${NC}"
echo ""

# Create output directory
mkdir -p "$OUTPUT_DIR"

# Build matrix
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "linux/arm/7"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
    "windows/arm64"
)

for PLATFORM in "${PLATFORMS[@]}"; do
    IFS='/' read -r -a PARTS <<< "$PLATFORM"
    GOOS="${PARTS[0]}"
    GOARCH="${PARTS[1]}"
    GOARM="${PARTS[2]:-}"
    
    OUTPUT_NAME="watchup-agent-${GOOS}-${GOARCH}"
    if [ -n "$GOARM" ]; then
        OUTPUT_NAME="watchup-agent-${GOOS}-armv${GOARM}"
    fi
    
    if [ "$GOOS" = "windows" ]; then
        OUTPUT_NAME="${OUTPUT_NAME}.exe"
    fi
    
    echo -e "${BLUE}Building for ${GOOS}/${GOARCH}${GOARM:+v$GOARM}...${NC}"
    
    env GOOS="$GOOS" GOARCH="$GOARCH" GOARM="$GOARM" CGO_ENABLED=0 \
        go build -ldflags="-s -w -X main.Version=${VERSION}" \
        -o "${OUTPUT_DIR}/${OUTPUT_NAME}" \
        cmd/agent/main.go cmd/agent/setup.go
    
    echo -e "${GREEN}✓ Built ${OUTPUT_NAME}${NC}"
done

echo ""
echo -e "${GREEN}Build complete! Binaries are in ${OUTPUT_DIR}/${NC}"
echo ""
echo "To create a release:"
echo "  git tag v${VERSION}"
echo "  git push origin v${VERSION}"