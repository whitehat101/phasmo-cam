package watch

import (
	"log"
	"regexp"

	"github.com/fsnotify/fsnotify"
)

var re = regexp.MustCompile(`SavedScreen\d.png$`)

// Watch ...
func Watch(basePath string, cb func(string)) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go watchLoop(watcher, cb)

	err = watcher.Add(basePath)
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan bool)
	<-done
}

func watchLoop(watcher *fsnotify.Watcher, cb func(string)) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if re.MatchString(event.Name) && event.Op&fsnotify.Write == fsnotify.Write {
				debounceTimer(event.Name, cb)
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("error:", err)
		}
	}
}
