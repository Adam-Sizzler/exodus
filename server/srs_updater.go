package server

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"cerberus-node/config"
)

const srsRefreshInterval = 12 * time.Hour

func (s *NodeServer) startSRSAutoUpdater() {
	if s == nil || s.Cfg == nil {
		return
	}

	s.Cfg.Logger.Info("SRS auto updater started", "interval", srsRefreshInterval.String())

	go func() {
		ticker := time.NewTicker(srsRefreshInterval)
		defer ticker.Stop()

		for range ticker.C {
			s.refreshSRSFromManifest()
		}
	}()
}

func (s *NodeServer) refreshSRSFromManifest() {
	manifestPath := filepath.Join(config.FixedSingboxDir, srsManifestName)

	lists, err := loadSRSManifest()
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			s.Cfg.Logger.Debug("SRS manifest not found, skipping periodic refresh", "path", manifestPath)
			return
		}
		s.Cfg.Logger.Warn("Failed to read SRS manifest for periodic refresh", "path", manifestPath, "error", err)
		return
	}
	if len(lists) == 0 {
		s.Cfg.Logger.Debug("SRS manifest is empty, skipping periodic refresh", "path", manifestPath)
		return
	}

	summary, syncErr := s.SyncSRSLists(lists)
	if syncErr != nil {
		s.Cfg.Logger.Warn("Periodic SRS refresh failed", "error", syncErr)
		return
	}

	s.Cfg.Logger.Info(
		"Periodic SRS refresh completed",
		"total", summary.Total,
		"configured", summary.Configured,
		"downloaded", summary.Downloaded,
		"failed", summary.Failed,
	)
}
