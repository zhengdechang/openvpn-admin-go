package cmd

import (
	"fmt"
	"log"

	"openvpn-admin-go/router"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

/// StartWebServer starts the Gin web server.
func StartWebServer() {
	r := gin.Default()

	// 注册路由
	api := r.Group("/api")
	{
		router.SetupUserRoutes(api)
		router.SetupManageRoutes(api)
		router.SetupServerRoutes(api)
		router.SetupClientRoutes(api)
		router.SetupLogRoutes(api)
	}

	// 在 goroutine 中启动 Web 服务器
	go func() {
		fmt.Println("Web 服务器启动在 :8085 端口")
		if err := r.Run(":8085"); err != nil {
			log.Printf("Web 服务器启动失败: %v", err) // Use log.Printf for goroutine
		}
	}()
}

var webCmd = &cobra.Command{
	Use:   "web",
	Short: "启动 Web 服务器",
	Run: func(cmd *cobra.Command, args []string) {
		if err := CoreInitializer(); err != nil {
			log.Fatalf("核心初始化失败: %v", err)
		}
		StartWebServer()
		// 保持程序运行
		select {}
	},
}

func init() {
	rootCmd.AddCommand(webCmd)
} 