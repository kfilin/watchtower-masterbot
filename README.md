
# WatchtowerMasterBot

🚀 **Multi-server Watchtower management via Telegram** - Phase 2 Complete ✅

A production-ready Telegram bot for managing multiple Watchtower instances with secure credential storage and real API integration.

## 🎯 Current Status: **Production-Ready Foundation**

### What Works Today

- **Multi-server management** with encrypted credential storage (AES-256)
- **Real Watchtower API integration** (v1.7.1 HTTP API)
- **Secure architecture** with memory-only token processing
- **Professional Telegram interface** with educational error handling
- **Adaptive API client** with runtime endpoint discovery

## 🏗️ Architecture Highlights

- **Go 1.21+** for performance and single-binary deployment
- **User-isolated server management** with thread-safe operations
- **Security-first implementation** with no persistent plaintext storage
- **Progressive enhancement** approach - core works today, advanced features ready

## 🚀 Quick Start

### Prerequisites

- Go 1.21+
- Telegram Bot Token from [@BotFather](https://t.me/botfather)
- Watchtower instance with HTTP API enabled

### Installation

```bash
# Clone the repository
git clone https://github.com/kfilin/watchtower-masterbot
cd watchtower-masterbot

# Set environment variables
export TELEGRAM_BOT_TOKEN="your_bot_token_from_botfather"
export ADMIN_USER_ID=304528450  # Your Telegram user ID
export ENCRYPTION_KEY="secure-encryption-key-change-in-production"

# Run the bot
go run main.go
```

### Docker (Coming Soon)

```bash
docker build -t watchtower-masterbot .
docker run -e TELEGRAM_BOT_TOKEN="your_token" watchtower-masterbot
```

## 📋 Key Commands

### Server Management

```text
/addserver <name> <url> <token>  - Add new Watchtower instance
/servers                         - List all configured servers  
/server <name>                   - Switch active server context
```

### Watchtower Commands

```text
/wt_update    - Trigger manual container updates
/wt_status    - Check Watchtower instance status
/wt_history   - View update timeline and results
/wt_metrics   - Performance statistics (v1.7+ required)
/wt_job       - Detailed job results (v1.7+ required)
```

## 🔧 Configuration

### Environment Variables

```bash
# Required
TELEGRAM_BOT_TOKEN=your_bot_token_here
ADMIN_USER_ID=304528450

# Optional (with defaults)
ENCRYPTION_KEY=default-key-change-in-production
PORT=8443
WEBHOOK_URL=your_webhook_url
```

### Adding Your First Server

1. Start chat with your bot in Telegram
2. Use `/addserver home https://your-watchtower-url your-token`
3. Switch with `/server home`
4. Trigger updates with `/wt_update`

## 🛡️ Security Features

- **AES-256 Encryption**: All tokens encrypted at rest
- **Memory-Only Processing**: Tokens decrypted only during API calls
- **User Isolation**: Each user manages their own servers
- **Input Validation**: All inputs sanitized and validated
- **No Shell Commands**: Pure HTTP API integration only

## 🏗️ Technical Architecture

Detailed in [docs/files.md](docs/files.md) and [docs/ARCH_DECISIONS.md](docs/ARCH_DECISIONS.md).

```text
watchtower-masterbot/
├── bot/
│   ├── bot.go           # Telegram integration & command routing
│   └── handlers.go      # Command implementations
├── servers/
│   ├── manager.go       # Multi-server management & encryption
│   └── types.go         # Data structures
├── internal/api/
│   └── watchtower_client.go  # Adaptive Watchtower API client
├── config/
│   └── config.go        # Environment configuration
└── main.go              # Application entry point
```

## 📚 Comprehensive Documentation

Complete knowledge base organized in Obsidian with:

- **Architectural decisions** and rationales
- **Security implementation** deep dive  
- **Technical challenge** solutions
- **Blog-ready content** and interview materials
- **Code patterns** and snippets library

## 🚧 Development Status

### ✅ Phase 2 Complete

- [x] Multi-server architecture with encrypted storage
- [x] Real Watchtower API integration
- [x] Comprehensive error handling and user guidance
- [x] Production-ready code quality
- [x] Complete documentation foundation

### 🔄 Phase 3 Planned

- [ ] Docker containerization and image building
- [ ] Kubernetes deployment manifests
- [ ] CI/CD pipeline implementation
- [ ] Enhanced features (update completion tracking)
- [ ] Advanced notifications and analytics

## 🤝 Contributing

See [docs/DEVELOPER.md](docs/DEVELOPER.md) for detailed setup and contribution guidelines.

1. Check the Backlog in [.agent/backlog.md](.agent/backlog.md).
2. Fork the repository
3. Create a feature branch (`git checkout -b feature/amazing-feature`)
4. Commit your changes (`git commit -m 'Add amazing feature'`)
5. Push to the branch (`git push origin feature/amazing-feature`)
6. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Watchtower](https://containrrr.dev/watchtower/) for container update automation
- [Go Telegram Bot API](https://github.com/go-telegram-bot-api/telegram-bot-api) for Telegram integration
- Built with security and user experience as primary concerns

---

**Ready for production deployment and enhanced features!** 🚀

## 🚀 Deployment

### Docker Deployment

```bash
# Build from source
docker build -t watchtower-masterbot .

# Or use docker-compose
docker-compose -f deploy/docker/docker-compose.yml up

Kubernetes Deployment
bash

# Apply Kubernetes manifests
kubectl apply -f deploy/kubernetes/k8s/

# Or use Helm
helm install watchtower-masterbot deploy/kubernetes/helm/watchtower-masterbot/

