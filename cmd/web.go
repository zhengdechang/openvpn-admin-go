package cmd

import (
	"fmt"
	"openvpn-admin-go/constants"
	"openvpn-admin-go/logging"
	"openvpn-admin-go/router"
	"os/exec"
	"time"

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
	// 检查服务状态
	cmd := exec.Command("systemctl", "is-active", constants.WebServiceName)
	if err := cmd.Run(); err == nil {
		fmt.Println("Web 服务已在运行")
		return
	}

	// 启动服务
	cmd = exec.Command("systemctl", "start", constants.WebServiceName)
	if err := cmd.Run(); err != nil {
		fmt.Printf("启动服务失败: %v\n请检查服务状态: systemctl status %s\n", err, constants.WebServiceName)
		return
	}

	// 等待服务启动
	time.Sleep(2 * time.Second)

	// 检查服务状态
	cmd = exec.Command("systemctl", "is-active", constants.WebServiceName)
	if err := cmd.Run(); err != nil {
		fmt.Printf("服务未正常运行: %v\n请检查服务日志: journalctl -u %s\n", err, constants.WebServiceName)
		return
	}

	fmt.Println("Web 服务启动成功")
}

func stopWebService() {
	// 检查服务状态
	cmd := exec.Command("systemctl", "is-active", constants.WebServiceName)
	if err := cmd.Run(); err != nil {
		fmt.Println("Web 服务未运行")
		return
	}

	// 停止服务
	cmd = exec.Command("systemctl", "stop", constants.WebServiceName)
	if err := cmd.Run(); err != nil {
		fmt.Printf("停止服务失败: %v\n请检查服务状态: systemctl status %s\n", err, constants.WebServiceName)
		return
	}

	// 等待服务完全停止
	time.Sleep(2 * time.Second)

	// 验证服务已停止
	cmd = exec.Command("systemctl", "is-active", constants.WebServiceName)
	if err := cmd.Run(); err == nil {
		fmt.Println("服务仍在运行")
		return
	}

	fmt.Println("Web 服务已停止")
}

func restartWebService() {
	// 重启服务
	cmd := exec.Command("systemctl", "restart", constants.WebServiceName)
	if err := cmd.Run(); err != nil {
		fmt.Printf("重启服务失败: %v\n请检查服务状态: systemctl status %s\n", err, constants.WebServiceName)
		return
	}

	// 等待服务启动
	time.Sleep(2 * time.Second)

	// 检查服务状态
	cmd = exec.Command("systemctl", "is-active", constants.WebServiceName)
	if err := cmd.Run(); err != nil {
		fmt.Printf("服务未正常运行: %v\n请检查服务日志: journalctl -u %s\n", err, constants.WebServiceName)
		return
	}

	fmt.Println("Web 服务重启成功")
}

func checkWebServiceStatus() {
	// 检查服务状态
	cmd := exec.Command("systemctl", "is-active", constants.WebServiceName)
	if err := cmd.Run(); err != nil {
		fmt.Println("Web 服务未运行")
	} else {
		fmt.Println("Web 服务正在运行")

		// 获取服务详细状态
		cmd = exec.Command("systemctl", "status", constants.WebServiceName, "--no-pager", "-l")
		if output, err := cmd.CombinedOutput(); err == nil {
			fmt.Printf("\n服务状态详情:\n%s\n", string(output))
		}
	}
	fmt.Println("\n按回车键返回...")
	fmt.Scanln()
}

func showWebServiceLogs() {
	// 使用journalctl查看systemd服务日志
	cmd := exec.Command("journalctl", "-u", constants.WebServiceName, "--no-pager", "-n", "50")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("获取服务日志失败: %v\n", err)
		fmt.Println("\n按回车键返回...")
		fmt.Scanln()
		return
	}

	// 显示日志内容
	fmt.Println("\n=== Web 服务日志 (最近50行) ===")
	fmt.Println(string(output))
	fmt.Println("\n按回车键返回...")
	fmt.Scanln()
}

// Web Server Functions - 实际的 Web 服务器实现

// runWebServer 实际的 Web 服务器实现
func runWebServer(port int) error {
	// 初始化核心服务
	if err := CoreInitializer(); err != nil {
		return fmt.Errorf("核心初始化失败: %v", err)
	}

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



func init() {
	// 添加 web-server 子命令
	webServerCmd := &cobra.Command{
		Use:   "web-server",
		Short: "运行 Web 服务器",
		Run: func(cmd *cobra.Command, args []string) {
			if err := runWebServer(webPort); err != nil {
				logging.Fatal("Web 服务器错误: %v", err)
			}
		},
	}
	webServerCmd.Flags().IntVarP(&webPort, "port", "p", 8085, "Web 服务器监听的端口")
	rootCmd.AddCommand(webServerCmd)
}
