package main

import (
	"database/sql"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Game struct {
	word     string
	guess    string
	attempts int
}

var (
	words  []string
	points map[string]int64
)

//This includes loading the users and the words.
func initGame() {
	wordsBytes, err := ioutil.ReadFile("words.txt")
	check(err)
	words = strings.Split(string(wordsBytes), "\n")
	points = make(map[string]int64)
	checkDbConnection()
	rows, err := db.Query("SELECT AccId, GamePts FROM Users")
	check(err)
	defer rows.Close()
	for rows.Next() {
		var id string
		var pt sql.NullInt64
		err = rows.Scan(&id, &pt)
		check(err)
		if pt.Valid {
			points[id] = pt.Int64
		} else {
			points[id] = 0
		}
	}
}

//This replaces the dashes with eventual matches.
func (self *Game) updateGuess(substr string) bool {

	//If this is a whole-word guess
	if len(substr) == len(self.word) {
		//Then check for it to match
		if substr == self.word {
			self.guess = self.word
			return true
		}
	}
	//Else it should be only a character, but let's check just in case.
	if len(substr) != 1 {
		//If the lenght of the substr isn't 1 (so it's not a char), then something went wrong. Pretend nothing has ever happened.
		return false
	}

	//If the word doesn't contain that char, return.
	if !strings.Contains(self.word, substr) {
		return false
	}

	//Loop over every letter of word and if it matches with that char, replace the corresponding dash in the guess.
	//I should review my variable naming. PF please help.

	for i := 1; i < len(self.word)-1; i++ {
		if string(self.word[i]) == substr {
			//That mess is because strings are immutable objects in Go.
			self.guess = self.guess[:i] + substr + self.guess[i+1:]
		}
	}
	return true
}

func newGame() *Game {
	rand.Seed(time.Now().UnixNano())
	//Randomly choose a word.
	word := strings.ToLower(words[rand.Intn(len(words))])
	//Replace the middle of the word with dashes
	guess := string(word[0]) + strings.Repeat("-", len(word)-2) + string(word[len(word)-1])
	playing = true

	//Easy win for Goo
	//fmt.Println(word)

	return &Game{word, guess, 0}
}

func flushPeopleDb() {
	checkDbConnection()
	for id, pts := range points {
		db.Exec("INSERT OR REPLACE INTO Users (AccId, Github, GamePts) VALUES ($1, (SELECT Github FROM Users WHERE AccId = $1), $2)", id, pts)
	}
}

//This is called once a message with gamePrefix is received
func gameMessage(msg *discordgo.Message) string {
	if !playing {
		game = newGame()
		session.ChannelMessageSend(msg.ChannelID, "You weren't playing so I started a game for you.")
	}
	if len(msg.Content) < 2 {
		//Nothing given, so message is "%"
		return "I need a letter or a guess, buddy. Not a percent sign."
	}
	if len(msg.Content) > 2 && len(msg.Content)-1 != len(game.word) {
		//The message is not a guess nor a single letter.
		return "You can give me only one letter or your guess."
	}
	if game.updateGuess(msg.Content[1:]) {
		if v, ok := points[msg.Author.ID]; ok {
			points[msg.Author.ID] = v + 1
		} else {
			points[msg.Author.ID] = 1
		}
	}
	if game.guess == game.word {
		playing = false
		flushPeopleDb()
		return "The word was " + game.word + "\nYou won!"
	}
	return game.guess
}
