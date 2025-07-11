package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"openvpn-admin-go/constants"
	"openvpn-admin-go/database"
	"openvpn-admin-go/logging"
	"openvpn-admin-go/model"
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
	// 先显示所有客户端列表
	fmt.Println("=== 当前客户端列表 ===")
	showClientList()

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
	// 先显示所有客户端列表
	fmt.Println("=== 当前客户端列表 ===")
	showClientList()

	username, err := getUsername()
	if err != nil {
		logging.Error("获取用户名失败: %v", err)
		return
	}

	// 首先从数据库查找用户
	var user model.User
	if err := database.DB.Where("name = ?", username).First(&user).Error; err != nil {
		fmt.Printf("数据库中未找到用户 %s: %v\n", username, err)
		return
	}

	// 暂停OpenVPN客户端
	if err := openvpn.PauseClient(username); err != nil {
		logging.Error("暂停OpenVPN客户端失败: %v", err)
		return
	}

	// 更新数据库状态
	user.IsPaused = true
	if err := database.DB.Save(&user).Error; err != nil {
		logging.Error("更新数据库状态失败: %v", err)
		// 尝试恢复OpenVPN状态
		openvpn.ResumeClient(username)
		return
	}

	fmt.Printf("客户端 %s 已暂停\n", username)
}

func ResumeClient() {
	// 先显示所有客户端列表
	fmt.Println("=== 当前客户端列表 ===")
	showClientList()

	username, err := getUsername()
	if err != nil {
		logging.Error("获取用户名失败: %v", err)
		return
	}

	// 首先从数据库查找用户
	var user model.User
	if err := database.DB.Where("name = ?", username).First(&user).Error; err != nil {
		fmt.Printf("数据库中未找到用户 %s: %v\n", username, err)
		return
	}

	// 恢复OpenVPN客户端
	if err := openvpn.ResumeClient(username); err != nil {
		logging.Error("恢复OpenVPN客户端失败: %v", err)
		return
	}

	// 更新数据库状态
	user.IsPaused = false
	if err := database.DB.Save(&user).Error; err != nil {
		logging.Error("更新数据库状态失败: %v", err)
		// 尝试重新暂停OpenVPN状态
		openvpn.PauseClient(username)
		return
	}

	fmt.Printf("客户端 %s 已恢复\n", username)
}

func ViewClientStatus() {
	// 先显示所有客户端列表
	fmt.Println("=== 当前客户端列表 ===")
	showClientList()

	username, err := getUsername()
	if err != nil {
		logging.Error("获取用户名失败: %v", err)
		return
	}

	// 首先从数据库查找用户
	var user model.User
	if err := database.DB.Where("name = ?", username).First(&user).Error; err != nil {
		fmt.Printf("数据库中未找到用户 %s: %v\n", username, err)
		return
	}

	fmt.Printf("=== 客户端 %s 详细状态 ===\n", username)
	fmt.Printf("用户ID: %s\n", user.ID)
	fmt.Printf("邮箱: %s\n", user.Email)
	fmt.Printf("角色: %s\n", user.Role)
	fmt.Printf("部门ID: %s\n", user.DepartmentID)
	fmt.Printf("固定IP: %s\n", user.FixedIP)
	fmt.Printf("子网: %s\n", user.Subnet)
	fmt.Printf("是否暂停: %t\n", user.IsPaused)
	fmt.Printf("创建时间: %s\n", user.CreatedAt.Format("2006-01-02 15:04:05"))

	// 获取OpenVPN连接状态
	status, err := openvpn.GetClientStatus(username)
	if err != nil {
		fmt.Printf("获取OpenVPN状态失败: %v\n", err)
		fmt.Println("--- OpenVPN连接状态: 离线 ---")
		return
	}

	if status == nil {
		fmt.Println("--- OpenVPN连接状态: 离线 ---")
		return
	}

	fmt.Println("--- OpenVPN连接状态: 在线 ---")
	fmt.Printf("连接时间: %s\n", status.ConnectedSince.Format("2006-01-02 15:04:05"))
	fmt.Printf("最后活动: %s\n", status.LastRef.Format("2006-01-02 15:04:05"))
	fmt.Printf("接收字节: %d\n", status.BytesReceived)
	fmt.Printf("发送字节: %d\n", status.BytesSent)
	fmt.Printf("虚拟地址: %s\n", status.VirtualAddress)
	fmt.Printf("真实地址: %s\n", status.RealAddress)
}

