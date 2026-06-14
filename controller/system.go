package controller

import (
	"bufio"
	"fmt"
	"math"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// SystemController 部署系统状态
type SystemController struct{}

// MemoryInfo 内存信息（单位：字节）
type MemoryInfo struct {
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	Available   uint64  `json:"available"`
	Buffers     uint64  `json:"buffers"`
	Cached      uint64  `json:"cached"`
	SwapTotal   uint64  `json:"swapTotal"`
	SwapUsed    uint64  `json:"swapUsed"`
	UsedPercent float64 `json:"usedPercent"`
}

// DiskInfo 磁盘信息（单位：字节）
type DiskInfo struct {
	Path        string  `json:"path"`
	Total       uint64  `json:"total"`
	Used        uint64  `json:"used"`
	Free        uint64  `json:"free"`
	UsedPercent float64 `json:"usedPercent"`
}

// SystemInfo 部署系统 + 服务的整体状态
type SystemInfo struct {
	Version       string     `json:"version"`         // 服务版本
	GoVersion     string     `json:"goVersion"`       // Go 运行时版本
	Hostname      string     `json:"hostname"`        // 主机名
	OS            string     `json:"os"`              // 操作系统 (GOOS)
	Arch          string     `json:"arch"`            // 架构 (GOARCH)
	KernelVersion string     `json:"kernelVersion"`   // 内核版本
	NumCPU        int        `json:"numCpu"`          // CPU 核心数
	CPUUsage      float64    `json:"cpuUsagePercent"` // CPU 使用率 %
	LoadAvg       []float64  `json:"loadAvg"`         // 平均负载 1/5/15 分钟
	Uptime        string     `json:"uptime"`          // 运行时间（可读）
	UptimeSeconds uint64     `json:"uptimeSeconds"`   // 运行时间（秒）
	LocalTime     string     `json:"localTime"`       // 本地时间
	Memory        MemoryInfo `json:"memory"`
	Disk          DiskInfo   `json:"disk"`
}

// AppVersion 返回服务版本：优先环境变量 APP_VERSION，缺省 "1.0.0"。
func AppVersion() string {
	if v := strings.TrimSpace(os.Getenv("APP_VERSION")); v != "" {
		return v
	}
	return "1.0.0"
}

// GetSystemInfo 返回部署系统状态（服务版本 + 主机 CPU/内存/磁盘/负载等）
func (c *SystemController) GetSystemInfo(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, collectSystemInfo())
}

func collectSystemInfo() SystemInfo {
	hostname, _ := os.Hostname()
	info := SystemInfo{
		Version:       AppVersion(),
		GoVersion:     runtime.Version(),
		Hostname:      hostname,
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		KernelVersion: readKernelVersion(),
		NumCPU:        runtime.NumCPU(),
		CPUUsage:      readCPUUsage(),
		LoadAvg:       readLoadAvg(),
		LocalTime:     time.Now().Format("2006-01-02 15:04:05"),
		Memory:        readMemInfo(),
		Disk:          diskUsage(diskPath()),
	}
	if up, ok := readUptime(); ok {
		info.UptimeSeconds = up
		info.Uptime = formatDuration(up)
	}
	return info
}

// diskPath 磁盘统计路径：环境变量 DISK_PATH，缺省 "/"
func diskPath() string {
	if p := strings.TrimSpace(os.Getenv("DISK_PATH")); p != "" {
		return p
	}
	return "/"
}

func readKernelVersion() string {
	b, err := os.ReadFile("/proc/sys/kernel/osrelease")
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(b))
}

func readLoadAvg() []float64 {
	out := []float64{}
	b, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return out
	}
	fields := strings.Fields(string(b))
	for i := 0; i < 3 && i < len(fields); i++ {
		if v, err := strconv.ParseFloat(fields[i], 64); err == nil {
			out = append(out, v)
		}
	}
	return out
}

func readUptime() (uint64, bool) {
	b, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, false
	}
	fields := strings.Fields(string(b))
	if len(fields) == 0 {
		return 0, false
	}
	sec, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, false
	}
	return uint64(sec), true
}

func readMemInfo() MemoryInfo {
	var mi MemoryInfo
	b, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return mi
	}
	vals := map[string]uint64{}
	sc := bufio.NewScanner(strings.NewReader(string(b)))
	for sc.Scan() {
		parts := strings.Fields(sc.Text())
		if len(parts) >= 2 {
			key := strings.TrimSuffix(parts[0], ":")
			if v, err := strconv.ParseUint(parts[1], 10, 64); err == nil {
				vals[key] = v * 1024 // kB -> bytes
			}
		}
	}
	mi.Total = vals["MemTotal"]
	mi.Free = vals["MemFree"]
	mi.Available = vals["MemAvailable"]
	mi.Buffers = vals["Buffers"]
	mi.Cached = vals["Cached"]
	mi.SwapTotal = vals["SwapTotal"]
	if vals["SwapTotal"] >= vals["SwapFree"] {
		mi.SwapUsed = vals["SwapTotal"] - vals["SwapFree"]
	}
	if mi.Available > 0 && mi.Total >= mi.Available {
		mi.Used = mi.Total - mi.Available
	} else if mi.Total >= mi.Free+mi.Buffers+mi.Cached {
		mi.Used = mi.Total - mi.Free - mi.Buffers - mi.Cached
	}
	if mi.Total > 0 {
		mi.UsedPercent = math.Round(float64(mi.Used)/float64(mi.Total)*1000) / 10
	}
	return mi
}

// readCPUUsage 采样 /proc/stat 两次（间隔 200ms）计算总体 CPU 使用率
func readCPUUsage() float64 {
	idle1, total1, ok1 := readCPUStat()
	if !ok1 {
		return 0
	}
	time.Sleep(200 * time.Millisecond)
	idle2, total2, ok2 := readCPUStat()
	if !ok2 {
		return 0
	}
	dt := total2 - total1
	di := idle2 - idle1
	if dt == 0 {
		return 0
	}
	usage := (1.0 - float64(di)/float64(dt)) * 100.0
	if usage < 0 {
		usage = 0
	}
	if usage > 100 {
		usage = 100
	}
	return math.Round(usage*10) / 10
}

func readCPUStat() (idle, total uint64, ok bool) {
	b, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0, false
	}
	for _, line := range strings.Split(string(b), "\n") {
		if strings.HasPrefix(line, "cpu ") {
			fields := strings.Fields(line)[1:]
			for i, f := range fields {
				v, _ := strconv.ParseUint(f, 10, 64)
				total += v
				if i == 3 { // idle 为第 4 个字段
					idle = v
				}
			}
			return idle, total, true
		}
	}
	return 0, 0, false
}

func formatDuration(sec uint64) string {
	d := sec / 86400
	h := (sec % 86400) / 3600
	m := (sec % 3600) / 60
	s := sec % 60
	parts := []string{}
	if d > 0 {
		parts = append(parts, fmt.Sprintf("%dd", d))
	}
	if h > 0 || d > 0 {
		parts = append(parts, fmt.Sprintf("%dh", h))
	}
	if m > 0 || h > 0 || d > 0 {
		parts = append(parts, fmt.Sprintf("%dm", m))
	}
	parts = append(parts, fmt.Sprintf("%ds", s))
	return strings.Join(parts, " ")
}
