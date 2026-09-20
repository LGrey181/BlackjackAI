package main

import (
	"BlackjackAI/blackjack"
	"fmt"

	"github.com/lgrey181/carddeck"
)

type basicAI struct {
	blackjackScore int
	seen           int
	decks          int
}

func (ai *basicAI) Bet(shuffled bool) int {
	if shuffled {
		ai.blackjackScore = 0
		ai.seen = 0
	}
	trueScore := ai.blackjackScore / ((ai.decks*52 - ai.seen) / 52)
	switch {
	case trueScore > 14:
		return 10000
	case trueScore > 8:
		return 5000
	default:
		return 1000
	}
}

func (ai *basicAI) Play(hand []carddeck.Card, dealer carddeck.Card) blackjack.Move {
	score := blackjack.Score(hand...)
	if len(hand) == 2 {
		if hand[0] == hand[1] {
			cardScore := blackjack.Score(hand[0])
			if cardScore >= 8 && cardScore != 10 {
				return blackjack.MoveSplit
			}
		}
		if (score == 10 || score == 11) && !blackjack.Soft(hand...) {
			return blackjack.MoveDouble
		}
	}
	dScore := blackjack.Score(dealer)
	if dScore >= 5 && dScore <= 6 {
		return blackjack.MoveStand
	}
	if score < 13 {
		return blackjack.MoveHit
	}
	return blackjack.MoveStand
}

func (ai *basicAI) Results(hand [][]carddeck.Card, dealer []carddeck.Card) {
	for _, card := range dealer {
		ai.count(card)
	}
	for _, h := range hand {
		for _, c := range h {
			ai.count(c)
		}
	}
}

func (ai *basicAI) count(card carddeck.Card) {
	score := blackjack.Score(card)
	switch {
	case score >= 10:
		ai.blackjackScore--
	case score <= 6:
		ai.blackjackScore++
	}
	ai.seen++
}

func main() {
	opts := blackjack.Options{
		NumDecks:        4,
		HandsNum:        50000,
		BlackjackPayout: 1.5,
	}
	game := blackjack.New(opts)
	winnings := game.Play(&basicAI{
		decks: 4,
	})
	fmt.Println(winnings)
}
