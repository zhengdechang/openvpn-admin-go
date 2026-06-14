package cmd

import (
	"fmt"

	"openvpn-admin-go/database"
	"openvpn-admin-go/logging"

	"github.com/spf13/cobra"
)

// migrateCmd 提供手动数据库迁移操作（up/down/status）。
// 注意：后端正常启动时会在 InitCore 中自动执行 up，此命令用于运维场景。
var migrateCmd = &cobra.Command{
	Use:   "migrate [up|down|status]",
	Short: "数据库迁移管理 (goose)",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := database.Init(); err != nil {
			logging.Fatal("数据库初始化失败: %v", err)
		}

		var err error
		switch args[0] {
		case "up":
			err = database.RunMigrations()
			if err == nil {
				fmt.Println("迁移已全部应用")
			}
		case "down":
			err = database.MigrateDown()
			if err == nil {
				fmt.Println("已回滚最近一次迁移")
			}
		case "status":
			err = database.MigrateStatus()
		default:
			logging.Fatal("未知的迁移操作: %s (可用: up|down|status)", args[0])
		}

		if err != nil {
			logging.Fatal("迁移操作失败: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
}
