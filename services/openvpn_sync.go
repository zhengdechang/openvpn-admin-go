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

	parsedClients, err := openvpn.ParseStatusLog(statusLogPath)
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
		dbOnlineUserMap[u.Name] = u // Assuming User.Name is the CommonName
	}

	processedUserNames := make(map[string]bool) // To track users processed in this cycle

	// --- Step 2: Process clients from the status log ---
	for _, clientStatus := range parsedClients {
		processedUserNames[clientStatus.CommonName] = true
		var user model.User
		// Try to find user by CommonName (assuming it maps to User.Name)
		result := db.Where("name = ?", clientStatus.CommonName).First(&user)

		now := time.Now()

		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("User with CommonName '%s' not found in DB. Consider creating or logging.", clientStatus.CommonName)
			// Optionally, create a new user here if that's desired behavior.
			// For now, we only update existing users.
			continue
		} else if result.Error != nil {
			log.Printf("Error fetching user '%s' from DB: %v", clientStatus.CommonName, result.Error)
			continue
		}

		// User found, update status
		user.IsOnline = true
		user.LastConnectionTime = &clientStatus.LastRef // LastRef from status log

		if err := db.Save(&user).Error; err != nil {
			log.Printf("Error updating user '%s' online status: %v", user.Name, err)
			continue // Skip to next client if user update fails
		}

		// Create or Update ClientLog
		var clientLog model.ClientLog
		logResult := db.Where("user_id = ? AND is_online = ?", user.ID, true).Last(&clientLog)

		currentTraffic := clientStatus.BytesReceived + clientStatus.BytesSent
		currentOnlineDuration := int64(0)
		if !clientStatus.ConnectedSince.IsZero() {
			currentOnlineDuration = int64(now.Sub(clientStatus.ConnectedSince).Seconds())
		}

		if logResult.Error == gorm.ErrRecordNotFound { // No active log, create new
			clientLog = model.ClientLog{
				UserID:             user.ID,
				IsOnline:           true,
				OnlineDuration:     currentOnlineDuration,
				TrafficUsage:       currentTraffic,
				LastConnectionTime: &clientStatus.ConnectedSince, // Log when this session started
				// CreatedAt will be set by GORM
			}
			if err := db.Create(&clientLog).Error; err != nil {
				log.Printf("Error creating new ClientLog for user '%s': %v", user.Name, err)
			} else {
				log.Printf("Created new ClientLog for connected user '%s'.", user.Name)
			}
		} else if logResult.Error == nil { // Active log found, update it
			clientLog.OnlineDuration = currentOnlineDuration // Update duration
			clientLog.TrafficUsage = currentTraffic          // Update with current session's total traffic
			// LastConnectionTime in ClientLog could mean two things:
			// 1. The start of this specific online session (clientStatus.ConnectedSince)
			// 2. The last time we saw this client (clientStatus.LastRef)
			// Let's use clientStatus.ConnectedSince for the log's LastConnectionTime to mark session start.
			// User.LastConnectionTime has LastRef.
			clientLog.LastConnectionTime = &clientStatus.ConnectedSince
			if err := db.Save(&clientLog).Error; err != nil {
				log.Printf("Error updating active ClientLog for user '%s': %v", user.Name, err)
			}
		} else { // Some other error fetching the client log
			log.Printf("Error fetching ClientLog for user '%s': %v", user.Name, logResult.Error)
		}
	}

	// --- Step 3: Process users who were in DB as online but are no longer in status log (disconnected) ---
	for _, dbUser := range dbOnlineUsers {
		if _, found := processedUserNames[dbUser.Name]; !found {
			// This user was online but is no longer in the status log -> disconnected
			log.Printf("User '%s' disconnected.", dbUser.Name)
			dbUser.IsOnline = false
			// LastConnectionTime for User model could be set to now, or keep the LastRef from previous cycle.
			// For now, we leave User.LastConnectionTime as it was (updated by LastRef when they were online).
			// If we want to mark exact disconnect time for User model:
			// disconnectedTime := time.Now()
			// dbUser.LastConnectionTime = &disconnectedTime
			if err := db.Save(&dbUser).Error; err != nil {
				log.Printf("Error updating user '%s' to offline: %v", dbUser.Name, err)
				continue
			}

			// Find their active ClientLog and mark it as offline
			var activeLog model.ClientLog
			if err := db.Where("user_id = ? AND is_online = ?", dbUser.ID, true).Last(&activeLog).Error; err == nil {
				activeLog.IsOnline = false
				// Finalize OnlineDuration: time since the log was created or since ConnectedSince
				// If activeLog.LastConnectionTime stores ConnectedSince:
				if activeLog.LastConnectionTime != nil {
					finalDuration := int64(time.Now().Sub(*activeLog.LastConnectionTime).Seconds())
					activeLog.OnlineDuration = finalDuration
				}
				// TrafficUsage is already the last known cumulative for that session.

				// Set the LastConnectionTime of the log to the actual disconnection time
				// This field in ClientLog now better represents "session ended at" or "last seen for this session"
				logDisconnectedTime := time.Now()
				activeLog.LastConnectionTime = &logDisconnectedTime

				if err := db.Save(&activeLog).Error; err != nil {
					log.Printf("Error finalizing ClientLog for disconnected user '%s': %v", dbUser.Name, err)
				} else {
					log.Printf("Finalized ClientLog for disconnected user '%s'. Duration: %d s, Traffic: %d bytes", dbUser.Name, activeLog.OnlineDuration, activeLog.TrafficUsage)
				}
			} else if err != gorm.ErrRecordNotFound {
				log.Printf("Error finding active ClientLog for disconnected user '%s': %v", dbUser.Name, err)
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

// Example of how it might be called in main.go (not part of this subtask's file changes)
/*
func main() {
    // ... other initializations ...
    database.Init() // Initialize DB
    db := database.DB

    // AutoMigrate models if you haven't already
    // db.AutoMigrate(&model.User{}, &model.ClientLog{})

    statusLog := "/var/log/openvpn/status.log" // Or from config
    syncInterval := 5 * time.Minute          // Or from config

    services.StartOpenVPNSyncService(db, statusLog, syncInterval)

    // ... start your web server etc. ...
	// select {} // Block main goroutine if needed
}
*/

// Helper function (if needed, e.g. for more complex duration calculation, but time.Sub.Seconds() is fine)
// func calculateDuration(startTime time.Time) int64 {
// 	return int64(time.Since(startTime).Seconds())
// }
