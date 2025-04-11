package cmd

import (
	"fmt"
	"os"

	"github.com/manifoldco/promptui"
)

func MainMenu() {
	for {
		prompt := promptui.Select{
			Label: "请选择操作",
			Items: []string{
				"服务器管理",
				"客户端管理",
				"退出",
			},
			Size: 3,
			Templates: &promptui.SelectTemplates{
				Label:    "{{ . }}",
				Active:   "➤ {{ . | cyan }}",
				Inactive: "  {{ . | white }}",
				Selected: "{{ . | green }}",
			},
			Keys: &promptui.SelectKeys{
				Prev:     promptui.Key{Code: promptui.KeyPrev, Display: "↑"},
				Next:     promptui.Key{Code: promptui.KeyNext, Display: "↓"},
				PageUp:   promptui.Key{Code: promptui.KeyBackward, Display: "←"},
				PageDown: promptui.Key{Code: promptui.KeyForward, Display: "→"},
			},
		}

		_, result, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt {
				fmt.Println("\n程序已退出")
				os.Exit(0)
			}
			fmt.Printf("选择失败: %v\n", err)
			continue
		}

		switch result {
		case "服务器管理":
			ServerMenu()
		case "客户端管理":
			ClientMenu()
		case "退出":
			fmt.Println("程序已退出")
			os.Exit(0)
		}
	}
}
