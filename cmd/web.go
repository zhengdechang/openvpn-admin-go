package cmd

import (
	"fmt"
	"log"
	"openvpn-admin-go/router"
	"github.com/gin-gonic/gin"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"
	"github.com/manifoldco/promptui"
	"os/exec"
	"github.com/spf13/cobra"
)

const pidFileName = "web.pid"

var webPort int

// WebMenu displays the web service management menu
func WebMenu() {
	for {
		fmt.Println("\n=== Web 服务管理 ===\n")
		fmt.Println("1.启动服务\n")
		fmt.Println("2.停止服务\n")
		fmt.Println("3.重启服务\n")
		fmt.Println("4.查看服务状态\n")
		fmt.Println("5.查看服务日志\n")
		fmt.Println("0.返回主菜单\n")
		fmt.Print("请选择操作 (0-5): ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "0":
			return
		case "1":
			startWebService()
		case "2":
			stopWebService()
		case "3":
			restartWebService()
		case "4":
			checkWebServiceStatus()
		case "5":
			showWebServiceLogs()
		default:
			fmt.Println("无效选择，请重试")
		}
	}
}

func startWebService() {
	prompt := promptui.Prompt{
		Label:   "请输入端口号",
		Default: "8085",
	}

	portStr, err := prompt.Run()
	if err != nil {
		fmt.Printf("输入失败: %v\n", err)
		return
	}

	port := 8085
	if portStr != "" {
		if p, err := strconv.Atoi(portStr); err == nil {
			port = p
		}
	}

	webPort = port
	if err := startWebServer(); err != nil {
		fmt.Printf("启动服务失败: %v\n", err)
	}
}

func stopWebService() {
	if err := stopWebServer(); err != nil {
		fmt.Printf("停止服务失败: %v\n", err)
	}
}

func restartWebService() {
	if err := restartWebServer(); err != nil {
		fmt.Printf("重启服务失败: %v\n", err)
	}
}

func checkWebServiceStatus() {
	pidDir := os.TempDir()
	pidFile := filepath.Join(pidDir, pidFileName)

	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Web 服务未运行")
		} else {
			fmt.Printf("检查服务状态失败: %v\n", err)
		}
		return
	}

	pid, err := strconv.Atoi(string(pidData))
	if err != nil {
		fmt.Printf("PID 文件内容无效: %v\n", err)
		return
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		fmt.Printf("获取进程信息失败: %v\n", err)
		return
	}

	if err := process.Signal(syscall.Signal(0)); err != nil {
		fmt.Println("Web 服务未运行")
	} else {
		fmt.Printf("Web 服务正在运行 (PID: %d)\n", pid)
	}
	fmt.Println("\n按回车键返回...")
	fmt.Scanln()
}

func showWebServiceLogs() {
	// 获取日志文件路径
	logFile := filepath.Join(os.TempDir(), "web.log")

	// 检查日志文件是否存在
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		fmt.Println("Web 服务日志文件不存在")
		fmt.Println("\n按回车键返回...")
		fmt.Scanln()
		return
	}

	// 读取并显示日志内容
	content, err := os.ReadFile(logFile)
	if err != nil {
		fmt.Printf("读取日志文件失败: %v\n", err)
		fmt.Println("\n按回车键返回...")
		fmt.Scanln()
		return
	}

	// 显示日志内容
	fmt.Println("\n=== Web 服务日志 ===")
	fmt.Println(string(content))
	fmt.Println("\n按回车键返回...")
	fmt.Scanln()
}

