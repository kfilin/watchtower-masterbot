package bot

import (
	"log"

	"watchtower-masterbot/config"
	"watchtower-masterbot/servers"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type WatchtowerBot struct {
	bot           *tgbotapi.BotAPI
	config        *config.Config
	serverManager *servers.ServerManager
}

func NewBot(cfg *config.Config, serverManager *servers.ServerManager) *WatchtowerBot {
	botAPI, err := tgbotapi.NewBotAPI(cfg.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	botAPI.Debug = true

	log.Printf("Authorized on account %s", botAPI.Self.UserName)

	return &WatchtowerBot{
		bot:           botAPI,
		config:        cfg,
		serverManager: serverManager,
	}
}

func (wb *WatchtowerBot) Start() error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := wb.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			wb.handleMessage(update.Message)
		}
	}

	return nil
}

func (wb *WatchtowerBot) handleMessage(message *tgbotapi.Message) {
	// Authorization check
	if message.From.ID != wb.config.AdminUserID {
		wb.sendMessage(message.Chat.ID, "❌ Unauthorized access")
		return
	}

	// Handle commands and button clicks
	switch message.Command() {
	case "start", "bot":
		wb.handleStart(message)
	case "add_server", "addserver": // Support both old and new during transition
		wb.handleAddServer(message)
	case "servers":
		wb.handleListServers(message)
	case "server":
		wb.handleSwitchServer(message)
	case "wt_update", "wtupdate": // Support both old and new
		wb.handleUpdate(message)
	default:
		// Handle button clicks and unknown commands
		wb.handleButtonClicks(message)
	}
}

func (wb *WatchtowerBot) handleStart(message *tgbotapi.Message) {
	welcomeText := `🚀 *Watchtower MasterBot*

*Manage multiple Watchtower instances from one place\!*

📋 *Available Commands:*

• /add\_server \- Add a new Watchtower server
• /servers \- List your managed servers  
• /server \- Switch active server context
• /wt\_update \- Trigger container updates

💡 *Quick Start:*
1\. Use /add\_server to add your first server
2\. Switch between servers with /server
3\. Trigger updates with /wt\_update

🔒 *Security:* All data encrypted with AES\-256`

	msg := tgbotapi.NewMessage(message.Chat.ID, welcomeText)
	msg.ParseMode = "MarkdownV2"
	msg.ReplyMarkup = wb.createMainMenu()
	wb.bot.Send(msg)
}

func (wb *WatchtowerBot) createMainMenu() tgbotapi.ReplyKeyboardMarkup {
	// Professional 2-column menu layout using the proper constructor
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🚀 Add Server"),
			tgbotapi.NewKeyboardButton("📋 Servers List"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("🔄 Switch Server"),
			tgbotapi.NewKeyboardButton("⚡ Update Containers"),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("ℹ️ Bot Info"),
		),
	)
}

func (wb *WatchtowerBot) handleButtonClicks(message *tgbotapi.Message) {
	text := message.Text

	switch text {
	case "🚀 Add Server":
		wb.handleAddServer(message)
	case "📋 Servers List":
		wb.handleListServers(message)
	case "🔄 Switch Server":
		wb.handleSwitchServer(message)
	case "⚡ Update Containers":
		wb.handleUpdate(message)
	case "ℹ️ Bot Info":
		wb.handleStart(message)
	default:
		wb.handleUnknownCommand(message)
	}
}

func (wb *WatchtowerBot) handleUnknownCommand(message *tgbotapi.Message) {
	helpText := `❓ *Unknown Command*

📋 *Available Commands:*
• /add_server - Add new Watchtower server
• /servers - List managed servers
• /server - Switch active server  
• /wt_update - Trigger container updates

💡 Click menu buttons or type commands directly`

	msg := tgbotapi.NewMessage(message.Chat.ID, helpText)
	msg.ParseMode = "Markdown"
	msg.ReplyMarkup = wb.createMainMenu()
	wb.bot.Send(msg)
}

func (wb *WatchtowerBot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"

	if _, sendErr := wb.bot.Send(msg); sendErr != nil {
		log.Printf("Error sending message: %v", sendErr)
	}
}
