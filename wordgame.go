package main

import (
	"io/ioutil"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

type Game struct {
	word     string
	guess    string
	attempts int
}

var (
	words []string
)

func loadWords() {
	wordsBytes, err := ioutil.ReadFile("words.txt")
	check(err)
	words = strings.Split(string(wordsBytes), "\n")
}

func (self *Game) updateGuess(substr string) {
	re := regexp.MustCompile(substr)
	indexes := re.FindAllStringIndex(self.word, -1)
	if indexes == nil {
		//Nope, wrong guess.
		return
	}
	for _, i := range indexes {
		self.guess = self.guess[:i[0]] + substr + self.guess[i[1]:]
	}
}

func newGame() *Game {
	rand.Seed(time.Now().UnixNano())
	word := strings.ToLower(words[rand.Intn(len(words))])
	guess := string(word[0]) + strings.Repeat("-", len(word)-2) + string(word[len(word)-1])
	playing = true
	return &Game{word, guess, 0}
}

//This should run only if playing == true, and thus game != nil.
//Returns the updated guess.
func handleGameMessage(msg string) string {
	if len(msg) < 2 {
		//Nothing given, so message is "%"
		return "I need a letter or a guess, buddy. Not a percent sign."
	}
	game.updateGuess(msg[1:])
	if game.guess == game.word {
		playing = false
		return "You won!"
	}
	return game.guess
}
