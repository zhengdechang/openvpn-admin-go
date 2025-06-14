package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
	"openvpn-admin-go/openvpn"
	"openvpn-admin-go/utils"
)

// CoreInitializer is a function variable to hold the core initialization logic from the main package.
var CoreInitializer func() error

// logCmd represents the log command
var logCmd = &cobra.Command{
	Use:   "log",
	Short: "Display OpenVPN status logs",
	Run: func(cmd *cobra.Command, args []string) {
		logPath := utils.GetOpenVPNStatusLogPath()
		if _, err := os.Stat(logPath); os.IsNotExist(err) {
			fmt.Printf("Log file not found at %s\n", logPath)
			return
		}

		content, err := os.ReadFile(logPath)
		if err != nil {
			log.Fatalf("Error reading log file %s: %v", logPath, err)
		}
		fmt.Println(string(content))
	},
}

var rootCmd = &cobra.Command{
	Use:   "openvpn-admin",
	Short: "OpenVPN 管理工具",
	Run: func(cmd *cobra.Command, args []string) {
		// 设置 Ctrl+C 处理
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		go func() {
			<-sigChan
			fmt.Println("\n程序已退出")
			os.Exit(0)
		}()

		// 加载配置
		cfg, err := openvpn.LoadConfig()
		if err != nil {
			fmt.Printf("加载配置失败: %v\n", err)
			return
		}

		// 显示主菜单
		MainMenu(cfg)
	},
}

func MainMenu(cfg *openvpn.Config) {
	for {
		fmt.Println("\n=== 欢迎使用 OpenVPN 管理程序 ===\n")
		fmt.Println("1.服务器管理    2.客户端管理 \n")
		fmt.Println("3.Web服务管理   4.查看配置 \n")
		fmt.Println("5.查看日志      0.退出程序 \n")
		fmt.Print("请选择操作 (0-5): ")

		var choice string
		fmt.Scanln(&choice)

		switch choice {
		case "0":
			fmt.Println("再见!")
			return
		case "1":
			ServerMenu()
		case "2":
			ClientMenu()
		case "3":
			WebMenu()
		case "4":
			ShowConfig(cfg)
		case "5":
			logCmd.Run(nil, nil)
		default:
			fmt.Println("无效选择，请重试")
		}
	}
}

func ShowConfig(cfg *openvpn.Config) {
	fmt.Println("\n当前配置:")
	fmt.Printf("服务器地址: %s\n", cfg.OpenVPNServerHostname)
	fmt.Printf("服务器端口: %d\n", cfg.OpenVPNPort)
	fmt.Printf("协议: %s\n", cfg.OpenVPNProto)
	fmt.Printf("服务器网络: %s\n", cfg.OpenVPNServerNetwork)
	fmt.Printf("子网掩码: %s\n", cfg.OpenVPNServerNetmask)
	fmt.Printf("客户端配置目录: %s\n", cfg.OpenVPNClientConfigDir)
	fmt.Printf("TLS版本: %s\n", cfg.OpenVPNTLSVersion)
	fmt.Printf("状态日志路径: %s\n", cfg.OpenVPNStatusLogPath)
	fmt.Printf("OpenVPN日志路径: %s\n", cfg.OpenVPNLogPath)
	fmt.Println("\n按回车键返回主菜单...")
	fmt.Scanln()
}



func Execute() {
	// webCmd is added to rootCmd in cmd/web.go's init()
	rootCmd.AddCommand(logCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
