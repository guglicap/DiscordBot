package main

import (
	"io/ioutil"
	"log"
	"strings"

	"database/sql"

	"fmt"

	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	gamePrefix string             //Game guesses prefix
	magicWord  string             //Keyword used to summon the bot
	game       *Game              //The current Game
	playing    bool               //Are we playing?
	token      string             //Bot Token
	session    *discordgo.Session //Discord session.
	db         *sql.DB            //DB interface
)

func main() {
	initialize()
	defer db.Close()
	discord, err := discordgo.New("Bot " + token) //This is assuming a Bot Token.
	session = discord
	discord.AddHandler(messageCreated)
	err = discord.Open()
	check(err)
	<-make(chan struct{}) //Prevent the program from exiting
}

//Initialize global vars
func initialize() {
	magicWord = "!bot "
	gamePrefix = "%"
	playing = false

	//Get the secrets from the file
	secrBytes, err := ioutil.ReadFile("secrets.txt")
	check(err)
	token = strings.Split(string(secrBytes), "\n")[0]

	db, err = sql.Open("sqlite3", "./db")
	check(err)

	initGame()
}

func messageCreated(s *discordgo.Session, msg *discordgo.MessageCreate) {

	//Check if the message has a pastebin link.
	link, id := hasPastebinLink(msg.Content)
	if len(link) != 0 && len(id) != 0 {
		lang := getPasteLanguage(link)
		if lang == "text" {
			lang = ""
		}
		raws := getPasteRaw(id)
		for _, raw := range raws {
			s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("```%s\n%s\n```", lang, raw))
		}
		return
	}

	//If not, let's see if it is a bot command.
	if strings.HasPrefix(msg.Content, magicWord) {
		if x, _ := s.Channel(msg.ChannelID); x.IsPrivate {
			sendMessage(s, msg, PrivateCommands, privateHelp)
		} else {
			sendMessage(s, msg, PublicCommands, publicHelp)
		}
		return
	}

	//If not, let's see if it is a game guess
	if strings.HasPrefix(msg.Content, gamePrefix) {
		reply := gameMessage(msg.Message)
		if len(reply) > 0 {
			s.ChannelMessageSend(msg.ChannelID, reply)
		}
	}

	//If not, bye.
}

//Takes the list of commands and the help function as parameters
func sendMessage(s *discordgo.Session, msg *discordgo.MessageCreate, commands map[string]Command, help func() string) {
	c := msg.ChannelID
	tokens := getTokens(c, msg.Content, help)
	if tokens == nil {
		return
	}
	if command, ok := commands[tokens[1]]; ok {
		reply := command.reply(msg.Message)
		if len(reply) > 0 {
			s.ChannelMessageSend(c, reply)
		}
	}
}
