package discord

import (
	"os"
	"path/filepath"

	"github.com/bwmarrin/discordgo"
)

// Factory for the Discord upload action
func Factory(token string, channelID string) (func(string) error, error) {
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}

	// Before using the message endpoint, you must connect to
	// and identify with a gateway at least once.
	err = s.Open()
	if err != nil {
		return nil, err
	}
	err = s.Close()
	if err != nil {
		return nil, err
	}

	cb := func(src string) error {
		f, err := os.Open(src)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = s.ChannelFileSend(channelID, filepath.Base(src), f)
		if err != nil {
			return err
		}
		return nil
	}

	return cb, nil
}
