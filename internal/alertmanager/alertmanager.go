package alertmanager

import (
	"encoding/json"
	"github.com/a-dev-mobile/server-alert-bot/internal/config"
	"github.com/a-dev-mobile/server-alert-bot/internal/models"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
	"strings"
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

// replaceUnsupportedChars заменяет неподдерживаемые символы на аналогичные
func replaceUnsupportedChars(input string) string {
	replacements := map[string]string{
		"-":  "–", // Замена дефиса на тире
		"*":  "•", // Замена звездочки на пуленаправляющий символ
		"_":  "‗", // Замена подчеркивания на двойное подчеркивание
		"[":  "⟦", // Замена открывающей квадратной скобки на другую скобку
		"]":  "⟧", // Замена закрывающей квадратной скобки на другую скобку
		"(":  "❨", // Замена открывающей круглой скобки на другую скобку
		")":  "❩", // Замена закрывающей круглой скобки на другую скобку
		"~":  "˜", // Замена тильды на другой символ
		"\\": "⧵", // Замена обратной косой черты на другой символ
	}
	for oldChar, newChar := range replacements {
		input = strings.ReplaceAll(input, oldChar, newChar)
	}
	return input
}

// FormatAlertMessage создает читаемое сообщение из оповещения
func FormatAlertMessage(alert models.Alert) string {
	var sb strings.Builder
	sb.WriteString("🚨 *Alert!* 🚨\n\n")
	sb.WriteString("*Summary:* " + replaceUnsupportedChars(alert.Annotations.Summary) + "\n")
	sb.WriteString("\n")
	sb.WriteString("*Instance:* " + replaceUnsupportedChars(alert.Labels.Instance) + "\n")
	sb.WriteString("*Job:* " + replaceUnsupportedChars(alert.Labels.Job) + "\n")
	sb.WriteString("*Severity:* " + alert.Labels.Severity + "\n")
	sb.WriteString("*State:* " + alert.State + "\n")
	sb.WriteString("*Active Since:* " + alert.ActiveAt + "\n")
	return sb.String()
}
