package main

import (
	"bufio"
	"fmt"
	// "log" // log.Fatalf replaced with fmt.Printf for env loading
	"os"
	"os/signal"
	"strings"
	"syscall"

	"openvpn-admin-go/cmd"
	// "openvpn-admin-go/common" // Moved to cmd/web.go
	// "openvpn-admin-go/constants" // Moved to cmd/web.go
	// "openvpn-admin-go/database" // Moved to cmd/web.go (or only migrate might be needed if CLI ops require it)
	// "openvpn-admin-go/model"    // Moved to cmd/web.go (or only migrate might be needed if CLI ops require it)
	// "openvpn-admin-go/openvpn" // Moved to cmd/web.go
	// "openvpn-admin-go/router" // Moved to cmd/web.go
	// "openvpn-admin-go/services" // Moved to cmd/web.go
	// "openvpn-admin-go/utils"    // Moved to cmd/web.go
	// "path/filepath" // Moved to cmd/web.go
	// "github.com/gin-gonic/gin" // Moved to cmd/web.go
)

// loadEnv 从.env文件加载环境变量
// This function is also present in cmd/web.go. Consider refactoring to a common utility.
// For now, keeping it here for any CLI commands that might need it,
// but ensuring the web command uses its own isolated version.
func loadEnv() error {
	file, err := os.Open(".env")
	if err != nil {
		// Attempt to load from parent directory if .env is not in current dir.
		// This helps when running `go run main.go web` from project root.
		if os.IsNotExist(err) {
			file, err = os.Open("../.env")
			if err != nil {
				return fmt.Errorf("无法打开.env文件 (尝试了当前目录和父目录): %v", err)
			}
		} else {
			return fmt.Errorf("无法打开.env文件: %v", err)
		}
	}
	defer file.Close()

	// For CLI, we might not need to be as strict with all vars as web server start.
	// Define essential vars for CLI if any, otherwise, this can be more relaxed.
	requiredVars := map[string]bool{
		// Example: "OPENVPN_CONFIG_DIR": true,
		// If no vars are strictly required for CLI commands to function at a basic level,
		// this map can be empty or checks can be warnings.
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
		if requiredVars[key] {
			loadedVars[key] = true
		}
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
		return fmt.Errorf("缺少CLI必需的环境变量: %s", strings.Join(missingVars, ", "))
	}
	return nil
}

func main() {
	// 加载环境变量
	if err := loadEnv(); err != nil {
		// For CLI, this might not be fatal. Web server (cmd/web.go) has its own loadEnv and will log.Fatalf if critical.
		fmt.Printf("警告: 加载 .env 文件 (main.go context) 失败: %v\n", err)
		fmt.Println("这可能影响部分CLI功能。`openvpn web` 命令会独立加载并严格检查其所需的环境变量。")
	}

	// 设置Ctrl+C信号处理 (remains for CLI commands)
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		fmt.Println("\n程序已退出")
		os.Exit(0)
	}()

	fmt.Println("OpenVPN 管理工具 (CLI Mode)")

	// // 获取当前工作目录 - Moved to cmd/web.go or not strictly needed for all CLI commands
	// dir, err := os.Getwd()
	// if err != nil {
	// 	log.Fatalf("获取工作目录失败: %v", err)
	// }
	// fmt.Printf("当前工作目录: %s\n", dir)

	// // 检查环境 - Moved to cmd/web.go, CLI commands might need selective checks or none
	// // The interactive environment check and installation is primarily for the initial setup,
	// // which is now better associated with the web server part or a dedicated setup command.
	// fmt.Println("正在检查运行环境...")
	// if err := cmd.CheckEnvironment(); err != nil {
	// 	fmt.Printf("环境检查失败: %v\n", err)
	// 	fmt.Println("是否自动安装所需环境？(y/n)")
	// 	var choice string
	// 	fmt.Scanln(&choice)
	// 	if choice == "y" || choice == "Y" {
	// 		if err := cmd.InstallEnvironment(); err != nil {
	// 			fmt.Printf("环境安装失败: %v\n", err)
	// 			fmt.Println("请确保您有足够的权限，软件源配置正确，网络连接稳定。")
	// 			fmt.Println("您可以手动检查并修复问题后重新运行程序。")
	// 			return
	// 		}
	// 		if err := cmd.CheckEnvironment(); err != nil {
	// 			fmt.Printf("环境检查仍然失败: %v\n", err)
	// 			fmt.Println("请手动检查并修复问题后重新运行程序。")
	// 			return
	// 		}
	// 	} else {
	// 		fmt.Println("请手动安装所需环境后重新运行程序。")
	// 		return
	// 	}
	// }
	// fmt.Println("环境检查通过")

	// // 初始化数据库 - All database setup, migration, seeding is moved to cmd/web.go
	// // if err := database.Init(); err != nil {
	// // 	log.Fatalf("数据库初始化失败: %v", err)
	// // }
	// // if err := database.Migrate(&model.User{}, &model.Department{}); err != nil {
	// // 	log.Fatalf("数据库迁移失败: %v", err)
	// // }
	// // Seed default superadmin user if not exists - Moved to cmd/web.go
	// func() {
	// 	// ... seeding logic ...
	// }()
	// // 确保数据库用户在 OpenVPN 客户端存在，不存在则自动创建 - Moved to cmd/web.go
	// func() {
	// 	// ... client creation logic ...
	// }()

	// // 启动 Web 服务器 - Moved to cmd/web.go
	// // r := gin.Default()
	// // api := r.Group("/api")
	// // {
	// // 	router.SetupUserRoutes(api)
	// // 	// ... other router setups
	// // }
	// // go func() {
	// // 	fmt.Println("Web 服务器启动在 :8085 端口")
	// // 	if err := r.Run(":8085"); err != nil {
	// // 		log.Fatal("Failed to start server:", err)
	// // 	}
	// // }()

	// // Start OpenVPN Data Synchronization Service - Moved to cmd/web.go
	// // statusLogPath := utils.GetOpenVPNStatusLogPath()
	// // syncInterval := utils.GetOpenVPNSyncInterval()
	// // log.Printf("Starting OpenVPN Sync Service: LogPath='%s', Interval=%s", statusLogPath, syncInterval)
	// // go services.StartOpenVPNSyncService(database.DB, statusLogPath, syncInterval)

	// 启动 Cobra 命令处理
	cmd.Execute()
}
