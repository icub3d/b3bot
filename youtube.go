package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/spf13/viper"
	"google.golang.org/api/googleapi/transport"
	youtube "google.golang.org/api/youtube/v3"
)

func getYoutubeVideos(s *discordgo.Session) error {
	// Make the connection.
	wc := &http.Client{
		Transport: &transport.APIKey{Key: viper.GetString("youtube.key")},
	}
	yt, err := youtube.New(wc)
	if err != nil {
		return fmt.Errorf("connecting to youtube: %v", err)
	}

	// Get our list of channels to search
	channels := viper.GetStringSlice("youtube.channels")
	for _, channel := range channels {
		// Get a list the channel caller.
		call := yt.Channels.List("contentDetails").Id(channel)
		r, err := call.Do()
		if err != nil {
			return fmt.Errorf("retreiving channel list: %v", err)
		}

		for _, item := range r.Items {
			// Get a list of videos (just the latest few).
			plCall := yt.PlaylistItems.List("snippet").
				PlaylistId(item.ContentDetails.RelatedPlaylists.Uploads).
				MaxResults(5)
			plr, err := plCall.Do()
			if err != nil {
				return fmt.Errorf("retreiving video list: %v", err)
			}
			for _, i := range plr.Items {
				// Get when it was published
				when, err := time.Parse(time.RFC3339, i.Snippet.PublishedAt)
				if err != nil {
					return fmt.Errorf("parsing time: %v", err)
				}
				// Show it if it was published in the last 5
				// minutes. This is based on the assumption we check
				// every 5 minutes.
				if when.After(time.Now().Add(-5 * time.Minute)) {
					msg := fmt.Sprintf("https://www.youtube.com/watch?v=%v",
						i.Snippet.ResourceId.VideoId)
					s.ChannelMessageSend("236232987086290944", msg)
				}
			}
		}
	}

	return nil
}
