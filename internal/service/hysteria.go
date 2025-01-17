package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Hysteria2Service struct{}

type Hysteria2Status struct {
	IsInstalled   bool   `json:"is_installed"`
	IsRunning     bool   `json:"is_running"`
	Version       string `json:"version"`
	BuildDate     string `json:"build_date,omitempty"`
	BuildType     string `json:"build_type,omitempty"`
	Platform      string `json:"platform,omitempty"`
	Architecture  string `json:"architecture,omitempty"`
	ServiceStatus string `json:"service_status,omitempty"` // 服务状态描述
	LastError     string `json:"last_error,omitempty"`     // 最后一次错误信息
	LoadState     string `json:"load_state,omitempty"`     // 加载状态
	ActiveState   string `json:"active_state,omitempty"`   // 活动状态
}

// 获取日志的选项
type LogOptions struct {
	Lines int    `json:"lines,omitempty"` // 返回的最大行数
	Since string `json:"since,omitempty"` // 从什么时间开始，如"5m", "2h"
	Level string `json:"level,omitempty"` // 日志级别过滤：info, error等
}

// 健康检查结果
type HealthCheck struct {
	IsRunning   bool   `json:"is_running"`
	PortOpen    bool   `json:"port_open"`
	ConfigValid bool   `json:"config_valid"`
	LastError   string `json:"last_error,omitempty"`
	CheckTime   string `json:"check_time"`
}

func NewHysteria2Service() *Hysteria2Service {
	return &Hysteria2Service{}
}

// 检查是否已安装
func (h *Hysteria2Service) IsInstalled() bool {
	cmd := exec.Command("which", "hysteria")
	err := cmd.Run()
	return err == nil
}

// 获取版本信息
func (h *Hysteria2Service) GetVersion() string {
	if !h.IsInstalled() {
		return ""
	}
	cmd := exec.Command("hysteria", "version")
	output, err := cmd.Output()
	if err != nil {
		return ""
	}

	// 解析输出找到版本号
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "Version:") {
			version := strings.TrimSpace(strings.TrimPrefix(line, "Version:"))
			return version
		}
	}
	return ""
}

// 获取服务详细状态
func (h *Hysteria2Service) getServiceStatus() (string, string, string, string) {
	cmd := exec.Command("systemctl", "status", "hysteria-server.service")
	output, _ := cmd.CombinedOutput()
	outputStr := string(output)

	// 解析状态输出
	var serviceStatus, lastError, loadState, activeState string

	lines := strings.Split(outputStr, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		switch {
		case strings.Contains(line, "Loaded:"):
			if strings.Contains(line, "loaded") {
				loadState = "loaded"
			} else {
				loadState = "not-loaded"
			}
		case strings.Contains(line, "Active:"):
			if strings.Contains(line, "active (running)") {
				activeState = "active"
				serviceStatus = "running"
			} else if strings.Contains(line, "failed") {
				activeState = "failed"
				serviceStatus = "failed"
			} else if strings.Contains(line, "inactive") {
				activeState = "inactive"
				serviceStatus = "stopped"
			}
		case strings.Contains(line, "FATAL") || strings.Contains(line, "ERROR") || strings.Contains(line, "Failed"):
			if lastError == "" { // 只获取第一个错误
				lastError = strings.TrimSpace(strings.Join(strings.Split(line, ":")[2:], ":"))
			}
		}
	}

	return serviceStatus, lastError, loadState, activeState
}

// 获取运行状态
func (h *Hysteria2Service) GetStatus() (*Hysteria2Status, error) {
	status := &Hysteria2Status{
		IsInstalled: h.IsInstalled(),
	}

	if status.IsInstalled {
		// 获取版本信息
		cmd := exec.Command("hysteria", "version")
		output, err := cmd.Output()
		if err == nil {
			lines := strings.Split(string(output), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				switch {
				case strings.HasPrefix(line, "Version:"):
					status.Version = strings.TrimSpace(strings.TrimPrefix(line, "Version:"))
				case strings.HasPrefix(line, "BuildDate:"):
					status.BuildDate = strings.TrimSpace(strings.TrimPrefix(line, "BuildDate:"))
				case strings.HasPrefix(line, "BuildType:"):
					status.BuildType = strings.TrimSpace(strings.TrimPrefix(line, "BuildType:"))
				case strings.HasPrefix(line, "Platform:"):
					status.Platform = strings.TrimSpace(strings.TrimPrefix(line, "Platform:"))
				case strings.HasPrefix(line, "Architecture:"):
					status.Architecture = strings.TrimSpace(strings.TrimPrefix(line, "Architecture:"))
				}
			}
		}

		// 获取服务详细状态
		serviceStatus, lastError, loadState, activeState := h.getServiceStatus()
		status.ServiceStatus = serviceStatus
		status.LastError = lastError
		status.LoadState = loadState
		status.ActiveState = activeState
		status.IsRunning = activeState == "active"
	}

	return status, nil
}

