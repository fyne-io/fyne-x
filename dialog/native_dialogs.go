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

// FileDialog is a file dialog that uses native file dialogs on Linux.
type FileDialog struct {
	*fyneDialog.FileDialog
	Title        string
	filter       storage.FileFilter
	callback     interface{}
	parentWindow fyne.Window
	save         bool
}

// NewFileOpen creates a new file dialog that uses native file dialogs on Linux.
func NewFileOpen(callback func(fyne.URIReadCloser, error), parent fyne.Window) *FileDialog {
	return &FileDialog{
		FileDialog: fyneDialog.NewFileOpen(callback, parent),
		Title:      "Select a file",
		callback:   callback,
		save:       false,
	}
}

func NewFileSave(callback func(fyne.URIWriteCloser, error), parent fyne.Window) *FileDialog {
	return &FileDialog{
		FileDialog: fyneDialog.NewFileSave(callback, parent),
		Title:      "Save file",
		callback:   callback,
		save:       true,
	}
}

// SaveFileDialog is a save file dialog that uses native file dialogs on Linux.
func ShowFileOpen(callback func(fyne.URIReadCloser, error), parent fyne.Window) {
	NewFileOpen(callback, parent).Show()
}

// SaveFileDialog is a save file dialog that uses native file dialogs on Linux. At this time, it does not support filters.
func (f *FileDialog) SetFilter(filter storage.FileFilter) {
	f.FileDialog.SetFilter(filter)
	f.filter = filter
}

// Show shows the file dialog.
func (f *FileDialog) Show() {
	wm := detectWM()
	log.Println("WM:", wm)
	switch wm {
	case wm_GNOME:
		if u, err := f.zenityFileDialog(); err != nil {
			log.Println("zenity error:", err)
			f.FileDialog.Show()
		} else {
			if !f.save {
				f.callback.(func(fyne.URIReadCloser, error))(u.(fyne.URIReadCloser), nil)
			} else {
				f.callback.(func(fyne.URIWriteCloser, error))(u.(fyne.URIWriteCloser), nil)
			}
		}
	case wm_KDE:
		if u, err := f.kdialogFileDialog(); err != nil {
			log.Println("kdialog error:", err)
			f.FileDialog.Show()
		} else {
			if !f.save {
				f.callback.(func(fyne.URIReadCloser, error))(u.(fyne.URIReadCloser), nil)
			} else {
				f.callback.(func(fyne.URIWriteCloser, error))(u.(fyne.URIWriteCloser), nil)
			}
		}
	default:
		f.FileDialog.Show()
	}

}

func (f *FileDialog) getReadCloser(path string) (fyne.URIReadCloser, error) {
	uri := repository.NewFileURI(strings.TrimSpace(path))
	return storage.Reader(uri)
}

func (f *FileDialog) getWriteCloser(path string) (fyne.URIWriteCloser, error) {
	uri := repository.NewFileURI(strings.TrimSpace(path))
	return storage.Writer(uri)
}

func (f *FileDialog) zenityFileDialog() (interface{}, error) {
	command := exec.Command("zenity", "--file-selection", "--title", f.Title, "--file-filter", "All files | *")
	out, err := command.Output()
	if err != nil {
		return nil, err
	}
	path := string(out)

	if f.save {
		return f.getWriteCloser(path)
	}

	return f.getReadCloser(path)
}

func (f *FileDialog) kdialogFileDialog() (interface{}, error) {
	command := exec.Command("kdialog", "--getopenfilename", f.Title, "All files | *")
	out, err := command.Output()
	if err != nil {
		return nil, err
	}
	path := string(out)
	if f.save {
		return f.getWriteCloser(path)
	}
	return f.getReadCloser(path)
}
