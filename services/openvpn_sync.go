package services

import (
	"log"
	"time"

	"gorm.io/gorm"
	"openvpn-admin-go/model"
	"openvpn-admin-go/openvpn" // Assuming the parser is in this package
)

// RunSyncCycle performs a single synchronization cycle of OpenVPN client statuses with the database.
func RunSyncCycle(db *gorm.DB, statusLogPath string) {
	log.Println("Running OpenVPN sync cycle...")

	parsedClients, _, err := openvpn.ParseStatusLog(statusLogPath)
	if err != nil {
		log.Printf("Error parsing OpenVPN status log: %v. Skipping sync cycle.", err)
		return
	}

	// --- Step 1: Fetch users currently marked as online in DB ---
	var dbOnlineUsers []model.User
	if err := db.Where("is_online = ?", true).Find(&dbOnlineUsers).Error; err != nil {
		log.Printf("Error fetching online users from DB: %v. Skipping sync cycle.", err)
		return
	}
	dbOnlineUserMap := make(map[string]model.User)
	for _, u := range dbOnlineUsers {
		dbOnlineUserMap[u.Name] = u
	}

	processedUserNames := make(map[string]bool)

	// --- Step 2: Process clients from the status log ---
	for _, clientStatus := range parsedClients {
		processedUserNames[clientStatus.CommonName] = true
		var user model.User
		result := db.Where("name = ?", clientStatus.CommonName).First(&user)

		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("User with CommonName '%s' not found in DB. (Log Username: '%s', Log RealAddress: '%s'). Consider creating or logging.", clientStatus.CommonName, clientStatus.Username, clientStatus.RealAddress)
			continue
		} else if result.Error != nil {
			log.Printf("Error fetching user with CommonName '%s' from DB: %v", clientStatus.CommonName, result.Error)
			continue
		}

		// Update all user status fields
		user.IsOnline = clientStatus.IsOnline
		user.RealAddress = clientStatus.RealAddress
		user.VirtualAddress = clientStatus.VirtualAddress
		user.BytesReceived = clientStatus.BytesReceived
		user.BytesSent = clientStatus.BytesSent
		user.OnlineDuration = clientStatus.OnlineDurationSeconds
		
		if !clientStatus.ConnectedSince.IsZero() {
			user.ConnectedSince = &clientStatus.ConnectedSince
		}
		if !clientStatus.LastRef.IsZero() {
			user.LastRef = &clientStatus.LastRef
			user.LastConnectionTime = &clientStatus.LastRef
		}

		if err := db.Save(&user).Error; err != nil {
			log.Printf("Error updating user '%s' status: %v", user.Name, err)
			continue
		}
	}

	// --- Step 3: Process users who were in DB as online but are no longer in status log (disconnected) ---
	now := time.Now()
	for _, dbUser := range dbOnlineUsers {
		if _, found := processedUserNames[dbUser.Name]; !found {
			log.Printf("User '%s' disconnected.", dbUser.Name)
			
			// Reset all online status fields
			dbUser.IsOnline = false
			dbUser.RealAddress = ""
			dbUser.VirtualAddress = ""
			dbUser.BytesReceived = 0
			dbUser.BytesSent = 0
			dbUser.OnlineDuration = 0
			dbUser.ConnectedSince = nil
			dbUser.LastRef = nil

			if err := db.Save(&dbUser).Error; err != nil {
				log.Printf("Error updating user '%s' to offline: %v", dbUser.Name, err)
				continue
			}
		}
	}
	log.Println("OpenVPN sync cycle finished.")
}

// StartOpenVPNSyncService launches a goroutine to periodically sync OpenVPN client statuses.
func StartOpenVPNSyncService(db *gorm.DB, statusLogPath string, interval time.Duration) {
	log.Printf("Starting OpenVPN Sync Service with interval %s. Log path: %s", interval, statusLogPath)
	go func() {
		// Run once immediately at start, then tick.
		RunSyncCycle(db, statusLogPath)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			RunSyncCycle(db, statusLogPath)
		}
	}()
}

