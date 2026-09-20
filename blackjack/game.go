package blackjack

import (
	"errors"

	"github.com/lgrey181/carddeck"
)

type state uint8

type Options struct {
	NumDecks        int
	HandsNum        int
	BlackjackPayout float64
}

func New(opts Options) Game {
	g := Game{
		state:    statePlayerTurn,
		dealerAI: &dealerAI{},
		balance:  0,
	}
	if opts.NumDecks == 0 {
		opts.NumDecks = 3
	}
	if opts.HandsNum == 0 {
		opts.HandsNum = 100
	}
	if opts.BlackjackPayout == 0.0 {
		opts.BlackjackPayout = 1.5
	}
	g.nDecks = opts.NumDecks
	g.nHands = opts.HandsNum
	g.blackjackPayout = opts.BlackjackPayout
	return g
}

const (
	stateBet state = iota
	statePlayerTurn
	stateDealerTurn
	stateHandOver
)

type Game struct {
	//unexported fields
	nDecks  int
	nHands  int
	balance int

	deck  []carddeck.Card
	state state

	player    []hand
	handInd   int
	playerBet int

	dealer          []carddeck.Card
	dealerAI        AI
	blackjackPayout float64
}

type hand struct {
	card []carddeck.Card
	bet  int
}

var (
	errbust = errors.New("Hand score exceeded 21!!")
)

func MoveHit(g *Game) error {
	hand := g.currentHand()
	var card carddeck.Card
	card, g.deck = draw(g.deck)
	*hand = append(*hand, card)
	if Score(*hand...) > 21 {
		return errbust
	}
	return nil
}

func MoveSplit(g *Game) error {
	cards := g.currentHand()
	if len(*cards) != 2 {
		return errors.New("You can only split on two of the same cards in your hand.")
	}
	if (*cards)[0].Rank != (*cards)[1].Rank {
		return errors.New("Both cards must have the same ranl to split.")
	}
	g.player = append(g.player, hand{
		card: []carddeck.Card{(*cards)[1]},
		bet:  g.player[g.handInd].bet,
	})
	g.player[g.handInd].card = (*cards)[:1]
	return nil
}

func MoveStand(g *Game) error {
	if g.state == stateDealerTurn {
		g.state++
		return nil
	}
	if g.state == statePlayerTurn {
		g.handInd++
		if g.handInd >= len(g.player) {
			g.state++
		}
		return nil
	}
	return errors.New("Invalid state")
}

func MoveDouble(g *Game) error {
	if len(*g.currentHand()) != 2 {
		return errors.New("Can only double on a hand with 2 cards")
	}
	g.playerBet *= 2
	MoveHit(g)
	return MoveStand(g)
}

func bet(g *Game, ai AI, shuffled bool) {
	bet := ai.Bet(shuffled)
	if bet < 100 {
		panic("Bet must be at least $100.")
	}
	g.playerBet = bet
}

func draw(cards []carddeck.Card) (carddeck.Card, []carddeck.Card) {
	return cards[0], cards[1:]
}

func (g *Game) Play(ai AI) int {
	g.deck = nil
	min := 52 * g.nDecks / 3
	for i := 0; i < g.nHands; i++ {
		shuffled := false
		if len(g.deck) < min {
			g.deck = carddeck.New(carddeck.Deck(g.nDecks), carddeck.Shuffle)
			shuffled = true
		}
		bet(g, ai, shuffled)
		deal(g)
		if Blackjack(g.dealer...) {
			endHand(g, ai)
			continue
		}
		for g.state == statePlayerTurn {
			hand := make([]carddeck.Card, len(*g.currentHand()))
			copy(hand, *g.currentHand())
			move := ai.Play(hand, g.dealer[0])
			err := move(g)
			switch err {
			case errbust:
				MoveStand(g)
			case nil:
				// noop
			default:
				panic(err)
			}
		}

		for g.state == stateDealerTurn {
			hand := make([]carddeck.Card, len(g.dealer))
			copy(hand, g.dealer)
			move := g.dealerAI.Play(hand, g.dealer[0])
			move(g)
		}

		endHand(g, ai)
	}
	return g.balance
}

func deal(g *Game) {
	g.handInd = 0
	playerHand := make([]carddeck.Card, 0, 5)
	g.dealer = make([]carddeck.Card, 0, 5)
	var individualCard carddeck.Card
	for i := 0; i < 2; i++ {
		individualCard, g.deck = draw(g.deck)
		playerHand = append(playerHand, individualCard)
		individualCard, g.deck = draw(g.deck)
		g.dealer = append(g.dealer, individualCard)
	}
	g.player = []hand{
		{
			card: playerHand,
			bet:  g.playerBet,
		},
	}
	g.state = statePlayerTurn
}

func Score(hand ...carddeck.Card) int {
	minScore := minScore(hand...)
	if minScore > 11 {
		return minScore
	}
	for _, c := range hand {
		if c.Rank == carddeck.Ace {
			return minScore + 10
		}
	}
	return minScore
}

func endHand(g *Game, playerAI AI) {
	dScore := Score(g.dealer...)
	dBlackjack := Blackjack(g.dealer...)
	allHands := make([][]carddeck.Card, len(g.player))

	for hi, hand := range g.player {
		cards := hand.card
		allHands[hi] = cards
		winnings := hand.bet
		pScore, pBlackjack := Score(cards...), Blackjack(cards...)
		switch {
		case pBlackjack && dBlackjack:
			winnings = 0
		case dBlackjack:
			winnings = -winnings
		case pBlackjack:
			winnings = int(float64(winnings) * g.blackjackPayout)
		case pScore > 21:
			//	fmt.Println("You busted.")
			winnings *= -1
		case dScore > 21:
		//	fmt.Println("Dealer busted.")
		case pScore > dScore:
		//	fmt.Println("You win!!")
		case dScore > pScore:
			//	fmt.Println("You lost :(")
			winnings *= -1
		case pScore == dScore:
			//	fmt.Println("It's a draw!")
			winnings = 0
		}
		g.balance += winnings
	}

	//fmt.Println()
	playerAI.Results(allHands, g.dealer)
	g.player = nil
	g.dealer = nil
}

func Soft(hand ...carddeck.Card) bool {
	minScore := minScore(hand...)
	score := Score(hand...)
	return minScore != score
}

func Blackjack(hand ...carddeck.Card) bool {
	return len(hand) == 2 && Score(hand...) == 21
}

func minScore(hand ...carddeck.Card) int {
	score := 0
	for _, c := range hand {
		score += min(int(c.Rank), 10)
	}

	return score
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (g *Game) currentHand() *[]carddeck.Card {
	switch g.state {
	case statePlayerTurn:
		return &g.player[g.handInd].card
	case stateDealerTurn:
		return &g.dealer
	default:
		panic("It isn't currently a player's turn")
	}
}
