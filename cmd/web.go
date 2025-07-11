package cmd

import (
	"bytes"
	"fmt"
	"openvpn-admin-go/constants"
	"openvpn-admin-go/database"
	"openvpn-admin-go/logging"
	"openvpn-admin-go/openvpn"
	"openvpn-admin-go/router"
	"openvpn-admin-go/services"
	"openvpn-admin-go/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

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
			fmt.Print("请输入Web服务端口 (默认8085): ")
			var portInput string
			fmt.Scanln(&portInput)
			port := 8085 // 默认端口
			if portInput != "" {
				if p, err := strconv.Atoi(portInput); err == nil && p > 0 && p <= 65535 {
					port = p
				} else {
					fmt.Println("端口号无效，使用默认端口8085")
				}
			}
			startWebService(port)
		case "2":
			stopWebService()
		case "3":
			fmt.Print("请输入Web服务端口 (默认8085): ")
			var portInput string
			fmt.Scanln(&portInput)
			port := 8085 // 默认端口
			if portInput != "" {
				if p, err := strconv.Atoi(portInput); err == nil && p > 0 && p <= 65535 {
					port = p
				} else {
					fmt.Println("端口号无效，使用默认端口8085")
				}
			}
			restartWebService(port)
		case "4":
			checkWebServiceStatus()
		case "5":
			showWebServiceLogs()
		default:
			fmt.Println("无效选择，请重试")
		}
	}
}



// installWebService 安装web服务的systemd服务文件
func installWebService(port int) error {
	// 检查是否以root权限运行
	if os.Geteuid() != 0 {
		return fmt.Errorf("请使用 sudo 运行程序以安装服务")
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

	// 模板数据
	data := map[string]interface{}{
		"WorkingDirectory": wd,
		"BinaryPath":       binaryPath,
		"Port":             port,
		"ConfigDirectory":  "/etc/openvpn",
	}

	// 解析模板
	templatePath := filepath.Join(wd, "template", "openvpn-web.j2")
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return fmt.Errorf("解析web服务模板失败: %v", err)
	}

	// 生成服务文件内容
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return fmt.Errorf("渲染web服务模板失败: %v", err)
	}

	// 写入systemd服务文件
	servicePath := "/etc/systemd/system/" + constants.WebServiceName
	if err := os.WriteFile(servicePath, buf.Bytes(), 0644); err != nil {
		return fmt.Errorf("写入服务文件失败: %v", err)
	}

	// 重新加载systemd配置
	reloadOutput := utils.ExecCommandWithResult("systemctl daemon-reload")
	if strings.Contains(reloadOutput, "Failed") || strings.Contains(reloadOutput, "failed") {
		return fmt.Errorf("重新加载systemd配置失败: %s", reloadOutput)
	}

	// 启用服务
	enableOutput := utils.ExecCommandWithResult(fmt.Sprintf("systemctl enable %s", constants.WebServiceName))
	if strings.Contains(enableOutput, "Failed") || strings.Contains(enableOutput, "failed") {
		return fmt.Errorf("启用web服务失败: %s", enableOutput)
	}

	fmt.Printf("Web服务已安装: %s\n", servicePath)
	return nil
}

func startWebService(port int) {
	// 检查服务文件是否存在以及端口是否匹配
	servicePath := "/etc/systemd/system/" + constants.WebServiceName
	needReinstall := false

	if _, err := os.Stat(servicePath); os.IsNotExist(err) {
		// 服务文件不存在，需要安装
		needReinstall = true
		fmt.Printf("Web服务未安装，正在安装（端口: %d）...\n", port)
	} else {
		// 服务文件存在，检查端口是否匹配
		content, err := os.ReadFile(servicePath)
		if err == nil {
			expectedPort := fmt.Sprintf("--port %d", port)
			if !strings.Contains(string(content), expectedPort) {
				needReinstall = true
				fmt.Printf("检测到端口变化，正在重新安装Web服务（端口: %d）...\n", port)
			}
		}
	}

	if needReinstall {
		if err := installWebService(port); err != nil {
			fmt.Printf("安装Web服务失败: %v\n", err)
			return
		}

		// 重新加载 systemd 配置
		fmt.Println("重新加载 systemd 配置...")
		reloadOutput := utils.ExecCommandWithResult("systemctl daemon-reload")
		if strings.Contains(reloadOutput, "Failed") || strings.Contains(reloadOutput, "failed") {
			fmt.Printf("重新加载 systemd 配置失败: %s\n", reloadOutput)
			return
		}
	}

	// 启动服务
	fmt.Printf("正在启动Web服务（端口: %d）...\n", port)
	utils.SystemctlStart(constants.WebServiceName)
}

func stopWebService() {
	// 停止服务
	fmt.Println("正在停止Web服务...")
	utils.SystemctlStop(constants.WebServiceName)
}

func restartWebService(port int) {
	// 重新安装服务（使用新端口）
	fmt.Printf("正在重新安装Web服务（端口: %d）...\n", port)
	if err := installWebService(port); err != nil {
		fmt.Printf("重新安装Web服务失败: %v\n", err)
		return
	}

	// 重新加载 systemd 配置
	fmt.Println("重新加载 systemd 配置...")
	utils.ExecCommandWithResult("systemctl daemon-reload")

	// 重启服务
	fmt.Println("正在重启Web服务...")
	utils.SystemctlRestart(constants.WebServiceName)
}

func checkWebServiceStatus() {
	// 获取服务状态
	fmt.Println("=== Web 服务状态 ===")
	statusOutput := utils.SystemctlStatus(constants.WebServiceName)
	if statusOutput != "" {
		fmt.Printf("%s\n", statusOutput)
	} else {
		fmt.Println("无法获取服务状态")
	}

	fmt.Println("\n按回车键返回...")
	fmt.Scanln()
}

func showWebServiceLogs() {
	// 使用journalctl查看systemd服务日志
	output := utils.ExecCommandWithResult(fmt.Sprintf("journalctl -u %s --no-pager -n 50", constants.WebServiceName))
	if output == "" {
		fmt.Println("获取服务日志失败或日志为空")
		fmt.Println("\n按回车键返回...")
		fmt.Scanln()
		return
	}

	// 显示日志内容
	fmt.Println("\n=== Web 服务日志 (最近50行) ===")
	fmt.Println(output)
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
	webServerCmd.Flags().IntVarP(&webPort, "port", "p", 8085, "Web 服务器监听的端口")
	rootCmd.AddCommand(webServerCmd)
}
