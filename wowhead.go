package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	customsearch "google.golang.org/api/customsearch/v1"
	"google.golang.org/api/googleapi/transport"
)

func init() {
	register(wowhead)
}

var wowhead = &cmd{
	name: "wowhead",
	help: " - query - search wowhead for the given query and return the top result",
	init: initWowhead,
	run:  runWowhead,
}

func initWowhead() {
}

func runWowhead(s *discordgo.Session, m *discordgo.MessageCreate, params []string) error {
	wc := &http.Client{
		Transport: &transport.APIKey{Key: viper.GetString("wowhead.key")},
	}
	css, err := customsearch.New(wc)
	if err != nil {
		return fmt.Errorf("creating search client: %v", err)
	}
	query := strings.Join(params, " ")
	search, err := css.Cse.List(query).Cx(viper.GetString("wowhead.csid")).Do()
	if err != nil {
		return fmt.Errorf("searching: %v", err)
	}

	if len(search.Items) < 1 {
		reply(s, m, "no results found")
		return nil
	}

	reply(s, m, search.Items[0].Link)
	return nil
}
