package cmd

import (
	"fmt"
	"log"

	"openvpn-admin-go/router"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strconv"
	"syscall"
	"time"

	"openvpn-admin-go/router"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
)

const pidFileName = "web.pid" // Or a more suitable path

var webPort int

// webCmd is the parent command for web server operations.
var webCmd = &cobra.Command{
	Use:   "web",
	Short: "管理 Web 服务器 (start, stop, restart)",
	Long:  `管理 Web 服务器. 使用 "web start" 启动, "web stop" 停止, "web restart" 重启.`,
	// No Run function for parent command
}

// startWebCmd starts the Gin web server.
var startWebCmd = &cobra.Command{
	Use:   "start",
	Short: "启动 Web 服务器",
	Run: func(cmd *cobra.Command, args []string) {
		if err := CoreInitializer(); err != nil {
			log.Fatalf("核心初始化失败: %v", err)
		}

		pidDir := os.TempDir() // Or a more persistent location like /var/run/yourapp/
		pidFile := filepath.Join(pidDir, pidFileName)

		// Check if PID file exists and process is running
		if pidData, err := os.ReadFile(pidFile); err == nil {
			if pid, err := strconv.Atoi(string(pidData)); err == nil {
				// Check if the process with this PID is actually running.
				// Sending signal 0 to a process checks if it exists without harming it.
				if process, _ := os.FindProcess(pid); process != nil && process.Signal(syscall.Signal(0)) == nil {
					log.Printf("Web 服务器似乎已在运行 (PID: %d, PID 文件: %s). 请先停止它.", pid, pidFile)
					return
				}
				log.Printf("发现陈旧的 PID 文件 (%s, PID: %d), 但进程未运行. 将覆盖 PID 文件.", pidFile, pid)
			} else {
				log.Printf("PID 文件 (%s) 内容无效, 将覆盖.", pidFile)
			}
		}

		// Setup Gin router
		r := gin.Default()
		api := r.Group("/api")
		{
			router.SetupUserRoutes(api)
			router.SetupManageRoutes(api)
			router.SetupServerRoutes(api)
			router.SetupClientRoutes(api)
			router.SetupLogRoutes(api)
		}

		serverAddr := fmt.Sprintf(":%d", webPort)

		// Channel to signal server start
		serverStarted := make(chan struct{})
		go func() {
			fmt.Printf("准备在 %s 端口启动 Web 服务器...\n", serverAddr)
			// Close channel once server is confirmed or fails to start
			// For r.Run, it blocks, so if it returns, it either failed or was stopped.
			// A more robust way would be to check if the port is listening.
			close(serverStarted)
			if err := r.Run(serverAddr); err != nil {
				// This log might not be seen if the main startWebCmd goroutine has exited after PID write
				log.Printf("Web 服务器错误: %v", err)
			}
		}()

		// Wait briefly for server to attempt start-up before writing PID
		// This is a simplification. True readiness should be confirmed.
		<-serverStarted
		fmt.Println("Web 服务器正在后台启动...")


		// Save PID of the current 'web start' command process
		// Note: This is the PID of the 'web start' command itself, not a detached Gin server.
		// For a simple management within this CLI, this is acceptable.
		currentPid := os.Getpid()
		if err := os.WriteFile(pidFile, []byte(strconv.Itoa(currentPid)), 0644); err != nil {
			log.Fatalf("无法写入 PID 文件 %s: %v", pidFile, err)
		}
		fmt.Printf("Web 服务器 '管理进程' PID: %d. PID 文件: %s\n", currentPid, pidFile)
		fmt.Println("使用 'web stop' 来停止此管理进程和服务器.")

		// Keep the 'start' command running until interrupted by SIGINT or SIGTERM
		// This allows the 'stop' command to find and kill this process.
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Wait for a signal
		receivedSignal := <-sigChan
		fmt.Printf("\n收到信号: %v. 正在关闭 Web 服务器管理进程...\n", receivedSignal)

		// Cleanup PID file on exit
		if err := os.Remove(pidFile); err != nil {
			if !os.IsNotExist(err) {
				log.Printf("关闭时删除 PID 文件失败: %v", err)
			}
		} else {
			log.Println("PID 文件已成功删除.")
		}
		// Perform any other necessary cleanup for the Gin server if possible/needed here
		// (e.g. context cancellation if the gin server was started with one)
	},
}

