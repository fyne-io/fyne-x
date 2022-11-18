package desktop

import (
	"io/ioutil"
	"os"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

func setup() (tmp, home string) {
	// create a false home directory
	var err error
	tmp, err = os.MkdirTemp("", "fyne-test-")
	if err != nil {
		panic(err)
	}
	home = os.Getenv("HOME")
	os.Setenv("HOME", tmp)

	// creat a false KDE configuration
	if err = os.MkdirAll(tmp+"/.config", 0755); err != nil {
		panic(err)
	}
	content := []byte("[General]\nwidgetStyle=GTK")
	ioutil.WriteFile(tmp+"/.config/kdeglobals", content, 0644)
	return
}

func teardown(tmp, home string) {
	os.RemoveAll(tmp)
	os.Setenv("HOME", home)
}
func TestKDETheme(t *testing.T) {
	tmp, home := setup()
	defer teardown(tmp, home)
	app := test.NewApp()
	app.Settings().SetTheme(NewKDETheme())
	win := app.NewWindow("Test")
	defer win.Close()
	win.Resize(fyne.NewSize(200, 200))
	win.SetContent(widget.NewLabel("Hello"))
}
