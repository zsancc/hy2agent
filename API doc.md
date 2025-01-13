# Hysteria2 Agent API 文档

## 基本信息
- 基础URL: `http://your-server:8080`
- 所有请求需要包含 API Key：
  ```
  Header: X-API-Key: your-api-key
  ```
- 所有响应均为 JSON 格式
- 错误响应格式：
  ```json
  {
      "error": "错误信息描述"
  }
  ```

## IP 访问控制
- 支持 IP 白名单和黑名单
- 白名单优先级高于黑名单
- 配置文件位置：`/etc/hy2agent/config.json`

## API 端点

### 系统状态

#### 获取系统状态
```http
GET /api/v1/status

Response 200:
{
    "cpu_usage": 25.5,
    "memory_usage": 60.2,
    "disk_usage": 45.8,
    "uptime": "5d 12h 30m",
    "load_average": {
        "1m": 0.5,
        "5m": 0.7,
        "15m": 0.6
    }
}
```

### 系统管理

#### 获取内存信息
```http
GET /api/v1/system/memory

Response 200:
{
    "total": 8589934592,
    "used": 4294967296,
    "free": 4294967296,
    "available": 6442450944,
    "usage_percent": 50,
    "swap_total": 2147483648,
    "swap_used": 1073741824,
    "swap_free": 1073741824
}
```

#### 获取磁盘信息
```http
GET /api/v1/system/disk

Response 200:
{
    "total": 1099511627776,
    "used": 549755813888,
    "free": 549755813888,
    "usage_percent": 50,
    "mount_point": "/",
    "fs_type": "ext4",
    "inodes_total": 65536000,
    "inodes_used": 32768000
}
```

#### 获取网络信息
```http
GET /api/v1/system/network

Response 200:
{
    "interfaces": [
        {
            "name": "eth0",
            "bytes_sent": 1000000,
            "bytes_recv": 2000000,
            "packets_sent": 10000,
            "packets_recv": 20000,
            "errors_in": 0,
            "errors_out": 0,
            "speed": "1000Mb/s"
        }
    ],
    "connections": {
        "tcp": 500,
        "udp": 200
    }
}
```

#### 获取系统信息
```http
GET /api/v1/system/info

Response 200:
{
    "hostname": "server1",
    "os": "Ubuntu 20.04",
    "kernel": "5.4.0-42-generic",
    "arch": "x86_64",
    "cpu_cores": 4,
    "cpu_model": "Intel(R) Xeon(R) CPU E5-2680 v3 @ 2.50GHz",
    "timezone": "UTC",
    "boot_time": "2024-01-12T00:00:00Z"
}
```

#### 系统控制
```http
POST /api/v1/system/reboot

Response 200:
{
    "message": "System is rebooting",
    "scheduled_time": "2024-01-12T12:00:00Z"
}

POST /api/v1/system/shutdown

Response 200:
{
    "message": "System is shutting down",
    "scheduled_time": "2024-01-12T12:00:00Z"
}
```

### Hysteria2 管理

#### 获取服务状态
```http
GET /api/v1/hysteria/status

Response 200:
{
    "is_installed": true,
    "is_running": true,
    "version": "v2.6.0",
    "build_date": "2024-01-12T00:10:09Z",
    "build_type": "release",
    "platform": "linux",
    "architecture": "amd64",
    "service_status": "running",
    "last_error": "",
    "load_state": "loaded",
    "active_state": "active",
    "memory_usage": 15360000,
    "cpu_usage": 0.5,
    "uptime": "2d 5h 30m"
}
```

#### 配置管理
```http
GET /api/v1/hysteria/config

Response 200:
{
    "config": "listen: :443\nauth:\n  type: password\n  password: your_password\nmasquerade:\n  type: proxy\n  proxy:\n    url: https://www.microsoft.com\n    rewriteHost: true"
}

PUT /api/v1/hysteria/config
Request:
{
    "config": "listen: :443\nauth:\n  type: password\n  password: new_password\nmasquerade:\n  type: proxy\n  proxy:\n    url: https://www.microsoft.com\n    rewriteHost: true"
}

Response 200:
{
    "message": "Config updated successfully",
    "backup_file": "config.yaml.bak.20240112150405"
}
```

