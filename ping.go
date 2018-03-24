package main

import (
	"github.com/bwmarrin/discordgo"
)

func init() {
	register(ping)
}

var ping = &cmd{
	name: "ping",
	help: "- get a pong response back",
	init: initPing,
	run:  runPing,
}

func initPing() {
}

func runPing(s *discordgo.Session, m *discordgo.MessageCreate, params []string) error {
	reply(s, m, "pong")
	return nil
}
