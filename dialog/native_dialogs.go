package dialog

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	fyneDialog "fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/storage/repository"
)

type wm uint8

const (
	wm_UNKNOWN wm = iota
	wm_GNOME
	wm_KDE
)

type FileSelector struct {
	Title        string
	Filters      []string
	callback     func(fyne.URIReadCloser, error)
	parentWindow fyne.Window
}

func detectWM() wm {
	// detect WM
	xdgCurrentDesktop := os.Getenv("XDG_CURRENT_DESKTOP")
	switch xdgCurrentDesktop {
	case "GNOME":
		return wm_GNOME
	case "KDE":
		return wm_KDE
	default:
		return wm_UNKNOWN
	}
}

func NewFileSelector(callback func(fyne.URIReadCloser, error), parent fyne.Window) *FileSelector {
	return &FileSelector{
		Title:    "Select a file",
		callback: callback,
	}
}

func ShowFileSelector(callback func(fyne.URIReadCloser, error), parent fyne.Window) {
	NewFileSelector(callback, parent).Show()
}

func (f *FileSelector) AddFilter(name string, extensions ...string) {
	f.Filters = append(f.Filters, name)
}

func (f *FileSelector) Show() {
	wm := detectWM()
	log.Println("WM:", wm)
	switch wm {
	case wm_GNOME:
		// use zenity
		//command := exec.Command("zenity", "--file-selection", "--title", f.Title)
		command := exec.Command("zenity", "--file-selection", "--title", f.Title, "--file-filter", "All files | *")
		out, err := command.Output()
		if err != nil {
			f.defaultDialog()
			return
		}
		path := string(out)
		u, err := f.getReadCloser(path)
		if err != nil {
			log.Println("Error:", err)
			return
		}
		f.callback(u, nil)
	case wm_KDE:
		// use kdialog
		//command := exec.Command("kdialog", "--getopenfilename", f.Title)
		command := exec.Command("kdialog", "--getopenfilename", f.Title, "All files | *")
		out, err := command.Output()
		if err != nil {
			f.defaultDialog()
			return
		}
		path := string(out)
		u, err := f.getReadCloser(path)
		if err != nil {
			log.Println("Error:", err)
			return
		}
		f.callback(u, nil)

	default:
		// use native dialog
		f.defaultDialog()
	}

}

func (f *FileSelector) defaultDialog() {
	d := fyneDialog.NewFileOpen(f.callback, f.parentWindow)
	d.Show()
}

func (f *FileSelector) getReadCloser(path string) (fyne.URIReadCloser, error) {
	uri := repository.NewFileURI(strings.TrimSpace(path))
	return storage.Reader(uri)
}
