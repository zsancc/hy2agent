package model

type SystemStatus struct {
    CPU     CPUInfo     `json:"cpu"`
    Memory  MemoryInfo  `json:"memory"`
    Disk    []DiskInfo  `json:"disk"`
    Network NetworkInfo `json:"network"`
    System  SystemInfo  `json:"system"`
}

type CPUInfo struct {
    Usage       float64   `json:"usage"`        // CPU总使用率
    CoreUsages  []float64 `json:"coreUsages"`   // 每个核心的使用率
    LoadAvg1    float64   `json:"loadAvg1"`     // 1分钟负载
    LoadAvg5    float64   `json:"loadAvg5"`     // 5分钟负载
    LoadAvg15   float64   `json:"loadAvg15"`    // 15分钟负载
}

type MemoryInfo struct {
    Total       uint64  `json:"total"`      // 总内存
    Used        uint64  `json:"used"`       // 已用内存
    Free        uint64  `json:"free"`       // 可用内存
    Cache       uint64  `json:"cache"`      // 缓存内存
}

type DiskInfo struct {
    Path        string  `json:"path"`       // 分区路径
    Total       uint64  `json:"total"`      // 总空间
    Used        uint64  `json:"used"`       // 已用空间
    Free        uint64  `json:"free"`       // 可用空间
    UsageRate   float64 `json:"usageRate"`  // 使用率
}

type NetworkInfo struct {
    UploadSpeed   uint64 `json:"uploadSpeed"`   // 实时上传速度
    DownloadSpeed uint64 `json:"downloadSpeed"` // 实时下载速度
    TotalUpload   uint64 `json:"totalUpload"`   // 总上传流量
    TotalDownload uint64 `json:"totalDownload"` // 总下载流量
}

type SystemInfo struct {
    IPv4        []string `json:"ipv4"`         // IPv4地址列表
    IPv6        []string `json:"ipv6"`         // IPv6地址列表
    OS          string   `json:"os"`           // 操作系统类型和版本
    Uptime      uint64   `json:"uptime"`       // 系统运行时间(秒)
    OnlineTime  uint64   `json:"onlineTime"`   // 在线时间(秒)
} 