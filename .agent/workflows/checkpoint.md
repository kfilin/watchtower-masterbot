---
description: How to save the current state of work (Checkpoint)
---

# 🏁 Checkpoint Workflow

This workflow defines how to save the project state, ensuring that the next agent or session picks up exactly where we left off.

## Trigger

- User types `/checkpoint` or explicitly asks to "save state".
- A major feature or refactor is completed.

## Steps

1. **Context Flush**
    - Stop current coding tasks.
    - Review the session's achievements.

2. **Update `.agent/last_session.md`**
    - **Date**: Today's date.
    - **Goal**: One-line summary of what was being worked on.
    - **Changes**: Bullet points of technical changes (files touched, decisions made).
    - **Next Steps**: Clear instructions for the next session.

3. **Update `.agent/Project-Hub.md`** (If applicable)
    - If the "Gold Standard" has changed (e.g., a new stable version), update the **Gold Standard Checkpoint** section.
    - Update **Current Development Status** checkboxes.

4. **Update `CHANGELOG.md`**
    - Add a new entry under `## [Unreleased]` or the new version number.
    - Follow [Keep a Changelog](https://keepachangelog.com/en/1.0.0/) format.

5. **Commit Changes (Optional but Recommended)**
    - If asked to commit: `git add .`, `git commit -m "checkpoint: <summary>"`

6. **Notify User**
    - "Checkpoint saved. Ready for handoff."
