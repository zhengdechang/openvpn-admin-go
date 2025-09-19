package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"openvpn-admin-go/cmd"
	"openvpn-admin-go/common"
	"openvpn-admin-go/constants"
	"openvpn-admin-go/database"
	"openvpn-admin-go/logging"
	"openvpn-admin-go/model"
	"openvpn-admin-go/openvpn"
	"path/filepath"

	"gorm.io/gorm"
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

	// 初始化日志系统
	logConfigPath := "config/logging.json"
	if err := logging.InitFromConfig(logConfigPath); err != nil {
		fmt.Printf("日志系统初始化失败: %v\n", err)
		// 使用默认配置继续运行
		defaultConfig := logging.Config{
			LogLevel:      logging.INFO,
			LogFilePath:   "logs/web.log",
			EnableAPI:     true,
			EnableConsole: false,
		}
		if err := logging.Init(defaultConfig); err != nil {
			return fmt.Errorf("日志系统初始化失败: %v", err)
		}
	}
	logging.Info("OpenVPN 管理工具启动中...")

	// 获取当前工作目录
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("获取工作目录失败: %v", err)
	}
	logging.Info("当前工作目录: %s", dir)

	// 检查环境
	fmt.Println("正在检查运行环境...")
	if err := cmd.CheckEnvironment(); err != nil {
		fmt.Printf("环境检查失败: %v\n", err)
		fmt.Println("尝试自动安装所需环境...")
		if errInstall := cmd.InstallEnvironment(); errInstall != nil {
			return fmt.Errorf("环境自动安装失败: %v\n请确保权限、软件源和网络连接正常，然后重试。", errInstall)
		}
		// 重新检查环境
		if errCheckAgain := cmd.CheckEnvironment(); errCheckAgain != nil {
			return fmt.Errorf("环境检查仍然失败: %v\n请手动检查并修复问题后重新运行程序。", errCheckAgain)
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
				logging.Error("默认超级管理员密码哈希失败: %v", errHash) // Log and continue
				return
			}
			super := model.User{
				Name:         "Super Admin",
				Email:        "superadmin@gmail.com",
				PasswordHash: hash,
				Role:         model.RoleSuperAdmin,
			}
			if errCreate := database.DB.Create(&super).Error; errCreate != nil {
				logging.Error("创建默认超级管理员失败: %v", errCreate) // Log and continue
			} else {
				logging.Info("已创建默认超级管理员: superadmin@gmail.com / superadmin")
			}
		}
	}()
	// 双向同步：确保数据库用户和 OpenVPN 客户端配置保持一致
	func() {
		// 第一步：确保数据库用户在 OpenVPN 客户端存在，不存在则自动创建
		var users []model.User
		if err := database.DB.Find(&users).Error; err != nil {
			logging.Error("查询用户列表失败: %v", err) // Log and continue
		} else {
			for _, u := range users {
				clientPath := filepath.Join(constants.ClientConfigDir, u.Name+".ovpn")
				if _, errStat := os.Stat(clientPath); os.IsNotExist(errStat) {
					if errCreate := openvpn.CreateClient(u.Name); errCreate != nil {
						logging.Error("创建 OpenVPN 客户端 %s 失败: %v", u.Name, errCreate) // Log and continue
					} else {
						logging.Info("为数据库用户 %s 创建了 OpenVPN 客户端配置", u.Name)
					}
				}
			}
		}

		// 第二步：检查 OpenVPN 客户端配置，如果数据库中没有对应用户则创建
		if _, err := os.Stat(constants.ClientConfigDir); os.IsNotExist(err) {
			logging.Warn("OpenVPN 客户端配置目录不存在: %s", constants.ClientConfigDir)
			return
		}

		files, err := os.ReadDir(constants.ClientConfigDir)
		if err != nil {
			logging.Error("读取 OpenVPN 客户端配置目录失败: %v", err)
			return
		}

		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".ovpn") {
				// 提取用户名（去掉.ovpn扩展名）
				userName := strings.TrimSuffix(file.Name(), ".ovpn")

				// 检查数据库中是否存在该用户（按用户名查找）
				var existingUser model.User
				if err := database.DB.Where("name = ?", userName).First(&existingUser).Error; err != nil {
					if err == gorm.ErrRecordNotFound {
						// 用户不存在，创建新用户
						hash, errHash := common.HashPassword("changeme123") // 默认密码
						if errHash != nil {
							logging.Error("为 OpenVPN 客户端 %s 生成默认密码哈希失败: %v", userName, errHash)
							continue
						}

						// 检查是否有固定IP配置
						fixedIP, errFixedIP := openvpn.GetClientFixedIP(userName)
						if errFixedIP != nil {
							logging.Warn("检查 OpenVPN 客户端 %s 固定IP配置失败: %v", userName, errFixedIP)
						}

						// 检查是否有子网配置
						subnet, errSubnet := openvpn.GetClientSubnet(userName)
						if errSubnet != nil {
							logging.Warn("检查 OpenVPN 客户端 %s 子网配置失败: %v", userName, errSubnet)
						}

						newUser := model.User{
							Name:         userName,
							Email:        userName + "@openvpn.local", // 生成默认邮箱
							PasswordHash: hash,
							Role:         model.RoleUser, // 默认角色为普通用户
							FixedIP:      fixedIP,        // 设置固定IP（如果有）
							Subnet:       subnet,         // 设置子网（如果有）
						}

						if errCreate := database.DB.Create(&newUser).Error; errCreate != nil {
							logging.Error("为 OpenVPN 客户端 %s 创建数据库用户失败: %v", userName, errCreate)
						} else {
							logging.Info("为 OpenVPN 客户端 %s 创建了数据库用户，默认密码: changeme123", userName)

							// 如果有固定IP配置，确保CCD配置正确
							if fixedIP != "" {
								if errSetIP := openvpn.SetClientFixedIP(userName, fixedIP); errSetIP != nil {
									logging.Error("为用户 %s 设置固定IP %s 失败: %v", userName, fixedIP, errSetIP)
								} else {
									logging.Info("用户 %s 已设置固定IP: %s", userName, fixedIP)
								}
							}

							// 如果有子网配置，确保CCD配置正确
							if subnet != "" {
								if errSetSubnet := openvpn.SetClientSubnet(userName, subnet); errSetSubnet != nil {
									logging.Error("为用户 %s 设置子网 %s 失败: %v", userName, subnet, errSetSubnet)
								} else {
									logging.Info("用户 %s 已设置子网: %s", userName, subnet)
								}
							}
						}
					} else {
						logging.Error("查询用户 %s 时发生错误: %v", userName, err)
					}
				} else {
					// 用户已存在，以数据库为准，确保CCD配置与数据库一致
					logging.Info("用户 %s 已存在，检查CCD配置是否与数据库一致", userName)

					// 如果数据库中有固定IP配置，确保CCD配置正确
					if existingUser.FixedIP != "" {
						if errSetIP := openvpn.SetClientFixedIP(userName, existingUser.FixedIP); errSetIP != nil {
							logging.Error("为用户 %s 设置固定IP %s 失败: %v", userName, existingUser.FixedIP, errSetIP)
						} else {
							logging.Info("为用户 %s 确保固定IP配置: %s", userName, existingUser.FixedIP)
						}
					}

					// 如果数据库中有子网配置，确保CCD配置正确
					if existingUser.Subnet != "" {
						if errSetSubnet := openvpn.SetClientSubnet(userName, existingUser.Subnet); errSetSubnet != nil {
							logging.Error("为用户 %s 设置子网 %s 失败: %v", userName, existingUser.Subnet, errSetSubnet)
						} else {
							logging.Info("为用户 %s 确保子网配置: %s", userName, existingUser.Subnet)
						}
					}
				}
			}
		}
	}()

	// OpenVPN 同步服务将在 web 服务启动时启动，避免数据库锁定问题
	logging.Info("核心初始化完成，OpenVPN 同步服务将在 web 服务启动时启动")

	return nil
}

func main() {
	// Assign the public functions to the variables in the cmd package.
	cmd.CoreInitializer = InitCore
	if err := cmd.CoreInitializer(); err != nil {
		logging.Fatal("核心初始化失败: %v", err)
	}
	cmd.Execute()
}
