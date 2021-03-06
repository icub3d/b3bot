package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

func init() {
	register(rio)
}

type raiderIOChar struct {
	Name     string
	Region   string
	Realm    string
	Profile  string `json:"profile_url"`
	Gear     map[string]float64
	MPScores map[string]float64 `json:"mythic_plus_scores"`
}

var rio = &cmd{
	name: "rio",
	help: " - character [realm] [region] - get the raider.io score and link for the given character. You can pass the realm and region as well (default eredar, us).",
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
		region = params[2]
		server = params[1]
		char = params[0]
	} else if len(params) == 2 {
		server = params[1]
		char = params[0]
	}

	u := fmt.Sprintf("https://raider.io/api/v1/characters/profile?region=%s&realm=%s&name=%s&fields=gear,mythic_plus_scores", region, server, char)
	resp, err := http.Get(u)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	c := raiderIOChar{}
	b := &bytes.Buffer{}
	_, err = b.ReadFrom(resp.Body)
	logrus.Infof("%s", b.String())
	dec := json.NewDecoder(b)
	err = dec.Decode(&c)
	if err != nil {
		return err
	}
	reply(s, m, fmt.Sprintf("%v: ilvl - %v, score - %v\n%v",
		c.Name, c.Gear["item_level_total"], c.MPScores["all"], c.Profile))

	return nil
}
