package blackjack

import (
	"fmt"

	"github.com/lgrey181/carddeck"
)

type AI interface {
	Bet(shuffled bool) int
	Play(hand []carddeck.Card, dealer carddeck.Card) Move
	Results(hand [][]carddeck.Card, dealer []carddeck.Card)
}

type Move func(g *Game) error

type dealerAI struct {
}

func (ai dealerAI) Bet(shuffled bool) int {
	return 1
}

func (ai dealerAI) Play(hand []carddeck.Card, dealer carddeck.Card) Move {
	dScore := Score(hand...)
	if dScore <= 16 || (dScore == 17 && Soft(hand...)) {
		return MoveHit
	} else {
		return MoveStand
	}
}

func (ai dealerAI) Results(hands [][]carddeck.Card, dealer []carddeck.Card) {
	//Nothing
}

func HumanAI() AI {
	return humanAI{}
}

type humanAI struct {
}

func (ai humanAI) Bet(shuffled bool) int {
	if shuffled {
		fmt.Println("The deck was just shuffled.")
	}
	fmt.Println("What would you like to bet?")
	var bet int
	fmt.Scanf("%d\n", &bet)
	return bet
}

func (ai humanAI) Play(hand []carddeck.Card, dealer carddeck.Card) Move {
	var input string
	for {
		fmt.Println("Player:", hand)
		fmt.Println("Dealer:", dealer)
		fmt.Println("What will you do? (h)it, (s)tand, (d)ouble, s(p)lit")
		fmt.Scanf("%s\n", &input)
		switch input {
		case "h":
			return MoveHit
		case "s":
			return MoveStand
		case "d":
			return MoveDouble
		case "p":
			return MoveSplit
		default:
			fmt.Println("Invalid option", input)
		}
	}
}

func (ai humanAI) Results(hands [][]carddeck.Card, dealer []carddeck.Card) {
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Player:")
	for _, h := range hands {
		fmt.Println("  ", h)
	}
	fmt.Println("Player:", hands)
	fmt.Println("Dealer", dealer)
}
