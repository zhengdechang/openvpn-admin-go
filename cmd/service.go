package cmd

import (
	"fmt"
	"os/exec"
)

const serviceName = "openvpn-server@server.service"

// startService starts the OpenVPN service.
func startService() {
	fmt.Println("Attempting to start OpenVPN service...")
	cmd := exec.Command("sudo", "systemctl", "start", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to start OpenVPN service: %v\nOutput: %s\n", err, string(output))
		return
	}
	fmt.Printf("OpenVPN service started successfully.\nOutput: %s\n", string(output))
}

// stopService stops the OpenVPN service.
func stopService() {
	fmt.Println("Attempting to stop OpenVPN service...")
	cmd := exec.Command("sudo", "systemctl", "stop", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to stop OpenVPN service: %v\nOutput: %s\n", err, string(output))
		return
	}
	fmt.Printf("OpenVPN service stopped successfully.\nOutput: %s\n", string(output))
}

// restartService restarts the OpenVPN service.
func restartService() {
	fmt.Println("Attempting to restart OpenVPN service...")
	cmd := exec.Command("sudo", "systemctl", "restart", serviceName)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to restart OpenVPN service: %v\nOutput: %s\n", err, string(output))
		return
	}
	fmt.Printf("OpenVPN service restarted successfully.\nOutput: %s\n", string(output))
}
