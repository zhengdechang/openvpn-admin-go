package cmd

import (
	"bytes"
	"fmt"
	"openvpn-admin-go/constants"
	"openvpn-admin-go/logging"
	"openvpn-admin-go/router"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
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

	// 创建openvpn-admin用户和组
	cmd := exec.Command("id", "openvpn-admin")
	if err := cmd.Run(); err != nil {
		// 用户不存在，创建用户
		cmd = exec.Command("useradd", "--system", "--no-create-home", "--shell", "/bin/false", "openvpn-admin")
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("创建openvpn-admin用户失败: %v", err)
		}
	}

	// 重新加载systemd配置
	cmd = exec.Command("systemctl", "daemon-reload")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("重新加载systemd配置失败: %v", err)
	}

	// 启用服务
	cmd = exec.Command("systemctl", "enable", constants.WebServiceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("启用web服务失败: %v", err)
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
	}

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

func restartWebService(port int) {
	// 先停止服务
	fmt.Println("正在停止Web服务...")
	cmd := exec.Command("systemctl", "stop", constants.WebServiceName)
	cmd.Run() // 忽略错误，因为服务可能未运行

	// 等待服务完全停止
	time.Sleep(2 * time.Second)

	// 重新安装服务（使用新端口）
	fmt.Printf("正在重新安装Web服务（端口: %d）...\n", port)
	if err := installWebService(port); err != nil {
		fmt.Printf("重新安装Web服务失败: %v\n", err)
		return
	}

	// 启动服务
	fmt.Println("正在启动Web服务...")
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

	fmt.Printf("Web 服务重启成功（端口: %d）\n", port)
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
		Use:   "web",
		Short: "运行 Web 服务器",
		Run: func(cmd *cobra.Command, args []string) {
			// 检查是否为开发模式
			if isDev := os.Getenv("DEV"); isDev == "true" || isDev == "1" {
				// 开发模式：直接运行web服务器
				if err := runWebServer(webPort); err != nil {
					logging.Fatal("Web 服务器错误: %v", err)
				}
			} else {
				// 生产模式：使用systemd服务
				startWebService(webPort)
			}
		},
	}
	webServerCmd.Flags().IntVarP(&webPort, "port", "p", 8085, "Web 服务器监听的端口")
	rootCmd.AddCommand(webServerCmd)
}
