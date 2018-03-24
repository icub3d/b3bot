package main

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
)

func init() {
	register(rio)
}

type raiderIOChar struct {
	Name     string
	Region   string
	Realm    string
	Profile  string `json:"profile_url"`
	Gear     map[string]int
	MPScores map[string]int `json:"mythic_plus_scores"`
}

var rio = &cmd{
	name: "rio",
	help: " - [region] [realm] character - get the raider.io score and link for the given character. You can pass the region and realm as well (default us, eredar).",
	init: initRio,
	run:  runRio,
}

func initRio() {
}

func runRio(s *discordgo.Session, m *discordgo.MessageCreate, params []string) error {
	if len(params) < 1 {
		return fmt.Errorf("at least one argument required; see !help rio")
	}

	region := "us"
	server := "eredar"
	char := params[0]
	if len(params) == 3 {
		region = params[0]
		server = params[1]
		char = params[2]
	} else if len(params) == 2 {
		server = params[0]
		char = params[1]
	}

	u := fmt.Sprintf("https://raider.io/api/v1/characters/profile?region=%s&realm=%s&name=%s&fields=gear,mythic_plus_scores", region, server, char)
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c := raiderIOChar{}
	dec := json.NewDecoder(resp.Body)
	err = dec.Decode(&c)
	if err != nil {
		return err
	}

	reply(s, m, fmt.Sprintf("%v: ilvl - %v, score - %v\n%v",
		c.Name, c.Gear["item_level_total"], c.MPScores["all"], c.Profile))

	return nil
}
