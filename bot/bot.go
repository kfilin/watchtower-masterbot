package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"watchtower-masterbot/config"
	"watchtower-masterbot/servers"
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

	// Handle commands
	switch message.Command() {
	case "start":
		wb.handleStart(message)
	case "addserver":
		wb.handleAddServer(message)
	case "servers":
		wb.handleListServers(message)
	case "server":
		wb.handleSwitchServer(message)
	case "wt_status":
		wb.handleStatus(message)
	case "wt_dashboard":
		wb.handleDashboard(message)
	case "wt_summary":
		wb.handleSummary(message)
	case "wt_update":
		wb.handleUpdate(message)
	default:
		wb.sendMessage(message.Chat.ID, 
				"🤖 *WatchtowerMasterBot*\n\n" +
				"*Server Management:*\n" +
				"`/addserver` - Add new server\n" +
				"`/servers` - List your servers\n" +
				"`/server` - Switch active server\n\n" +
				"*Watchtower Commands:*\n" +
				"`/wt_status` - Container status\n" +
				"`/wt_dashboard` - Overview\n" +
				"`/wt_summary` - Update history\n" +
				"`/wt_update` - Manual updates")
	}
}

func (wb *WatchtowerBot) handleStart(message *tgbotapi.Message) {
	msg := `🚀 *WatchtowerMasterBot Started!*

*Available Commands:*
/addserver - Add a new Watchtower server
/servers - List your servers  
/server - Switch active server
/wt_status - Check container status

*Get Started:*
1. Use /addserver to add your first Watchtower instance
2. Switch between servers with /server <name>
3. Use /wt_status to check container status`

	wb.sendMessage(message.Chat.ID, msg)
}

func (wb *WatchtowerBot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"

	if _, sendErr := wb.bot.Send(msg); sendErr != nil {
		log.Printf("Error sending message: %v", sendErr)
	}
}
