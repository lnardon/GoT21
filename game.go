package main

import (
	"fmt"
	"math/rand"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type model struct {
	playerHand []int
	dealerHand []int
	playerTurn bool
	gameOver   bool
	message    string
	width      int
	height     int
}

func InitialModel() model {
	w, h, _ := term.GetSize(int(os.Stdout.Fd()))
	return model{
		playerHand: []int{drawCard(), drawCard()},
		dealerHand: []int{drawCard()},
		playerTurn: true,
		gameOver:   false,
		message:    "Your turn! Press 'h' to hit, 's' to stand.",
		width:      w,
		height:     h,
	}
}

func (m model) Init() tea.Cmd {
	rand.Int()
	return nil
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

type message string

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	case tea.KeyMsg:
		switch msg.String() {
		case "h":
			if !m.playerTurn {
				break
			}

			m.playerHand = append(m.playerHand, drawCard())
			if handValue(m.playerHand) > 21 {
				m.message = "You busted! Dealer wins. Press 'a' to play again or 'q' to quit."
				m.playerTurn = false
			}
		case "s":
			if !m.playerTurn {
				break
			}

			m.playerTurn = false
			for handValue(m.dealerHand) < 17 {
				m.dealerHand = append(m.dealerHand, drawCard())
			}
			winner := determineWinner(m)
			m.message = fmt.Sprintf("%s\n%s", winner, "Press 'a' to play again or 'q' to quit.")
		case "q":
			fmt.Print("\033[H\033[2J")
			return m, tea.Quit
		case "a":
			m = InitialModel()
		}
	}
	return m, nil
}

func determineWinner(m model) string {
	playerScore := handValue(m.playerHand)
	dealerScore := handValue(m.dealerHand)

	if playerScore > 21 {
		return "You busted! Dealer wins."
	}

	if dealerScore > 21 {
		return "Dealer busted! You win."
	}

	if playerScore > dealerScore {
		return "You win!"
	}

	if dealerScore > playerScore {
		return "Dealer wins!"
	}
	return "It's a tie!"
}

func renderCards(hand []int) string {
	cardStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		Width(6).
		Padding(1).
		Align(lipgloss.Center)

	cards := []string{}
	for _, card := range hand {
		cards = append(cards, cardStyle.Render(fmt.Sprintf("%d", card)))
	}

	w, _, _ := term.GetSize(int(os.Stdout.Fd()))
	t := lipgloss.JoinHorizontal(lipgloss.Top, cards...)
	container := lipgloss.NewStyle().
		Width(w).
		Align(lipgloss.Center).
		Render(fmt.Sprintf("%s", t))
	return container
}

func (m model) View() string {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#FAFAFA")).
		Background(lipgloss.Color("#7D56F4")).
		Padding(1, 2).
		MarginBottom(1).
		Width(m.width)

	return lipgloss.JoinVertical(lipgloss.Left,
		headerStyle.Render("Got21"),
		"Dealer's Hand:",
		renderCards(m.dealerHand),
		"\n\n\n\n",
		"Your Hand:",
		renderCards(m.playerHand),
		"",
		m.message,
	)
}
