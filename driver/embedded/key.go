package fynex

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/embedded"
)

func triggerKey(key uint16, dir embedded.KeyDirection, queue chan embedded.Event) {
	name := fyne.KeyName("")
	switch key {
	case 1:
		name = fyne.KeyUp
	case 2:
		name = fyne.KeyDown
	case 3:
		name = fyne.KeyRight
	case 4:
		name = fyne.KeyLeft
	case 27: // esc
		name = fyne.KeyEscape
	case 13: // ret
		name = fyne.KeyReturn
	case 8: // backspace
		name = fyne.KeyBackspace
	case 9: // tab
		name = fyne.KeyTab
	case ' ':
		name = fyne.KeySpace
	default:
		if key > ' ' && key < '~' {
			name = fyne.KeyName(rune(key))
			if key >= 'a' && key <= 'z' {
				name = fyne.KeyName(rune(key) - 'a' + 'A')
			}
		}
	}
	queue <- &embedded.KeyEvent{Name: name, Direction: dir}

	// visible characters
	if key >= ' ' && key <= '~' {
		if dir == embedded.KeyReleased {
			return
		}

		queue <- &embedded.CharacterEvent{Rune: rune(key)}
	}
}
