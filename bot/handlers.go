package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (wb *WatchtowerBot) handleAddServer(message *tgbotapi.Message) {
	args := strings.Fields(message.CommandArguments())

	if len(args) < 3 {
		// Show usage example when no arguments provided
		wb.sendMessage(message.Chat.ID,
			"## 📥 Add a Watchtower Server\n\n"+
				"**Usage:** `/addserver <nickname> <watchtower_url> <token>`\n\n"+
				"**Examples:**\n"+
				"`/addserver home https://watchtower.local your_token_here`\n"+
				"`/addserver vps https://watchtower.yourserver.com abc123token`\n\n"+
				"**Parameters:**\n"+
				"• **nickname** - Easy-to-remember name (e.g., home, vps, work)\n"+
				"• **watchtower_url** - Your Watchtower API endpoint\n"+
				"• **token** - Watchtower authentication token")
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
		"✅ Server **%s** added successfully!\n\n"+
			"🌐 **URL:** `%s`\n"+
			"🔑 **Token:** `%s`\n\n"+
			"Use `/server %s` to switch to this server or `/wt_status` to check status",
		nickname, watchtowerURL, "••••••••", nickname)

	wb.sendMessage(message.Chat.ID, response)
}

func (wb *WatchtowerBot) handleListServers(message *tgbotapi.Message) {
	serverList, err := wb.serverManager.ListServers(message.From.ID)
	if err != nil {
		wb.sendMessage(message.Chat.ID, "❌ No servers configured. Use **/addserver** to add your first server.")
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
		// Make server names non-clickable by removing backticks
		response.WriteString(fmt.Sprintf("%s %s\n", indicator, server))

		// Limit to 10 servers for clean display
		if i >= 9 {
			response.WriteString("\n... and more")
			break
		}
	}

	response.WriteString("\n\n*Quick Actions:*\n")
	response.WriteString("• Use `/server <name>` to switch active server\n")
	response.WriteString("• Use `/wt_status` to check current server status\n")
	response.WriteString("• Use `/addserver` to add more servers")

	wb.sendMessage(message.Chat.ID, response.String())
}

func (wb *WatchtowerBot) handleSwitchServer(message *tgbotapi.Message) {
	targetServer := strings.TrimSpace(message.CommandArguments())

	if targetServer == "" {
		// Show current server if no argument provided
		currentServer, err := wb.serverManager.GetCurrentServer(message.From.ID)
		if err != nil {
			wb.sendMessage(message.Chat.ID,
				"❌ No active server.\n\n"+
					"Use **/addserver** to add your first Watchtower server.")
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
			"Use **/wt_status** to check container status or **/servers** to switch again.", targetServer))
}

// Watchtower command with real API integration
func (wb *WatchtowerBot) handleStatus(message *tgbotapi.Message) {
	currentServer, err := wb.serverManager.GetCurrentServer(message.From.ID)
	if err != nil {
		wb.sendMessage(message.Chat.ID,
			"❌ No active server configured.\n\n"+
				"Use **/addserver** to add your first Watchtower server.")
		return
	}

	wb.sendMessage(message.Chat.ID,
		fmt.Sprintf("🔄 *Checking Watchtower status...*\n\n"+
			"🌐 Server: `%s`\n"+
			"📡 URL: %s",
			currentServer.Nickname, currentServer.WatchtowerURL))

	// Get API client
	client, err := wb.serverManager.GetAPIClient(message.From.ID)
	if err != nil {
		wb.sendMessage(message.Chat.ID,
			fmt.Sprintf("❌ Failed to create API client: %v", err))
		return
	}

	// Get basic status
	status, err := client.GetStatus()
	if err != nil {
		wb.sendMessage(message.Chat.ID,
			fmt.Sprintf("❌ Failed to get status: %v", err))
		return
	}

	response := fmt.Sprintf("📊 *Watchtower Status*\n\n"+
		"🌐 Server: `%s`\n"+
		"🔄 Version: %s\n"+
		"📈 Status: %s\n\n"+
		"⚠️  *Note:* Watchtower API doesn't provide container status.\n"+
		"Use **/wt_update** to trigger manual update and see results.",
		currentServer.Nickname, status.Version, status.Status)

	wb.sendMessage(message.Chat.ID, response)
}

// Update handleUpdate to use the real API
func (wb *WatchtowerBot) handleUpdate(message *tgbotapi.Message) {
	currentServer, err := wb.serverManager.GetCurrentServer(message.From.ID)
	if err != nil {
		wb.sendMessage(message.Chat.ID,
			"❌ No active server configured.\n\n"+
				"Use **/addserver** to add your first Watchtower server.")
		return
	}

	wb.sendMessage(message.Chat.ID,
		fmt.Sprintf("🚀 *Triggering manual update...*\n\n"+
			"🌐 Server: `%s`\n"+
			"📡 URL: %s\n\n"+
			"This may take a few minutes...",
			currentServer.Nickname, currentServer.WatchtowerURL))

	// Get API client and trigger real update
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

	// Format the response
	var response strings.Builder
	response.WriteString("✅ *Update Triggered Successfully!*\n\n")

	if len(updateResponse.Updated) > 0 {
		response.WriteString("🔄 *Updated Containers:*\n")
		for _, container := range updateResponse.Updated {
			response.WriteString(fmt.Sprintf("• %s\n", container))
		}
		response.WriteString("\n")
	}

	if len(updateResponse.Failed) > 0 {
		response.WriteString("❌ *Failed Containers:*\n")
		for _, container := range updateResponse.Failed {
			response.WriteString(fmt.Sprintf("• %s\n", container))
		}
		response.WriteString("\n")
	}

	if len(updateResponse.Updated) == 0 && len(updateResponse.Failed) == 0 {
		response.WriteString("📋 No containers were updated or failed.\n")
		response.WriteString("All containers are up to date! ✅")
	}

	wb.sendMessage(message.Chat.ID, response.String())
}

func (wb *WatchtowerBot) handleDashboard(message *tgbotapi.Message) {
	currentServer, err := wb.serverManager.GetCurrentServer(message.From.ID)
	if err != nil {
		wb.sendMessage(message.Chat.ID,
			"❌ No active server configured.\n\n"+
				"Use **/addserver** to add your first Watchtower server.")
		return
	}

	wb.sendMessage(message.Chat.ID,
		fmt.Sprintf("📊 *Dashboard - Coming Soon*\n\n"+
			"🌐 Server: `%s`\n"+
			"📡 URL: %s\n\n"+
			"*Feature Status:* 🚧 Under Development\n\n"+
			"Use **/wt_status** for current status or **/wt_update** to trigger updates.",
			currentServer.Nickname, currentServer.WatchtowerURL))
}

func (wb *WatchtowerBot) handleSummary(message *tgbotapi.Message) {
	currentServer, err := wb.serverManager.GetCurrentServer(message.From.ID)
	if err != nil {
		wb.sendMessage(message.Chat.ID,
			"❌ No active server configured.\n\n"+
				"Use **/addserver** to add your first Watchtower server.")
		return
	}

	wb.sendMessage(message.Chat.ID,
		fmt.Sprintf("📈 *Update Summary - Coming Soon*\n\n"+
			"🌐 Server: `%s`\n"+
			"📡 URL: %s\n\n"+
			"*Feature Status:* 🚧 Under Development\n\n"+
			"Use **/wt_update** to trigger updates and see recent changes.",
			currentServer.Nickname, currentServer.WatchtowerURL))
}
