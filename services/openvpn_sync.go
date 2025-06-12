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

	// Note: logUpdateTime is parsed but not directly used here yet, as clientStatus.IsOnline and
	// clientStatus.OnlineDurationSeconds are expected to be correctly calculated by ParseStatusLog.
	// It's available if more complex logic involving the log's timestamp is needed later.
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
		dbOnlineUserMap[u.Name] = u // User.Name should store the derived username
	}

	processedUserNames := make(map[string]bool) // To track users processed in this cycle (using derived username)

	// --- Step 2: Process clients from the status log ---
	for _, clientStatus := range parsedClients {
		processedUserNames[clientStatus.CommonName] = true // Use CommonName
		var user model.User
		// Find user by CommonName
		result := db.Where("name = ?", clientStatus.CommonName).First(&user)

		if result.Error == gorm.ErrRecordNotFound {
			log.Printf("User with CommonName '%s' not found in DB. (Log Username: '%s', Log RealAddress: '%s'). Consider creating or logging.", clientStatus.CommonName, clientStatus.Username, clientStatus.RealAddress)
			// For now, we only update existing users.
			continue
		} else if result.Error != nil {
			log.Printf("Error fetching user with CommonName '%s' from DB: %v", clientStatus.CommonName, result.Error)
			continue
		}

		// User found, update status
		user.IsOnline = clientStatus.IsOnline
		if !clientStatus.LastRef.IsZero() {
			user.LastConnectionTime = &clientStatus.LastRef
		} else {
			// Retain existing LastConnectionTime or set to nil if it must be strictly from LastRef
			// For now, if LastRef is zero, we don't update it.
		}

		if err := db.Save(&user).Error; err != nil {
			log.Printf("Error updating user '%s' online status: %v", user.Name, err)
			continue // Skip to next client if user update fails
		}

		// Create or Update ClientLog
		var clientLog model.ClientLog
		// Attempt to find an existing *active* log for this user.
		// An active log is one that is marked IsOnline = true.
		logResult := db.Where("user_id = ? AND is_online = ?", user.ID, true).Last(&clientLog)

		currentTraffic := clientStatus.BytesReceived + clientStatus.BytesSent

		if logResult.Error == gorm.ErrRecordNotFound { // No active log, create new
			newLog := model.ClientLog{
				UserID:         user.ID,
				IsOnline:       clientStatus.IsOnline,
				RealAddress:    clientStatus.RealAddress,
				OnlineDuration: clientStatus.OnlineDurationSeconds, // This is time.Since(ConnectedSince) from parser
				TrafficUsage:   currentTraffic,
			}
			if !clientStatus.ConnectedSince.IsZero() {
				newLog.LastConnectionTime = &clientStatus.ConnectedSince // Log session start time
			}
			// CreatedAt will be set by GORM or BeforeCreate hook
			if err := db.Create(&newLog).Error; err != nil {
				log.Printf("Error creating new ClientLog for user '%s': %v", user.Name, err)
			} else {
				log.Printf("Created new ClientLog for user '%s'. RealAddress: %s", user.Name, newLog.RealAddress)
			}
		} else if logResult.Error == nil { // Active log found, update it
			clientLog.IsOnline = clientStatus.IsOnline // Update online status (could have been reconnected)
			clientLog.RealAddress = clientStatus.RealAddress // Update real address
			clientLog.OnlineDuration = clientStatus.OnlineDurationSeconds // Update duration
			clientLog.TrafficUsage = currentTraffic         // Update with current session's total traffic

			// clientLog.LastConnectionTime should ideally remain the start of this session (ConnectedSince).
			// If clientStatus.ConnectedSince represents the true start of the current continuous session, update it.
			// This assumes clientStatus.ConnectedSince is stable for an ongoing session.
			if !clientStatus.ConnectedSince.IsZero() &&
			   (clientLog.LastConnectionTime == nil || (*clientLog.LastConnectionTime != clientStatus.ConnectedSince && clientStatus.IsOnline)) {
				// Update if ConnectedSince changed, e.g., after a quick reconnect that didn't create a new log entry.
				// Only update if IsOnline is true, to ensure this is for current session.
				clientLog.LastConnectionTime = &clientStatus.ConnectedSince
			}

			if err := db.Save(&clientLog).Error; err != nil {
				log.Printf("Error updating active ClientLog for user '%s': %v", user.Name, err)
			}
		} else { // Some other error fetching the client log
			log.Printf("Error fetching ClientLog for user '%s': %v", user.Name, logResult.Error)
		}
	}

	// --- Step 3: Process users who were in DB as online but are no longer in status log (disconnected) ---
	now := time.Now() // Define 'now' for disconnected client processing
	for _, dbUser := range dbOnlineUsers {
		if _, found := processedUserNames[dbUser.Name]; !found { // User.Name is the derived username
			// This user was online but is no longer in the status log -> disconnected
			log.Printf("User '%s' (Username: %s) disconnected.", dbUser.Name, dbUser.Name)
			dbUser.IsOnline = false
			// We don't update dbUser.LastConnectionTime here; it reflects the last time they were seen (LastRef).
			// If it needs to be the disconnect time:
			// disconnectedTimestamp := now
			// dbUser.LastConnectionTime = &disconnectedTimestamp
			if err := db.Save(&dbUser).Error; err != nil {
				log.Printf("Error updating user '%s' to offline: %v", dbUser.Name, err)
				continue
			}

			// Find their active ClientLog and mark it as offline
			var activeLog model.ClientLog
			// Search for a log that was marked as online for this user.
			if err := db.Where("user_id = ? AND is_online = ?", dbUser.ID, true).Last(&activeLog).Error; err == nil {
				activeLog.IsOnline = false

				// Finalize OnlineDuration based on its stored ConnectedSince (session start).
				if activeLog.LastConnectionTime != nil { // This field stored ConnectedSince for the session
					// Ensure duration is not negative if clocks are skewed or ConnectedSince was in future.
					// Also, using 'now' (disconnect detection time) as end of session.
					sessionEndTime := now
					if sessionEndTime.Before(*activeLog.LastConnectionTime) {
						activeLog.OnlineDuration = 0 // Or log error, ConnectedSince is unexpectedly in future
					} else {
						activeLog.OnlineDuration = int64(sessionEndTime.Sub(*activeLog.LastConnectionTime).Seconds())
					}
				} else {
					// If LastConnectionTime (ConnectedSince) was nil, duration might be suspect or 0.
					// This case should be rare if logs are created correctly.
					activeLog.OnlineDuration = 0
				}
				// TrafficUsage is already the last known cumulative for that session.

				// Set the LastConnectionTime of the log to the actual disconnection time (when we noticed it).
				// This field in ClientLog now represents "session ended at".
                logDisconnectedTime := now
                activeLog.LastConnectionTime = &logDisconnectedTime // Overwrite ConnectedSince with disconnect time

				if err := db.Save(&activeLog).Error; err != nil {
					log.Printf("Error finalizing ClientLog for disconnected user '%s': %v", dbUser.Name, err)
				} else {
					log.Printf("Finalized ClientLog for disconnected user '%s'. Duration: %d s, Traffic: %d bytes, RealAddress: %s", dbUser.Name, activeLog.OnlineDuration, activeLog.TrafficUsage, activeLog.RealAddress)
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
