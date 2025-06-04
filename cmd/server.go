package cmd

import (
	"fmt"
	"math/rand"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"openvpn-admin-go/constants"
	"openvpn-admin-go/openvpn"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "管理 OpenVPN 服务器",
	Run: func(cmd *cobra.Command, args []string) {
		ServerMenu()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}

// 统一使用标准服务名称
const serviceName = "openvpn-server@server.service"

func ServerMenu() {
	// 加载配置
	cfg, err := openvpn.LoadConfig()
	if err != nil {
		fmt.Printf("加载配置失败: %v\n", err)
		return
	}

	menuItems := []string{
		"启动服务器",
		"停止服务器",
		"重启服务器",
		"查看服务器状态",
		"更新服务器配置",
		"返回主菜单",
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
		case "启动服务器":
			startServer(cfg)
		case "停止服务器":
			stopServer()
		case "重启服务器":
			restartServer(cfg)
		case "查看服务器状态":
			checkServerStatus()
		case "更新服务器配置":
			if err := UpdateConfig(); err != nil {
				fmt.Printf("更新配置失败: %v\n", err)
			}
		case "返回主菜单":
			return
		}
	}
}

func startServer(cfg *openvpn.Config) {
	// 生成配置文件
	configContent, err := cfg.GenerateServerConfig()
	if err != nil {
		fmt.Printf("生成配置文件失败: %v\n", err)
		return
	}
	configPath := constants.ServerConfigPath
	fmt.Println(configPath)
	// 写入配置文件
	if err := WriteConfig(configContent); err != nil {
		fmt.Printf("写入配置文件失败: %v\n", err)
		return
	}

	// 启动服务
	cmd := exec.Command("sudo", "systemctl", "start", constants.ServiceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("启动失败: %v\n输出: %s\n", err, string(output))
		return
	}
	fmt.Println("✅ 服务已启动")
}

func stopServer() {
	cmd := exec.Command("sudo", "systemctl", "stop", constants.ServiceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("停止失败: %v\n输出: %s\n", err, string(output))
		return
	}
	fmt.Println("✅ 服务已停止")
}

func restartServer(cfg *openvpn.Config) error {
	// 停止服务
	cmd := exec.Command("systemctl", "stop", constants.ServiceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("停止服务失败: %v\n输出: %s", err, string(output))
	}

	// 等待服务完全停止
	time.Sleep(2 * time.Second)

	// 启动服务
	cmd = exec.Command("systemctl", "start", constants.ServiceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("启动服务失败: %v\n输出: %s", err, string(output))
	}

	return nil
}

func checkServerStatus() {
	cmd := exec.Command("sudo", "systemctl", "status", constants.ServiceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("获取状态失败: %v\n输出: %s\n", err, string(output))
		return
	}
	fmt.Println(string(output))
}

func updatePort(cfg *openvpn.Config) error {
	// 读取当前配置文件
	configPath := constants.ServerConfigPath
	configContent, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 获取当前端口
	currentPort := "1194" // 默认端口
	lines := strings.Split(string(configContent), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "port ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				currentPort = parts[1]
			}
			break
		}
	}

	// 生成随机端口 (1024-65535)
	rand.Seed(time.Now().UnixNano())
	randomPort := rand.Intn(64511) + 1024

	fmt.Printf("当前端口: %s\n", currentPort)
	fmt.Printf("随机端口: %d\n", randomPort)
	fmt.Print("请输入新的端口号 (直接按回车使用随机端口): ")

	var input string
	fmt.Scanln(&input)

	var newPort string
	if input == "" {
		newPort = strconv.Itoa(randomPort)
	} else {
		port, err := strconv.Atoi(input)
		if err != nil {
			return fmt.Errorf("输入失败: %v", err)
		}
		if port < 1 || port > 65535 {
			return fmt.Errorf("端口号必须在 1-65535 之间")
		}
		newPort = input
	}

	// 更新配置
	cfg.OpenVPNPort, _ = strconv.Atoi(newPort)

	// 保存配置后需要重新生成服务配置
	if err := openvpn.SaveConfig(cfg); err != nil {
		return fmt.Errorf("保存配置失败: %v", err)
	}

	// 生成新的服务端配置
	if err := openvpn.UpdatePort(cfg.OpenVPNPort); err != nil {
		return fmt.Errorf("更新端口失败: %v", err)
	}

	fmt.Printf("端口已更新为 %s\n", newPort)
	return nil
}

