package storage

import (
	"log/slog"
	"os"
	"time"
)

func cleanupService(ticker *time.Ticker) {
	for range ticker.C {
		slog.Debug("Running cleanup...")

		expired := ListExpired()
		for _, file := range expired {
			slog.Debug("Deleting expired file", "id", file.ID)
			Delete(file.ID)

			if err := os.Remove(Path(file.ID)); err != nil {
				slog.Error("Failed to delete file on cleanup", "id", file.ID, "error", err)
				continue
			}
		}

		slog.Info("Finished cleanup", "removed", len(expired))
	}
}

func StartCleanupService(cleanupInterval time.Duration, storagePath string) {
	slog.Info("Starting up cleanup service...", "interval", cleanupInterval)

	ticker := time.NewTicker(cleanupInterval)
	go cleanupService(ticker)
}
