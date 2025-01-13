#!/bin/bash

VERSION="v1.0.0"
PLATFORMS=("linux/amd64" "linux/arm64" "linux/arm")

for PLATFORM in "${PLATFORMS[@]}"; do
    OS=${PLATFORM%/*}
    ARCH=${PLATFORM#*/}
    
    echo "Building for $OS/$ARCH..."
    GOOS=$OS GOARCH=$ARCH go build -o "build/hy2agent-$OS-$ARCH"
    
    if [ $? -eq 0 ]; then
        echo "Successfully built for $OS/$ARCH"
    else
        echo "Failed to build for $OS/$ARCH"
        exit 1
    fi
done

# 创建发布包
cd build
for FILE in *; do
    tar czf "$FILE.tar.gz" "$FILE"
    sha256sum "$FILE.tar.gz" >> checksums.txt
done 