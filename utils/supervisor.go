package utils

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// SupervisorctlStart 启动 supervisor 服务
func SupervisorctlStart(name string) {
	if err := supervisorctlBase(name, "start"); err != nil {
		fmt.Println(Red(fmt.Sprintf("启动%s失败: %v", name, err)))
	} else {
		fmt.Println(Green(fmt.Sprintf("启动%s成功!", name)))
	}
}

// SupervisorctlStop 停止 supervisor 服务
func SupervisorctlStop(name string) {
	if err := supervisorctlBase(name, "stop"); err != nil {
		fmt.Println(Red(fmt.Sprintf("停止%s失败: %v", name, err)))
	} else {
		fmt.Println(Green(fmt.Sprintf("停止%s成功!", name)))
	}
}

// SupervisorctlRestart 重启 supervisor 服务
func SupervisorctlRestart(name string) {
	if err := supervisorctlBase(name, "restart"); err != nil {
		fmt.Println(Red(fmt.Sprintf("重启%s失败: %v", name, err)))
	} else {
		fmt.Println(Green(fmt.Sprintf("重启%s成功!", name)))
	}
}

// SupervisorctlStatus 获取 supervisor 服务状态
func SupervisorctlStatus(name string) string {
	out, err := supervisorctlBaseWithOutput(name, "status")
	if err != nil {
		return fmt.Sprintf("获取%s状态失败: %v", name, err)
	}
	return out
}

// SupervisorctlReload 重新加载 supervisor 配置
func SupervisorctlReload() error {
	cmd := exec.Command("supervisorctl", "reread")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("重新读取配置失败: %v", err)
	}
	
	cmd = exec.Command("supervisorctl", "update")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("更新配置失败: %v", err)
	}
	
	return nil
}

// StartSupervisord 启动 supervisord 守护进程
func StartSupervisord(configPath string) error {
	// 检查 supervisord 是否已经运行
	if IsSupervisordRunning() {
		fmt.Println("supervisord 已经在运行")
		return nil
	}
	
	var cmd *exec.Cmd
	if configPath != "" {
		cmd = exec.Command("supervisord", "-c", configPath)
	} else {
		cmd = exec.Command("supervisord")
	}
	
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动 supervisord 失败: %v", err)
	}
	
	fmt.Println(Green("supervisord 启动成功!"))
	return nil
}

// StopSupervisord 停止 supervisord 守护进程
func StopSupervisord() error {
	cmd := exec.Command("supervisorctl", "shutdown")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("停止 supervisord 失败: %v", err)
	}
	
	fmt.Println(Green("supervisord 停止成功!"))
	return nil
}

// IsSupervisordRunning 检查 supervisord 是否正在运行
func IsSupervisordRunning() bool {
	cmd := exec.Command("supervisorctl", "status")
	err := cmd.Run()
	return err == nil
}

// GetAllServiceStatus 获取所有服务状态
func GetAllServiceStatus() string {
	out, err := exec.Command("supervisorctl", "status").CombinedOutput()
	if err != nil {
		return fmt.Sprintf("获取服务状态失败: %v", err)
	}
	return string(out)
}

// supervisorctlBase 执行 supervisorctl 命令的基础函数
func supervisorctlBase(name, operate string) error {
	cmd := exec.Command("supervisorctl", operate, name)
	return cmd.Run()
}

// supervisorctlBaseWithOutput 执行 supervisorctl 命令并返回输出
func supervisorctlBaseWithOutput(name, operate string) (string, error) {
	cmd := exec.Command("supervisorctl", operate, name)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// IsServiceRunning 检查特定服务是否正在运行
func IsServiceRunning(serviceName string) bool {
	out, err := supervisorctlBaseWithOutput(serviceName, "status")
	if err != nil {
		return false
	}
	return strings.Contains(out, "RUNNING")
}

// GetServiceLogs 获取服务日志
func GetServiceLogs(serviceName string, lines int) (string, error) {
	var cmd *exec.Cmd
	if lines > 0 {
		cmd = exec.Command("supervisorctl", "tail", "-"+fmt.Sprintf("%d", lines), serviceName)
	} else {
		cmd = exec.Command("supervisorctl", "tail", serviceName)
	}
	
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("获取日志失败: %v", err)
	}
	
	return string(out), nil
}

// FollowServiceLogs 实时跟踪服务日志
func FollowServiceLogs(serviceName string) error {
	cmd := exec.Command("supervisorctl", "tail", "-f", serviceName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

// CheckSupervisorInstalled 检查 supervisor 是否已安装
func CheckSupervisorInstalled() bool {
	_, err := exec.LookPath("supervisorctl")
	return err == nil
}

// GetSupervisorVersion 获取 supervisor 版本
func GetSupervisorVersion() (string, error) {
	cmd := exec.Command("supervisorctl", "version")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}
