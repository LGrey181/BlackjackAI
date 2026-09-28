# BlackjackAI

BlackjackAI is a Go program for simulating blackjack hands against a dealer. It includes a basic player AI that places bets, chooses when to hit or stand, and can split or double down when its rules allow.

## Functionality

- Simulates many blackjack hands using one or more decks.
- Tracks the player balance and applies blackjack payouts.
- Uses a card-counting score to adjust the AI's bet size as cards are revealed.
- Supports hit, stand, split, and double-down moves.
- Handles aces, soft hands, blackjacks, dealer busts, pushes, and player busts.
- Reports the final balance after the simulation completes.

## Run

Run the simulation from the project root with the standard Go toolchain. The default executable configuration uses four decks and simulates 50,000 hands before printing the final winnings.

## Package API

The `blackjack` package provides game construction through configurable options, game execution with a player AI, hand scoring, blackjack detection, soft-hand detection, and move functions for the supported player actions.