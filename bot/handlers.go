package bot

import (
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (wb *WatchtowerBot) handleAddServer(message *tgbotapi.Message) {
	// Handle button click
	if message.Text == "🚀 Add Server" {
		wb.sendMessage(message.Chat.ID,
			"📥 *Add a Watchtower Server*\n\n"+
				"*Usage:* `/add_server <nickname> <watchtower_url> <token>`\n\n"+
				"*Examples:*\n"+
				"• `/add_server home https://watchtower.local your_token_here`\n"+
				"• `/add_server vps https://watchtower.yourserver.com abc123token`\n\n"+
				"*Parameters:*\n"+
				"• *nickname* - Easy-to-remember name (home, vps, work)\n"+
				"• *watchtower_url* - Your Watchtower API endpoint\n"+
				"• *token* - Watchtower authentication token")
		return
	}

	args := strings.Fields(message.CommandArguments())

	if len(args) < 3 {
		wb.sendMessage(message.Chat.ID,
			"📥 *Add a Watchtower Server*\n\n"+
				"*Usage:* `/add_server <nickname> <watchtower_url> <token>`\n\n"+
				"*Examples:*\n"+
				"• `/add_server home https://watchtower.local your_token_here`\n"+
				"• `/add_server vps https://watchtower.yourserver.com abc123token`\n\n"+
				"*Parameters:*\n"+
				"• *nickname* - Easy-to-remember name (home, vps, work)\n"+
				"• *watchtower_url* - Your Watchtower API endpoint\n"+
				"• *token* - Watchtower authentication token")
		return
	}

	nickname := args[0]
	watchtowerURL := args[1]
	token := args[2]

	if !strings.HasPrefix(watchtowerURL, "http") {
		watchtowerURL = "https://" + watchtowerURL
	}

	err := wb.serverManager.AddServer(message.From.ID, nickname, watchtowerURL, token)
	if err != nil {
		wb.sendMessage(message.Chat.ID, fmt.Sprintf("❌ Error adding server: %v", err))
		return
	}

	response := fmt.Sprintf(
		"✅ *Server %s added successfully!*\n\n"+
			"🌐 *URL:* `%s`\n"+
			"🔑 *Token:* `%s`\n\n"+
			"Use `/server %s` to switch to this server or `/servers` to see all servers",
		nickname, watchtowerURL, "••••••••", nickname)

	wb.sendMessage(message.Chat.ID, response)
}

func (wb *WatchtowerBot) handleListServers(message *tgbotapi.Message) {
	serverList, err := wb.serverManager.ListServers(message.From.ID)
	if err != nil {
		wb.sendMessage(message.Chat.ID, "❌ No servers configured. Use /add_server to add your first server.")
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
		// Fixed alignment - consistent spacing with monospace
		response.WriteString(fmt.Sprintf("%s `%s`\n", indicator, server))

		if i >= 9 {
			response.WriteString("\n... and more")
			break
		}
	}

	response.WriteString("\n\n🚀 *Quick Actions:*\n")
	response.WriteString("• Use `/server <name>` to switch active server\n")
	response.WriteString("• Use `/wt_update` to trigger container updates\n")
	response.WriteString("• Use `/add_server` to add more servers")

	wb.sendMessage(message.Chat.ID, response.String())
}

func (wb *WatchtowerBot) handleSwitchServer(message *tgbotapi.Message) {
	// Handle button click
	if message.Text == "🔄 Switch Server" {
		currentServer, err := wb.serverManager.GetCurrentServer(message.From.ID)
		if err != nil {
			wb.sendMessage(message.Chat.ID,
				"🔄 *Switch Active Server*\n\n"+
					"*Usage:* `/server <server_name>`\n\n"+
					"*Example:* `/server home`\n\n"+
					"Use `/servers` to see your available servers.")
			return
		}

		wb.sendMessage(message.Chat.ID,
			fmt.Sprintf("📍 *Current Server:* `%s`\n\n"+
				"*To switch servers:*\n"+
				"`/server <server_name>`\n\n"+
				"Use `/servers` to see all available servers.",
				currentServer.Nickname))
		return
	}

	targetServer := strings.TrimSpace(message.CommandArguments())

	if targetServer == "" {
		currentServer, err := wb.serverManager.GetCurrentServer(message.From.ID)
		if err != nil {
			wb.sendMessage(message.Chat.ID,
				"❌ No active server.\n\n"+
					"Use **/add_server** to add your first Watchtower server.")
			return
		}

		wb.sendMessage(message.Chat.ID,
			fmt.Sprintf("📍 *Current Server:* `%s`\n🌐 *URL:* %s\n\n"+
				"Use `/server <name>` to switch servers or `/servers` to see all servers.",
				currentServer.Nickname, currentServer.WatchtowerURL))
		return
	}

	err := wb.serverManager.SwitchServer(message.From.ID, targetServer)
	if err != nil {
		wb.sendMessage(message.Chat.ID,
			fmt.Sprintf("❌ Server `%s` not found.\n\n"+
				"Use **/servers** to see your available servers.", targetServer))
		return
	}

	wb.sendMessage(message.Chat.ID,
		fmt.Sprintf("🔄 *Now managing:* `%s`\n\n"+
			"Use **/wt_update** to trigger updates or **/servers** to switch again.", targetServer))
}

func (wb *WatchtowerBot) handleUpdate(message *tgbotapi.Message) {
	currentServer, err := wb.serverManager.GetCurrentServer(message.From.ID)
	if err != nil {
		wb.sendMessage(message.Chat.ID,
			"❌ No active server configured.\n\n"+
				"Use **/add_server** to add your first Watchtower server.")
		return
	}

	// Send immediate feedback
	wb.sendMessage(message.Chat.ID,
		fmt.Sprintf("🚀 *Triggering container update...*\n\n"+
			"🌐 Server: `%s`\n"+
			"📡 URL: %s\n\n"+
			"⏱️ *This may take 2-5 minutes...*\n"+
			"I'll notify you when complete.",
			currentServer.Nickname, currentServer.WatchtowerURL))

	client, err := wb.serverManager.GetAPIClient(message.From.ID)
	if err != nil {
		wb.sendMessage(message.Chat.ID,
			fmt.Sprintf("❌ Failed to create API client: %v", err))
		return
	}

	updateResponse, err := client.TriggerUpdate()
	if err != nil {
		wb.sendMessage(message.Chat.ID,
			fmt.Sprintf("❌ Failed to trigger update: %v", err))
		return
	}

	// Build response EXACTLY as requested
	var response strings.Builder
	response.WriteString("✅ *Update Triggered Successfully!*\n\n")
	response.WriteString("📋 *Result:* Watchtower is checking for updates\n\n")

	// Show container counts and names exactly as requested
	updatedCount := len(updateResponse.Updated)
	failedCount := len(updateResponse.Failed)

	response.WriteString(fmt.Sprintf("🔄 *Containers updated:* `%d`\n", updatedCount))

	if updatedCount > 0 {
		containerNames := strings.Join(updateResponse.Updated, ", ")
		response.WriteString(fmt.Sprintf("📦 *Container(s) name:* `%s`\n", containerNames))
	} else {
		response.WriteString("📦 *Container(s) name:* `none`\n")
	}

	if failedCount > 0 {
		response.WriteString(fmt.Sprintf("\n❌ *Containers failed:* `%d`\n", failedCount))
		failedNames := strings.Join(updateResponse.Failed, ", ")
		response.WriteString(fmt.Sprintf("💥 *Failed container(s):* `%s`\n", failedNames))
	}

	response.WriteString("\n🔍 *Use `/servers` to manage your servers*")

	wb.sendMessage(message.Chat.ID, response.String())
}

func (wb *WatchtowerBot) handleTerminal(message *tgbotapi.Message) {
	if wb.webAppURL == "" {
		wb.sendMessage(message.Chat.ID, "❌ *Retro Terminal* is not configured.\n\n"+
			"Please set `WEBAPP_URL` in your `.env` file.\n\n"+
			"💡 *Note:* This must be a public HTTPS URL pointing to this bot's port (default :8080) at the `/terminal` path.\n\n"+
			"Example: `WEBAPP_URL=https://your-tunnel.ngrok-free.app/terminal`")
		return
	}

	msg := tgbotapi.NewMessage(message.Chat.ID, "📟 *ACCESSING MASTER CONTROL...*\n\nWelcome to the specialized retro terminal interface.")
	msg.ParseMode = "Markdown"

	// Create Inline Keyboard with WebApp button
	markup := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonWebApp("🔓 Open Terminal", tgbotapi.WebAppInfo{URL: wb.webAppURL}),
		),
	)
	msg.ReplyMarkup = markup

	if _, err := wb.API.Send(msg); err != nil {
		log.Printf("❌ Failed to send terminal message: %v", err)
	}
}