// stopWebCmd stops the Gin web server.
var stopWebCmd = &cobra.Command{
	Use:   "stop",
	Short: "停止 Web 服务器",
	Run: func(cmd *cobra.Command, args []string) {
		pidDir := os.TempDir()
		pidFile := filepath.Join(pidDir, pidFileName)

		pidData, err := os.ReadFile(pidFile)
		if err != nil {
			if os.IsNotExist(err) {
				log.Println("PID 文件未找到. 服务器可能未运行.")
				// Optionally, try to find process by name if PID file is missing but server might be running.
				// This is more complex and platform-dependent. For now, rely on PID file.
				return
			}
			log.Fatalf("读取 PID 文件失败 %s: %v", pidFile, err)
		}

		pid, err := strconv.Atoi(string(pidData))
		if err != nil {
			log.Fatalf("PID 文件内容无效 %s: %v", pidFile, err)
		}

		process, err := os.FindProcess(pid)
		if err != nil {
			// This error means os.FindProcess failed, not that the process doesn't exist.
			// However, on Unix, FindProcess always succeeds for valid PIDs.
			log.Printf("获取进程句柄失败 (PID: %d): %v. 可能进程已不存在.", pid, err)
			// Attempt to remove stale PID file
			if rmErr := os.Remove(pidFile); rmErr == nil {
				log.Printf("已删除陈旧/无效的 PID 文件: %s", pidFile)
			}
			return
		}

		// Send SIGTERM to the process
		log.Printf("正在向 PID %d 发送 SIGTERM 信号...", pid)
		if err := process.Signal(syscall.SIGTERM); err != nil {
			// Check if the process was already dead
			// Note: Error checking for Signal can be tricky.
			// If Signal returns "os: process already finished", it's a clean case.
			// Other errors might mean the process exists but we can't signal it (permissions?),
			// or it's a different kind of issue.
			if checkErr := process.Signal(syscall.Signal(0)); checkErr != nil { // Check if actually running
				log.Printf("进程 (PID: %d) 似乎未在运行: %v", pid, checkErr)
				if rmErr := os.Remove(pidFile); rmErr == nil {
					log.Printf("已删除陈旧的 PID 文件: %s", pidFile)
				} else if !os.IsNotExist(rmErr) {
					log.Printf("尝试删除陈旧 PID 文件失败 %s: %v", pidFile, rmErr)
				}
			} else {
				// Process is running but Signal(SIGTERM) failed for another reason
				log.Fatalf("发送 SIGTERM 信号失败 (PID: %d): %v", pid, err)
			}
		} else {
			log.Printf("已成功向 PID %d 发送 SIGTERM 信号.", pid)
			// The terminated process's shutdown hook (in startWebCmd) should remove the PID file.
			// We can add a small delay and check/remove it here as a fallback.
			time.Sleep(100 * time.Millisecond)
			if _, statErr := os.Stat(pidFile); statErr == nil {
				if rmErr := os.Remove(pidFile); rmErr == nil {
					log.Println("由停止命令删除的 PID 文件.")
				} else {
					log.Printf("停止命令删除 PID 文件失败: %v", rmErr)
				}
			}
		}
	},
}

// restartWebCmd restarts the Gin web server.
var restartWebCmd = &cobra.Command{
	Use:   "restart",
	Short: "重启 Web 服务器",
	Run: func(cmd *cobra.Command, args []string) {
		log.Println("正在重启 Web 服务器...")

		// Call stop logic
		log.Println("执行停止操作...")
		stopWebCmd.Run(cmd, args)

		// Wait for a bit to ensure the old process has time to shut down
		log.Println("等待服务器关闭 (2 秒)...")
		time.Sleep(2 * time.Second)

		// Call start logic
		log.Println("执行启动操作...")
		startWebCmd.Run(cmd, args)
		log.Println("Web 服务器重启过程已启动.")
	},
}

func init() {
	// Add flags and subcommands here
	webCmd.PersistentFlags().IntVarP(&webPort, "port", "p", 8085, "Web 服务器监听的端口")

	webCmd.AddCommand(startWebCmd)
	webCmd.AddCommand(stopWebCmd)
	webCmd.AddCommand(restartWebCmd)

	rootCmd.AddCommand(webCmd)
} 