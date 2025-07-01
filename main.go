package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"openvpn-admin-go/cmd"
	"openvpn-admin-go/common"
	"openvpn-admin-go/constants"
	"openvpn-admin-go/database"
	"openvpn-admin-go/model"
	"openvpn-admin-go/openvpn"
	"openvpn-admin-go/services" // Added for OpenVPN Sync Service
	"openvpn-admin-go/utils"    // Added for config utils
	"path/filepath"
)

// InitCore initializes core application services.
func InitCore() error {
	// 不再需要加载环境变量，配置将从 JSON 文件或常量中加载

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
	cfg, err := openvpn.LoadConfig()
	if err != nil {
		// If config loading fails here, it's a significant issue as other parts might also fail.
		// However, InitCore is already designed to return an error.
		return fmt.Errorf("无法加载 OpenVPN 配置以启动同步服务: %v", err)
	}
	statusLogPath := cfg.OpenVPNStatusLogPath      // Use configured path
	syncInterval := utils.GetOpenVPNSyncInterval() // Assuming this handles its own config or defaults
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
