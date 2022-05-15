package desktop

import (
	"errors"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	ft "fyne.io/fyne/v2/theme"
)

// KDE theme is based on the KDE or Plasma theme.
type KDE struct {
	variant    fyne.ThemeVariant
	bgColor    color.Color
	fgColor    color.Color
	fontConfig string
	fontSize   float32
	font       fyne.Resource
}

// NewKDETheme returns a new KDE theme.
func NewKDETheme() fyne.Theme {
	kde := &KDE{
		variant: ft.VariantDark,
	}
	kde.decodeTheme()
	log.Println(kde)

	return kde
}

// Color returns the color for the specified name.
//
// Implements: fyne.Theme
func (k *KDE) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {
	switch name {
	case ft.ColorNameBackground:
		return k.bgColor
	case ft.ColorNameForeground:
		return k.fgColor
	}
	return ft.DefaultTheme().Color(name, k.variant)
}

// Icon returns the icon for the specified name.
//
// Implements: fyne.Theme
func (k *KDE) Icon(i fyne.ThemeIconName) fyne.Resource {
	return ft.DefaultTheme().Icon(i)
}

// Font returns the font for the specified name.
//
// Implements: fyne.Theme
func (k *KDE) Font(s fyne.TextStyle) fyne.Resource {
	if k.font != nil {
		return k.font
	}
	return ft.DefaultTheme().Font(s)
}

// Size returns the size of the font for the specified text style.
//
// Implements: fyne.Theme
func (k *KDE) Size(s fyne.ThemeSizeName) float32 {
	if s == ft.SizeNameText {
		return k.fontSize
	}
	return ft.DefaultTheme().Size(s)
}

// decodeTheme initialize the theme.
func (k *KDE) decodeTheme() {
	k.loadScheme()
	k.setFont()
}

func (k *KDE) loadScheme() error {
	// the theme name is declared in ~/.config/kdedefaults/kdeglobals
	// in the ini section [General] as "ColorScheme" entry
	homedir, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	content, err := ioutil.ReadFile(filepath.Join(homedir, ".config/kdeglobals"))
	if err != nil {
		return err
	}

	section := ""
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "[") {
			section = strings.ReplaceAll(line, "[", "")
			section = strings.ReplaceAll(section, "]", "")
		}
		if section == "Colors:Window" {
			if strings.HasPrefix(line, "BackgroundNormal=") {
				k.bgColor = k.parseColor(strings.ReplaceAll(line, "BackgroundNormal=", ""))
			}
			if strings.HasPrefix(line, "ForegroundNormal=") {
				k.fgColor = k.parseColor(strings.ReplaceAll(line, "ForegroundNormal=", ""))
			}
		}
		if section == "General" {
			if strings.HasPrefix(line, "font=") {
				k.fontConfig = strings.ReplaceAll(line, "font=", "")
			}
		}
	}

	return errors.New("Unable to find the KDE color scheme")
}

func (k *KDE) parseColor(col string) color.Color {
	// the color is in the form r,g,b,
	// we need to convert it to a color.Color

	// split the string
	cols := strings.Split(col, ",")
	// convert the string to int
	r, _ := strconv.Atoi(cols[0])
	g, _ := strconv.Atoi(cols[1])
	b, _ := strconv.Atoi(cols[2])

	// convert the int to a color.Color
	return color.RGBA{uint8(r), uint8(g), uint8(b), 0xff}
}

func (k *KDE) setFont() {

	if k.fontConfig == "" {
		log.Println("WTF")
		return
	}
	// the font is in the form "fontname,size,...", so we can split it
	font := strings.Split(k.fontConfig, ",")
	name := font[0]
	size, _ := strconv.ParseFloat(font[1], 32)
	k.fontSize = float32(size)

	// we need to load the font, Gnome struct has got some nice methods
	fontpath, err := getFontPath(name)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(fontpath)
	if filepath.Ext(fontpath) == ".ttf" {
		font, err := ioutil.ReadFile(fontpath)
		if err != nil {
			log.Println(err)
			return
		}
		k.font = fyne.NewStaticResource(fontpath, font)
	} else {
		font, err := converToTTF(fontpath)
		if err != nil {
			log.Println(err)
			return
		}
		k.font = fyne.NewStaticResource(fontpath, font)
	}
}
