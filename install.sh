#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 检查是否为 root 用户
if [ "$EUID" -ne 0 ]; then 
    echo -e "${RED}请使用 root 权限运行此脚本${NC}"
    exit 1
fi

# 检查系统类型和架构
OS=$(uname -s)
ARCH=$(uname -m)
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64)
        ARCH="arm64"
        ;;
    armv7l)
        ARCH="arm"
        ;;
    *)
        echo -e "${RED}不支持的系统架构: $ARCH${NC}"
        exit 1
        ;;
esac

# 检查包管理器
if command -v apt >/dev/null 2>&1; then
    PKG_MANAGER="apt"
    PKG_UPDATE="apt update"
    PKG_INSTALL="apt install -y"
elif command -v yum >/dev/null 2>&1; then
    PKG_MANAGER="yum"
    PKG_UPDATE="yum makecache"
    PKG_INSTALL="yum install -y"
else
    echo -e "${RED}不支持的包管理器${NC}"
    exit 1
fi

# 安装依赖
echo -e "${YELLOW}正在安装依赖...${NC}"
$PKG_UPDATE
$PKG_INSTALL curl wget systemd

# 创建目录
mkdir -p /etc/hy2agent
mkdir -p /usr/local/bin

# 下载对应版本的二进制文件
VERSION="v1.0.0"
DOWNLOAD_URL="https://github.com/yourusername/hy2agent/releases/download/${VERSION}/hy2agent-${OS}-${ARCH}"

echo -e "${YELLOW}正在下载 hy2agent...${NC}"
wget -O /usr/local/bin/hy2agent $DOWNLOAD_URL || {
    echo -e "${RED}下载失败${NC}"
    exit 1
}

# 设置执行权限
chmod +x /usr/local/bin/hy2agent

# 创建 systemd 服务文件
cat > /etc/systemd/system/hy2agent.service << EOF
[Unit]
Description=Hysteria2 Agent Service
After=network.target

[Service]
Type=simple
ExecStart=/usr/local/bin/hy2agent
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
EOF

# 重载 systemd 并启动服务
systemctl daemon-reload
systemctl enable hy2agent
systemctl start hy2agent

# 等待服务启动
sleep 2

# 检查服务状态
if systemctl is-active --quiet hy2agent; then
    echo -e "${GREEN}hy2agent 安装成功！${NC}"
    # 显示 API Key
    API_KEY=$(grep -o '"api_key": "[^"]*' /etc/hy2agent/config.json | cut -d'"' -f4)
    echo -e "${GREEN}API Key: ${YELLOW}$API_KEY${NC}"
    echo -e "${GREEN}请妥善保管 API Key${NC}"
else
    echo -e "${RED}hy2agent 安装失败，请检查日志${NC}"
    journalctl -u hy2agent --no-pager -n 50
    exit 1
fi

# 显示使用说明
echo -e "\n${GREEN}使用说明：${NC}"
echo -e "1. API 文档：http://your-server:8080/docs"
echo -e "2. 配置文件位置：/etc/hy2agent/config.json"
echo -e "3. 服务控制："
echo -e "   - 启动：systemctl start hy2agent"
echo -e "   - 停止：systemctl stop hy2agent"
echo -e "   - 重启：systemctl restart hy2agent"
echo -e "   - 状态：systemctl status hy2agent"
echo -e "4. 日志查看：journalctl -u hy2agent -f" 