// 安装Hysteria2
func (h *Hysteria2Service) Install() (string, error) {
	// 安装命令
	installCmd := exec.Command("bash", "-c", "curl -fsSL https://get.hy2.sh/ | bash")
	output, err := installCmd.CombinedOutput()
	if err != nil {
		return string(output), err
	}

	// 设置开机自启
	enableCmd := exec.Command("systemctl", "enable", "hysteria-server.service")
	if err := enableCmd.Run(); err != nil {
		return string(output), err
	}

	return string(output), nil
}

// 卸载Hysteria2
func (h *Hysteria2Service) Uninstall() (string, error) {
	cmd := exec.Command("bash", "-c", "curl -fsSL https://get.hy2.sh/ | bash -s -- --remove")
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// 更新Hysteria2
func (h *Hysteria2Service) Update() (string, error) {
	cmd := exec.Command("bash", "-c", "curl -fsSL https://get.hy2.sh/ | bash")
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// 获取配置
func (h *Hysteria2Service) GetConfig() (string, error) {
	data, err := os.ReadFile("/etc/hysteria/config.yaml")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// 备份配置
func (h *Hysteria2Service) BackupConfig() (string, error) {
	// 读取当前配置
	data, err := os.ReadFile("/etc/hysteria/config.yaml")
	if err != nil {
		return "", err
	}

	// 生成备份文件名（带时间戳）
	backupPath := fmt.Sprintf("/etc/hysteria/config.yaml.bak.%s",
		time.Now().Format("20060102150405"))

	// 写入备份文件
	if err := os.WriteFile(backupPath, data, 0644); err != nil {
		return "", err
	}

	return backupPath, nil
}

// 修改配置时自动备份
func (h *Hysteria2Service) UpdateConfig(config string) error {
	// 先备份当前配置
	if _, err := h.BackupConfig(); err != nil {
		return fmt.Errorf("failed to backup config: %v", err)
	}

	// 写入新配置
	if err := ioutil.WriteFile("/etc/hysteria/config.yaml", []byte(config), 0644); err != nil {
		return err
	}

	return h.Restart()
}

// 获取日志
func (h *Hysteria2Service) GetLogs(opts *LogOptions) (string, error) {
	args := []string{"--no-pager", "-u", "hysteria-server.service"}

	if opts != nil {
		if opts.Lines > 0 {
			args = append(args, "-n", fmt.Sprintf("%d", opts.Lines))
		}
		if opts.Since != "" {
			// 确保时间格式正确，例如 "5m"（5分钟）, "2h"（2小时）
			if strings.HasSuffix(opts.Since, "m") || strings.HasSuffix(opts.Since, "h") {
				args = append(args, fmt.Sprintf("--since='%s ago'", opts.Since))
			}
		}
		if opts.Level != "" {
			// 日志级别映射
			switch strings.ToLower(opts.Level) {
			case "error":
				args = append(args, "-p", "err")
			case "warning", "warn":
				args = append(args, "-p", "warning")
			case "info":
				args = append(args, "-p", "info")
			case "debug":
				args = append(args, "-p", "debug")
			}
		}
	}

	// 使用bash -c来执行命令，因为--since参数需要shell解释
	cmdStr := fmt.Sprintf("journalctl %s", strings.Join(args, " "))
	cmd := exec.Command("bash", "-c", cmdStr)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %v, output: %s", err, string(output))
	}

	return string(output), nil
}

// 启动服务
func (h *Hysteria2Service) Start() error {
	// 执行启动命令
	cmd := exec.Command("systemctl", "start", "hysteria-server.service")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to start service")
	}

	// 等待一小段时间让服务状态更新
	time.Sleep(time.Second)

	// 获取详细状态
	status, _ := h.GetStatus()
	if status != nil {
		if status.ServiceStatus != "running" {
			if status.LastError != "" {
				return fmt.Errorf("failed to start service: %s", status.LastError)
			}
			return fmt.Errorf("failed to start service: service is in %s state", status.ServiceStatus)
		}
	}

	return nil
}

// 停止服务
func (h *Hysteria2Service) Stop() error {
	// 执行停止命令
	cmd := exec.Command("systemctl", "stop", "hysteria-server.service")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to stop service")
	}

	// 等待一小段时间让服务状态更新
	time.Sleep(time.Second)

	// 只检查是否已停止
	cmd = exec.Command("systemctl", "is-active", "hysteria-server.service")
	output, _ := cmd.Output()
	status := strings.TrimSpace(string(output))

	switch status {
	case "inactive":
		return nil
	case "failed":
		return fmt.Errorf("service is in failed state")
	case "unknown":
		return fmt.Errorf("service state is unknown")
	default:
		return fmt.Errorf("failed to stop service: current state is %s", status)
	}
}

