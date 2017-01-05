package main

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

func privateHelp() string {
	result := "Syntax:\n  !bot <command>\n  Commands:\n"
	for k, v := range PrivateCommands {
		result += "\t" + k + ":\t" + v.desc + "\n"
	}
	return result
}

func githubPrivate(msg *discordgo.Message) string {
	accID := msg.Author.ID
	if len(accID) < 1 {
		return ""
	}
	//checkDbConnection()
	tokens := strings.Fields(msg.Content)
	if len(tokens) < 3 {
		return "Not enough Tokens"
	}
	if len(tokens[2]) < 1 {
		return "Link null"
	}
	_, err := db.Exec("INSERT OR REPLACE INTO Users (AccId, Github, GamePts) VALUES ($1, $2, (SELECT GamePts FROM Users WHERE AccId = $1))", accID, tokens[2])
	check(err)
	return "Yes, sir."
}
