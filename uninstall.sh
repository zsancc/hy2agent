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

echo -e "${YELLOW}开始卸载 hy2agent...${NC}"

# 停止服务
echo -e "${YELLOW}停止 hy2agent 服务...${NC}"
systemctl stop hy2agent
systemctl disable hy2agent

# 删除服务文件
echo -e "${YELLOW}删除服务文件...${NC}"
rm -f /etc/systemd/system/hy2agent.service
systemctl daemon-reload

# 删除程序文件
echo -e "${YELLOW}删除程序文件...${NC}"
rm -f /usr/local/bin/hy2agent

# 询问是否删除配置文件
read -p "是否删除配置文件？(y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo -e "${YELLOW}删除配置文件...${NC}"
    rm -rf /etc/hy2agent
    echo -e "${GREEN}配置文件已删除${NC}"
else
    echo -e "${YELLOW}保留配置文件在 /etc/hy2agent${NC}"
fi

# 检查是否还有残留进程
if pgrep hy2agent > /dev/null; then
    echo -e "${YELLOW}关闭残留进程...${NC}"
    pkill hy2agent
fi

echo -e "${GREEN}hy2agent 卸载完成！${NC}"

# 显示清理建议
echo -e "\n${YELLOW}建议清理：${NC}"
echo "1. 检查 journalctl 日志："
echo "   journalctl --vacuum-time=1d"
echo "2. 如果不再需要，可以删除安装目录："
echo "   rm -rf /etc/hy2agent" 