package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	// Initialize config
	viper.SetConfigName(".b3bot")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("reading config: %v", err)
	}

	// Setup logging.
	level, err := logrus.ParseLevel(viper.GetString("logrus.level"))
	if err != nil {
		logrus.Fatalf("parsing logrus.level from config: %v", err)
	}
	logrus.SetLevel(level)

	// init our cmds.
	initAll()

	// Create a new Discord session using the provided bot token.
	dg, err := discordgo.New("Bot " + viper.GetString("discord.token"))
	if err != nil {
		logrus.Fatalf("creating discord session: %v", err)
	}

	// Register the messageCreate func as a callback for MessageCreate events.
	dg.AddHandler(messageCreate)

	// Open a websocket connection to Discord and begin listening.
	err = dg.Open()
	if err != nil {
		logrus.Fatalf("opening discord connection: %v", err)
	}

	go cron(dg)

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	dg.Close()
}

func cron(s *discordgo.Session) {
	// TODO since we don't have any real open resources at this point,
	// I'm just quiting here. If we ever change that, we'll need to
	// exit clenanly.
	for {
		select {
		// Check once every 5 minutes
		case <-time.After(5 * time.Minute):
			getYoutubeVideos(s)
		}
	}
}

// This function will be called (due to AddHandler above) every time a new
// message is created on any channel that the autenticated bot has access to.
func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Ignore all messages created by the bot itself
	// This isn't required in this specific example but it's a good practice.
	if m.Author.ID == s.State.User.ID || !strings.HasPrefix(m.Content, "!") {
		return
	}

	parts := strings.Split(m.Content, " ")
	if len(parts) < 1 {
		logrus.Debugf("split message too small: %v", m.Content)
		return
	}

	c := strings.TrimPrefix(parts[0], "!")
	cmd := get(c)
	if cmd == nil {
		logrus.Debugf("unknown command: %v", c)
		reply(s, m, fmt.Sprintf("unknown command '%v': use !help for commands", c))
		return
	}

	parts = parts[1:]
	if err := cmd.run(s, m, parts); err != nil {
		logrus.Errorf("running command '%v': %v", c, err)
		reply(s, m, fmt.Sprintf("running command '%v' failed: %v", c, err))
	}
}

func reply(s *discordgo.Session, m *discordgo.MessageCreate, message string) {
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("@%s, %s", m.Author.Username, message))
}
