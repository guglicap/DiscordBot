package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"math/rand"
	"regexp"
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

func loadWords() {
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

func (self *Game) updateGuess(substr string) bool {
	if len(substr) == len(self.word) {
		substr = substr[1 : len(substr)-1]
	}
	re := regexp.MustCompile(substr)
	indexes := re.FindAllStringIndex(self.word[1:len(self.word)-1], -1)
	if indexes == nil {
		//Nope, wrong guess.
		return false
	}
	for _, i := range indexes {
		x := i[0] + 1
		y := i[1] + 1
		self.guess = self.guess[:x] + self.word[x:y] + self.guess[y:]
	}
	return true
}

func newGame() *Game {
	rand.Seed(time.Now().UnixNano())
	word := strings.ToLower(words[rand.Intn(len(words))])
	guess := string(word[0]) + strings.Repeat("-", len(word)-2) + string(word[len(word)-1])
	playing = true
	return &Game{word, guess, 0}
}

func flushPeopleDb() {
	checkDbConnection()
	for id, pts := range points {
		db.Exec("INSERT INTO Users (AccId, GamePts) VALUES (?, ?) ON DUPLICATE KEY UPDATE GamePts=VALUES(GamePts)", id, pts)
	}
}

//This is just another function of type replyString
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
		fmt.Println(len(msg.Content))
		fmt.Println(game.word)
		return "You can give me only one letter or your guess."
	}
	if game.updateGuess(regexp.QuoteMeta(msg.Content[1:])) {
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
