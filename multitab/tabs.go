package multitab

import (
	"fyne.io/fyne/v2"
)

type Tab struct {
	title   string
	content fyne.CanvasObject
	icon    fyne.Resource
	tooltip string
}
