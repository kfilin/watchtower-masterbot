# 📂 File Manifest

This document provides a map of the project structure and the purpose of each file.

## 🧠 Agent Brain (`.agent/`)

* **`Collaboration-Blueprint.md`**: The definitive guide on AI-Human collaboration protocols.
* **`Project-Hub.md`**: The single source of truth for project status, vision, and tech stack.
* **`last_session.md`**: A log of the most recent work session for context continuity.
* **`workflows/`**: Standard operating procedures (e.g., `checkpoint.md`).
* **`backlog.md`**: The project to-do list and idea parking lot.

## 🤖 Core Application (`/`)

* **`main.go`**: The application entry point. Initializes config, server manager, and starts the bot.
* **`main_test.go`**: Integration tests for the main application flow.
* **`go.mod` / `go.sum`**: Go module definitions and dependency checksums.
* **`Dockerfile`**: Instructions for building the container image.

## 📦 Bot Logic (`bot/`)

* **`bot.go`**: Initializes the Telegram bot API and sets up the update loop.
* **`handlers.go`**: Contains the command handlers (e.g., `/start`, `/addserver`, `/wt_update`).
* **`metrics.go`**: Handles internal metrics collection (if applicable).

## ⚙️ Configuration (`config/`)

* **`config.go`**: Loads and validates environment variables (`TELEGRAM_BOT_TOKEN`, `ENCRYPTION_KEY`, etc.).
* **`config_test.go`**: Unit tests for configuration loading.

## 🏥 Health Checks (`health/`)

* **`health.go`**: Implements health check endpoints (e.g., for Kubernetes probes).
* **`health_test.go`**: Tests for health check logic.

## 🔌 Internal API (`internal/api/`)

* **`watchtower_client.go`**: The HTTP client responsible for communicating with Watchtower instances. Handles API version detection and authentication.

## 🖥️ Server Management (`servers/`)

* **`manager.go`**: The core domain logic. Manages the list of Watchtower servers, handles AES encryption of tokens, and provides thread-safe access.
* **`types.go`**: Defines the `Server` struct and other domain models.

## 🚀 Deployment (`deploy/`)

* **`docker/`**: Docker Compose files.
* **`kubernetes/`**: K8s manifests (Deployments, Services, RBAC, Secrets) and Helm charts.

## 📜 Scripts (`scripts/`)

* **`final_verification.sh`**: End-to-end verification script.
* **`test_with_real_credentials.sh`**: Script for testing with live credentials (use with caution).
