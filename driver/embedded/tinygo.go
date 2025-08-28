//go:build tinygo || noos

package fynex

import (
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/embedded"
)

const noKey = uint16(0xFFFF)

type tinygo struct {
	queue chan embedded.Event
}

func NewTinyGoDriver() embedded.Driver {
	return &tinygo{queue: make(chan embedded.Event)}
}

func (t *tinygo) Queue() chan embedded.Event {
	return t.queue
}

func (t *tinygo) Run() {
	go runEvents(t.queue)
	fyne.CurrentApp().Run()
}

func runEvents(queue chan embedded.Event) {
	key := noKey
	for {
		time.Sleep(time.Millisecond * 10) // don't poll too fast

		newKey := nextKey()
		if newKey == key {
			continue
		}

		if newKey == noKey {
			triggerKey(key, embedded.KeyReleased, queue)
			key = newKey
		} else {
			if newKey == 0x100 { // escape
				break
			}

			typed := mapKey(newKey)
			triggerKey(typed, embedded.KeyPressed, queue)
			key = newKey
		}
	}

	fyne.Do(fyne.CurrentApp().Quit)
}

func mapKey(key uint16) uint16 {
	// TODO handle shift...
	if key >= 'A' && key <= 'Z' {
		return key + 'a' - 'A'
	}

	switch key {
	case 0x100:
		return 27 // esc
	case 0x101:
		return 13 // ret
	case 0x103:
		return 8 // backspace
	case 0x102:
		return 9 // tab
	case 0x105:
		return 127 // delete

	case 265:
		return 1
	case 264:
		return 2
	case 263:
		return 4
	case 262:
		return 3
	}

	return key
}
