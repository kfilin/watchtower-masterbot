# 🏛️ Architectural Decision Records (ADR)

This document records the significant architectural decisions made during the project's evolution.

## ADR-001: Memory-Only Credential Storage

* **Status**: Accepted
* **Date**: 2024-XX-XX
* **Context**: Users need to store API tokens for their Watchtower instances. Storing them in a database creates a liability/attack surface.
* **Decision**: We will store tokens in memory only, encrypted with AES-256. If the bot restarts, users must re-authenticate/re-add servers (or we rely on a secure persistence layer in Phase 3).
* **Consequence**: Higher security, but lower convenience (statelessness).

## ADR-002: Direct HTTP Client over SDK

* **Status**: Accepted
* **Context**: Watchtower's API is simple and effectively stable.
* **Decision**: Use a custom Go HTTP client (`internal/api/watchtower_client.go`) rather than a heavy third-party SDK.
* **Consequence**: Fewer dependencies, tighter control over request context and timeouts.

## ADR-003: Agentic "Brain" Structure

* **Status**: Accepted
* **Date**: 2026-01-30
* **Context**: Need to maintain context across AI coding sessions and mirror the successful `massage-bot` workflow.
* **Decision**: Adopt the `.agent` directory structure with `Collaboration-Blueprint.md` and `Project-Hub.md`.
* **Consequence**: Better continuity for AI-assisted development.
