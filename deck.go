package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"math/rand"
)

const CardsInDeck = 52

const (
	Heart = iota
	Spade
	Club
	Diamond
)

var Suits = []int{Heart, Spade, Club, Diamond}

const (
	Ace = iota
	Two
	Three
	Four
	Five
	Six
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
)

var Faces = []int{Ace, Two, Three, Four, Five, Six, Seven, Eight, Nine, Ten, Jack, Queen, King}
var FaceNames = []string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Jack", "Queen", "King"}
var FaceValue = []int{11, 2, 3, 4, 5, 6, 7, 8, 9, 10, 10, 10, 10}

type Card struct {
	Suit int
	Face int
}

type Deck struct {
	Cards []Card
}

func NewDeck(shuffle bool) Deck {
	cards := make([]Card, CardsInDeck)
	shuff := rand.Perm(CardsInDeck)
	i := 0

	for _, suit := range Suits {
		for _, face := range Faces {
			card := Card{Suit: suit, Face: face}

			if shuffle {
				cards[shuff[i]] = card
			} else {
				cards[i] = card
			}
			i = i + 1
		}
	}
	return Deck{Cards: cards}
}

type Hand struct {
	Cards []Card
	Aces  int // number of aces
}

func (h Hand) Value() int {
	total := 0
	for _, card := range h.Cards {
		total += FaceValue[card.Face]
	}

	aces := h.Aces
	for {
		if aces == 0 || total <= 21 {
			break
		}

		aces = aces - 1
		total = total - 10
	}
	return total
}

func (h *Hand) CountAces() {
	aces := 0
	for _, card := range h.Cards {
		if card.Face == Ace {
			aces++
		}
	}
	h.Aces = aces
}

func (h *Hand) AddCard(c Card) {
	if c.Face == Ace {
		h.Aces++
	}
	h.Cards = append(h.Cards, c)
}

type Game struct {
	Deck   Deck
	Dealer Hand
	Player Hand
}

func NewGame(shuffle bool) Game {
	deck := NewDeck(shuffle)
	dealer := Hand{}
	player := Hand{}

	return Game{Deck: deck, Dealer: dealer, Player: player}
}

func (g *Game) Deal(player bool) {
	top := g.Deck.Cards[0]
	if player {
		g.Player.AddCard(top)
	} else {
		g.Dealer.AddCard(top)
	}
	g.Deck.Cards = g.Deck.Cards[1:]
}

func (g *Game) Setup() {
	g.Deal(true)
	g.Deal(false)
	g.Deal(true)
	g.Deal(false)

	for {
		if g.Player.Value() >= 12 {
			break
		}
		g.Deal(true)
	}
}

func (g Game) String() string {
	return fmt.Sprintf("Game:\n\tDeck: %d cards\n\tPlayer %d cards (%d value)\n\tDealer %d cards (%d value)",
		len(g.Deck.Cards), len(g.Player.Cards), g.Player.Value(), len(g.Dealer.Cards), g.Dealer.Value())

}

type Reward int

const (
	Lose Reward = iota - 1
	Draw
	Win
)

func (g Game) Outcome() Reward {
	dealerValue := g.Dealer.Value()
	playerValue := g.Player.Value()

	if playerValue > 21 {
		return Lose
	} else if dealerValue > 21 {
		return Win
	} else if playerValue > dealerValue {
		return Win
	} else if dealerValue > playerValue {
		return Lose
	} else {
		return Draw
	}
}

func (g Game) State() State {
	playerTotal := g.Player.Value()
	dealerFace := g.Dealer.Cards[0].Face

	return State{PlayerTotal: playerTotal,
		DealerFace: dealerFace,
		Aces:       g.Player.Aces}
}

type Action int

const (
	Stay Action = iota
	Hit
)

var Actions = []Action{Stay, Hit}
var ActionNames = []string{"Stay", "Hit"}

type State struct {
	PlayerTotal int
	Aces        int
	DealerFace  int
}

