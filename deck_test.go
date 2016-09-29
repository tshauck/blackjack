package main

import (
	"testing"
)

func TestCardValue(t *testing.T) {
	var cardTests = []struct {
		hand  Hand
		value int
	}{
		{Hand{Cards: []Card{Card{Face: Ace}, Card{Face: Ace}}}, 12},
		{Hand{Cards: []Card{Card{Face: Ace}, Card{Face: Ace}, Card{Face: King}}}, 12},
		{Hand{Cards: []Card{Card{Face: Ace}, Card{Face: King}}}, 21},
		{Hand{Cards: []Card{Card{Face: King}, Card{Face: King}}}, 20},
	}

	for _, ct := range cardTests {
		ct.hand.CountAces()
		if ct.value != ct.hand.Value() {
			t.Errorf("Hand value not equal, hand: %v, hand total: %d, value: %d.",
				ct.hand, ct.hand.Value(), ct.value)
		}
	}
}

func TestNewDeckHas52(t *testing.T) {
	deck := NewDeck(false)

	if len(deck.Cards) != 52 {
		t.Errorf("Wrong number of cards in deck: %d", len(deck.Cards))
	}
}

func TestGameDeal(t *testing.T) {

	game := NewGame(false)
	if len(game.Deck.Cards) != 52 {
		t.Errorf("Wrong number of cards in deck: %d", len(game.Deck.Cards))
	}

	game.Deal(true)
	if len(game.Deck.Cards) != 51 {
		t.Errorf("Wrong number of cards in deck: %d", len(game.Deck.Cards))
	}
	if len(game.Player.Cards) != 1 {
		t.Errorf("Wrong number of cards in player hand: %d", len(game.Player.Cards))
	}

	game.Deal(false)
	if len(game.Deck.Cards) != 50 {
		t.Errorf("Wrong number of cards in deck: %d", len(game.Deck.Cards))
	}
	if len(game.Dealer.Cards) != 1 {
		t.Errorf("Wrong number of cards in player hand: %d", len(game.Player.Cards))
	}

}

func TestGameSetup(t *testing.T) {
	game := NewGame(false)
	game.Setup()

	if len(game.Dealer.Cards) != 2 {
		t.Errorf("Dealer should have 2 cards, they have %d", len(game.Dealer.Cards))
	}

	if game.Player.Value() < 12 {
		t.Errorf("player should have at least a hand of twelve or more, they have %d",
			game.Player.Value())
	}
}

func TestGameOutcome(t *testing.T) {
	game := NewGame(false)
	game.Setup()

	// Test setup is a Win.
	if !(game.Outcome() == Win) {
		t.Errorf("player should have won they did not outcome: %s", game.Outcome())
	}

	// Test that Game is a Win (dealer > 21)
	game = Game{
		Deck:   NewDeck(false),
		Dealer: Hand{Cards: []Card{Card{Suit: Heart, Face: King}, Card{Suit: Heart, Face: King}, Card{Suit: Heart, Face: King}}},
		Player: Hand{Cards: []Card{Card{Suit: Heart, Face: King}, Card{Suit: Heart, Face: King}}},
	}

	// Test setup is a Win.
	if !(game.Outcome() == Win) {
		t.Errorf("player should have won they did not outcome: %s", game.Outcome())
	}

	// Test that Game is a Draw.
	game = Game{
		Deck:   NewDeck(false),
		Dealer: Hand{Cards: []Card{Card{Suit: Heart, Face: King}, Card{Suit: Heart, Face: King}}},
		Player: Hand{Cards: []Card{Card{Suit: Heart, Face: King}, Card{Suit: Heart, Face: King}}},
	}

	if !(game.Outcome() == Draw) {
		t.Errorf("player should have won they did not outcome: %s", game.Outcome())
	}

	// Test that Game is a Lose.
	game = Game{
		Deck:   NewDeck(false),
		Dealer: Hand{Cards: []Card{Card{Suit: Heart, Face: King}, Card{Suit: Heart, Face: King}}},
		Player: Hand{Cards: []Card{Card{Suit: Heart, Face: Six}, Card{Suit: Heart, Face: King}}},
	}

	if !(game.Outcome() == Lose) {
		t.Errorf("player should have won they did not outcome: %s", game.Outcome())
	}

	// Test that Game is a Lose (> 21)
	game = Game{
		Deck:   NewDeck(false),
		Dealer: Hand{Cards: []Card{Card{Suit: Heart, Face: King}, Card{Suit: Heart, Face: King}}},
		Player: Hand{Cards: []Card{Card{Suit: Heart, Face: King}, Card{Suit: Heart, Face: King}, Card{Suit: Heart, Face: King}}},
	}

	if !(game.Outcome() == Lose) {
		t.Errorf("player should have won they did not outcome: %s", game.Outcome())
	}
}
