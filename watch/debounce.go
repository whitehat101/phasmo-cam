package watch

import "time"

type void struct{}

var (
	debounce = map[string](chan void){}
	cancel   = void{}
)

func debounceTimer(name string, cb func(string)) {
	if _, ok := debounce[name]; !ok {
		debounce[name] = make(chan void)
	}
	select {
	case debounce[name] <- cancel:
	default:
	}

	go func() {
		select {
		case <-debounce[name]:
		case <-time.After(250 * time.Millisecond):
			cb(name)
		}
	}()
}
