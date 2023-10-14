package main

import (
	"encoding/json"
	"net/http"

	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	xwidget "fyne.io/x/fyne/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Completion Entry Demo")

	entry := xwidget.NewCompletionEntry([]string{})
	entry.OnChanged = func(s string) {
		// completion start for text length >= 3
		if len(s) < 3 {
			entry.HideCompletion()
			return
		}

		// Make a search on wikipedia
		resp, err := http.Get(
			"https://en.wikipedia.org/w/api.php?action=opensearch&search=" + entry.Text,
		)
		if err != nil {
			entry.HideCompletion()
			return
		}

		// Get the list of possible completion
		var results [][]string
		json.NewDecoder(resp.Body).Decode(&results)

		// no results
		if len(results) == 0 {
			entry.HideCompletion()
			return
		}

		// then show them
		entry.SetOptions(results[1])
		entry.ShowCompletion()
	}

	w.SetContent(container.NewBorder(entry, nil, nil, nil, widget.NewLabel("Enter 3 chars to start search...")))
	w.ShowAndRun()
}
