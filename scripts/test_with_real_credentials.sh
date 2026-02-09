#!/bin/bash
cd "$(dirname "$0")/.."

echo "🔐 Testing with actual credentials..."

# Check if credentials are set
if [ -z "$TELEGRAM_BOT_TOKEN" ]; then
    echo "❌ TELEGRAM_BOT_TOKEN environment variable is not set"
    echo "💡 Get it from @BotFather on Telegram"
    exit 1
fi

if [ -z "$ADMIN_USER_ID" ]; then
    echo "❌ ADMIN_USER_ID environment variable is not set" 
    echo "💡 Get your user ID from @userinfobot on Telegram"
    exit 1
fi

echo "✅ Credentials found, testing application..."

# Build the application
echo "🔨 Building application..."
go build -o watchtower-masterbot .

# Test run
echo "🚀 Starting application with real credentials..."
TELEGRAM_BOT_TOKEN="$TELEGRAM_BOT_TOKEN" \
ADMIN_USER_ID="$ADMIN_USER_ID" \
HEALTH_PORT="8080" \
ENCRYPTION_KEY="test-encryption-key-$(date +%s)" \
./watchtower-masterbot
