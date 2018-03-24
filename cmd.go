package main

import (
	"github.com/bwmarrin/discordgo"
)

type cmd struct {
	name string
	help string
	init func()
	run  func(s *discordgo.Session, m *discordgo.MessageCreate, params []string) error
}

var cmds = map[string]*cmd{}

func register(c *cmd) {
	cmds[c.name] = c
}

func get(name string) *cmd {
	return cmds[name]
}

func list() []string {
	names := []string{}
	for name := range cmds {
		names = append(names, name)
	}
	return names
}

func initAll() {
	for _, cmd := range cmds {
		cmd.init()
	}
}
