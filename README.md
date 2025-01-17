# Hysteria2 Agent

一个用于管理 Hysteria2 服务的 RESTful API，用于hy2web管理面板（待开发）。

## 功能特点

- 完整的 Hysteria2 服务管理
  - 安装/卸载/更新
  - 启动/停止/重启
  - 配置管理和备份
  - 日志查询
  - 状态监控
- 系统管理功能
  - CPU/内存/磁盘监控
  - 网络状态监控
  - 系统信息查询
- 安全特性
  - API Key 认证
  - 访问控制
    - IP 白名单：支持 IPv4 和 IPv6
  - 配置自动备份

## 快速安装

### 一键安装脚本
```bash
# 使用 curl
bash -c "$(curl -fsSL https://raw.githubusercontent.com/zsancc/hy2agent/main/install.sh)"

# 或使用 wget
bash -c "$(wget -qO- https://raw.githubusercontent.com/zsancc/hy2agent/main/install.sh)"
```

### 一键卸载
```bash
# 使用 curl
bash -c "$(curl -fsSL https://raw.githubusercontent.com/zsancc/hy2agent/main/uninstall.sh)"

# 或使用 wget
bash -c "$(wget -qO- https://raw.githubusercontent.com/zsancc/hy2agent/main/uninstall.sh)"
```

## 从源码编译安装

1. 下载项目
```bash
git clone https://github.com/zsancc/hy2agent.git
cd hy2agent
```

2. 编译
```bash
go build -o hy2agent
```

3. 安装服务
```bash
sudo mv hy2agent /usr/local/bin/
sudo mkdir -p /etc/hy2agent
```

4. 创建服务文件
```bash
sudo nano /etc/systemd/system/hy2agent.service
```

添加以下内容：
```ini
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
```

5. 启动服务
```bash
sudo systemctl daemon-reload
sudo systemctl enable hy2agent
sudo systemctl start hy2agent
```

## 配置

配置文件位置：`/etc/hy2agent/config.json`

```json
{
    "api_key": "your-api-key",
    "ip_whitelist": ["192.168.1.100"]
}
```

首次启动时会自动生成随机 API Key。

## 使用示例

1. 获取服务状态
```bash
curl -H "X-API-Key: your-api-key" http://localhost:8080/api/v1/hysteria/status
```

2. 更新配置
```bash
curl -H "X-API-Key: your-api-key" -X PUT \
     -H "Content-Type: application/json" \
     -d '{"config":"your-config-here"}' \
     http://localhost:8080/api/v1/hysteria/config
```

3. 查看日志
```bash
curl -H "X-API-Key: your-api-key" \
     "http://localhost:8080/api/v1/hysteria/logs?lines=100&level=error"
```

更多 API 详情请参考 [API 文档](API%20doc.md)。

## 安全建议

1. 使用强密码作为 API Key
2. 访问控制
    - IP 白名单：支持 IPv4 和 IPv6
3. HTTPS 配置：
   - 方式一：安装时选择自动配置 HTTPS（使用 acme.sh）
     - 自动申请 Let's Encrypt 证书
     - 配置每日自动续签
     - 证书续签后自动重启服务
   - 方式二：使用 Nginx/Caddy 等反向代理
4. 定期备份配置文件

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
