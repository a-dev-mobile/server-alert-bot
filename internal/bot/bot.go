package bot

import (
	"github.com/a-dev-mobile/server-alert-bot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"golang.org/x/exp/slog"

)

// BotService структура для хранения состояния и конфигурации бота
type BotService struct {
	Config *config.Config
	Logger *slog.Logger
	BotAPI *tgbotapi.BotAPI
}

// NewBotService создает новый экземпляр BotService
func NewBotService(cfg *config.Config, logger *slog.Logger, botAPI *tgbotapi.BotAPI) *BotService {
	return &BotService{
		Config: cfg,
		Logger: logger,
		BotAPI: botAPI,
	}
}

// SendStartupMessage отправляет сообщение при запуске бота
func (s *BotService) SendStartupMessage() {
	startupMessage := "🤖 Бот запущен и мониторит систему."
	msg := tgbotapi.NewMessage(s.Config.TelegramChatID, startupMessage)
	msg.ParseMode = "Markdown"
	if _, err := s.BotAPI.Send(msg); err != nil {
		s.Logger.Error("Ошибка при отправке сообщения о запуске:", err)
	}
}

// SendServiceDownMessage отправляет сообщение при недоступности сервиса
func (s *BotService) SendServiceDownMessage(service string) {
	msg := tgbotapi.NewMessage(s.Config.TelegramChatID, "❌ Сервис недоступен: "+service)
	msg.ParseMode = "Markdown"
	if _, err := s.BotAPI.Send(msg); err != nil {
		s.Logger.Error("Ошибка при отправке сообщения о недоступности сервиса:", err)
	}
}

// SendTelegramMessage отправляет сообщение через Telegram
func (s *BotService) SendTelegramMessage(bot *tgbotapi.BotAPI, message string) {
    msg := tgbotapi.NewMessage(s.Config.TelegramChatID, message)
    msg.ParseMode = "Markdown"  // Ensure Markdown parsing is enabled
    if _, err := bot.Send(msg); err != nil {
        s.Logger.Error("Ошибка отправки сообщения:", err)
    }
}

// SendNoAlertsMessage отправляет сообщение об отсутствии тревог
func (s *BotService) SendNoAlertsMessage(bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(s.Config.TelegramChatID, "✅ Все системы функционируют стабильно.")
	msg.ParseMode = "Markdown"
	if _, err := bot.Send(msg); err != nil {
		s.Logger.Error("Ошибка при отправке сообщения об отсутствии тревог:", err)
	}
}

// SendAlertResolvedMessage отправляет сообщение об устранении всех тревог
func (s *BotService) SendAlertResolvedMessage(bot *tgbotapi.BotAPI) {
	msg := tgbotapi.NewMessage(s.Config.TelegramChatID, "🛠️ Все проблемы были успешно устранены.")
	msg.ParseMode = "Markdown"
	if _, err := bot.Send(msg); err != nil {
		s.Logger.Error("Ошибка при отправке сообщения о решении проблем:", err)
	}
}
