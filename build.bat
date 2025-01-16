@echo off
setlocal

set VERSION=v1.0.0
set PLATFORMS=linux/amd64 linux/arm64

:: 创建构建目录
if not exist build mkdir build

:: 遍历平台进行构建
for %%p in (%PLATFORMS%) do (
    for /f "tokens=1,2 delims=/" %%a in ("%%p") do (
        echo Building for %%a/%%b...
        
        :: 设置环境变量并构建
        set GOOS=%%a
        set GOARCH=%%b
        go build -o build/hy2agent-%%a-%%b
        
        if errorlevel 1 (
            echo Failed to build for %%a/%%b
            exit /b 1
        ) else (
            echo Successfully built for %%a/%%b
        )
    )
)

:: 进入构建目录
cd build

:: 创建压缩包和校验和
for %%f in (*) do (
    if not "%%~xf"==".gz" if not "%%~xf"==".txt" (
        tar -czf "%%f.tar.gz" "%%f"
        sha256sum "%%f.tar.gz" >> checksums.txt
    )
)

:: 返回上级目录
cd ..

echo Build complete! Check the build directory for outputs.
dir /b build

endlocal 