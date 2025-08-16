package game

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"executive-chef/internal/deck"
	"executive-chef/internal/player"
)

// Turn represents a single turn in the game.
type Turn struct {
	Deck   *deck.Deck
	Player *player.Player
}

// DraftPhase performs the drafting phase of a turn. Ten cards are revealed
// and the player may draft three of them.
func (t *Turn) DraftPhase() {
	reveal := t.Deck.Draw(10)
	reader := bufio.NewReader(os.Stdin)
	for picks := 0; picks < 3 && len(reveal) > 0; picks++ {
		fmt.Println("Choose an ingredient to draft:")
		for i, card := range reveal {
			fmt.Printf("%d) %s (%s)\n", i+1, card.Name, card.Role)
		}
		fmt.Print("Selection: ")
		input, _ := reader.ReadString('\n')
		n, err := strconv.Atoi(strings.TrimSpace(input))
		if err != nil || n < 1 || n > len(reveal) {
			fmt.Println("Invalid selection")
			picks--
			continue
		}
		chosen := reveal[n-1]
		t.Player.Add(chosen)
		reveal = append(reveal[:n-1], reveal[n:]...)
	}
}