func updateServerIP(cfg *openvpn.Config) error {
	prompt := promptui.Prompt{
		Label: "请输入新的服务器地址",
		Validate: func(input string) error {
			if len(strings.TrimSpace(input)) == 0 {
				return fmt.Errorf("服务器地址不能为空")
			}
			return nil
		},
	}

	newIP, err := prompt.Run()
	if err != nil {
		if strings.Contains(err.Error(), "^C") {
			fmt.Println("\n操作已取消")
			return nil
		}
		fmt.Printf("输入失败: %v\n", err)
		return err
	}

	// 更新配置
	cfg.OpenVPNServerHostname = newIP

	// 保存配置后需要重新生成服务配置
	if err := openvpn.SaveConfig(cfg); err != nil {
		fmt.Printf("保存配置失败: %v\n", err)
		return err
	}

	// 添加配置重载
	reloadCmd := exec.Command("sudo", "systemctl", "daemon-reload")
	if output, err := reloadCmd.CombinedOutput(); err != nil {
		fmt.Printf("配置重载失败: %v\n输出: %s\n", err, string(output))
		return err
	}

	if err := openvpn.RestartServer(); err != nil {
		return fmt.Errorf("重启服务失败: %v", err)
	}
	fmt.Printf("服务器地址已更新为 %s\n", newIP)
	return nil
}

// updateServerIPAndMask 修改服务器IP和子网掩码
func updateServerIPAndMask() error {
	// 读取当前配置
	configContent, err := os.ReadFile(constants.ServerConfigPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 获取当前IP和子网掩码
	currentIP := "10.8.0.0"        // 默认IP
	currentMask := "255.255.255.0" // 默认子网掩码
	lines := strings.Split(string(configContent), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "server ") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				ipParts := strings.Split(parts[1], " ")
				if len(ipParts) >= 2 {
					currentIP = ipParts[0]
					currentMask = ipParts[1]
				}
			}
			break
		}
	}

	// 提示输入新IP和子网掩码（CIDR格式）
	fmt.Printf("当前服务器IP: %s/%s\n", currentIP, currentMask)
	fmt.Print("请输入新IP和子网掩码 (格式: 10.8.0.0/24): ")
	var input string
	fmt.Scanln(&input)

	// 解析CIDR格式
	parts := strings.Split(input, "/")
	if len(parts) != 2 {
		return fmt.Errorf("请输入有效的CIDR格式 (例如: 10.8.0.0/24)")
	}

	ip := parts[0]
	maskBits, err := strconv.Atoi(parts[1])
	if err != nil || maskBits < 0 || maskBits > 32 {
		return fmt.Errorf("请输入有效的子网掩码位数 (0-32)")
	}

	// 验证IP
	ipParts := strings.Split(ip, ".")
	if len(ipParts) != 4 {
		return fmt.Errorf("请输入有效的IP地址")
	}
	for _, part := range ipParts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 || num > 255 {
			return fmt.Errorf("请输入有效的IP地址")
		}
	}

	// 转换掩码位数到点分十进制
	mask := net.CIDRMask(maskBits, 32)
	maskStr := fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])

	// 更新配置文件
	lines = strings.Split(string(configContent), "\n")
	for i, line := range lines {
		if strings.HasPrefix(line, "server ") {
			lines[i] = fmt.Sprintf("server %s %s", ip, maskStr)
			break
		}
	}

	// 写入新配置
	if err := os.WriteFile(constants.ServerConfigPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	fmt.Printf("服务器IP和子网掩码已更新为: %s %s\n", ip, maskStr)
	return nil
}

// generateTLSKey 生成tls-auth密钥
func generateTLSKey(cfg *openvpn.Config) error {
	// 检查密钥文件是否已存在
	if _, err := os.Stat(constants.ServerTLSKeyPath); err == nil {
		fmt.Printf("TLS密钥文件已存在: %s\n", constants.ServerTLSKeyPath)
		return nil
	}

	// 生成tls-auth密钥
	cmd := exec.Command("openvpn", "--genkey", "secret", constants.ServerTLSKeyPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("生成TLS密钥失败: %v\n输出: %s", err, string(output))
	}

	// 设置适当的权限
	if err := os.Chmod(constants.ServerTLSKeyPath, 0600); err != nil {
		return fmt.Errorf("设置TLS密钥文件权限失败: %v", err)
	}

	fmt.Printf("TLS密钥已生成: %s\n", constants.ServerTLSKeyPath)
	return nil
}

