package main

//GUILD ID 227330229520171008
//ROLE ID 238008384203390976
import (
	"log"
	"strings"

	"database/sql"

	"github.com/bwmarrin/discordgo"
	_ "github.com/go-sql-driver/mysql"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

var (
	magicWord string
	session   *discordgo.Session
	db        *sql.DB
)

func main() {
	magicWord = "!bot"
	discord, err := discordgo.New("Bot " + "MjI5MzAyMTQ5Nzk5MDg0MDMy.CshR2w.PWImFjIlo028OGVMn0kFaVMTMFk")
	check(err)
	discord.AddHandler(messageCreated)
	session = discord
	db, err = sql.Open("mysql", "goo:Gucefalo1@tcp(127.0.0.1:3306)/discordbot")
	check(err)
	defer db.Close()
	err = discord.Open()
	check(err)
	<-make(chan struct{})
}

func messageCreated(s *discordgo.Session, msg *discordgo.MessageCreate) {
	if strings.HasPrefix(msg.Content, magicWord) {
		if x, _ := s.Channel(msg.ChannelID); x.IsPrivate {
			privateMessage(s, msg)
		} else {
			publicMessage(s, msg)
		}
	}
}

func privateMessage(s *discordgo.Session, msg *discordgo.MessageCreate) {
	c := msg.ChannelID
	tokens := strings.Fields(msg.Content)
	if len(tokens) < 2 {
		s.ChannelMessageSend(c, privateHelp())
		return
	}
	if tokens[1] == "help" {
		s.ChannelMessageSend(c, privateHelp())
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
	tokens := strings.Fields(msg.Content)
	if len(tokens) < 2 {
		s.ChannelMessageSend(c, "Yes?")
		return
	}
	if tokens[1] == "help" {
		s.ChannelMessageSend(c, publicHelp())
		return
	}
	if command, ok := PublicCommands[tokens[1]]; ok {
		reply := command.reply(msg.Message)
		if len(reply) > 0 {
			s.ChannelMessageSend(c, reply)
		}
	}
}
