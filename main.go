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
	secrets    []string           //Bot Token and DB Auth
	session    *discordgo.Session //Discord session.
	db         *sql.DB            //DB interface
)

func main() {
	initialize()
	defer db.Close()
	discord, err := discordgo.New("Bot " + secrets[0])
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
	secrBytes, err := ioutil.ReadFile("secrets.txt")
	check(err)
	secrets = strings.Split(string(secrBytes), "\n")
	db, err = sql.Open("sqlite3", "./db")
	check(err)
	loadWords()
}

func messageCreated(s *discordgo.Session, msg *discordgo.MessageCreate) {

	//Check if the message has a pastebin link.
	link, id := hasPastebinLink(msg.Content)
	if len(link) != 0 && len(id) != 0 {
		lang := getPasteLanguage(link)
		if lang == "text" {
			lang = ""
		}
		raw := getPasteRaw(id)
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("```%s\n%s\n```", lang, raw))
		return
	}

	//If not, let's see if it is a bot command.
	if strings.HasPrefix(msg.Content, magicWord) {
		if x, _ := s.Channel(msg.ChannelID); x.IsPrivate {
			privateMessage(s, msg)
		} else {
			publicMessage(s, msg)
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

func privateMessage(s *discordgo.Session, msg *discordgo.MessageCreate) {
	c := msg.ChannelID
	tokens := getTokens(c, msg.Content, privateHelp)
	if tokens == nil {
		return
	}
	if command, ok := PrivateCommands[tokens[1]]; ok {
		reply := command.reply(msg.Message)
		if len(reply) > 0 {
			s.ChannelMessageSend(c, reply)
		}
	}
}

func publicMessage(s *discordgo.Session, msg *discordgo.MessageCreate) {
	c := msg.ChannelID
	tokens := getTokens(c, msg.Content, publicHelp)
	if tokens == nil {
		return
	}
	if command, ok := PublicCommands[tokens[1]]; ok {
		reply := command.reply(msg.Message)
		if len(reply) > 0 {
			s.ChannelMessageSend(c, reply)
		}
	}
}
