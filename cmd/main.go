package cmd

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	"openvpn-admin-go/openvpn"
	"openvpn-admin-go/utils"
)

// CoreInitializer is a function variable to hold the core initialization logic from the main package.
var CoreInitializer func() error

// WebServerStarter is a function variable to hold the web server starting logic from the main package.
var WebServerStarter func()

// webCmd represents the web command
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "Start the REST API web server",
	Run: func(cmd *cobra.Command, args []string) {
		if CoreInitializer == nil || WebServerStarter == nil {
			log.Fatal("Core services not initialized. Please ensure the main application sets them up.")
		}
		log.Println("Initializing core services for web server...")
		if err := CoreInitializer(); err != nil {
			log.Fatalf("Core initialization failed: %v", err)
		}
		log.Println("Starting web server...")
		WebServerStarter()
		log.Println("Web server started. Command will remain active to keep the server running.")
		// Block indefinitely to keep the command running
		select {}
	},
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the service (Placeholder)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("start command called (Full system service functionality not yet implemented.)")
	},
}

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop the service (Placeholder)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("stop command called (Full system service functionality not yet implemented.)")
	},
}

// restartCmd represents the restart command
var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the service (Placeholder)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("restart command called (Full system service functionality not yet implemented.)")
	},
}

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
	menuItems := []string{
		"服务器管理",
		"客户端管理",
		"退出",
	}

	prompt := promptui.Select{
		Label: "请选择操作",
		Items: menuItems,
		Size:  len(menuItems),
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "➤ {{ . | cyan }}",
			Inactive: "  {{ . | white }}",
			Selected: "{{ . | green }}",
		},
		HideSelected: true,
	}

	for {
		_, result, err := prompt.Run()
		if err != nil {
			if strings.Contains(err.Error(), "^C") {
				fmt.Println("\n程序已退出")
				os.Exit(0)
			}
			fmt.Printf("选择失败: %v\n", err)
			return
		}

		switch result {
		case "服务器管理":
			ServerMenu()
		case "客户端管理":
			ClientMenu()
		case "退出":
			return
		}
	}
}

func Execute() {
	rootCmd.AddCommand(webCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(logCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
