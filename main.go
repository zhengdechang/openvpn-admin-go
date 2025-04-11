package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"openvpn-admin-go/cmd"
)

func main() {
	// 设置Ctrl+C信号处理
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n程序已退出")
		os.Exit(0)
	}()

	fmt.Println("OpenVPN 管理工具启动中...")
	
	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("获取工作目录失败: %v", err)
	}
	fmt.Printf("当前工作目录: %s\n", dir)

	// 检查环境
	fmt.Println("正在检查运行环境...")
	if err := cmd.CheckEnvironment(); err != nil {
		fmt.Printf("环境检查失败: %v\n", err)
		fmt.Println("是否自动安装所需环境？(y/n)")
		var choice string
		fmt.Scanln(&choice)
		if choice == "y" || choice == "Y" {
			if err := cmd.InstallEnvironment(); err != nil {
				fmt.Printf("环境安装失败: %v\n", err)
				fmt.Println("请确保您有足够的权限，软件源配置正确，网络连接稳定。")
				fmt.Println("您可以手动检查并修复问题后重新运行程序。")
				return
			}
			// 重新检查环境
			if err := cmd.CheckEnvironment(); err != nil {
				fmt.Printf("环境检查仍然失败: %v\n", err)
				fmt.Println("请手动检查并修复问题后重新运行程序。")
				return
			}
		} else {
			fmt.Println("请手动安装所需环境后重新运行程序。")
			return
		}
	}

	fmt.Println("环境检查通过")
	
	// 启动主菜单
	cmd.Execute()
}
