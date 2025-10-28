package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (wb *WatchtowerBot) handleAddServer(message *tgbotapi.Message) {
	args := strings.Fields(message.CommandArguments())

	if len(args) < 3 {
		wb.sendMessage(message.Chat.ID,
			"❌ Usage: /addserver <nickname> <watchtower_url> <token>\n\n"+
				"Example:\n"+
				"`/addserver home https://wt.home.lan token123`\n"+
				"`/addserver vps https://watchtower.myvps.net token456`")
		return
	}

	nickname := args[0]
	watchtowerURL := args[1]
	token := args[2]

	// Validate URL format
	if !strings.HasPrefix(watchtowerURL, "http") {
		watchtowerURL = "https://" + watchtowerURL
	}

	err := wb.serverManager.AddServer(message.From.ID, nickname, watchtowerURL, token)
	if err != nil {
		wb.sendMessage(message.Chat.ID, fmt.Sprintf("❌ Error adding server: %v", err))
		return
	}

	response := fmt.Sprintf(
		"✅ Server `%s` added successfully!\n\n"+
			"🌐 URL: `%s`\n"+
			"🔑 Token: `%s`\n\n"+
			"Use `/server %s` to switch to this server",
		nickname, watchtowerURL, "••••••••", nickname)

	wb.sendMessage(message.Chat.ID, response)
}

func (wb *WatchtowerBot) handleListServers(message *tgbotapi.Message) {
	serverList, err := wb.serverManager.ListServers(message.From.ID)
	if err != nil {
		wb.sendMessage(message.Chat.ID, "❌ No servers configured. Use /addserver to add your first server.")
		return
	}

	currentServer, err := wb.serverManager.GetCurrentServer(message.From.ID)
	currentServerName := ""
	if err == nil {
		currentServerName = currentServer.Nickname
	}

	var response strings.Builder
	response.WriteString("📋 *Your Watchtower Servers*\n\n")

	for i, server := range serverList {
		indicator := "  "
		if server == currentServerName {
			indicator = "📍"
		}
		response.WriteString(fmt.Sprintf("%s `%s`\n", indicator, server))

		// Limit to 10 servers for clean display
		if i >= 9 {
			response.WriteString("\n... and more")
			break
		}
	}

	response.WriteString("\n\nUse `/server <name>` to switch active server")

	wb.sendMessage(message.Chat.ID, response.String())
}

func (wb *WatchtowerBot) handleSwitchServer(message *tgbotapi.Message) {
	targetServer := strings.TrimSpace(message.CommandArguments())

	if targetServer == "" {
		// Show current server if no argument provided
		currentServer, err := wb.serverManager.GetCurrentServer(message.From.ID)
		if err != nil {
			wb.sendMessage(message.Chat.ID, "❌ No active server. Use /addserver to add your first server.")
			return
		}

		wb.sendMessage(message.Chat.ID,
			fmt.Sprintf("📍 Current server: `%s`\n🌐 %s",
				currentServer.Nickname, currentServer.WatchtowerURL))
		return
	}

	err := wb.serverManager.SwitchServer(message.From.ID, targetServer)
	if err != nil {
		wb.sendMessage(message.Chat.ID,
			fmt.Sprintf("❌ Server `%s` not found. Use /servers to see your available servers.", targetServer))
		return
	}

	wb.sendMessage(message.Chat.ID,
		fmt.Sprintf("🔄 Now managing: `%s`\n\nUse `/wt-status` to check container status.", targetServer))
}

// Watchtower command stubs
func (wb *WatchtowerBot) handleStatus(message *tgbotapi.Message) {
	wb.sendMessage(message.Chat.ID, "🔄 Status command - implementing Watchtower API integration next...")
}

func (wb *WatchtowerBot) handleDashboard(message *tgbotapi.Message) {
	wb.sendMessage(message.Chat.ID, "📊 Dashboard command - implementing multi-format output next...")
}

func (wb *WatchtowerBot) handleSummary(message *tgbotapi.Message) {
	wb.sendMessage(message.Chat.ID, "📈 Summary command - coming soon...")
}

func (wb *WatchtowerBot) handleUpdate(message *tgbotapi.Message) {
	wb.sendMessage(message.Chat.ID, "🔄 Update command - manual trigger coming soon...")
}
