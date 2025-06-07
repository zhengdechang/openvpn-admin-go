package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"openvpn-admin-go/constants"
)

// showLogs displays OpenVPN logs.
// It can show the main log or status log, with an option to follow the main log.
func showLogs(args []string) {
	targetLog := constants.OpenVPNLogPath // Default to main log
	logTypeName := "main OpenVPN log"
	follow := false
	tailLines := "50" // Default number of lines for cat/tail non-follow

	// Parse arguments
	isStatusLog := false
	for _, arg := range args {
		lowArg := strings.ToLower(arg)
		if lowArg == "-f" || lowArg == "--follow" {
			follow = true
		}
		if lowArg == "status" {
			isStatusLog = true
		}
	}

	if isStatusLog {
		targetLog = constants.ServerStatusLogPath
		logTypeName = "OpenVPN status log"
	}

	if follow {
		fmt.Printf("Following %s from: %s (Ctrl+C to stop)\n", logTypeName, targetLog)
		if isStatusLog {
			fmt.Println("Note: Status log is often overwritten rather than appended. 'Follow' might show limited updates.")
		}
	} else {
		fmt.Printf("Displaying last %s lines of %s from: %s\n", tailLines, logTypeName, targetLog)
	}
	fmt.Println("---") // Separator

	var cmd *exec.Cmd
	if follow {
		cmd = exec.Command("sudo", "tail", "-f", targetLog)
	} else {
		cmd = exec.Command("sudo", "tail", "-n", tailLines, targetLog)
	}

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if err != nil {
		// Don't use log.Fatalf here as it will exit the main program.
		// Cobra will typically print the error from the RunE function if we return it.
		// For a simple Println, this is fine for now.
		fmt.Printf("Error executing command to show logs: %v\n", err)
		if exitErr, ok := err.(*exec.ExitError); ok {
			fmt.Printf("Command error output: %s\n", string(exitErr.Stderr))
		}
	}
}
