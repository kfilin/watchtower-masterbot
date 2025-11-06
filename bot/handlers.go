package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (wb *WatchtowerBot) handleAddServer(message *tgbotapi.Message) {
	args := strings.Fields(message.CommandArguments())

	if len(args) < 3 {
		// Clean Markdown without headers or complex formatting
		wb.sendMessage(message.Chat.ID,
			"📥 *Add a Watchtower Server*\n\n"+
				"*Usage:* `/addserver <nickname> <watchtower_url> <token>`\n\n"+
				"*Examples:*\n"+
				"• `/addserver home https://watchtower.local your_token_here`\n"+
				"• `/addserver vps https://watchtower.yourserver.com abc123token`\n\n"+
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
			"Use `/server %s` to switch to this server or `/wt_history` to check status",
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
		response.WriteString(fmt.Sprintf("%s %s\n", indicator, server))

		if i >= 9 {
			response.WriteString("\n... and more")
			break
		}
	}

	response.WriteString("\n\n*Quick Actions:*\n")
	response.WriteString("• Use `/server <name>` to switch active server\n")
	response.WriteString("• Use `/wt_history` to check current server status\n")
	response.WriteString("• Use `/addserver` to add more servers")

	wb.sendMessage(message.Chat.ID, response.String())
}

func (wb *WatchtowerBot) handleSwitchServer(message *tgbotapi.Message) {
	targetServer := strings.TrimSpace(message.CommandArguments())

	if targetServer == "" {
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
			"Use **/wt_history** to check update history or **/servers** to switch again.", targetServer))
}

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

	response.WriteString("\n\n💡 *Use `/wt_history` to monitor progress*")

	wb.sendMessage(message.Chat.ID, response.String())
}

func (wb *WatchtowerBot) handleHistory(message *tgbotapi.Message) {
	currentServer, err := wb.serverManager.GetCurrentServer(message.From.ID)
	if err != nil {
		wb.sendMessage(message.Chat.ID,
			"❌ No active server configured.\n\n"+
				"Use **/addserver** to add your first Watchtower server.")
		return
	}

	wb.sendMessage(message.Chat.ID,
		fmt.Sprintf("📊 *Checking update history...*\n\n"+
			"🌐 Server: `%s`\n"+
			"📡 URL: %s",
			currentServer.Nickname, currentServer.WatchtowerURL))

	client, err := wb.serverManager.GetAPIClient(message.From.ID)
	if err != nil {
		wb.sendMessage(message.Chat.ID,
			fmt.Sprintf("❌ Failed to create API client: %v", err))
		return
	}

	updateJobs, err := client.GetUpdateJobs(5)
	if err != nil {
		// Provide specific guidance based on error type
		errorMsg := err.Error()

		if strings.Contains(errorMsg, "not available") {
			wb.sendMessage(message.Chat.ID,
				"ℹ️  *Update History Not Available*\n\n"+
					"Your Watchtower instance doesn't support update history.\n\n"+
					"*To enable this feature:*\n"+
					"• Update to Watchtower v1.7+\n"+
					"• Add `--http-api` to your Watchtower command\n"+
					"• Ensure API is properly configured\n\n"+
					"*What works now:*\n"+
					"• `/wt_trigger` - Manual updates\n"+
					"• `/wt_update` - Legacy command")
		} else if strings.Contains(errorMsg, "timeout") {
			wb.sendMessage(message.Chat.ID,
				"⏱️  *Server Timeout*\n\n"+
					"Watchtower server is not responding.\n\n"+
					"*Troubleshooting:*\n"+
					"• Check if Watchtower is running\n"+
					"• Verify your URL and token\n"+
					"• Ensure HTTP API is enabled\n"+
					"• Check server connectivity")
		} else {
			wb.sendMessage(message.Chat.ID,
				fmt.Sprintf("❌ Failed to get update history: %v", err))
		}
		return
	}

	if len(updateJobs) == 0 {
		wb.sendMessage(message.Chat.ID,
			"📭 *No Update History*\n\n"+
				"No update jobs found in history.\n\n"+
				"This could mean:\n"+
				"• Watchtower hasn't run any updates yet\n"+
				"• Updates were run before history was enabled\n"+
				"• All containers are up to date\n\n"+
				"Use **/wt_trigger** to run an update now.")
		return
	}

	// Display the history (existing code)
	var response strings.Builder
	response.WriteString(fmt.Sprintf("🕐 *Update History - %s*\n\n", currentServer.Nickname))

	for _, job := range updateJobs {
		var statusEmoji string
		switch job.State {
		case "failed":
			statusEmoji = "❌"
		case "running":
			statusEmoji = "🔄"
		default:
			statusEmoji = "✅"
		}

		jobIDShort := job.ID
		if len(jobIDShort) > 8 {
			jobIDShort = jobIDShort[:8] + "..."
		}

		startTime := job.Started.Format("Jan 02 15:04")
		var endTime string
		if job.Ended.IsZero() {
			endTime = "Ongoing"
		} else {
			endTime = job.Ended.Format("Jan 02 15:04")
		}

		response.WriteString(fmt.Sprintf("%s *Job %s*\n", statusEmoji, jobIDShort))
		response.WriteString(fmt.Sprintf("   Status: `%s`\n", job.State))
		response.WriteString(fmt.Sprintf("   Started: `%s`\n", startTime))
		response.WriteString(fmt.Sprintf("   Ended: `%s`\n", endTime))

		if len(job.Results) > 0 {
			successCount := 0
			failedCount := 0
			for _, result := range job.Results {
				if result.Status == "success" {
					successCount++
				} else {
					failedCount++
				}
			}
			response.WriteString(fmt.Sprintf("   Results: ✅ %d | ❌ %d\n", successCount, failedCount))
		}
		response.WriteString("\n")
	}

	response.WriteString("💡 *Use `/wt_job <job_id>` for detailed container results*")

	wb.sendMessage(message.Chat.ID, response.String())
}

