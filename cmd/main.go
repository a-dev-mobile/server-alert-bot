package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/a-dev-mobile/server-alert-bot/internal/alertmanager"
	"github.com/a-dev-mobile/server-alert-bot/internal/bot"
	"github.com/a-dev-mobile/server-alert-bot/internal/config"
	"github.com/a-dev-mobile/server-alert-bot/internal/logging"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/slog"
)

func main() {
	cfg, lg := initializeApp()

	// Создание объекта бота с использованием токена
	botAPI, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		log.Fatalf("Panic due to failure in creating bot API: %s", err)
	}

	botAPI.Debug = false
	lg.Info("Authorized on account", "username", botAPI.Self.UserName)

	// Проверка сервисов при запуске
	if err := checkServices(cfg, lg); err != nil {
		botService := bot.NewBotService(cfg, lg, botAPI)
		botService.SendServiceDownMessage(err.Error())
		log.Fatalf("Service check failed: %s", err)
	}

	// Отправка стартового сообщения при запуске бота
	botService := bot.NewBotService(cfg, lg, botAPI)
	botService.SendStartupMessage()

	// Настройка таймера для регулярного опроса
	ticker := time.NewTicker(cfg.PollingInterval)
	defer ticker.Stop()

	// Мапа для отслеживания времени последней отправки уведомлений
	lastSent := make(map[string]time.Time)
	// Переменная для отслеживания активных тревог
	var alertsActive bool = false

	// Основной цикл опроса
	for range ticker.C {
		alerts := alertmanager.FetchAlerts(cfg, lg)
		if len(alerts.Data.Alerts) > 0 {
			if !alertsActive {
				alertsActive = true // Флаг активных тревог устанавливается в true
				lg.Debug("Активное состояние оповещений изменено", "active", alertsActive)
			}
			for _, alert := range alerts.Data.Alerts {
				// уникальный ключ для каждого оповещения
				key := alert.Labels.Alertname + alert.Labels.Instance
				lastTime, exists := lastSent[key]
				// Проверка, нужно ли отправить повторное уведомление
				if !exists || time.Since(lastTime) >= cfg.RepeatNotificationInterval {
					message := alertmanager.FormatAlertMessage(alert)
					botService.SendTelegramMessage(botAPI, message)
					lastSent[key] = time.Now()
					lg.Debug("Alert sent", "alert", message)
				}
			}
		} else if alertsActive {
			alertsActive = false                        // Сброс флага активных тревог
			lastSent = make(map[string]time.Time)       // Сброс времени последней отправки
			botService.SendAlertResolvedMessage(botAPI) // Сообщение об устранении тревог
			botService.SendNoAlertsMessage(botAPI)      // Сообщение о стабильной работе системы
			lg.Debug("No more active alerts")
		}
	}
}

func initializeApp() (*config.Config, *slog.Logger) {
	cfg := getConfigOrFail()
	lg := logging.SetupLogger(cfg)
	return cfg, lg
}

func getConfigOrFail() *config.Config {
	cfg, err := config.LoadConfig("../config", "config.yaml")
	if err != nil {
		log.Fatalf("Error loading config: %s", err)
	}
	return cfg
}

func checkServices(cfg *config.Config, lg *slog.Logger) error {
	services := []string{cfg.AlertManagerURL}
	for _, url := range services {
		resp, err := http.Get(url)
		if err != nil || resp.StatusCode != 200 {
			lg.Error("Service is down", "url", url, "status", resp.Status)
			return fmt.Errorf("service at %s is down", url)
		}
	}
	return nil
}
