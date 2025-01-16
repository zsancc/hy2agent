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

# 检查是否已安装
if [ -f "/usr/local/bin/hy2agent" ]; then
    echo -e "${YELLOW}检测到已安装 hy2agent${NC}"
    echo -e "请选择操作："
    echo -e "1. 重新安装"
    echo -e "2. 更新"
    echo -e "3. 退出"
    read -p "请输入选项 (1-3): " choice
    
    case $choice in
        1)
            echo -e "${YELLOW}开始重新安装...${NC}"
            # 停止服务
            systemctl stop hy2agent
            systemctl disable hy2agent
            ;;
        2)
            echo -e "${YELLOW}开始更新...${NC}"
            # 备份配置
            if [ -f "/etc/hy2agent/config.json" ]; then
                cp /etc/hy2agent/config.json /etc/hy2agent/config.json.bak
                echo -e "${GREEN}已备份配置文件到 /etc/hy2agent/config.json.bak${NC}"
            fi
            # 停止服务
            systemctl stop hy2agent
            ;;
        3)
            echo -e "${YELLOW}退出安装${NC}"
            exit 0
            ;;
        *)
            echo -e "${RED}无效选项${NC}"
            exit 1
            ;;
    esac
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

# 配置端口和白名单
read -p "请输入服务端口 (默认: 8080): " PORT
PORT=${PORT:-8080}

# 配置 HTTPS
echo -e "\n${YELLOW}配置 HTTPS...${NC}"

# 安装 acme.sh
echo -e "${YELLOW}安装 acme.sh...${NC}"
curl https://get.acme.sh | sh

# 设置自动更新
~/.acme.sh/acme.sh --upgrade --auto-upgrade

# 配置域名
read -p "请输入域名: " DOMAIN
if [ -z "$DOMAIN" ]; then
    echo -e "${RED}域名不能为空${NC}"
    exit 1
fi

# 配置邮箱
read -p "请输入邮箱 (用于证书申请): " EMAIL
if [ -z "$EMAIL" ]; then
    echo -e "${RED}邮箱不能为空${NC}"
    exit 1
fi

# 设置默认CA
~/.acme.sh/acme.sh --set-default-ca --server letsencrypt

# 申请证书
echo -e "${YELLOW}开始申请证书...${NC}"
if ! apply_cert $DOMAIN; then
    echo -e "${RED}证书申请失败${NC}"
    echo -e "${YELLOW}请检查：${NC}"
    echo "1. 域名是否正确解析到本机IP"
    echo "2. 是否有足够的权限"
    exit 1
fi

# 安装证书
mkdir -p /etc/hy2agent/cert
~/.acme.sh/acme.sh --install-cert -d $DOMAIN --ecc \
    --key-file /etc/hy2agent/cert/private.key \
    --fullchain-file /etc/hy2agent/cert/cert.pem
    
# 配置自动续签服务
cat > /etc/systemd/system/acme-renew.service << EOF
[Unit]
Description=Acme certificate renewal service
After=network.target

[Service]
Type=oneshot
ExecStart=/root/.acme.sh/acme.sh --cron --home /root/.acme.sh
ExecStartPost=/bin/systemctl restart hy2agent
EOF

# 配置自动续签定时器
cat > /etc/systemd/system/acme-renew.timer << EOF
[Unit]
Description=Acme certificate renewal timer

[Timer]
OnCalendar=*-*-* 00:00:00
RandomizedDelaySec=3600
Persistent=true

[Install]
WantedBy=timers.target
EOF

# 启用自动续签服务
systemctl daemon-reload
systemctl enable acme-renew.timer
systemctl start acme-renew.timer

# 配置服务启动参数
EXTRA_ARGS="-port $PORT -tls -cert /etc/hy2agent/cert/cert.pem -key /etc/hy2agent/cert/private.key"

echo -e "${GREEN}已配置证书自动续签，每天凌晨随机时间检查并续签${NC}"

# 配置访问控制
echo -e "\n${YELLOW}访问控制配置：${NC}"

# 获取本机公网IP
echo -e "${YELLOW}正在获取本机公网IP...${NC}"
IPV4=$(curl -s4 4.ipw.cn)
IPV6=$(curl -s6 6.ipw.cn)