// 重启服务
func (h *Hysteria2Service) Restart() error {
	// 执行重启命令
	cmd := exec.Command("systemctl", "restart", "hysteria-server.service")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to restart service")
	}

	// 等待一小段时间让服务状态更新
	time.Sleep(time.Second)

	// 获取详细状态
	status, _ := h.GetStatus()
	if status != nil {
		if status.ServiceStatus != "running" {
			if status.LastError != "" {
				return fmt.Errorf("failed to restart service: %s", status.LastError)
			}
			return fmt.Errorf("failed to restart service: service is in %s state", status.ServiceStatus)
		}
	}

	return nil
}

// 执行健康检查
func (h *Hysteria2Service) CheckHealth() (*HealthCheck, error) {
	health := &HealthCheck{
		CheckTime: time.Now().Format(time.RFC3339),
	}

	// 检查服务状态
	status, _ := h.GetStatus()
	if status != nil {
		health.IsRunning = status.IsRunning
		health.LastError = status.LastError
	}

	// 检查配置文件是否有效
	if _, err := h.GetConfig(); err == nil {
		health.ConfigValid = true
	}

	// TODO: 添加端口检查逻辑

	return health, nil
}

// 获取可用版本列表
func (h *Hysteria2Service) GetAvailableVersions() ([]string, error) {
	// 从GitHub API获取版本列表
	cmd := exec.Command("curl", "-s", "https://api.github.com/repos/apernet/hysteria/releases")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	// 解析JSON响应
	var releases []struct {
		TagName string `json:"tag_name"`
	}
	if err := json.Unmarshal(output, &releases); err != nil {
		return nil, fmt.Errorf("failed to parse GitHub response: %v", err)
	}

	// 提取版本号
	versions := make([]string, 0, len(releases))
	for _, release := range releases {
		versions = append(versions, release.TagName)
	}

	return versions, nil
}

// 安装指定版本
func (h *Hysteria2Service) InstallVersion(version string) error {
	cmd := exec.Command("bash", "-c", fmt.Sprintf("curl -fsSL https://get.hy2.sh/ | bash -s -- --version %s", version))
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install version %s: %v", version, err)
	}
	return nil
}

// 获取配置备份列表
func (h *Hysteria2Service) GetConfigBackups() ([]string, error) {
	// 读取/etc/hysteria目录下的所有备份文件
	files, err := os.ReadDir("/etc/hysteria")
	if err != nil {
		return nil, err
	}

	// 筛选出备份文件
	backups := make([]string, 0)
	for _, file := range files {
		if strings.HasPrefix(file.Name(), "config.yaml.bak.") {
			backups = append(backups, file.Name())
		}
	}

	// 按时间倒序排序（最新的在前）
	sort.Slice(backups, func(i, j int) bool {
		return backups[i] > backups[j]
	})

	return backups, nil
}

// 恢复配置备份
func (h *Hysteria2Service) RestoreConfig(backup string) error {
	// 安全检查：确保文件名是备份文件
	if !strings.HasPrefix(backup, "config.yaml.bak.") {
		return fmt.Errorf("invalid backup file name")
	}

	backupPath := filepath.Join("/etc/hysteria", backup)
	configPath := "/etc/hysteria/config.yaml"

	// 检查备份文件是否存在
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("backup file not found: %s", backup)
	}

	// 读取备份文件
	data, err := os.ReadFile(backupPath)
	if err != nil {
		return fmt.Errorf("failed to read backup file: %v", err)
	}

	// 写入配置文件
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to restore config: %v", err)
	}

	// 重启服务以应用新配置
	return h.Restart()
}
