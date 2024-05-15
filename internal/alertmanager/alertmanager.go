package alertmanager

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/a-dev-mobile/server-alert-bot/internal/config"
	"github.com/a-dev-mobile/server-alert-bot/internal/models"
	"golang.org/x/exp/slog"
)

// FetchAlerts извлекает текущие оповещения из AlertManager
func FetchAlerts(cfg *config.Config, lg *slog.Logger) models.AlertData {
	var alertData models.AlertData

	resp, err := http.Get(cfg.AlertManagerURL)
	if err != nil {
		lg.Error("Failed to fetch alerts", "error", err)
		return alertData
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		lg.Error("AlertManager returned non-OK status", "status", resp.Status)
		return alertData
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		lg.Error("Error reading response body:", err)
		return alertData
	}

	if err := json.Unmarshal(body, &alertData); err != nil {
		lg.Error("Error decoding response JSON:", err)
		return alertData
	}

	return alertData
}

// FormatAlertMessage создает читаемое сообщение из оповещения
func FormatAlertMessage(alert models.Alert) string {
	return alert.Annotations.Summary
}
