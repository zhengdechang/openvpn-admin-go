package cmd

import (
	"fmt"

	"openvpn-admin-go/utils"

	"github.com/spf13/cobra"
)

var (
	supervisorMainOnly  bool
	supervisorService   string
	supervisorPort      int
	supervisorAutostart bool
)

var supervisorConfigCmd = &cobra.Command{
	Use:   "supervisor-config",
	Short: "配置或更新 supervisor 服务",
	RunE: func(cmd *cobra.Command, args []string) error {
		if !supervisorMainOnly && supervisorService == "" {
			return fmt.Errorf("请通过 --main-only 或 --service 指定需要执行的操作")
		}

		if supervisorMainOnly {
			if err := utils.InstallSupervisorMainConfig(); err != nil {
				return fmt.Errorf("安装 supervisor 主配置失败: %w", err)
			}
		}

		if supervisorService != "" && !utils.IsSupervisorConfigExists() {
			if err := utils.InstallSupervisorMainConfig(); err != nil {
				return fmt.Errorf("安装 supervisor 主配置失败: %w", err)
			}
		}

		reloadIfRunning := func() error {
			if utils.IsSupervisordRunning() {
				if err := utils.SupervisorctlReload(); err != nil {
					return fmt.Errorf("重新加载 supervisor 配置失败: %w", err)
				}
			}
			return nil
		}

		switch supervisorService {
		case "":
			// no-op
		case "web", "api":
			config := utils.ServiceConfig{
				Port:      supervisorPort,
				AutoStart: supervisorAutostart,
			}
			if err := utils.InstallWebServiceConfig(config); err != nil {
				return fmt.Errorf("配置 API 服务失败: %w", err)
			}
			if err := reloadIfRunning(); err != nil {
				return err
			}
		case "openvpn", "backend":
			if err := utils.InstallOpenVPNServiceConfig(supervisorAutostart); err != nil {
				return fmt.Errorf("配置 OpenVPN 服务失败: %w", err)
			}
			if err := reloadIfRunning(); err != nil {
				return err
			}
		case "frontend":
			if err := utils.InstallFrontendServiceConfig(supervisorAutostart); err != nil {
				return fmt.Errorf("配置前端服务失败: %w", err)
			}
			if err := reloadIfRunning(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("未知服务类型: %s", supervisorService)
		}

		return nil
	},
}

func init() {
	supervisorConfigCmd.Flags().BoolVar(&supervisorMainOnly, "main-only", false, "安装或刷新 supervisor 主配置")
	supervisorConfigCmd.Flags().StringVar(&supervisorService, "service", "", "配置的服务类型 (web|api|openvpn|backend|frontend)")
	supervisorConfigCmd.Flags().IntVar(&supervisorPort, "port", 8085, "API 服务端口")
	supervisorConfigCmd.Flags().BoolVar(&supervisorAutostart, "autostart", false, "是否在 supervisord 中启用自动启动")
	rootCmd.AddCommand(supervisorConfigCmd)
}
