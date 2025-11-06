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
	case "wt_history":
		wb.handleHistory(message)
		//	case "wt_trigger":
		//		wb.handleTrigger(message)
	case "wt_metrics":
		wb.handleMetrics(message)
	case "wt_job":
		wb.handleJob(message)
	case "wt_dashboard":
		wb.handleDashboard(message)
	case "wt_summary":
		wb.handleSummary(message)
	case "wt_update":
		wb.handleUpdate(message)
	default:
		wb.sendMessage(message.Chat.ID,
			"🤖 *WatchtowerMasterBot v2.0*\n\n"+
				"*Server Management:*\n"+
				"`/addserver` - Add new server\n"+
				"`/servers` - List your servers\n"+
				"`/server` - Switch active server\n\n"+
				"*Update Commands:*\n"+
				"`/wt_update` - Trigger container updates\n\n"+
				"*Advanced Features (v1.7+ required):*\n"+
				"`/wt_history` - Update timeline & results\n"+
				"`/wt_metrics` - Performance statistics\n"+
				"`/wt_job` - Detailed job results\n\n"+
				"💡 *Advanced features require Watchtower v1.7+ with HTTP API*")
	}
}

func (wb *WatchtowerBot) handleStart(message *tgbotapi.Message) {
	msg := `🚀 *WatchtowerMasterBot v2.0 Started!*

*Server Management:*
/addserver - Add a new Watchtower server
/servers - List your servers  
/server - Switch active server

"*Update Commands:*\n"+
"/wt_update - Trigger container updates\n"+
"/wt_trigger - Alias for updates\n\n"+
"*Advanced (v1.7+ required):*\n"+
"/wt_history - Update timeline & results\n"+
"/wt_metrics - Statistics & performance\n"+
"/wt_job - Detailed job results"


*Get Started:*
1. Use /addserver to add your first Watchtower instance
2. Switch between severs with /server <name>
3. Use /wt_history to see update status and results`

	wb.sendMessage(message.Chat.ID, msg)
}

func (wb *WatchtowerBot) sendMessage(chatID int64, text string) {
	msg := tgbotapi.NewMessage(chatID, text)
	msg.ParseMode = "Markdown"

	if _, sendErr := wb.bot.Send(msg); sendErr != nil {
		log.Printf("Error sending message: %v", sendErr)
	}
}
