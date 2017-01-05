package main

import "strings"

func checkDbConnection() {
	err := db.Ping()
	check(err)
}

func getTokens(c string, msg string, help func() string) []string {
	tokens := strings.Fields(msg)
	if len(tokens) < 2 {
		session.ChannelMessageSend(c, "Yes?")
		return nil
	}
	if tokens[1] == "help" {
		session.ChannelMessageSend(c, help())
		return nil
	}
	return tokens
}