// 更新服务器配置
func UpdateConfig() error {
	// 读取当前配置文件
	configContent, err := os.ReadFile(constants.ServerConfigPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 显示当前配置
	fmt.Println("\n当前配置:")
	fmt.Println(string(configContent))

	// 选择要修改的配置项
	prompt := promptui.Select{
		Label: "请选择要修改的配置项",
		Items: []string{
			"修改端口",
			"修改服务器地址",
			"修改服务器IP和子网掩码",
			"修改OpenVPN路由",
			"返回",
		},
		Templates: &promptui.SelectTemplates{
			Label:    "{{ . }}",
			Active:   "➤ {{ . | cyan }}",
			Inactive: "  {{ . | white }}",
			Selected: "{{ . | green }}",
		},
		HideSelected: true,
	}

	_, result, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("选择失败: %v", err)
	}

	// 加载配置
	cfg, err := openvpn.LoadConfig()
	if err != nil {
		return fmt.Errorf("加载配置失败: %v", err)
	}

	switch result {
	case "修改端口":
		if err := updatePort(cfg); err != nil {
			return fmt.Errorf("修改端口失败: %v", err)
		}
		return openvpn.UpdateServerConfig()
	case "修改服务器地址":
		if err := updateServerIP(cfg); err != nil {
			return fmt.Errorf("修改服务器地址失败: %v", err)
		}
		return openvpn.UpdateServerConfig()
	case "修改服务器IP和子网掩码":
		if err := updateServerIPAndMask(); err != nil {
			return fmt.Errorf("修改服务器IP和子网掩码失败: %v", err)
		}
		return openvpn.UpdateServerConfig()
	case "修改OpenVPN路由":
		if err := updateRoute(); err != nil {
			return fmt.Errorf("修改OpenVPN路由失败: %v", err)
		}
		return openvpn.UpdateServerConfig()
	case "返回":
		return nil
	}

	return nil
}

func updateRoute() error {
	// 读取当前配置
	configContent, err := os.ReadFile(constants.ServerConfigPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 显示当前路由配置
	fmt.Println("\n当前路由配置:")
	lines := strings.Split(string(configContent), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "push \"route ") {
			fmt.Println(line)
		}
	}

	// 选择操作
	prompt := promptui.Select{
		Label: "请选择操作",
		Items: []string{
			"添加路由",
			"删除路由",
			"返回",
		},
	}

	_, result, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("选择失败: %v", err)
	}

	switch result {
	case "添加路由":
		return addRoute(lines)
	case "删除路由":
		return deleteRoute(lines)
	case "返回":
		return nil
	}

	return nil
}

func addRoute(lines []string) error {
	// 提示输入新路由
	fmt.Print("请输入要添加的路由 (格式: 10.10.100.0/23,10.10.98.0/23): ")
	var input string
	fmt.Scanln(&input)

	// 分割多个路由
	routes := strings.Split(input, ",")
	var newRoutes []string

	for _, route := range routes {
		route = strings.TrimSpace(route)
		if route == "" {
			continue
		}

		// 解析CIDR格式
		parts := strings.Split(route, "/")
		if len(parts) != 2 {
			return fmt.Errorf("请输入有效的CIDR格式 (例如: 10.10.100.0/23)")
		}

		ip := parts[0]
		maskBits, err := strconv.Atoi(parts[1])
		if err != nil || maskBits < 0 || maskBits > 32 {
			return fmt.Errorf("请输入有效的子网掩码位数 (0-32)")
		}

		// 验证IP
		ipParts := strings.Split(ip, ".")
		if len(ipParts) != 4 {
			return fmt.Errorf("请输入有效的IP地址")
		}
		for _, part := range ipParts {
			num, err := strconv.Atoi(part)
			if err != nil || num < 0 || num > 255 {
				return fmt.Errorf("请输入有效的IP地址")
			}
		}

		// 转换掩码位数到点分十进制
		mask := net.CIDRMask(maskBits, 32)
		maskStr := fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])

		// 添加新路由
		routeLine := fmt.Sprintf("push \"route %s %s\"", ip, maskStr)
		newRoutes = append(newRoutes, routeLine)
	}

	// 添加新路由到配置
	lines = append(lines, newRoutes...)

	// 写入新配置
	if err := os.WriteFile(constants.ServerConfigPath, []byte(strings.Join(lines, "\n")), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	fmt.Println("路由已添加:")
	for _, route := range newRoutes {
		fmt.Println(route)
	}
	return nil
}

