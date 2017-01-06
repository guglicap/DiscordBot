package main

import (
	"log"
	"strings"
)

func checkDbConnection() {
	err := db.Ping()
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func help(commands map[string]Command) string {
	result := "Syntax:\n  !bot <command>\n  Commands:\n"
	for k, v := range commands {
		if len(v.desc) > 0 {
			result += "\t" + k + ":\t" + v.desc + "\n"
		}
	}
	return result
}

func getTokens(c string, msg string) []string {
	tokens := strings.Fields(msg)
	if len(tokens) < 2 {
		session.ChannelMessageSend(c, "Yes?")
		return nil
	}
	return tokens
}
