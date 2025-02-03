package main

import (
	"fmt"
	"math/rand"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	playerHand []int
	dealerHand []int
	playerTurn bool
	gameOver   bool
	message    string
}

func initialModel() model {
	return model{
		playerHand: []int{drawCard(), drawCard()},
		dealerHand: []int{drawCard()},
		playerTurn: true,
		gameOver:   false,
		message:    "Your turn! Press 'h' to hit, 's' to stand.",
	}
}

func drawCard() int {
	return rand.Intn(10) + 2
}

func handValue(hand []int) int {
	sum, aces := 0, 0
	for _, card := range hand {
		sum += card
		if card == 11 {
			aces++
		}
	}
	for sum > 21 && aces > 0 {
		sum -= 10
		aces--
	}
	return sum
}

func checkWinner(m model) string {
	playerScore := handValue(m.playerHand)
	dealerScore := handValue(m.dealerHand)

	if playerScore > 21 {
		return "You busted! Dealer wins."
	} else if dealerScore > 21 {
		return "Dealer busted! You win."
	} else if !m.playerTurn {
		if playerScore > dealerScore {
			return "You win!"
		} else if dealerScore > playerScore {
			return "Dealer wins!"
		} else {
			return "It's a tie!"
		}
	}
	return "Would you like to play [a]gain or [quit]?"
}

type message string

func (m model) Init() tea.Cmd {
	rand.Int()
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.gameOver {
			return m, tea.Quit
		}
		switch msg.String() {
		case "h":
			if !m.playerTurn {
				break
			}

			m.playerHand = append(m.playerHand, drawCard())
			if handValue(m.playerHand) > 21 {
				m.message = checkWinner(m)
				m.message = "Would you like to play [a]gain or [q]uit?"
			}
		case "s":
			m.playerTurn = false
			for handValue(m.dealerHand) < 20 && handValue(m.playerHand) > handValue(m.dealerHand) {
				m.dealerHand = append(m.dealerHand, drawCard())
			}
			m.message = checkWinner(m)
		case "q":
			return m, tea.Quit
		case "a":
			m = initialModel()
		}
	}
	return m, nil
}

func renderCards(hand []int) string {
	cards := []string{}
	for _, card := range hand {
		cards = append(cards, fmt.Sprintf("[ %d ]", card))
	}
	return strings.Join(cards, " ")
}

func (m model) View() string {
	return fmt.Sprintf(
		"\n ♣ Dealer's Hand: %s (Value: %d)\n\nYour Hand: %s (Value: %d)\n\n%s\n\n\n\n\nPress 'h' to hit, 's' to stand, 'q' to quit.\n",
		renderCards(m.dealerHand), handValue(m.dealerHand),
		renderCards(m.playerHand), handValue(m.playerHand),
		m.message,
	)
}

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running game:", err)
	}
}
