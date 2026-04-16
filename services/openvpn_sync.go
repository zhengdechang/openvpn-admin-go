package services

import (
	"context"
	"sync"
	"time"

	"openvpn-admin-go/logging"
	"openvpn-admin-go/model"
	"openvpn-admin-go/openvpn"

	"gorm.io/gorm"
)

// createNotification writes a notification record best-effort (outside transaction).
// Errors are logged but never propagate — a lost notification is not worth stopping the sync.
func createNotification(db *gorm.DB, notifType model.NotificationType, userName, realIP, virtualIP string) {
	n := model.Notification{
		Type:      notifType,
		UserName:  userName,
		RealIP:    realIP,
		VirtualIP: virtualIP,
	}
	if err := db.Create(&n).Error; err != nil {
		logging.Error("Failed to create notification for user '%s' (%s): %v", userName, notifType, err)
	}
}

// RunSyncCycle performs a single synchronization cycle of OpenVPN client statuses with the database.
func RunSyncCycle(db *gorm.DB, statusLogPath string) {
	logging.Info("Running OpenVPN sync cycle...")

	parsedClients, _, err := openvpn.ParseStatusLog(statusLogPath)
	if err != nil {
		logging.Error("Error parsing OpenVPN status log: %v. Skipping sync cycle.", err)
		return
	}

	// Step 1: Fetch users currently marked as online in DB
	var dbOnlineUsers []model.User
	if err := db.Where("is_online = ?", true).Find(&dbOnlineUsers).Error; err != nil {
		logging.Error("Error fetching online users from DB: %v. Skipping sync cycle.", err)
		return
	}
	dbOnlineUserMap := make(map[string]model.User)
	for _, u := range dbOnlineUsers {
		dbOnlineUserMap[u.Name] = u
	}

	processedUserNames := make(map[string]bool)
	// newlyConnectedThisCycle tracks users that connected in this cycle (for flapping debounce).
	// If a user both connects and disconnects in the same cycle, we skip the "connected" notification
	// and only emit the "disconnected" one.
	newlyConnectedThisCycle := make(map[string]struct {
		realIP    string
		virtualIP string
	})

	// Step 2: Process clients from the status log (batch update in transaction)
	if err := db.Transaction(func(tx *gorm.DB) error {
		for _, clientStatus := range parsedClients {
			processedUserNames[clientStatus.CommonName] = true
			var user model.User
			result := tx.Where("name = ?", clientStatus.CommonName).First(&user)

			if result.Error == gorm.ErrRecordNotFound {
				logging.Warn("User '%s' not found in DB, skipping sync.", clientStatus.CommonName)
				continue
			} else if result.Error != nil {
				logging.Error("Error fetching user '%s' from DB: %v", clientStatus.CommonName, result.Error)
				continue
			}

			// Detect new connection: was offline before this cycle, now online
			wasOffline := !user.IsOnline
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

			if err := tx.Save(&user).Error; err != nil {
				logging.Error("Error updating user '%s' status: %v", user.Name, err)
			}

			// Record that this user newly connected this cycle (for debounce check in Step 3)
			if wasOffline && clientStatus.IsOnline {
				newlyConnectedThisCycle[clientStatus.CommonName] = struct {
					realIP    string
					virtualIP string
				}{clientStatus.RealAddress, clientStatus.VirtualAddress}
			}
		}
		return nil
	}); err != nil {
		logging.Error("Sync cycle transaction failed: %v", err)
		return
	}

	// Emit "connected" notifications OUTSIDE the transaction, best-effort.
	// Skip users that also disconnected this cycle (flapping debounce).
	for userName, info := range newlyConnectedThisCycle {
		if _, alsoDisconnected := dbOnlineUserMap[userName]; !alsoDisconnected {
			// Was offline before AND still in status log → genuine new connection
			createNotification(db, model.NotificationTypeConnected, userName, info.realIP, info.virtualIP)
			logging.Info("Notification: user '%s' connected from %s", userName, info.realIP)
		}
		// If alsoDisconnected (was online before → implies it was already tracked as online, which
		// contradicts wasOffline check above — this branch is unreachable; left for clarity)
	}

	// Step 3: Mark disconnected users offline (batch in transaction)
	if err := db.Transaction(func(tx *gorm.DB) error {
		for _, dbUser := range dbOnlineUsers {
			if _, found := processedUserNames[dbUser.Name]; !found {
				logging.Info("User '%s' disconnected.", dbUser.Name)
				dbUser.IsOnline = false
				dbUser.RealAddress = ""
				dbUser.VirtualAddress = ""
				dbUser.BytesReceived = 0
				dbUser.BytesSent = 0
				dbUser.OnlineDuration = 0
				dbUser.ConnectedSince = nil
				dbUser.LastRef = nil

				if err := tx.Save(&dbUser).Error; err != nil {
					logging.Error("Error updating user '%s' to offline: %v", dbUser.Name, err)
				}
			}
		}
		return nil
	}); err != nil {
		logging.Error("Disconnect sync transaction failed: %v", err)
	}

	// Emit "disconnected" notifications OUTSIDE the transaction, best-effort.
	for _, dbUser := range dbOnlineUsers {
		if _, found := processedUserNames[dbUser.Name]; !found {
			createNotification(db, model.NotificationTypeDisconnected, dbUser.Name, dbUser.RealAddress, dbUser.VirtualAddress)
			logging.Info("Notification: user '%s' disconnected", dbUser.Name)
		}
	}

	logging.Info("OpenVPN sync cycle finished.")
}

// StartOpenVPNSyncService 启动 OpenVPN 状态同步服务，支持 context 取消和 WaitGroup 优雅退出
func StartOpenVPNSyncService(ctx context.Context, wg *sync.WaitGroup, db *gorm.DB, statusLogPath string, interval time.Duration) {
	logging.Info("Starting OpenVPN Sync Service with interval %s. Log path: %s", interval, statusLogPath)
	wg.Add(1)
	go func() {
		defer wg.Done()
		RunSyncCycle(db, statusLogPath)
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				logging.Info("OpenVPN Sync Service stopping...")
				return
			case <-ticker.C:
				RunSyncCycle(db, statusLogPath)
			}
		}
	}()
}