type Agent struct {
	Q     map[State]map[Action]int
	Visit map[State]map[Action]int
	L     *log.Logger
}

func (a Agent) SavePolicy(fileName string) error {
	var data []map[string]interface{}
	for state, actions := range a.Q {
		for action, q := range actions {
			datum := map[string]interface{}{
				"playerTotal":    state.PlayerTotal,
				"aces":           state.Aces,
				"dealerFaceName": FaceNames[state.DealerFace],
				"dealerFace":     state.DealerFace,
				"action":         ActionNames[action],
				"q":              q,
				"visits":         a.Visit[state][action],
			}
			data = append(data, datum)
		}
	}
	dataJson, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(fileName, dataJson, 0644)
}

func (a *Agent) Policy(s State) Action {
	hitReward := float32(a.Q[s][Hit]) / float32(a.Visit[s][Hit]+1) // Use smoothing
	stayReward := float32(a.Q[s][Stay]) / float32(a.Visit[s][Stay]+1)

	if hitReward > stayReward {
		return Hit
	} else {
		return Stay
	}
}

func NewAgent(l *log.Logger) Agent {

	var q map[State]map[Action]int
	q = make(map[State]map[Action]int)

	var v map[State]map[Action]int
	v = make(map[State]map[Action]int)

	for _, nAces := range []int{0, 1, 2, 3, 4} {
		for playerTotal := 12; playerTotal < 22; playerTotal++ {
			for _, face := range Faces {

				var qa = make(map[Action]int)
				var va = make(map[Action]int)

				for _, action := range Actions {
					qa[action] = 0
					va[action] = 0
				}

				a := State{PlayerTotal: playerTotal, Aces: nAces, DealerFace: face}
				q[a] = qa
				v[a] = va
			}
		}
	}

	return Agent{Q: q, Visit: v, L: l}
}

type GameEvent struct {
	State  State // Might change this to a Game
	Action Action
}

func (a *Agent) PlayGames(nGames int) {
	for nGame := 0; nGame < nGames; nGame++ {
		a.L.WithFields(log.Fields{"nGame": nGame}).Info("Starting Game.")

		game := NewGame(true)
		game.Setup()
		var gameEvents []GameEvent

		// Play the player's hand according to the policy.
		for {
			state := game.State()
			action := a.Policy(state)
			gameEvents = append(gameEvents, GameEvent{State: state, Action: action})
			if action == Hit {
				game.Deal(true)
			} else {
				break
			}
		}

		// Play the dealer's hand.
		for {
			if game.Dealer.Value() >= 17 {
				break
			}
			game.Deal(false)
		}

		outcome := game.Outcome()
		a.L.WithFields(log.Fields{"aces": game.Player.Aces, "nGame": nGame, "reward": outcome, "playerTotal": game.Player.Value(), "dealerValue": game.Dealer.Value()}).Info("Finished Game.")
		a.UpdatePolicy(gameEvents, game.Outcome())
	}
}

func (a *Agent) UpdatePolicy(gameEvents []GameEvent, reward Reward) {
	for _, gameEvent := range gameEvents {
		if gameEvent.State.PlayerTotal > 21 {
			a.L.Debugf("player went bust with a score of %d", gameEvent.State.PlayerTotal)
			continue
		}
		a.Q[gameEvent.State][gameEvent.Action] += int(reward)
		a.Visit[gameEvent.State][gameEvent.Action]++

		a.L.WithFields(log.Fields{
			"playerTotal": gameEvent.State.PlayerTotal,
			"dealerFace":  gameEvent.State.DealerFace,
			"aces":        gameEvent.State.Aces,
			"action":      gameEvent.Action,
			"visits":      a.Visit[gameEvent.State][gameEvent.Action],
			"q":           a.Q[gameEvent.State][gameEvent.Action],
			"reward":      reward,
		}).Info("updated action state values")
	}
}