func ListClients() {
	// 首先从数据库获取所有用户
	fmt.Println("=== 所有客户端列表 ===")

	// 获取数据库中的所有用户
	users, err := getAllUsersFromDB()
	if err != nil {
		logging.Error("获取数据库用户列表失败: %v", err)
		return
	}

	// 获取OpenVPN状态日志中的连接状态
	liveStatuses, err := openvpn.GetAllClientStatuses()
	if err != nil {
		logging.Warn("获取OpenVPN状态失败: %v", err)
		liveStatuses = []openvpn.ClientStatus{} // 继续显示数据库中的用户
	}

	// 创建状态映射
	statusMap := make(map[string]openvpn.ClientStatus)
	for _, status := range liveStatuses {
		statusMap[status.CommonName] = status
	}

	if len(users) == 0 {
		fmt.Println("数据库中没有找到任何用户")
		return
	}

	fmt.Println("----------------------------------------")
	fmt.Printf("%-20s %-10s %-15s %-15s %-20s\n", "用户名", "状态", "虚拟地址", "真实地址", "连接时间")
	fmt.Println("----------------------------------------")

	for _, user := range users {
		status := "离线"
		virtualAddr := "-"
		realAddr := "-"
		connectedTime := "-"

		if liveStatus, exists := statusMap[user.Name]; exists {
			status = "在线"
			virtualAddr = liveStatus.VirtualAddress
			realAddr = liveStatus.RealAddress
			if !liveStatus.ConnectedSince.IsZero() {
				connectedTime = liveStatus.ConnectedSince.Format("15:04:05")
			}
		} else if user.IsPaused {
			status = "已暂停"
		}

		fmt.Printf("%-20s %-10s %-15s %-15s %-20s\n",
			user.Name,
			status,
			virtualAddr,
			realAddr,
			connectedTime)
	}
	fmt.Println("----------------------------------------")
	fmt.Printf("总计: %d 个用户，其中 %d 个在线\n", len(users), len(liveStatuses))
}

// getAllUsersFromDB 从数据库获取所有用户
func getAllUsersFromDB() ([]model.User, error) {
	var users []model.User
	if err := database.DB.Order("created_at desc").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("查询用户失败: %v", err)
	}
	return users, nil
}

// showClientList 显示简化的客户端列表
func showClientList() {
	// 获取数据库中的所有用户
	users, err := getAllUsersFromDB()
	if err != nil {
		fmt.Printf("获取用户列表失败: %v\n", err)
		return
	}

	// 获取OpenVPN状态日志中的连接状态
	liveStatuses, err := openvpn.GetAllClientStatuses()
	if err != nil {
		liveStatuses = []openvpn.ClientStatus{} // 继续显示数据库中的用户
	}

	// 创建状态映射
	statusMap := make(map[string]openvpn.ClientStatus)
	for _, status := range liveStatuses {
		statusMap[status.CommonName] = status
	}

	if len(users) == 0 {
		fmt.Println("没有找到任何用户")
		return
	}

	fmt.Printf("%-20s %-10s %-15s\n", "用户名", "状态", "虚拟地址")
	fmt.Println("----------------------------------------")

	for _, user := range users {
		status := "离线"
		virtualAddr := "-"

		if liveStatus, exists := statusMap[user.Name]; exists {
			status = "在线"
			virtualAddr = liveStatus.VirtualAddress
		} else if user.IsPaused {
			status = "已暂停"
		}

		fmt.Printf("%-20s %-10s %-15s\n",
			user.Name,
			status,
			virtualAddr)
	}
	fmt.Println("----------------------------------------")
	fmt.Printf("总计: %d 个用户\n", len(users))
	fmt.Println()
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
