package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func publicHelp() string {
	result := "Syntax:\n  !bot <command>\n  Commands:\n"
	for k, v := range PublicCommands {
		if len(v.desc) > 0 {
			result += "\t" + k + ":\t" + v.desc + "\n"
		}
	}
	return result
}

func me(msg *discordgo.Message) string {
	return "You're " + msg.Author.Username
}

func bot(msg *discordgo.Message) string {
	return "I am the bot. Bleep Bloop."
}

func kappa(msg *discordgo.Message) string {
	err := sendFile(msg.ChannelID, "kappa.png")
	if err != nil {
		return err.Error()
	}
	return ""
}

func sendFile(c, filename string) error {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	session.ChannelFileSend(c, filename, bytes.NewBuffer(b))
	return nil
}

func whois(msg *discordgo.Message) string {
	tokens := strings.Fields(msg.Content)
	if len(tokens) < 3 {
		return "Who?"
	}
	if id, ok := NametoID[strings.ToLower(tokens[2])]; ok {
		if desc, ok := IDtoWhois[id]; ok {
			return desc
		}
	}
	return "I don't know him. Sorry."
}

func githubPublic(msg *discordgo.Message) string {
	tokens := strings.Fields(msg.Content)
	if len(tokens) < 3 {
		return "Who?"
	}
	var link string
	if id, ok := NametoID[strings.ToLower(tokens[2])]; ok {
		err := db.QueryRow("SELECT Github FROM Users WHERE AccId=?", id).Scan(&link)
		if err != nil {
			return ""
		}
	}
	if len(link) < 1 {
		return "Hmm, I don't know."
	}
	return "There you go\n" + link
}

func source(msg *discordgo.Message) string {
	return "I'm here!\n" + "https://github.com/guglicap/DiscordBot"
}

func startGame() string {
	if playing {
		return "You're already playing!"
	}
	game = newGame()
	playing = true
	return game.guess
}

func getGamePoints(tokens []string) string {
	var points int
	if id, ok := NametoID[strings.ToLower(tokens[3])]; ok {
		err := db.QueryRow("SELECT GamePts FROM Users WHERE AccId=?", id).Scan(&points)
		if err != nil {
			return ""
		}
	}
	return fmt.Sprintf("%s has %d points", tokens[3], points)
}

func gameCmd(msg *discordgo.Message) string {
	tokens := strings.Fields(msg.Content)
	if len(tokens) == 2 {
		return startGame()
	}
	if len(tokens) == 4 {
		return getGamePoints(tokens)
	}
	return ""
}
