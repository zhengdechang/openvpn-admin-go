package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"openvpn-admin-go/cmd"
	"bufio"
	"strings"
)

// loadEnv 从.env文件加载环境变量
func loadEnv() error {
	file, err := os.Open(".env") 
	if err != nil {
		return fmt.Errorf("无法打开.env文件: %v", err)
	}
	defer file.Close()

	requiredVars := map[string]bool{
		"OPENVPN_PORT":            true,
		"OPENVPN_PROTO":           true,
		"OPENVPN_SERVER_NETWORK":  true,
		"OPENVPN_SERVER_NETMASK":  true,
		"OPENVPN_SERVER_HOSTNAME": true,
		"OPENVPN_SERVER_IP":       true,
	}

	loadedVars := make(map[string]bool)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// 跳过注释和空行
		if strings.HasPrefix(line, "#") || len(strings.TrimSpace(line)) == 0 {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// 设置环境变量
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("设置环境变量失败 %s: %v", key, err)
		}

		loadedVars[key] = true
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取.env文件失败: %v", err)
	}

	// 检查必需的环境变量是否都已加载
	var missingVars []string
	for varName := range requiredVars {
		if !loadedVars[varName] {
			missingVars = append(missingVars, varName)
		}
	}

	if len(missingVars) > 0 {
		return fmt.Errorf("缺少必需的环境变量: %s", strings.Join(missingVars, ", "))
	}

	return nil
}

func main() {
	// 加载环境变量
	if err := loadEnv(); err != nil {
		log.Fatalf("错误: 加载环境变量失败: %v\n请确保.env文件存在且包含所有必需的配置项。", err)
	}

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