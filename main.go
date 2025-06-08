package main

import (
   "bufio"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"openvpn-admin-go/cmd"
   "openvpn-admin-go/database"
   "openvpn-admin-go/model"
   "openvpn-admin-go/common"
   "openvpn-admin-go/openvpn"
   "openvpn-admin-go/constants"
   "path/filepath"
   "openvpn-admin-go/services" // Added for OpenVPN Sync Service
   "openvpn-admin-go/utils"    // Added for config utils
)

// loadEnv 从.env文件加载环境变量
func loadEnv() error {
	// Check if .env file exists. If not, copy from .env.example
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Println(".env file not found. Attempting to copy from .env.example...")
		sourceFile, err := os.Open(".env.example")
		if err != nil {
			return fmt.Errorf("无法打开 .env.example 文件: %v", err)
		}
		defer sourceFile.Close()

		destinationFile, err := os.Create(".env")
		if err != nil {
			return fmt.Errorf("无法创建 .env 文件: %v", err)
		}
		defer destinationFile.Close()

		_, err = bufio.NewReader(sourceFile).WriteTo(destinationFile)
		if err != nil {
			return fmt.Errorf("无法从 .env.example 复制到 .env: %v", err)
		}
		log.Println(".env file copied successfully from .env.example.")
	}

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

// InitCore initializes core application services.
func InitCore() error {
	// 加载环境变量
	if err := loadEnv(); err != nil {
		return fmt.Errorf("加载环境变量失败: %v\n请确保.env文件存在且包含所有必需的配置项。", err)
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
		return fmt.Errorf("获取工作目录失败: %v", err)
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
			if errInstall := cmd.InstallEnvironment(); errInstall != nil {
				return fmt.Errorf("环境安装失败: %v\n请确保您有足够的权限，软件源配置正确，网络连接稳定。\n您可以手动检查并修复问题后重新运行程序。", errInstall)
			}
			// 重新检查环境
			if errCheckAgain := cmd.CheckEnvironment(); errCheckAgain != nil {
				return fmt.Errorf("环境检查仍然失败: %v\n请手动检查并修复问题后重新运行程序。", errCheckAgain)
			}
		} else {
			return fmt.Errorf("请手动安装所需环境后重新运行程序")
		}
	}

	fmt.Println("环境检查通过")

	// 初始化数据库
	if err := database.Init(); err != nil {
		return fmt.Errorf("数据库初始化失败: %v", err)
	}
	if err := database.Migrate(&model.User{}, &model.Department{}); err != nil {
		return fmt.Errorf("数据库迁移失败: %v", err)
	}
	// Seed default superadmin user if not exists
	func() {
		var existing model.User
		if err := database.DB.Where("email = ?", "superadmin@gmail.com").First(&existing).Error; err != nil {
			hash, errHash := common.HashPassword("superadmin")
			if errHash != nil {
				log.Printf("默认超级管理员密码哈希失败: %v", errHash) // Log and continue
				return
			}
			super := model.User{
				Name:         "Super Admin",
				Email:        "superadmin@gmail.com",
				PasswordHash: hash,
				Role:         model.RoleSuperAdmin,
			}
			if errCreate := database.DB.Create(&super).Error; errCreate != nil {
				log.Printf("创建默认超级管理员失败: %v", errCreate) // Log and continue
			} else {
				log.Println("已创建默认超级管理员: superadmin@gmail.com / superadmin")
			}
		}
	}()
	// 确保数据库用户在 OpenVPN 客户端存在，不存在则自动创建
	func() {
		var users []model.User
		if err := database.DB.Find(&users).Error; err != nil {
			log.Printf("查询用户列表失败: %v", err) // Log and continue
		} else {
			for _, u := range users {
				clientPath := filepath.Join(constants.ClientConfigDir, u.ID+".ovpn")
				if _, errStat := os.Stat(clientPath); os.IsNotExist(errStat) {
					if errCreate := openvpn.CreateClient(u.ID); errCreate != nil {
						log.Printf("创建 OpenVPN 客户端 %s 失败: %v", u.ID, errCreate) // Log and continue
					}
				}
			}
		}
	}()

	// Start OpenVPN Data Synchronization Service
	statusLogPath := utils.GetOpenVPNStatusLogPath()
	syncInterval := utils.GetOpenVPNSyncInterval()
	log.Printf("Starting OpenVPN Sync Service: LogPath='%s', Interval=%s", statusLogPath, syncInterval)
	go services.StartOpenVPNSyncService(database.DB, statusLogPath, syncInterval)

	return nil
}



func main() {
	// Assign the public functions to the variables in the cmd package.
	cmd.CoreInitializer = InitCore
	if err := cmd.CoreInitializer(); err != nil {
		log.Fatalf("核心初始化失败: %v", err)
	}
	cmd.Execute()
}
