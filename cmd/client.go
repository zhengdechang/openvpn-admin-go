package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"openvpn-admin-go/constants"
	"openvpn-admin-go/logging"
	"openvpn-admin-go/openvpn"

	"github.com/manifoldco/promptui"
)

func ClientMenu() {
	for {
		fmt.Println("\n=== 客户端管理 ===\n")
		fmt.Println("1. 创建客户端\n")
		fmt.Println("2. 删除客户端\n")
		fmt.Println("3. 暂停客户端\n")
		fmt.Println("4. 恢复客户端\n")
		fmt.Println("5. 查看客户端状态\n")
		fmt.Println("6. 查看所有客户端\n")
		fmt.Println("0. 返回主菜单\n")
		fmt.Print("请选择操作 (0-6): ")

		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Printf("读取输入失败: %v\n", err)
			continue
		}

		// 移除输入中的空白字符
		input = strings.TrimSpace(input)
		if input == "" {
			continue
		}

		// 检查是否是Ctrl+C
		if input == "^C" {
			fmt.Println("\n程序已退出")
			os.Exit(0)
		}

		choice, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("无效的选择，请输入数字")
			continue
		}

		switch choice {
		case 0:
			return
		case 1:
			if err := CreateClient(); err != nil {
				fmt.Printf("创建客户端失败: %v\n", err)
			}
		case 2:
			DeleteClient()
		case 3:
			PauseClient()
		case 4:
			ResumeClient()
		case 5:
			ViewClientStatus()
		case 6:
			ListClients()
		default:
			fmt.Println("无效的选择，请输入0-6之间的数字")
		}
	}
}

func CreateClient() error {
	// 获取用户名
	username, err := getUsername()
	if err != nil {
		return err
	}

	// 调用openvpn包中的函数创建客户端
	if err := openvpn.CreateClient(username); err != nil {
		return fmt.Errorf("创建客户端失败: %v", err)
	}

	return nil
}

func DeleteClient() {
	username, err := getUsername()
	if err != nil {
		logging.Error("获取用户名失败: %v", err)
		return
	}

	if err := openvpn.DeleteClient(username); err != nil {
		logging.Error("删除客户端失败: %v", err)
	} else {
		fmt.Printf("客户端 %s 删除成功\n", username)
	}
}

func PauseClient() {
	username, err := getUsername()
	if err != nil {
		logging.Error("获取用户名失败: %v", err)
		return
	}

	if err := openvpn.PauseClient(username); err != nil {
		logging.Error("暂停客户端失败: %v", err)
	} else {
		fmt.Printf("客户端 %s 已暂停\n", username)
	}
}

func ResumeClient() {
	username, err := getUsername()
	if err != nil {
		logging.Error("获取用户名失败: %v", err)
		return
	}

	if err := openvpn.ResumeClient(username); err != nil {
		logging.Error("恢复客户端失败: %v", err)
	} else {
		fmt.Printf("客户端 %s 已恢复\n", username)
	}
}

func ViewClientStatus() {
	username, err := getUsername()
	if err != nil {
		logging.Error("获取用户名失败: %v", err)
		return
	}

	status, err := openvpn.GetClientStatus(username)
	if err != nil {
		logging.Error("获取客户端状态失败: %v", err)
		return
	}

	if status == nil {
		fmt.Printf("客户端 %s 不存在\n", username)
		return
	}

	fmt.Printf("客户端 %s 状态:\n", status.CommonName)
	fmt.Printf("连接时间: %s\n", status.ConnectedSince.Format("2006-01-02 15:04:05"))
	fmt.Printf("最后活动: %s\n", status.LastRef.Format("2006-01-02 15:04:05"))
	fmt.Printf("接收字节: %d\n", status.BytesReceived)
	fmt.Printf("发送字节: %d\n", status.BytesSent)
	fmt.Printf("虚拟地址: %s\n", status.VirtualAddress)
	fmt.Printf("真实地址: %s\n", status.RealAddress)
}

