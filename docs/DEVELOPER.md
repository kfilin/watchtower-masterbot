# 👩‍💻 Developer Guide

Welcome to the **Watchtower Masterbot** development team! This guide will help you set up your environment and understand our development workflow.

## 🛠️ Prerequisites

* **Go**: Version 1.21 or higher.
* **Docker**: For containerized testing.
* **Watchtower**: At least one running instance of Watchtower (v1.7.1+) with HTTP API enabled for testing.

## 🏗️ Local Setup Workflow

See `.agent/workflows/local-dev.md` (ToDo) for a step-by-step guide.

1. **Clone the Repo**:

    ```bash
    git clone https://github.com/kfilin/watchtower-masterbot
    cd watchtower-masterbot
    ```

2. **Set Environment Variables**:
    Create a `.env` file or export variables:

    ```bash
    export TELEGRAM_BOT_TOKEN="your_test_token"
    export ADMIN_USER_ID=123456789
    export ENCRYPTION_KEY="developer-key-32-bytes-long-exact!!!!!"
    ```

3. **Run the Bot**:

    ```bash
    go run main.go
    ```

## 🧪 Testing

We use standard Go testing.

* **Run Unit Tests**:

    ```bash
    go test ./...
    ```

* **Run Integration Tests**:
    (Requires configured environment)

    ```bash
    scripts/test_with_real_credentials.sh
    ```

## 📦 Project Structure

We follow a standard Go project layout:

* `cmd/` is omitted as `main.go` is in root for simplicity in this phase.
* `internal/` for private application code.
* `pkg/` (or specific folders like `bot`, `servers`) for library code.

See `docs/files.md` for a complete file map.

## 🤝 Contribution Guidelines

1. **Check the Backlog**: See `.agent/backlog.md`.
2. **Create a Branch**: `feature/your-feature-name`.
3. **Commit**: Use descriptive messages.
4. **Pull Request**: Request review from the maintainer.

## 🔒 Security

* **Never commit tokens**.
* **Always use the `servers/manager.go` encryption methods** when handling user credentials.
* **Memory Only**: Avoid writing sensitive data to disk.
