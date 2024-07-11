
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	
	"openvpn-admin-go/util"
	"openvpn-admin-go/constants"
	"openvpn-admin-go/common"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "openvpn",
	Run: func(cmd *cobra.Command, args []string) {
		mainMenu()
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func check() {
	if !util.IsExists(constants.OpenVPNConfigPath) {
		fmt.Println("本机未安装openvpn, 正在自动安装...")
	}
}

func mainMenu() {
	check()
exit:
	for {
		fmt.Println()
		fmt.Println(util.Cyan("欢迎使用openvpn管理程序"))
		fmt.Println()
		menuList := []string{"openvpn管理", "用户管理", "安装管理", "web管理", "查看配置", "生成json"}
		switch util.LoopInput("请选择: ", menuList, false) {
		case 1:
			fmt.Println("openvpn管理")
		case 2:
			fmt.Println(common.GetConfig().Port)
		case 3:
			common.UpdatePort(1198)
		default:
			break exit
		}
	}
}
