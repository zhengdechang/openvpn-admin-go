package cmd

import (
	"fmt"
	"openvpn-admin-go/constants"
	"openvpn-admin-go/database"
	"openvpn-admin-go/logging"
	"openvpn-admin-go/openvpn"
	"openvpn-admin-go/router"
	"openvpn-admin-go/services"
	"openvpn-admin-go/utils"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

var webPort int

// WebMenu displays the web service management menu
func WebMenu() {
	for {
		fmt.Println()
		fmt.Println("=== Web 服务管理 ===")
		fmt.Println()
		fmt.Println("1.启动 Web 服务（前后端）")
		fmt.Println("2.停止 Web 服务（前后端）")
		fmt.Println("3.查看 API 服务状态")
		fmt.Println("4.查看 API 服务日志")
		fmt.Println("5.查看前端状态")
		fmt.Println("6.查看前端日志")
		fmt.Println("7.查看所有服务状态")
		fmt.Println("0.返回主菜单")
		fmt.Print("请选择操作 (0-7): ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "0":
			return
		case "1":
			frontendPort := promptFrontendPort()
			startWebServices(frontendPort)
		case "2":
			stopWebServices()
		case "3":
			checkAPIServiceStatus()
		case "4":
			showAPIServiceLogs()
		case "5":
			checkFrontendStatus()
		case "6":
			showFrontendLogs()
		case "7":
			printAllServiceStatus()
		default:
			fmt.Println("无效选择，请重试")
		}
	}
}

func promptFrontendPort() int {
	fmt.Printf("请输入前端 (Nginx) 监听端口 (默认%d): ", constants.DefaultFrontendPort)
	var portInput string
	fmt.Scanln(&portInput)

	port := constants.DefaultFrontendPort
	if portInput != "" {
		if p, err := strconv.Atoi(portInput); err == nil && p > 0 && p <= 65535 {
			port = p
		} else {
			fmt.Printf("端口号无效，使用默认端口%d\n", constants.DefaultFrontendPort)
		}
	}

	return port
}

func startWebServices(frontendPort int) {
	backendPort := constants.DefaultWebAPIPort
	fmt.Printf("后端 API 端口固定为 %d\n", backendPort)
	startAPIService(backendPort)

	if err := startFrontendService(frontendPort); err != nil {
		fmt.Printf("启动前端失败: %v\n", err)
		return
	}

	fmt.Println("前后端服务已启动")
}

func stopWebServices() {
	stopFrontendService()
	stopAPIService()
}

// installAPIService 安装 API 服务的supervisor配置
func installAPIService(port int) error {
	// 检查 supervisor 是否已安装
	if !utils.CheckSupervisorInstalled() {
		return fmt.Errorf("supervisor 未安装，请先安装 supervisor")
	}

	// 获取当前工作目录
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %v", err)
	}

	// 获取二进制文件路径
	binaryPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取可执行文件路径失败: %v", err)
	}

	// 安装 supervisor 主配置文件（如果不存在）
	if !utils.IsSupervisorConfigExists() {
		if err := utils.InstallSupervisorMainConfig(); err != nil {
			return fmt.Errorf("安装 supervisor 主配置失败: %v", err)
		}
	}

	// 创建 API 服务配置
	apiConfig := utils.ServiceConfig{
		BinaryPath:       binaryPath,
		WorkingDirectory: wd,
		Port:             port,
		DBPath:           "/app/data/db.sqlite3",
		OpenVPNConfigDir: "/etc/openvpn",
		AutoStart:        false, // 默认不自动启动
	}

	// 安装 API 服务配置
	if err := utils.InstallWebServiceConfig(apiConfig); err != nil {
		return fmt.Errorf("安装 API 服务配置失败: %v", err)
	}

	// 启动 supervisord（如果未运行）
	if !utils.IsSupervisordRunning() {
		if err := utils.StartSupervisord(""); err != nil {
			return fmt.Errorf("启动 supervisord 失败: %v", err)
		}
	} else {
		// 重新加载配置
		if err := utils.SupervisorctlReload(); err != nil {
			return fmt.Errorf("重新加载 supervisor 配置失败: %v", err)
		}
	}

	fmt.Printf("API 服务已安装到 supervisor\n")
	return nil
}

