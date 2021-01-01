package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/whitehat101/phasmo-cam/discord"
	"github.com/whitehat101/phasmo-cam/watch"
)

type action func(string) error

const (
	timestamp    = "PhasmoCam_2006-01-02T15_04_05Z0700.png"
	fakeHomePath = `%UserProfile%\Pictures\Phasmophobia`
)

var (
	token        = flag.String("discord-token", "", "Discord bot token")
	channelID    = flag.String("discord-channelID", "", "Discord channel ID")
	basePath     = flag.String("phasmo-dir", `C:\Program Files (x86)\Steam\steamapps\common\Phasmophobia`, "Phasmophobia install dir")
	savePath     = flag.String("save-dir", fakeHomePath, "camera backup dir")
	afterActions = []action{}
)

func init() {
	flag.Parse()
	if *savePath == fakeHomePath {
		dir, _ := os.UserHomeDir()
		*savePath = filepath.Join(dir, "Pictures", "Phasmophobia")
	}
}

func main() {
	if *token != "" && *channelID != "" {
		log.Println("Connecting to Discord")
		if cb, err := discord.Factory(*token, *channelID); err == nil {
			afterActions = append(afterActions, cb)
		} else {
			log.Fatal(err)
		}
	}

	matches, _ := filepath.Glob(filepath.Join(*basePath, "SavedScreen*.png"))
	for _, match := range matches {
		uploadCallback(match)
	}
	watch.Watch(*basePath, uploadCallback)
}

func uploadCallback(source string) {
	var t time.Time
	var len int64
	if s, err := os.Stat(source); err == nil {
		t = s.ModTime()
		len = s.Size()
	} else {
		t = time.Now()
	}
	dest := filepath.Join(*savePath, t.Format(timestamp))
	if s, err := os.Stat(dest); err == nil && s.Size() == len {
		log.Println("SKIP", source)
		return
	}
	if err := copy(source, dest); err != nil {
		log.Fatal(err)
	}
	log.Println("BACKUP", source)
	for _, cb := range afterActions {
		if err := cb(dest); err != nil {
			log.Println(err)
		}
	}
}

// Copy the src file to dst.
// Any existing file will be overwritten and will not copy file attributes.
func copy(src string, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
