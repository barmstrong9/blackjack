package main

import (
	"fmt"
	"strings"

	"github.com/barmstrong9/brandon/gophercises/deck"
)

//Hand allows the user and dealer to have cards
type Hand []deck.Card

func (h Hand) String() string {
	strs := make([]string, len(h))
	for i := range h {
		strs[i] = h[i].String()
	}
	return strings.Join(strs, ", ")
}

//DealerString hides the dealers second card
func (h Hand) DealerString() string {
	return h[0].String() + ", **HIDDEN**"
}

//Score lets us have aces be worth 11
func (h Hand) Score() int {
	minScore := h.MinScore()
	if minScore > 11 {
		return minScore
	}
	for _, c := range h {
		if c.Rank == deck.Ace {
			//Currently, ace = 1, we change it to 11, only if the score is 11 or below
			return minScore + 10
		}
	}
	return minScore
}

//MinScore allows us to figure out the minimum score
func (h Hand) MinScore() int {
	score := 0
	for _, c := range h {
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

//Shuffle shuffles the cards before the game
func Shuffle(gs GameState) GameState {
	ret := clone(gs)
	ret.Deck = deck.NewDeck(deck.MultiDeck(3), deck.Shuffle)
	return ret
}

//Deal deals the cards, going player, dealer, player, dealer.
func Deal(gs GameState) GameState {
	ret := clone(gs)
	ret.Player = make(Hand, 0, 5)
	ret.Dealer = make(Hand, 0, 5)
	var card deck.Card
	for i := 0; i < 2; i++ {
		card, ret.Deck = draw(ret.Deck)
		ret.Player = append(ret.Player, card)
		card, ret.Deck = draw(ret.Deck)
		ret.Dealer = append(ret.Dealer, card)
	}
	ret.State = StatePlayerTurn
	return ret
}

//Stand ends the turn of the player or dealer
func Stand(gs GameState) GameState {
	ret := clone(gs)
	ret.State++
	return ret
}

//Hit gives the user another card.
func Hit(gs GameState) GameState {
	ret := clone(gs)
	hand := ret.CurrentPlayer()
	var card deck.Card
	card, ret.Deck = draw(ret.Deck)
	*hand = append(*hand, card)
	if hand.Score() > 21 {
		return Stand(ret)
	}
	return ret
}

//EndHand compares the scores of the user and dealer to see if either bust and who won
func EndHand(gs GameState) GameState {
	ret := clone(gs)
	pScore, dScore := ret.Player.Score(), ret.Dealer.Score()
	fmt.Println()
	fmt.Println("==FINAL HANDS==")
	fmt.Println("Your Cards:", ret.Player, "\nYour Score:", pScore)
	fmt.Println("Dealer's Cards:", ret.Dealer, "\nDealer's Score:", dScore)
	switch {
	case pScore > 21:
		fmt.Println("You Busted, You Lost!")
	case dScore > 21:
		fmt.Println("Dealer Busted, You Win!")
	case pScore > dScore:
		fmt.Println("You Win!")
	case dScore > pScore:
		fmt.Println("You Lost! Try Again")
	case dScore == pScore:
		fmt.Println("Draw")
	}
	fmt.Println()
	ret.Player = nil
	ret.Dealer = nil
	return ret
}

func main() {
	var gs GameState
	gs = Shuffle(gs)
	fmt.Println("\nWelcome To Black Jack")
	gs = Deal(gs)
	if gs.Dealer.Score() == 21 || gs.Player.Score() == 21 {
		fmt.Println()
		fmt.Println("Your Cards:", gs.Player)
		fmt.Println("Dealer's Cards:", gs.Dealer)
		fmt.Println("***BLACKJACK***")
		gs = Deal(gs)

	}
	var input string
	for gs.State == StatePlayerTurn {
		fmt.Println("\nYour Cards:", gs.Player)
		fmt.Println("Dealer's Cards:", gs.Dealer.DealerString())

		fmt.Println("What will you do? (H)it, (S)tand")
		fmt.Scanf("%s\n", &input)
		lowerInput := strings.ToLower(input)
		switch lowerInput {
		case "h":
			gs = Hit(gs)
		case "s":
			gs = Stand(gs)
		default:
			fmt.Println("Invalid Option:", input)
		}
	}
	for gs.State == StateDealerTurn {
		if gs.Dealer.Score() <= 16 || (gs.Dealer.Score() == 17 && gs.Dealer.MinScore() != 17) {
			gs = Hit(gs)
		} else {
			gs = Stand(gs)
		}
	}
	gs = EndHand(gs)
}

func draw(cards []deck.Card) (deck.Card, []deck.Card) {
	return cards[0], cards[1:]
}

//State allows for transition between player and dealer
type State int8

const (
	//StatePlayerTurn shows that it is currently the users turn. Index = 0
	StatePlayerTurn State = iota
	//StateDealerTurn shows that it is currently the dealers turn. Index = 1
	StateDealerTurn
	//StateHandOver shows that it is no longer the user or dealers turn, the current game should end. Index = 2
	StateHandOver
)

//GameState allows information to be cloned from one GameState to another
type GameState struct {
	Deck   []deck.Card
	State  State
	Player Hand
	Dealer Hand
}

//CurrentPlayer is a fallback to see if there are any errors in the code
func (gs *GameState) CurrentPlayer() *Hand {
	switch gs.State {
	case StatePlayerTurn:
		return &gs.Player
	case StateDealerTurn:
		return &gs.Dealer
	default:
		panic("it currently isn't any players turn")
	}
}

func clone(gs GameState) GameState {
	ret := GameState{
		Deck:   make([]deck.Card, len(gs.Deck)),
		State:  gs.State,
		Player: make(Hand, len(gs.Player)),
		Dealer: make(Hand, len(gs.Dealer)),
	}
	copy(ret.Deck, gs.Deck)
	copy(ret.Player, gs.Player)
	copy(ret.Dealer, gs.Dealer)
	return ret
}