func ListClients() {
	statuses, err := openvpn.GetAllClientStatuses()
	if err != nil {
		logging.Error("获取客户端列表失败: %v", err)
		return
	}

	if len(statuses) == 0 {
		fmt.Println("没有找到任何客户端")
		return
	}

	fmt.Println("客户端列表:")
	fmt.Println("----------------------------------------")
	fmt.Printf("%-20s %-20s %-15s %-15s\n", "用户名", "连接时间", "虚拟地址", "真实地址")
	fmt.Println("----------------------------------------")

	for _, status := range statuses {
		fmt.Printf("%-20s %-20s %-15s %-15s\n",
			status.CommonName,
			status.ConnectedSince.Format("2006-01-02 15:04:05"),
			status.VirtualAddress,
			status.RealAddress)
	}
	fmt.Println("----------------------------------------")
}

func getUsername() (string, error) {
	prompt := promptui.Prompt{
		Label: "请输入用户名",
		Validate: func(input string) error {
			if len(strings.TrimSpace(input)) == 0 {
				return fmt.Errorf("用户名不能为空")
			}
			return nil
		},
	}

	username, err := prompt.Run()
	if err != nil {
		fmt.Printf("输入失败: %v\n", err)
		return "", err
	}

	return username, nil
}

func getServerIP() string {
	// 读取服务器配置文件获取IP
	config, err := os.ReadFile(constants.ServerConfigPath)
	if err != nil {
		return "your-server-ip" // 默认值
	}

	lines := strings.Split(string(config), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "server ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				return parts[1]
			}
		}
	}

	return "your-server-ip" // 默认值
}

func readFile(path string) string {
	content, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(content)
}

// 读取服务器配置
func readServerConfig() (string, error) {
	// 检查服务端配置文件是否存在
	if _, err := os.Stat(constants.ServerConfigPath); os.IsNotExist(err) {
		return "", fmt.Errorf("服务端配置文件不存在: %s", constants.ServerConfigPath)
	}

	// 读取服务器配置
	config, err := os.ReadFile(constants.ServerConfigPath)
	if err != nil {
		return "", fmt.Errorf("读取服务器配置失败: %v", err)
	}
	return string(config), nil
}

// 检查客户端配置目录
func checkClientConfigDir() error {
	// 检查客户端配置目录
	if _, err := os.Stat(constants.ClientConfigDir); os.IsNotExist(err) {
		return fmt.Errorf("客户端配置目录不存在: %s", constants.ClientConfigDir)
	}
	return nil
}

// 创建客户端配置目录
func createClientConfigDir() error {
	// 创建客户端配置目录
	if err := os.MkdirAll(constants.ClientConfigDir, 0755); err != nil {
		return fmt.Errorf("创建客户端配置目录失败: %v", err)
	}
	return nil
}

// 写入客户端配置文件
func writeClientConfig(username, config string) error {
	// 检查客户端配置目录
	if err := checkClientConfigDir(); err != nil {
		return err
	}

	// 写入客户端配置文件
	configPath := filepath.Join(constants.ClientConfigDir, username+".ovpn")
	if err := os.WriteFile(configPath, []byte(config), 0644); err != nil {
		return fmt.Errorf("写入客户端配置文件失败: %v", err)
	}
	return nil
}

// 删除客户端配置文件
func deleteClientConfig(username string) error {
	// 删除客户端配置文件
	configPath := filepath.Join(constants.ClientConfigDir, username+".ovpn")
	if err := os.Remove(configPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("删除客户端配置文件失败: %v", err)
	}
	return nil
}

// 检查客户端配置文件是否存在
func checkClientConfig(username string) error {
	// 检查客户端配置文件是否存在
	configPath := filepath.Join(constants.ClientConfigDir, username+".ovpn")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return fmt.Errorf("客户端配置文件不存在: %s", configPath)
	}

	// 检查客户端证书文件是否存在
	clientCertPath := filepath.Join(constants.ClientConfigDir, username+".crt")
	if _, err := os.Stat(clientCertPath); os.IsNotExist(err) {
		return fmt.Errorf("客户端证书文件不存在: %s", clientCertPath)
	}

	// 检查客户端密钥文件是否存在
	clientKeyPath := filepath.Join(constants.ClientConfigDir, username+".key")
	if _, err := os.Stat(clientKeyPath); os.IsNotExist(err) {
		return fmt.Errorf("客户端密钥文件不存在: %s", clientKeyPath)
	}

	return nil
}