#### 日志查询
```http
GET /api/v1/hysteria/logs?lines=100&since=5m&level=error

Response 200:
{
    "logs": "2024-01-12 12:00:00 [ERROR] Failed to bind port: address already in use\n2024-01-12 12:00:05 [ERROR] Authentication failed from 1.2.3.4:12345",
    "total_lines": 2,
    "filter_info": {
        "lines": 100,
        "since": "5m",
        "level": "error"
    }
}
```

#### 服务控制
```http
POST /api/v1/hysteria/install

Response 200:
{
    "message": "Hysteria2 installed successfully",
    "version": "v2.6.0",
    "output": "Installation completed..."
}

POST /api/v1/hysteria/uninstall

Response 200:
{
    "message": "Hysteria2 uninstalled successfully",
    "output": "Uninstallation completed..."
}

POST /api/v1/hysteria/update

Response 200:
{
    "message": "Hysteria2 updated successfully",
    "old_version": "v2.5.0",
    "new_version": "v2.6.0",
    "output": "Update completed..."
}

POST /api/v1/hysteria/start

Response 200:
{
    "message": "Hysteria2 service started",
    "status": "running"
}

POST /api/v1/hysteria/stop

Response 200:
{
    "message": "Hysteria2 service stopped",
    "status": "stopped"
}

POST /api/v1/hysteria/restart

Response 200:
{
    "message": "Hysteria2 service restarted",
    "status": "running"
}
```

#### 配置备份
```http
GET /api/v1/hysteria/config/backups

Response 200:
{
    "backups": [
        {
            "filename": "config.yaml.bak.20240112150405",
            "size": 1024,
            "created_at": "2024-01-12T15:04:05Z"
        },
        {
            "filename": "config.yaml.bak.20240112140305",
            "size": 1024,
            "created_at": "2024-01-12T14:03:05Z"
        }
    ]
}

POST /api/v1/hysteria/config/restore
Request:
{
    "backup": "config.yaml.bak.20240112150405"
}

Response 200:
{
    "message": "Config restored successfully",
    "restored_from": "config.yaml.bak.20240112150405",
    "service_restarted": true
}
```

### 访问控制管理

#### IP 白名单
```http
GET /api/v1/config/whitelist

Response 200:
{
    "whitelist": [
        {
            "ip": "192.168.1.100",
            "added_at": "2024-01-12T00:00:00Z",
            "comment": "Office IP"
        },
        {
            "ip": "10.0.0.1",
            "added_at": "2024-01-12T00:00:00Z",
            "comment": "Home IP"
        }
    ]
}

PUT /api/v1/config/whitelist
Request:
{
    "ips": ["192.168.1.100", "10.0.0.1"],
    "comments": {
        "192.168.1.100": "Office IP",
        "10.0.0.1": "Home IP"
    }
}

Response 200:
{
    "message": "Whitelist updated",
    "updated_at": "2024-01-12T12:00:00Z"
}
```

#### IP 黑名单
```http
GET /api/v1/config/blacklist

Response 200:
{
    "blacklist": [
        {
            "ip": "1.2.3.4",
            "added_at": "2024-01-12T00:00:00Z",
            "reason": "Suspicious activity"
        },
        {
            "ip": "5.6.7.8",
            "added_at": "2024-01-12T00:00:00Z",
            "reason": "Multiple failed attempts"
        }
    ]
}

PUT /api/v1/config/blacklist
Request:
{
    "ips": ["1.2.3.4", "5.6.7.8"],
    "reasons": {
        "1.2.3.4": "Suspicious activity",
        "5.6.7.8": "Multiple failed attempts"
    }
}

Response 200:
{
    "message": "Blacklist updated",
    "updated_at": "2024-01-12T12:00:00Z"
}
```

## 错误码说明

- 200: 请求成功
- 400: 请求参数错误
- 401: 认证失败（API Key 无效）
- 403: 访问被拒绝（IP 不在白名单或在黑名单中）
- 404: 资源不存在
- 500: 服务器内部错误

## 注意事项

1. 所有时间相关的字段均使用 ISO 8601 格式
2. 文件大小单位为字节
3. CPU 和内存使用率为百分比（0-100）
4. 配置文件使用 YAML 格式
5. 日志级别包括：error, warning, info, debug
