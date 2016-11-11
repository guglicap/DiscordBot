package main

import (
	"bytes"
	"io/ioutil"
	"strings"

	"github.com/bwmarrin/discordgo"
)

/*COLORS
  DEFAULT: 0,
  AQUA: 1752220,
  GREEN: 3066993,
  BLUE: 3447003,
  PURPLE: 10181046,
  GOLD: 15844367,
  ORANGE: 15105570,
  RED: 15158332,
  GREY: 9807270,
  DARKER_GREY: 8359053,
  NAVY: 3426654,
  DARK_AQUA: 1146986,
  DARK_GREEN: 2067276,
  DARK_BLUE: 2123412,
  DARK_PURPLE: 7419530,
  DARK_GOLD: 12745742,
  DARK_ORANGE: 11027200,
  DARK_RED: 10038562,
  DARK_GREY: 9936031,
  LIGHT_GREY: 12370112,
  DARK_NAVY: 2899536
*/

func publicHelp() string {
	result := "Syntax:\n  !bot <command>\n  Commands:\n"
	for k, v := range PublicCommands {
		result += "\t" + k + ":\t" + v.desc + "\n"
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

func code(msg *discordgo.Message) string {
	err := sendFile(msg.ChannelID, "main.go")
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
	discord.ChannelFileSend(c, filename, bytes.NewBuffer(b))
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

func startGame(msg *discordgo.Message) string {
	if playing {
		return "You're already playing!"
	}
	game = newGame()
	playing = true
	return game.guess
}