// Web Server Functions
func startWebServer() error {
	if err := CoreInitializer(); err != nil {
		return fmt.Errorf("核心初始化失败: %v", err)
	}

	pidDir := os.TempDir()
	pidFile := filepath.Join(pidDir, pidFileName)
	logFile := filepath.Join(pidDir, "web.log")

	// Check if PID file exists and process is running
	if pidData, err := os.ReadFile(pidFile); err == nil {
		if pid, err := strconv.Atoi(string(pidData)); err == nil {
			if process, _ := os.FindProcess(pid); process != nil && process.Signal(syscall.Signal(0)) == nil {
				return fmt.Errorf("Web 服务器似乎已在运行 (PID: %d, PID 文件: %s). 请先停止它", pid, pidFile)
			}
			log.Printf("发现陈旧的 PID 文件 (%s, PID: %d), 但进程未运行. 将覆盖 PID 文件.", pidFile, pid)
		} else {
			log.Printf("PID 文件 (%s) 内容无效, 将覆盖.", pidFile)
		}
	}

	// 创建日志文件
	logF, err := os.Create(logFile)
	if err != nil {
		return fmt.Errorf("创建日志文件失败: %v", err)
	}
	defer logF.Close()

	// 创建新进程
	cmd := exec.Command(os.Args[0], "web-server", "--port", strconv.Itoa(webPort))
	cmd.Stdout = logF
	cmd.Stderr = logF

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("启动 Web 服务器进程失败: %v", err)
	}

	// 保存新进程的 PID
	if err := os.WriteFile(pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644); err != nil {
		cmd.Process.Kill() // 清理失败的进程
		return fmt.Errorf("无法写入 PID 文件 %s: %v", pidFile, err)
	}

	fmt.Printf("Web 服务器正在启动 (PID: %d)...\n", cmd.Process.Pid)
	fmt.Printf("日志文件位置: %s\n", logFile)

	// 等待一段时间确保服务器启动
	time.Sleep(2 * time.Second)

	// 检查进程是否还在运行
	if err := cmd.Process.Signal(syscall.Signal(0)); err != nil {
		return fmt.Errorf("Web 服务器启动失败")
	}

	fmt.Printf("Web 服务器已成功启动 (PID: %d)\n", cmd.Process.Pid)
	fmt.Printf("PID 文件位置: %s\n", pidFile)

	return nil
}

// 实际的 Web 服务器实现
func runWebServer(port int) error {
	// Setup Gin router
	r := gin.Default()
	api := r.Group("/api")
	{
		router.SetupUserRoutes(api)
		router.SetupManageRoutes(api)
		router.SetupServerRoutes(api)
		router.SetupClientRoutes(api)
		router.SetupLogRoutes(api)
	}

	serverAddr := fmt.Sprintf(":%d", port)
	fmt.Printf("Web 服务器正在监听 %s...\n", serverAddr)
	return r.Run(serverAddr)
}

func stopWebServer() error {
	pidDir := os.TempDir()
	pidFile := filepath.Join(pidDir, pidFileName)

	pidData, err := os.ReadFile(pidFile)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("PID 文件未找到. 服务器可能未运行")
		}
		return fmt.Errorf("读取 PID 文件失败 %s: %v", pidFile, err)
	}

	pid, err := strconv.Atoi(string(pidData))
	if err != nil {
		return fmt.Errorf("PID 文件内容无效 %s: %v", pidFile, err)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		if rmErr := os.Remove(pidFile); rmErr == nil {
			log.Printf("已删除陈旧/无效的 PID 文件: %s", pidFile)
		}
		return fmt.Errorf("获取进程句柄失败 (PID: %d): %v. 可能进程已不存在", pid, err)
	}

	log.Printf("正在向 PID %d 发送 SIGTERM 信号...", pid)
	if err := process.Signal(syscall.SIGTERM); err != nil {
		if checkErr := process.Signal(syscall.Signal(0)); checkErr != nil {
			log.Printf("进程 (PID: %d) 似乎未在运行: %v", pid, checkErr)
			if rmErr := os.Remove(pidFile); rmErr == nil {
				log.Printf("已删除陈旧的 PID 文件: %s", pidFile)
			}
		} else {
			return fmt.Errorf("发送 SIGTERM 信号失败 (PID: %d): %v", pid, err)
		}
	} else {
		log.Printf("已成功向 PID %d 发送 SIGTERM 信号.", pid)
		time.Sleep(100 * time.Millisecond)
		if _, statErr := os.Stat(pidFile); statErr == nil {
			if rmErr := os.Remove(pidFile); rmErr == nil {
				log.Println("由停止命令删除的 PID 文件.")
			}
		}
	}
	fmt.Println("Web 服务已停止")
	return nil
}

func restartWebServer() error {
	log.Println("正在重启 Web 服务器...")
	
	if err := stopWebServer(); err != nil {
		log.Printf("停止服务时出错: %v", err)
	}
	
	log.Println("等待服务器关闭 (2 秒)...")
	time.Sleep(2 * time.Second)
	
	if err := startWebServer(); err != nil {
		return fmt.Errorf("启动服务时出错: %v", err)
	}
	
	log.Println("Web 服务器重启过程已启动.")
	return nil
}

func init() {
	// 添加 web-server 子命令
	webServerCmd := &cobra.Command{
		Use:   "web-server",
		Short: "运行 Web 服务器",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runWebServer(webPort); err != nil {
				log.Fatalf("Web 服务器错误: %v", err)
			}
		},
	}
	webServerCmd.Flags().IntVarP(&webPort, "port", "p", 8085, "Web 服务器监听的端口")
	rootCmd.AddCommand(webServerCmd)
} 