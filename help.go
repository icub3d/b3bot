package main

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(help)
}

var help = &cmd{
	name: "help",
	help: " - [cmd] - get help on commands this bot can perform",
	init: initHelp,
	run:  runHelp,
}

func initHelp() {
}

func runHelp(s *discordgo.Session, m *discordgo.MessageCreate, params []string) error {
	if len(params) < 1 {
		cmds := list()
		reply(s, m, fmt.Sprintf("cmds: %v\nfor details, use !help [cmd]",
			strings.Join(cmds, " ")))
		return nil
	}

	cmd := get(params[0])
	if cmd == nil {
		reply(s, m, fmt.Sprintf("unknown command '%v': use !help for commands", params[0]))
		return nil
	}

	reply(s, m, fmt.Sprintf("%s%s", cmd.name, cmd.help))
	return nil
}
