package main

import "github.com/bwmarrin/discordgo"

type replyString (func(*discordgo.Message) string)

//Command is a struct. Holds the help description and the function callback for the action to take.
type Command struct {
	desc  string
	reply replyString
}

//PublicCommands is an array of commands that can be used in a public channel.
var PublicCommands = map[string]Command{
	"me": {
		desc:  "Tells you who you are.",
		reply: me,
	},
	"you": {
		desc:  "Tells you who I am.",
		reply: bot,
	},
	"whois": {
		desc:  "Use followed by a name.",
		reply: whois,
	},
	"kappa": {
		desc:  "You know this one. Unless you're PF.",
		reply: kappa,
	},
	"github": {
		desc:  "Use followed by a name, links you to that person's github.",
		reply: githubPublic,
	},
}

//PrivateCommands is a list of commands that can be used in a private channel.
var PrivateCommands = map[string]Command{
	"github": {
		desc:  "Set your github link",
		reply: githubPrivate,
	},
}