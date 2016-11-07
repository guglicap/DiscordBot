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

func checkDbConnection() {
	err := db.Ping()
	check(err)
}

func githubPrivate(msg *discordgo.Message) string {
	accID := msg.Author.ID
	if len(accID) < 1 {
		return ""
	}
	checkDbConnection()
	tokens := strings.Fields(msg.Content)
	if len(tokens) < 3 {
		return "Not enough Tokens"
	}
	if len(tokens[2]) < 1 {
		return "Link null"
	}
	_, err := db.Exec("INSERT INTO Users (AccId, Github) VALUES (?, ?) ON DUPLICATE KEY UPDATE Github=VALUES(Github)", accID, tokens[2])
	check(err)
	return "Yes, sir."
}
