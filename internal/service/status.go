package service

import (
    "github.com/shirou/gopsutil/v3/cpu"
    "github.com/shirou/gopsutil/v3/load"
    "github.com/shirou/gopsutil/v3/mem"
    "github.com/shirou/gopsutil/v3/disk"
    gopsnet "github.com/shirou/gopsutil/v3/net"
    "github.com/shirou/gopsutil/v3/host"
    "agentapi/internal/model"
    "time"
    "net"
)

type StatusService struct{}

func NewStatusService() *StatusService {
    return &StatusService{}
}

func (s *StatusService) GetSystemStatus() (*model.SystemStatus, error) {
    status := &model.SystemStatus{}
    
    // 获取CPU信息
    cpuInfo, err := s.getCPUInfo()
    if err != nil {
        return nil, err
    }
    status.CPU = cpuInfo
    
    // 获取内存信息
    memInfo, err := s.GetMemoryInfo()
    if err != nil {
        return nil, err
    }
    status.Memory = memInfo
    
    // 获取磁盘信息
    diskInfo, err := s.GetDiskInfo()
    if err != nil {
        return nil, err
    }
    status.Disk = diskInfo
    
    // 获取网络信息
    netInfo, err := s.GetNetworkInfo()
    if err != nil {
        return nil, err
    }
    status.Network = netInfo
    
    // 获取系统信息
    sysInfo, err := s.GetSystemInfo()
    if err != nil {
        return nil, err
    }
    status.System = sysInfo
    
    return status, nil
}

func (s *StatusService) getCPUInfo() (model.CPUInfo, error) {
    var info model.CPUInfo
    
    // 获取CPU使用率
    percent, err := cpu.Percent(time.Second, false)
    if err != nil {
        return info, err
    }
    info.Usage = percent[0]
    
    // 获取每个核心使用率
    perCPU, err := cpu.Percent(time.Second, true)
    if err != nil {
        return info, err
    }
    info.CoreUsages = perCPU
    
    // 获取负载信息
    loadAvg, err := load.Avg()
    if err != nil {
        return info, err
    }
    info.LoadAvg1 = loadAvg.Load1
    info.LoadAvg5 = loadAvg.Load5
    info.LoadAvg15 = loadAvg.Load15
    
    return info, nil
}

func (s *StatusService) GetMemoryInfo() (model.MemoryInfo, error) {
    var info model.MemoryInfo
    
    memStat, err := mem.VirtualMemory()
    if err != nil {
        return info, err
    }
    
    info.Total = memStat.Total
    info.Used = memStat.Used
    info.Free = memStat.Free
    info.Cache = memStat.Cached
    
    return info, nil
}

func (s *StatusService) GetDiskInfo() ([]model.DiskInfo, error) {
    var diskInfos []model.DiskInfo
    
    partitions, err := disk.Partitions(false)
    if err != nil {
        return nil, err
    }
    
    for _, partition := range partitions {
        usage, err := disk.Usage(partition.Mountpoint)
        if err != nil {
            continue
        }
        
        diskInfo := model.DiskInfo{
            Path:      partition.Mountpoint,
            Total:     usage.Total,
            Used:      usage.Used,
            Free:      usage.Free,
            UsageRate: usage.UsedPercent,
        }
        diskInfos = append(diskInfos, diskInfo)
    }
    
    return diskInfos, nil
}

func (s *StatusService) GetNetworkInfo() (model.NetworkInfo, error) {
    var info model.NetworkInfo
    
    // 获取网络IO统计
    netStats, err := gopsnet.IOCounters(false)
    if err != nil {
        return info, err
    }
    
    if len(netStats) > 0 {
        info.TotalUpload = netStats[0].BytesSent
        info.TotalDownload = netStats[0].BytesRecv
        
        // 计算速度需要两个时间点的数据
        time.Sleep(time.Second)
        newStats, err := gopsnet.IOCounters(false)
        if err == nil && len(newStats) > 0 {
            info.UploadSpeed = newStats[0].BytesSent - netStats[0].BytesSent
            info.DownloadSpeed = newStats[0].BytesRecv - netStats[0].BytesRecv
        }
    }
    
    return info, nil
}

func (s *StatusService) GetSystemInfo() (model.SystemInfo, error) {
    var info model.SystemInfo
    
    // 获取主机信息
    hostInfo, err := host.Info()
    if err != nil {
        return info, err
    }
    
    info.OS = hostInfo.Platform + " " + hostInfo.PlatformVersion
    info.Uptime = hostInfo.Uptime
    
    // 获取IP地址
    interfaces, err := net.Interfaces()
    if err != nil {
        return info, err
    }

    for _, iface := range interfaces {
        // 跳过回环接口和非活动接口
        if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
            continue
        }

        addrs, err := iface.Addrs()
        if err != nil {
            continue
        }

        for _, addr := range addrs {
            // 尝试将地址转换为IP
            ipNet, ok := addr.(*net.IPNet)
            if !ok {
                continue
            }
            ip := ipNet.IP
            if ip.IsLoopback() {
                continue
            }

            // 区分IPv4和IPv6
            if ip4 := ip.To4(); ip4 != nil {
                info.IPv4 = append(info.IPv4, ip4.String())
            } else if ip6 := ip.To16(); ip6 != nil {
                info.IPv6 = append(info.IPv6, ip6.String())
            }
        }
    }
    
    return info, nil
}

// 其他获取信息的方法实现... 