if [ ! -z "$IPV4" ]; then
    echo -e "${GREEN}检测到IPv4地址: ${YELLOW}$IPV4${NC}"
fi

if [ ! -z "$IPV6" ]; then
    echo -e "${GREEN}检测到IPv6地址: ${YELLOW}$IPV6${NC}"
fi

echo "请输入允许访问的IP地址（多个用空格分隔）"
echo -e "${YELLOW}提示: 本机公网IP已自动添加到白名单${NC}"
read -p "IP地址 (例如: 192.168.1.100 10.0.0.1): " WEB_IPS

# 处理输入为数组
IPS=(${WEB_IPS})

# 构建JSON数组
IP_LIST="["

# 添加本机IP
[ ! -z "$IPV4" ] && IP_LIST="$IP_LIST\"$IPV4\","
[ ! -z "$IPV6" ] && IP_LIST="$IP_LIST\"$IPV6\","

# 添加本地回环地址
IP_LIST="$IP_LIST\"127.0.0.1\","

# 添加用户输入的IP
for ip in "${IPS[@]}"; do
    IP_LIST="$IP_LIST\"$ip\","
done

# 移除最后一个逗号并关闭数组
IP_LIST="${IP_LIST%,}]"

# 创建初始配置文件
cat > /etc/hy2agent/config.json << EOF
{
    "api_key": "$(openssl rand -hex 32)",
    "ip_whitelist": $IP_LIST
}
EOF

# 检查防火墙
if command -v ufw >/dev/null 2>&1; then
    echo -e "${YELLOW}检测到 UFW 防火墙${NC}"
    echo -e "正在添加防火墙规则..."
    ufw allow $PORT/tcp
    echo -e "${GREEN}已添加 UFW 规则: $PORT/tcp${NC}"
elif command -v firewall-cmd >/dev/null 2>&1; then
    echo -e "${YELLOW}检测到 FirewallD 防火墙${NC}"
    echo -e "正在添加防火墙规则..."
    firewall-cmd --permanent --add-port=$PORT/tcp
    firewall-cmd --reload
    echo -e "${GREEN}已添加 FirewallD 规则: $PORT/tcp${NC}"
else
    echo -e "${YELLOW}未检测到支持的防火墙，请手动放通 $PORT 端口${NC}"
fi

# 下载对应版本的二进制文件
VERSION="v1.0.0"
DOWNLOAD_URL="https://github.com/zsancc/hy2agent/releases/download/${VERSION}/hy2agent-${OS}-${ARCH}"

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
ExecStart=/usr/local/bin/hy2agent $EXTRA_ARGS
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
    echo -e "${GREEN}IP 白名单: ${YELLOW}$WEB_IPS${NC}"
    echo -e "${GREEN}服务端口: ${YELLOW}$PORT${NC}"
    
    # 如果是更新，恢复配置
    if [ "$choice" = "2" ] && [ -f "/etc/hy2agent/config.json.bak" ]; then
        mv /etc/hy2agent/config.json.bak /etc/hy2agent/config.json
        echo -e "${GREEN}已恢复配置文件${NC}"
        systemctl restart hy2agent
    fi
else
    echo -e "${RED}hy2agent 安装失败，请检查日志${NC}"
    journalctl -u hy2agent --no-pager -n 50
    exit 1
fi

# 显示使用说明
echo -e "\n${GREEN}使用说明：${NC}"
echo -e "1. API 文档：https://$DOMAIN/docs"
echo -e "\n${YELLOW}证书信息：${NC}"
echo -e "- 证书位置：/etc/hy2agent/cert/"
echo -e "- 自动续签：每天凌晨随机时间"
echo -e "- 查看续签状态：systemctl status acme-renew.timer"
echo -e "- 手动续签：~/.acme.sh/acme.sh --cron --home ~/.acme.sh"
echo -e "2. 配置文件位置：/etc/hy2agent/config.json"
echo -e "3. 防火墙端口：$PORT/tcp 需要放通"
echo -e "3. 服务控制："
echo -e "   - 启动：systemctl start hy2agent"
echo -e "   - 停止：systemctl stop hy2agent"
echo -e "   - 重启：systemctl restart hy2agent"
echo -e "   - 状态：systemctl status hy2agent"
echo -e "4. 日志查看：journalctl -u hy2agent -f" 