func (wb *WatchtowerBot) handleMetrics(message *tgbotapi.Message) {
	currentServer, err := wb.serverManager.GetCurrentServer(message.From.ID)
	if err != nil {
		wb.sendMessage(message.Chat.ID,
			"❌ No active server configured.\n\n"+
				"Use **/addserver** to add your first Watchtower server.")
		return
	}

	wb.sendMessage(message.Chat.ID,
		fmt.Sprintf("📈 *Fetching metrics...*\n\n"+
			"🌐 Server: `%s`\n"+
			"📡 URL: %s",
			currentServer.Nickname, currentServer.WatchtowerURL))

	client, err := wb.serverManager.GetAPIClient(message.From.ID)
	if err != nil {
		wb.sendMessage(message.Chat.ID,
			fmt.Sprintf("❌ Failed to create API client: %v", err))
		return
	}

	metrics, err := client.GetMetrics()
	if err != nil {
		if strings.Contains(err.Error(), "not available") {
			wb.sendMessage(message.Chat.ID,
				"ℹ️  *Metrics Not Available*\n\n"+
					"Your Watchtower instance doesn't expose metrics.\n\n"+
					"*To enable metrics:*\n"+
					"• Update to Watchtower v1.7+\n"+
					"• Add `--http-api` and `--metrics` flags\n"+
					"• Configure Prometheus endpoint\n\n"+
					"Metrics provide insights into update frequency and success rates.")
		} else {
			wb.sendMessage(message.Chat.ID,
				fmt.Sprintf("❌ Failed to get metrics: %v", err))
		}
		return
	}

	var response strings.Builder
	response.WriteString(fmt.Sprintf("📊 *Metrics - %s*\n\n", currentServer.Nickname))

	updateCount := "0"
	failureCount := "0"
	lastUpdate := "Never"

	for key, value := range metrics.Data {
		switch {
		case strings.Contains(key, "watchtower_updated"):
			updateCount = value
		case strings.Contains(key, "watchtower_failed"):
			failureCount = value
		case strings.Contains(key, "watchtower_last_run"):
			if value != "0" {
				lastUpdate = value
			}
		}
	}

	response.WriteString(fmt.Sprintf("✅ *Total Updates:* `%s`\n", updateCount))
	response.WriteString(fmt.Sprintf("❌ *Total Failures:* `%s`\n", failureCount))
	response.WriteString(fmt.Sprintf("🕐 *Last Run:* `%s`\n\n", lastUpdate))

	response.WriteString("💡 *Metrics show overall Watchtower performance*\n")
	response.WriteString("Use `/wt_history` for recent update details")

	wb.sendMessage(message.Chat.ID, response.String())
}