func deleteRoute(lines []string) error {
	// 收集所有路由
	var routes []string
	for _, line := range lines {
		if strings.HasPrefix(line, "push \"route ") {
			routes = append(routes, line)
		}
	}

	if len(routes) == 0 {
		fmt.Println("没有可删除的路由")
		return nil
	}

	// 选择要删除的路由
	prompt := promptui.Select{
		Label: "请选择要删除的路由",
		Items: routes,
	}

	index, _, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("选择失败: %v", err)
	}

	// 删除选中的路由
	var newLines []string
	routeIndex := 0
	for _, line := range lines {
		if strings.HasPrefix(line, "push \"route ") {
			if routeIndex != index {
				newLines = append(newLines, line)
			}
			routeIndex++
		} else {
			newLines = append(newLines, line)
		}
	}

	// 写入新配置
	if err := os.WriteFile(constants.ServerConfigPath, []byte(strings.Join(newLines, "\n")), 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %v", err)
	}

	fmt.Println("路由已删除:", routes[index])
	return nil
}

// 读取配置文件
func readConfig() (string, error) {
	// 读取配置文件
	config, err := os.ReadFile(constants.ServerConfigPath)
	if err != nil {
		return "", fmt.Errorf("读取配置文件失败: %v", err)
	}
	return string(config), nil
}

// 检查配置文件是否存在
func checkConfig() error {
	// 检查配置文件是否存在
	if _, err := os.Stat(constants.ServerConfigPath); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s", constants.ServerConfigPath)
	}
	return nil
}

// 重启OpenVPN服务
func restartService() error {
	// 重启OpenVPN服务
	cmd := exec.Command("systemctl", "restart", constants.ServiceName)
	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("重启服务失败: %v\n输出: %s", err, string(output))
	}
	return nil
}

func RestartService() error {
	fmt.Println("正在检查 OpenVPN 所需文件...")

	// 检查配置文件
	configFile := constants.ServerConfigPath
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		return fmt.Errorf("配置文件不存在: %s\n请确保配置文件已正确生成", configFile)
	}

	// 检查证书文件
	if err := CheckCertFiles(); err != nil {
		fmt.Println("\n请执行以下步骤:")
		fmt.Println("1. 确保file目录下包含所有必要的证书文件")
		fmt.Println("2. 使用root权限运行程序")
		fmt.Println("3. 选择自动安装环境选项")
		fmt.Println("4. 程序会自动复制证书文件到正确位置")
		return fmt.Errorf("证书文件缺失")
	}

	// 先停止服务
	cmd := exec.Command("systemctl", "stop", constants.ServiceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("停止服务失败: %v\n请检查服务状态: systemctl status %s", err, constants.ServiceName)
	}

	// 等待服务完全停止
	time.Sleep(2 * time.Second)

	// 启动服务
	cmd = exec.Command("systemctl", "start", constants.ServiceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("启动服务失败: %v\n请检查服务状态: systemctl status %s", err, constants.ServiceName)
	}

	// 等待服务启动
	time.Sleep(2 * time.Second)

	// 检查服务状态
	cmd = exec.Command("systemctl", "is-active", constants.ServiceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("服务未正常运行: %v\n请检查服务日志: journalctl -u %s", err, constants.ServiceName)
	}

	fmt.Println("服务重启成功")
	return nil
}

func StopService() error {
	// 检查服务状态
	cmd := exec.Command("systemctl", "is-active", constants.ServiceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("服务未运行")
	}

	// 停止服务
	cmd = exec.Command("systemctl", "stop", constants.ServiceName)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("停止服务失败: %v\n请检查服务状态: systemctl status %s", err, constants.ServiceName)
	}

	// 等待服务完全停止
	time.Sleep(2 * time.Second)

	// 验证服务已停止
	cmd = exec.Command("systemctl", "is-active", constants.ServiceName)
	if err := cmd.Run(); err == nil {
		return fmt.Errorf("服务仍在运行")
	}

	return nil
}
