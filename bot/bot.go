package bot

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/kfilin/watchtower-masterbot/servers"
)

// WatchtowerBot matches the receiver name in your handlers.go
type WatchtowerBot struct {
	API           *tgbotapi.BotAPI
	AdminID       int64
	serverManager *servers.ServerManager
	webAppURL     string
}

// GetManager returns the internal ServerManager
func (wb *WatchtowerBot) GetManager() *servers.ServerManager {
	return wb.serverManager
}

// NewBot initializes the bot without panicking
func NewBot(token string, adminID int64, encryptionKey string, webAppURL string) (*WatchtowerBot, error) {
	if token == "" {
		return nil, fmt.Errorf("TELEGRAM_BOT_TOKEN is missing")
	}

	api, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, fmt.Errorf("failed to authenticate with Telegram: %w", err)
	}

	api.Debug = false

	// Initialize the ServerManager with encryption key
	mgr := servers.NewManager(encryptionKey)

	return &WatchtowerBot{
		API:           api,
		AdminID:       adminID,
		serverManager: mgr,
		webAppURL:     webAppURL,
	}, nil
}

// Start begins the update loop
func (wb *WatchtowerBot) Start() {
	log.Printf("🤖 Authorized on account %s", wb.API.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := wb.API.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Security Check
		if wb.AdminID != 0 && update.Message.From.ID != wb.AdminID {
			log.Printf("🔒 Security: Ignored message from unauthorized user %d (%s)",
				update.Message.From.ID, update.Message.From.UserName)
			continue
		}

		log.Printf("📥 Received message: '%s' from %s", update.Message.Text, update.Message.From.UserName)
		wb.Handle(update)
	}
}

// Handle dispatches commands to methods defined in handlers.go
func (wb *WatchtowerBot) Handle(update tgbotapi.Update) {
	msg := update.Message
	cmd := msg.Command()
	text := msg.Text

	switch {
	case cmd == "start":
		wb.showMainMenu(msg.Chat.ID)
	case cmd == "add_server" || text == "🚀 Add Server":
		wb.handleAddServer(msg)
	case cmd == "servers" || text == "📋 List Servers":
		wb.handleListServers(msg)
	case cmd == "server" || text == "🔄 Switch Server":
		wb.handleSwitchServer(msg)
	case cmd == "wt_update":
		wb.handleUpdate(msg)
	case cmd == "terminal":
		wb.handleTerminal(msg)
	default:
		// Unknown command, show menu
		log.Printf("❓ Unknown command: %s", cmd)
		wb.showMainMenu(msg.Chat.ID)
	}
}

// sendMessage is a helper used by handlers.go
func (wb *WatchtowerBot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"

	// Create persistent keyboard menu
	keyboard := tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🚀 Add Server"),
			tgbotapi.NewKeyboardButton("🔄 Switch Server"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("📋 List Servers"),
		),
	)
	msg.ReplyMarkup = keyboard

	if _, err := wb.API.Send(msg); err != nil {
		log.Printf("❌ Telegram API Error: %v | ChatID: %d | Text: [%s]", err, chatID, text)
	} else {
		log.Printf("📤 Sent message to chat %d", chatID)
	}
}

// showMainMenu displays the welcome message
func (wb *WatchtowerBot) showMainMenu(chatID int64) {
	log.Printf("📱 Showing main menu to chat %d", chatID)
	text := "🚀 *Watchtower MasterBot*\n\n" +
		"Manage multiple Watchtower instances from one place!\n\n" +
		"📋 *Available Commands:*\n\n" +
		"• `/add_server` - Add a new Watchtower server\n" +
		"• `/servers` - List your managed servers\n" +
		"• `/server` - Switch active server context\n" +
		"• `/wt_update` - Trigger container updates\n" +
		"• `/terminal` - 📟 Access Advanced Terminal\n\n" +
		"💡 *Quick Start:*\n" +
		"1. Use `/add_server` to add your first server\n" +
		"2. Switch between servers with `/server`\n" +
		"3. Trigger updates with `/wt_update`\n\n" +
		"🔒 *Security:* All data encrypted with AES-256"
	wb.sendMessage(chatID, text)
}
