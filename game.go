package main

import (
	"fmt"
	"strings"

	"github.com/shaurya947/gophercises-deck"
)

type Hand []deck.Card

type Player struct {
	Hand
}

type Dealer struct {
	Hand
}

type GameState uint8

const (
	Deal GameState = iota
	PlayerTurn
	DealerTurn
	End
)

type Game struct {
	Deck []deck.Card
	Player
	Dealer
	GameState
}

func NewGame() *Game {
	return &Game{
		Deck:      deck.New(deck.Decks(3), deck.Shuffle),
		Player:    Player{},
		Dealer:    Dealer{},
		GameState: Deal,
	}
}

func (g *Game) Start() {
	for {
		switch g.GameState {
		case Deal:
			g.dealCards()
		case PlayerTurn:
			g.promptPlayer()
		case DealerTurn:
			g.showDealerHand()
		case End:
			g.declareWinner()
			fallthrough
		default:
			return
		}
	}
}

func (g *Game) dealCards() {
	fmt.Println("Welcome to simple Blackjack! Dealing cards...")
	g.dealToPlayer()
	g.dealToDealer()
	g.dealToPlayer()
	g.dealToDealer()
	g.GameState = PlayerTurn
}

func (g *Game) promptPlayer() {
	fmt.Println("Player's turn")
	var input string

	for {
		fmt.Println("Type h for hit, or s for stand, and press enter")
		fmt.Scanln(&input)
		input = strings.ToLower(strings.TrimSpace(input))
		if input == "h" || input == "s" {
			break
		}
		fmt.Println("Invalid input!")
	}

	if input == "h" {
		g.dealToPlayer()
		if getHighestSafeHandScore(g.Player.Hand) > 21 {
			g.GameState = End
		}
	} else {
		g.GameState = DealerTurn
	}
}

func (g *Game) showDealerHand() {
	fmt.Println("Dealer's turn")
	fmt.Printf("Dealer's second card: %s\n", g.Dealer.Hand[1].Rank.String())
	g.GameState = End
}

func (g *Game) declareWinner() {
	playerScore := getHighestSafeHandScore(g.Player.Hand)
	dealerScore := getHighestSafeHandScore(g.Dealer.Hand)

	if playerScore > 21 {
		fmt.Println("Player bust, dealer wins")
	} else if playerScore > dealerScore {
		fmt.Println("Player wins")
	} else if playerScore == dealerScore {
		fmt.Println("Draw")
	} else {
		fmt.Println("Dealer wins")
	}
}

func (g *Game) dealToPlayer() {
	g.Player.Hand = append(g.Player.Hand, g.dealOne())
	i := len(g.Player.Hand)
	fmt.Printf("Player card %d: %s\n", i, g.Player.Hand[i-1].Rank.String())
}

func (g *Game) dealToDealer() {
	g.Dealer.Hand = append(g.Dealer.Hand, g.dealOne())
	i := len(g.Dealer.Hand)
	rank := g.Dealer.Hand[i-1].Rank.String()
	if i == 2 {
		rank = "hidden"
	}
	fmt.Printf("Dealer card %d: %s\n", i, rank)
}

func (g *Game) dealOne() deck.Card {
	dealt := g.Deck[0]
	g.Deck = g.Deck[1:]
	return dealt
}

func getHighestSafeHandScore(h Hand) int {
	possibleScores := [2]int{}

	for _, card := range h {
		switch card.Rank {
		case deck.Jack:
			fallthrough
		case deck.Queen:
			fallthrough
		case deck.King:
			fallthrough
		case deck.Ten:
			possibleScores[0] += 10
			if possibleScores[1] != 0 {
				possibleScores[1] += 10
			}
		case deck.Ace:
			possibleScores[0] += 1
			if possibleScores[1] == 0 {
				possibleScores[1] = possibleScores[0] + 11
			} else {
				possibleScores[1] += 1
			}
		default:
			possibleScores[0] += int(card.Rank)
			if possibleScores[1] != 0 {
				possibleScores[1] += int(card.Rank)
			}
		}
	}

	if possibleScores[1] == 0 {
		return possibleScores[0]
	}

	if possibleScores[1] <= 21 {
		return possibleScores[1]
	}

	return possibleScores[0]
}