func (wb *WatchtowerBot) handleJob(message *tgbotapi.Message) {
	jobID := strings.TrimSpace(message.CommandArguments())

	if jobID == "" {
		wb.sendMessage(message.Chat.ID,
			"🔍 *Job Details*\n\n"+
				"**Usage:** `/wt_job <job_id>`\n\n"+
				"**Example:** `/wt_job abc123def456`\n\n"+
				"Get job IDs from `/wt_history` command\n\n"+
				"💡 *Note:* This requires Watchtower v1.7+ with HTTP API enabled")
		return
	}

	currentServer, err := wb.serverManager.GetCurrentServer(message.From.ID)
	if err != nil {
		wb.sendMessage(message.Chat.ID,
			"❌ No active server configured.\n\n"+
				"Use **/addserver** to add your first Watchtower server.")
		return
	}

	wb.sendMessage(message.Chat.ID,
		fmt.Sprintf("🔍 *Fetching job details...*\n\n"+
			"🌐 Server: `%s`\n"+
			"📋 Job ID: `%s`",
			currentServer.Nickname, jobID))

	client, err := wb.serverManager.GetAPIClient(message.From.ID)
	if err != nil {
		wb.sendMessage(message.Chat.ID,
			fmt.Sprintf("❌ Failed to create API client: %v", err))
		return
	}

	job, err := client.GetUpdateJob(jobID)
	if err != nil {
		if strings.Contains(err.Error(), "not available") {
			wb.sendMessage(message.Chat.ID,
				"ℹ️  *Job Details Not Available*\n\n"+
					"Your Watchtower instance doesn't support detailed job tracking.\n\n"+
					"*To enable this feature:*\n"+
					"• Update to Watchtower v1.7+\n"+
					"• Add `--http-api` to your configuration\n\n"+
					"Job details show individual container update results.")
		} else {
			wb.sendMessage(message.Chat.ID,
				fmt.Sprintf("❌ Failed to get job details: %v", err))
		}
		return
	}

	// Display job details (existing code)
	var response strings.Builder
	response.WriteString(fmt.Sprintf("📋 *Job Details - %s*\n\n", jobID))

	var statusEmoji string
	switch job.State {
	case "failed":
		statusEmoji = "❌"
	case "running":
		statusEmoji = "🔄"
	default:
		statusEmoji = "✅"
	}

	response.WriteString(fmt.Sprintf("%s *Status:* `%s`\n", statusEmoji, job.State))
	response.WriteString(fmt.Sprintf("🕐 *Started:* `%s`\n", job.Started.Format("Jan 02, 2006 15:04:05")))

	if job.Ended.IsZero() {
		response.WriteString("🕐 *Ended:* `Ongoing`\n")
	} else {
		response.WriteString(fmt.Sprintf("🕐 *Ended:* `%s`\n", job.Ended.Format("Jan 02, 2006 15:04:05")))
	}

	if len(job.Results) > 0 {
		response.WriteString(fmt.Sprintf("\n📊 *Results (%d containers):*\n\n", len(job.Results)))

		successCount := 0
		for _, result := range job.Results {
			resultEmoji := "✅"
			if result.Status != "success" {
				resultEmoji = "❌"
			} else {
				successCount++
			}

			response.WriteString(fmt.Sprintf("%s *%s*\n", resultEmoji, result.Container))
			response.WriteString(fmt.Sprintf("   Status: `%s`\n", result.Status))
			if result.Error != "" {
				response.WriteString(fmt.Sprintf("   Error: `%s`\n", result.Error))
			}
			response.WriteString("\n")
		}

		response.WriteString(fmt.Sprintf("📈 *Summary:* ✅ %d/%d successful\n", successCount, len(job.Results)))
	} else {
		response.WriteString("\n📊 *Results:* No container results available\n")
	}

	wb.sendMessage(message.Chat.ID, response.String())
}

func (wb *WatchtowerBot) handleStatus(message *tgbotapi.Message) {
	wb.sendMessage(message.Chat.ID,
		"🔄 *Command Updated*\n\n"+
			"The `/wt_status` command has been replaced with `/wt_history` which shows actual update results!\n\n"+
			"**New Commands:**\n"+
			"• `/wt_history` - Recent update jobs & results\n"+
			"• `/wt_trigger` - Manual update trigger\n"+
			"• `/wt_metrics` - Update statistics\n"+
			"• `/wt_job` - Detailed job results")
}

func (wb *WatchtowerBot) handleDashboard(message *tgbotapi.Message) {
	wb.sendMessage(message.Chat.ID,
		"📊 *Dashboard Features Now Available!*\n\n"+
			"All dashboard functionality is now available through:\n\n"+
			"• `/wt_history` - Update timeline & status\n"+
			"• `/wt_metrics` - Performance statistics\n"+
			"• `/wt_job` - Detailed container results\n\n"+
			"No need for separate dashboard command!")
}

func (wb *WatchtowerBot) handleSummary(message *tgbotapi.Message) {
	wb.sendMessage(message.Chat.ID,
		"📈 *Summary Features Enhanced!*\n\n"+
			"Update summaries are now available through:\n\n"+
			"• `/wt_history` - Recent update overview\n"+
			"• `/wt_metrics` - Historical performance\n"+
			"• `/wt_job` - Detailed per-job summaries\n\n"+
			"These commands provide more detailed and actionable information!")
}
