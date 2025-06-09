package openvpn

import (
	"fmt"
	"os"
	"path/filepath"
	"strings" // Required for netmask validation if added, or general string ops
	// "log" // For logging errors if necessary
)

// SetClientFixedIP creates or updates a client-specific configuration file (CCD)
// to assign a fixed IP address to a client.
// commonName is typically the user's ID.
// ipAddress is the fixed IP to assign.
// This function will fetch the serverNetmask from the main OpenVPN configuration.
func SetClientFixedIP(commonName string, ipAddress string) error {
	if commonName == "" {
		return fmt.Errorf("commonName cannot be empty")
	}
	if ipAddress == "" {
		return fmt.Errorf("ipAddress cannot be empty")
	}

	// TODO: Add IP address format validation for ipAddress

	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load OpenVPN configuration: %w", err)
	}

	if cfg.OpenVPNClientConfigDir == "" {
		return fmt.Errorf("OpenVPNClientConfigDir is not set in the configuration")
	}
	if cfg.OpenVPNServerNetmask == "" {
		return fmt.Errorf("OpenVPNServerNetmask is not set in the configuration")
	}

	// Ensure the CCD directory exists
	err = os.MkdirAll(cfg.OpenVPNClientConfigDir, 0755) // rwxr-xr-x
	if err != nil {
		return fmt.Errorf("failed to create client config directory '%s': %w", cfg.OpenVPNClientConfigDir, err)
	}

	ccdFilePath := filepath.Join(cfg.OpenVPNClientConfigDir, commonName)
	content := fmt.Sprintf("ifconfig-push %s %s\n", ipAddress, cfg.OpenVPNServerNetmask)

	err = os.WriteFile(ccdFilePath, []byte(content), 0644) // rw-r--r--
	if err != nil {
		return fmt.Errorf("failed to write client fixed IP config file '%s': %w", ccdFilePath, err)
	}

	// log.Printf("Successfully set fixed IP for %s to %s in %s", commonName, ipAddress, ccdFilePath)
	return nil
}

// RemoveClientFixedIP removes the client-specific configuration file (CCD)
// for the given commonName, effectively removing their fixed IP assignment.
func RemoveClientFixedIP(commonName string) error {
	if commonName == "" {
		return fmt.Errorf("commonName cannot be empty")
	}

	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load OpenVPN configuration: %w", err)
	}

	if cfg.OpenVPNClientConfigDir == "" {
		// If the directory isn't set, there's nothing to remove.
		// Depending on desired strictness, this could be an error or a silent success.
		// log.Printf("OpenVPNClientConfigDir is not set, skipping removal for %s", commonName)
		return nil
	}

	ccdFilePath := filepath.Join(cfg.OpenVPNClientConfigDir, commonName)

	// Check if the file exists before trying to remove it
	if _, err := os.Stat(ccdFilePath); os.IsNotExist(err) {
		// File doesn't exist, so consider it successfully "removed"
		// log.Printf("CCD file for %s does not exist at %s, no action needed.", commonName, ccdFilePath)
		return nil
	} else if err != nil {
		// Other error during stat
		return fmt.Errorf("failed to check client fixed IP config file '%s': %w", ccdFilePath, err)
	}

	err = os.Remove(ccdFilePath)
	if err != nil {
		return fmt.Errorf("failed to remove client fixed IP config file '%s': %w", ccdFilePath, err)
	}

	// log.Printf("Successfully removed fixed IP config for %s from %s", commonName, ccdFilePath)
	return nil
}

// GetClientFixedIP reads the fixed IP address from a client-specific configuration file.
// It returns the IP address string if found, or an empty string if not found or on error.
func GetClientFixedIP(commonName string) (string, error) {
	if commonName == "" {
		return "", fmt.Errorf("commonName cannot be empty")
	}

	cfg, err := LoadConfig()
	if err != nil {
		return "", fmt.Errorf("failed to load OpenVPN configuration: %w", err)
	}

	if cfg.OpenVPNClientConfigDir == "" {
		// log.Printf("OpenVPNClientConfigDir is not set, cannot get fixed IP for %s", commonName)
		return "", nil // Or return an error: fmt.Errorf("OpenVPNClientConfigDir is not set")
	}

	ccdFilePath := filepath.Join(cfg.OpenVPNClientConfigDir, commonName)

	if _, err := os.Stat(ccdFilePath); os.IsNotExist(err) {
		// File doesn't exist, no fixed IP assigned
		return "", nil
	} else if err != nil {
		return "", fmt.Errorf("failed to check client fixed IP config file '%s': %w", ccdFilePath, err)
	}

	contentBytes, err := os.ReadFile(ccdFilePath)
	if err != nil {
		return "", fmt.Errorf("failed to read client fixed IP config file '%s': %w", ccdFilePath, err)
	}

	content := string(contentBytes)
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, "ifconfig-push") {
			parts := strings.Fields(trimmedLine)
			// Expected: "ifconfig-push" "ip_address" "netmask"
			if len(parts) == 3 {
				// TODO: Validate that parts[1] is a valid IP address
				return parts[1], nil
			}
			// log.Printf("Malformed ifconfig-push line in %s for %s: %s", ccdFilePath, commonName, line)
		}
	}

	// log.Printf("No ifconfig-push directive found in %s for %s", ccdFilePath, commonName)
	return "", nil // No ifconfig-push line found
}