func startAPIService(port int) {
	// 检查 API 服务配置是否存在
	needReinstall := false

	if !utils.IsWebServiceConfigExists() {
		// 配置文件不存在，需要安装
		needReinstall = true
		fmt.Printf("API 服务配置未安装，正在安装（端口: %d）...\n", port)
	} else {
		// 配置文件存在，检查端口是否匹配
		content, err := os.ReadFile(constants.SupervisorWebConfigPath)
		if err == nil {
			expectedPort := fmt.Sprintf("--port %d", port)
			if !strings.Contains(string(content), expectedPort) {
				needReinstall = true
				fmt.Printf("检测到端口变化，正在重新安装 API 服务（端口: %d）...\n", port)
			}
		}
	}

	if needReinstall {
		if err := installAPIService(port); err != nil {
			fmt.Printf("安装 API 服务失败: %v\n", err)
			return
		}
	}

	// 确保 supervisord 正在运行
	if !utils.IsSupervisordRunning() {
		fmt.Println("启动 supervisord...")
		if err := utils.StartSupervisord(""); err != nil {
			fmt.Printf("启动 supervisord 失败: %v\n", err)
			return
		}
	}

	// 启动 API 服务
	fmt.Printf("正在启动 API 服务（端口: %d）...\n", port)
	utils.SupervisorctlStart(constants.SupervisorWebServiceName)
}

func stopAPIService() {
	// 停止服务
	fmt.Println("正在停止 API 服务...")
	utils.SupervisorctlStop(constants.SupervisorWebServiceName)
}

func checkAPIServiceStatus() {
	// 获取服务状态
	fmt.Println("=== API 服务状态 ===")
	statusOutput := utils.SupervisorctlStatus(constants.SupervisorWebServiceName)
	if statusOutput != "" {
		fmt.Printf("%s\n", statusOutput)
	} else {
		fmt.Println("无法获取服务状态")
	}

	promptReturn()
}

func showAPIServiceLogs() {
	// 使用 supervisor 查看服务日志
	fmt.Println("\n=== API 服务日志 (最近50行) ===")
	output, err := utils.GetServiceLogs(constants.SupervisorWebServiceName, 50)
	if err != nil {
		fmt.Printf("获取服务日志失败: %v\n", err)
		promptReturn()
		return
	}

	if output == "" {
		fmt.Println("日志为空")
	} else {
		fmt.Println(output)
	}

	promptReturn()
}

func startFrontendService(port int) error {
	if !utils.IsSupervisorConfigExists() {
		if err := utils.InstallSupervisorMainConfig(); err != nil {
			return fmt.Errorf("安装 supervisor 主配置失败: %v", err)
		}
	}

	if err := utils.UpdateNginxListenPort(constants.NginxConfigPath, port); err != nil {
		return fmt.Errorf("更新前端监听端口失败: %v", err)
	}

	if !utils.IsFrontendServiceConfigExists() {
		fmt.Println("前端服务配置未安装，正在安装...")
		if err := utils.InstallFrontendServiceConfig(false); err != nil {
			return fmt.Errorf("安装前端服务配置失败: %v", err)
		}
	}

	if !utils.IsSupervisordRunning() {
		fmt.Println("启动 supervisord...")
		if err := utils.StartSupervisord(""); err != nil {
			return fmt.Errorf("启动 supervisord 失败: %v", err)
		}
	} else {
		if err := utils.SupervisorctlReload(); err != nil {
			return fmt.Errorf("重新加载 supervisor 配置失败: %v", err)
		}
	}

	fmt.Printf("正在启动前端 (Nginx)（端口: %d）...\n", port)
	utils.SupervisorctlStart(constants.SupervisorFrontendServiceName)
	return nil
}

