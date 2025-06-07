package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	// "os/signal" // Not needed for web server alone
	// "syscall" // Not needed for web server alone

	// Keep imports from original main.go that are relevant to web server
	"openvpn-admin-go/common"
	"openvpn-admin-go/constants"
	"openvpn-admin-go/database"
	"openvpn-admin-go/model"
	"openvpn-admin-go/openvpn"
	"openvpn-admin-go/router"
	"openvpn-admin-go/services"
	"openvpn-admin-go/utils"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// loadEnv loads environment variables from .env file.
// This is a copy of the function in root main.go.
// Consider refactoring to a common utility package later.
func loadEnv() error {
	file, err := os.Open(".env")
	if err != nil {
		// If .env is not found, try to load from parent directory, useful for `go run main.go web`
		if os.IsNotExist(err) {
			file, err = os.Open("../.env")
			if err != nil {
				return fmt.Errorf("无法打开.env文件: %v", err)
			}
		} else {
			return fmt.Errorf("无法打开.env文件: %v", err)
		}
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
		if strings.HasPrefix(line, "#") || len(strings.TrimSpace(line)) == 0 {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])
		if err := os.Setenv(key, value); err != nil {
			return fmt.Errorf("设置环境变量失败 %s: %v", key, err)
		}
		loadedVars[key] = true
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("读取.env文件失败: %v", err)
	}

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

// startWebServer initializes and starts the Gin web server.
func startWebServer() {
	fmt.Println("Starting web server...")

	// Load environment variables
	if err := loadEnv(); err != nil {
		log.Fatalf("错误: 加载环境变量失败: %v\n请确保.env文件存在且包含所有必需的配置项。", err)
	}

	// Get current working directory (useful for context, though not strictly used by server itself yet)
	dir, err := os.Getwd()
	if err != nil {
		log.Printf("获取工作目录失败: %v (non-fatal for web server)", err)
	} else {
		fmt.Printf("当前工作目录 (web server context): %s\n", dir)
	}

	// Check and install environment if necessary
	// This uses CheckEnvironment and InstallEnvironment from the cmd package (environment.go)
	fmt.Println("正在检查运行环境 (web server context)...")
	if err := CheckEnvironment(); err != nil {
		fmt.Printf("环境检查失败 (web server context): %v\n", err)
		// Attempt to install environment if check fails (non-interactive for server context)
		fmt.Println("尝试自动安装所需环境 (web server context)...")
		if errInstall := InstallEnvironment(); errInstall != nil {
			fmt.Printf("环境安装失败 (web server context): %v\n", errInstall)
			fmt.Println("请确保您有足够的权限，软件源配置正确，网络连接稳定。")
			fmt.Println("您可以手动检查并修复问题后重新运行程序。")
			// For a web command, we might want to exit if env setup fails critically
			// For now, we'll log and proceed, but this might need adjustment
		} else {
			// Re-check environment after installation attempt
			if errRecheck := CheckEnvironment(); errRecheck != nil {
				fmt.Printf("环境检查仍然失败 (web server context): %v\n", errRecheck)
				fmt.Println("请手动检查并修复问题后重新运行程序。")
				// Exit or handle error more gracefully
			} else {
				fmt.Println("环境已成功安装并通过检查 (web server context)。")
			}
		}
	} else {
		fmt.Println("环境检查通过 (web server context)。")
	}

	// Initialize database
	if err := database.Init(); err != nil {
		log.Fatalf("数据库初始化失败: %v", err)
	}
	if err := database.Migrate(&model.User{}, &model.Department{}); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}

	// Seed default superadmin user if not exists
	func() {
		var existing model.User
		if err := database.DB.Where("email = ?", "superadmin@gmail.com").First(&existing).Error; err != nil {
			hash, err := common.HashPassword("superadmin")
			if err != nil {
				log.Printf("默认超级管理员密码哈希失败: %v", err)
				return
			}
			super := model.User{
				Name:         "Super Admin",
				Email:        "superadmin@gmail.com",
				PasswordHash: hash,
				Role:         model.RoleSuperAdmin,
			}
			if err := database.DB.Create(&super).Error; err != nil {
				log.Printf("创建默认超级管理员失败: %v", err)
			} else {
				log.Println("已创建默认超级管理员: superadmin@gmail.com / superadmin")
			}
		}
	}()

	// Ensure database users have OpenVPN clients
	func() {
		var users []model.User
		if err := database.DB.Find(&users).Error; err != nil {
			log.Printf("查询用户列表失败: %v", err)
		} else {
			for _, u := range users {
				clientPath := filepath.Join(constants.ClientConfigDir, u.ID+".ovpn")
				if _, err := os.Stat(clientPath); os.IsNotExist(err) {
					// Ensure the user ID is valid (e.g. not empty) before creating client
					if u.ID == "" {
						log.Printf("用户ID为空，跳过创建OpenVPN客户端: User Name %s, Email %s", u.Name, u.Email)
						continue
					}
					if err := openvpn.CreateClient(u.ID); err != nil {
						log.Printf("创建 OpenVPN 客户端 %s 失败: %v", u.ID, err)
					}
				}
			}
		}
	}()

	// Start Gin server
	r := gin.Default()
	api := r.Group("/api")
	{
		router.SetupUserRoutes(api)
		router.SetupManageRoutes(api)
		router.SetupServerRoutes(api)
		router.SetupClientRoutes(api)
		router.SetupLogRoutes(api)
	}

	// Start OpenVPN Data Synchronization Service
	statusLogPath := utils.GetOpenVPNStatusLogPath()
	syncInterval := utils.GetOpenVPNSyncInterval()
	log.Printf("Starting OpenVPN Sync Service (from web command): LogPath='%s', Interval=%s", statusLogPath, syncInterval)
	go services.StartOpenVPNSyncService(database.DB, statusLogPath, syncInterval)

	fmt.Println("Web 服务器准备启动在 :8085 端口...")
	if err := r.Run(":8085"); err != nil {
		log.Fatalf("启动 Web 服务器失败: %v", err)
	}
	fmt.Println("Web 服务器已停止.") // This line will likely not be reached if Run is blocking and successful
}
