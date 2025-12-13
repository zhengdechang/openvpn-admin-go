package utils

import (
	"fmt"
	"os"
	"regexp"
)

// UpdateNginxListenPort 更新 Nginx 配置中的监听端口
func UpdateNginxListenPort(configPath string, port int) error {
	if port <= 0 || port > 65535 {
		return fmt.Errorf("无效的端口号: %d", port)
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取 Nginx 配置失败: %v", err)
	}

	re := regexp.MustCompile(`listen\s+\d+;`)
	if !re.Match(content) {
		return fmt.Errorf("未在 %s 中找到 listen 指令", configPath)
	}

	updated := re.ReplaceAll(content, []byte(fmt.Sprintf("listen %d;", port)))
	if err := os.WriteFile(configPath, updated, 0644); err != nil {
		return fmt.Errorf("写入 Nginx 配置失败: %v", err)
	}

	return nil
}