func stopFrontendService() {
	fmt.Println("正在停止前端 (Nginx)...")
	utils.SupervisorctlStop(constants.SupervisorFrontendServiceName)
}

func checkFrontendStatus() {
	fmt.Println("=== 前端服务状态 ===")
	statusOutput := utils.SupervisorctlStatus(constants.SupervisorFrontendServiceName)
	if statusOutput != "" {
		fmt.Printf("%s\n", statusOutput)
	} else {
		fmt.Println("无法获取服务状态")
	}

	promptReturn()
}

func showFrontendLogs() {
	fmt.Println("\n=== 前端服务日志 (最近50行) ===")
	output, err := utils.GetServiceLogs(constants.SupervisorFrontendServiceName, 50)
	if err != nil {
		fmt.Printf("获取服务日志失败: %v\n", err)
		promptReturn()
		return
	}

	if output == "" {
		fmt.Println("日志为空")
	} else {
		fmt.Println(output)
	}

	promptReturn()
}

func printAllServiceStatus() {
	fmt.Println("\n=== 所有服务状态 ===")
	allStatus := utils.GetAllServiceStatus()
	fmt.Printf("%s\n", allStatus)
	promptReturn()
}

func promptReturn() {
	fmt.Println("\n按回车键返回...")
	fmt.Scanln()
}

// Web Server Functions - 实际的 Web 服务器实现

// runWebServer 实际的 Web 服务器实现
func runWebServer(port int) error {
	// 启动 OpenVPN 同步服务（核心已在 main 中初始化）
	cfg, err := openvpn.LoadConfig()
	if err != nil {
		return fmt.Errorf("无法加载 OpenVPN 配置以启动同步服务: %v", err)
	}
	statusLogPath := cfg.OpenVPNStatusLogPath
	syncInterval := utils.GetOpenVPNSyncInterval()
	logging.Info("Starting OpenVPN Sync Service: LogPath='%s', Interval=%s", statusLogPath, syncInterval)
	go services.StartOpenVPNSyncService(database.DB, statusLogPath, syncInterval)

	// Setup Gin router
	r := gin.Default()

	// 添加日志中间件
	r.Use(logging.GinLoggingMiddleware())

	api := r.Group("/api")
	{
		router.SetupHealthRoutes(api)
		router.SetupUserRoutes(api)
		router.SetupManageRoutes(api)
		router.SetupServerRoutes(api)
		router.SetupClientRoutes(api)
		router.SetupLogRoutes(api)
	}

	serverAddr := fmt.Sprintf(":%d", port)
	logging.Info("Web 服务器正在监听 %s...", serverAddr)
	fmt.Printf("Web 服务器正在监听 %s...\n", serverAddr)
	return r.Run(serverAddr)
}

// isRunningInSystemd 检查是否在 systemd 服务中运行
func isRunningInSystemd() bool {
	// 检查是否有 systemd 相关的环境变量
	if os.Getenv("INVOCATION_ID") != "" {
		return true
	}
	// 检查父进程是否为 systemd
	if _, err := os.Stat("/run/systemd/system"); err == nil {
		return true
	}
	return false
}

func init() {
	// 添加 web 子命令
	webServerCmd := &cobra.Command{
		Use:   "web",
		Short: "运行 Web 服务器",
		Run: func(cmd *cobra.Command, args []string) {
			// 从命令行参数获取端口
			port, _ := cmd.Flags().GetInt("port")
			// 直接运行 web 服务器
			if err := runWebServer(port); err != nil {
				logging.Fatal("Web 服务器错误: %v", err)
			}
		},
	}
	webServerCmd.Flags().IntVarP(&webPort, "port", "p", constants.DefaultWebAPIPort, "Web 服务器监听的端口")
	rootCmd.AddCommand(webServerCmd)
}
