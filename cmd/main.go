package cmd

import (
	"fmt"
	"os"
	// "os/signal"
	// "strings"
	// "syscall"

	// "github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
	// "openvpn-admin-go/openvpn"
)

var openvpnCmd = &cobra.Command{
	Use:   "openvpn",
	Short: "OpenVPN management tool",
	// Run: func(cmd *cobra.Command, args []string) {
	// // 设置 Ctrl+C 处理
	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	// go func() {
	// 	<-sigChan
	// 	fmt.Println("\n程序已退出")
	// 	os.Exit(0)
	// }()

	// // 加载配置
	// cfg, err := openvpn.LoadConfig()
	// if err != nil {
	// 	fmt.Printf("加载配置失败: %v\n", err)
	// 	return
	// }

	// // 显示主菜单
	// MainMenu(cfg)
	// },
}

// func MainMenu(cfg *openvpn.Config) {
// 	menuItems := []string{
// 		"服务器管理",
// 		"客户端管理",
// 		"退出",
// 	}

// 	prompt := promptui.Select{
// 		Label: "请选择操作",
// 		Items: menuItems,
// 		Size:  len(menuItems),
// 		Templates: &promptui.SelectTemplates{
// 			Label:    "{{ . }}",
// 			Active:   "➤ {{ . | cyan }}",
// 			Inactive: "  {{ . | white }}",
// 			Selected: "{{ . | green }}",
// 		},
// 		HideSelected: true,
// 	}

// 	for {
// 		_, result, err := prompt.Run()
// 		if err != nil {
// 			if strings.Contains(err.Error(), "^C") {
// 				fmt.Println("\n程序已退出")
// 				os.Exit(0)
// 			}
// 			fmt.Printf("选择失败: %v\n", err)
// 			return
// 		}

// 		switch result {
// 		case "服务器管理":
// 			ServerMenu()
// 		case "客户端管理":
// 			ClientMenu()
// 		case "退出":
// 			return
// 		}
// 	}
// }

func Execute() {
	if err := openvpnCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Define subcommands
	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start OpenVPN server",
		Run: func(cmd *cobra.Command, args []string) {
			startService()
		},
	}

	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop OpenVPN server",
		Run: func(cmd *cobra.Command, args []string) {
			stopService()
		},
	}

	var restartCmd = &cobra.Command{
		Use:   "restart",
		Short: "Restart OpenVPN server",
		Run: func(cmd *cobra.Command, args []string) {
			restartService()
		},
	}

	var logCmd = &cobra.Command{
		Use:   "log",
		Short: "View OpenVPN logs (e.g., `log -f` to follow, `log status` for status log)",
		Run: func(cmd *cobra.Command, args []string) {
			showLogs(args)
		},
	}

	var webCmd = &cobra.Command{
		Use:   "web",
		Short: "Start web UI",
		Run: func(cmd *cobra.Command, args []string) {
			startWebServer()
		},
	}

	// Add subcommands to root command
	openvpnCmd.AddCommand(startCmd, stopCmd, restartCmd, logCmd, webCmd)
}
