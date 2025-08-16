# AGENTS

## Architecture Overview

This project implements a small drafting game.

- **Game Core** (`internal/game`)
  - `Turn` orchestrates the two main phases:
    - **Draft phase**: the deck reveals ingredients and the player drafts a subset.
    - **Design phase**: the player combines drafted ingredients into dishes.
  - Events and actions are exchanged with the UI through channels.

- **Domain Packages** (`internal/ingredient`, `internal/dish`, `internal/player`, etc.)
  encapsulate gameplay entities.

## UI

- The terminal UI lives in `internal/ui` and is built with
  [Bubble Tea](https://github.com/charmbracelet/bubbletea).
- The UI runs in a single persistent program and is divided into **modes**.
  Each mode handles rendering and input for a game phase.
    - `draftMode`: displays revealed ingredients and lets the player choose picks.
    - `designMode`: shows drafted ingredients, lets the player name dishes,
      select ingredients, and finish designing.
- Game events are sent to the UI program via a channel, and player actions are
  sent back to the game through another channel.

## Checks

Run the following to verify changes compile:

```bash
go mod tidy
go build ./...
```
