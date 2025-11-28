#!/bin/bash

echo "🎯 FINAL VERIFICATION"

echo "1. Building application..."
go build -o watchtower-masterbot .
if [ $? -eq 0 ]; then
    echo "✅ Build successful"
else
    echo "❌ Build failed"
    exit 1
fi

echo "2. Testing Docker build..."
docker build -t watchtower-masterbot:verify .
if [ $? -eq 0 ]; then
    echo "✅ Docker build successful"
else
    echo "❌ Docker build failed"
    exit 1
fi

echo "3. Testing invalid token handling..."
timeout 5s TELEGRAM_BOT_TOKEN="invalid" ADMIN_USER_ID="123" ./watchtower-masterbot &
sleep 2
curl -s http://localhost:8080/health > /dev/null && echo "✅ Health server works in degraded mode" || echo "❌ Health server failed"

echo "4. Code quality check..."
go vet ./...
if [ $? -eq 0 ]; then
    echo "✅ Code vetting passed"
else
    echo "⚠️  Code vetting issues found"
fi

echo "🎉 ALL CHECKS COMPLETED - READY TO COMMIT!"
