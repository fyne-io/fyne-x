package theme

import (
	"encoding/json"
	"image/color"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
	ft "fyne.io/fyne/v2/theme"
)

const gjsScript = `
let gtkVersion = Number(ARGV[0] || 4);
imports.gi.versions.Gtk = gtkVersion + ".0";

const { Gtk } = imports.gi;
if (gtkVersion === 3) {
  Gtk.init(null);
} else {
  Gtk.init();
}

const colors = {
  viewbg: [],
  viewfg: [],
  background: [],
  foreground: [],
  borders: [],
};

const win = new Gtk.Window();
const ctx = win.get_style_context();

let [ok, bg] = [false, null];

[ok, bg] = ctx.lookup_color("theme_base_color");
if (!ok) {
  [ok, bg] = ctx.lookup_color("view_bg_color");
}
colors.viewbg = [bg.red, bg.green, bg.blue, bg.alpha];

[ok, bg] = ctx.lookup_color("theme_text_color");
if (!ok) {
  [ok, bg] = ctx.lookup_color("view_fg_color");
}
colors.viewfg = [bg.red, bg.green, bg.blue, bg.alpha];

[ok, bg] = ctx.lookup_color("theme_bg_color");
if (!ok) {
  [ok, bg] = ctx.lookup_color("window_bg_color");
}
colors.background = [bg.red, bg.green, bg.blue, bg.alpha];

[ok, bg] = ctx.lookup_color("theme_fg_color");
if (!ok) {
  [ok, bg] = ctx.lookup_color("window_fg_color");
}
colors.foreground = [bg.red, bg.green, bg.blue, bg.alpha];

[ok, bg] = ctx.lookup_color("borders");
if (!ok) {
  [ok, bg] = ctx.lookup_color("unfocused_borders");
}
colors.borders = [bg.red, bg.green, bg.blue, bg.alpha];

print(JSON.stringify(colors));
`

// Gnome theme, based on the Gnome desktop environment preferences.
type Gnome struct {
	bgColor         color.Color
	fgColor         color.Color
	viewBgColor     color.Color
	viewFgColor     color.Color
	borderColor     color.Color
	fontScaleFactor float32
	font            fyne.Resource
	fontSize        float32
	variant         fyne.ThemeVariant
}

// NewGnomeTheme returns a new Gnome theme based on the given gtk version. If gtkVersion is -1,
// the theme will try to determine the higher Gtk version available for the current GtkTheme.
func NewGnomeTheme(gtkVersion int) fyne.Theme {
	gnome := &Gnome{}
	if gtkVersion == -1 {
		gtkVersion = gnome.getGTKVersion()
	}
	gnome.decodeTheme(gtkVersion, theme.VariantDark)

	return gnome
}

// Color returns the color for the given color name
//
// Implements: fyne.Theme
func (gnome *Gnome) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {

	switch name {
	case theme.ColorNameBackground:
		return gnome.bgColor
	case theme.ColorNameForeground:
		return gnome.fgColor
	case theme.ColorNameButton, theme.ColorNameInputBackground:
		return gnome.viewBgColor
	default:
		return ft.DefaultTheme().Color(name, gnome.variant)
	}
}

// Icon returns the icon for the given name.
//
// Implements: fyne.Theme
func (g *Gnome) Icon(i fyne.ThemeIconName) fyne.Resource {
	return ft.DefaultTheme().Icon(i)
}

// Font returns the font for the given name.
//
// Implements: fyne.Theme
func (g *Gnome) Font(s fyne.TextStyle) fyne.Resource {
	return ft.DefaultTheme().Font(s)
}

// Size returns the size for the given name. It will scale the detected Gnome font size
// by the Gnome font factor.
//
// Implements: fyne.Theme
func (g *Gnome) Size(s fyne.ThemeSizeName) float32 {
	switch s {
	case theme.SizeNameText:
		return g.fontScaleFactor * g.fontSize
	}
	return ft.DefaultTheme().Size(s)
}

func (gnome *Gnome) decodeTheme(gtkVersion int, variant fyne.ThemeVariant) {
	// default
	gnome.bgColor = theme.DefaultTheme().Color(theme.ColorNameBackground, variant)
	gnome.fgColor = theme.DefaultTheme().Color(theme.ColorNameForeground, variant)
	gnome.fontSize = theme.DefaultTheme().Size(theme.SizeNameText)
	wg := sync.WaitGroup{}

	// make things faster in concurrent mode
	wg.Add(3)
	go func() {
		gnome.getColors(gtkVersion)
		gnome.setVariant()
		wg.Done()
	}()
	go func() {
		gnome.getFont()
		wg.Done()
	}()
	go func() {
		gnome.fontScale()
		wg.Done()
	}()
	wg.Wait()
}

func (gnome *Gnome) getColors(gtkVersion int) {

	// we will call gjs to get the colors
	gjs, err := exec.LookPath("gjs")
	if err != nil {
		log.Println(err)
		return
	}

	// create a temp file to store the colors
	f, err := ioutil.TempFile("", "fyne-theme-gnome-")
	if err != nil {
		log.Println(err)
		return
	}
	defer os.Remove(f.Name())

	// write the script to the temp file
	_, err = f.WriteString(gjsScript)
	if err != nil {
		log.Println(err)
		return
	}

	// run the script
	cmd := exec.Command(gjs, f.Name(), strconv.Itoa(gtkVersion))
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err, string(out))
		return
	}

	// decode json to apply to the gnome theme
	colors := struct {
		WindowBGcolor []float32 `json:"background,-"`
		WindowFGcolor []float32 `json:"foreground,-"`
		ViewBGcolor   []float32 `json:"viewbg,-"`
		ViewFGcolor   []float32 `json:"viewfg,-"`
		Borders       []float32 `json:"borders,-"`
	}{}
	err = json.Unmarshal(out, &colors)
	if err != nil {
		log.Println(err)
		return
	}

	// convert the colors to fyne colors
	gnome.bgColor = color.RGBA{
		R: uint8(colors.WindowBGcolor[0] * 255),
		G: uint8(colors.WindowBGcolor[1] * 255),
		B: uint8(colors.WindowBGcolor[2] * 255),
		A: uint8(colors.WindowBGcolor[3] * 255)}

	gnome.fgColor = color.RGBA{
		R: uint8(colors.WindowFGcolor[0] * 255),
		G: uint8(colors.WindowFGcolor[1] * 255),
		B: uint8(colors.WindowFGcolor[2] * 255),
		A: uint8(colors.WindowFGcolor[3] * 255)}

	gnome.borderColor = color.RGBA{
		R: uint8(colors.Borders[0] * 255),
		G: uint8(colors.Borders[1] * 255),
		B: uint8(colors.Borders[2] * 255),
		A: uint8(colors.Borders[3] * 255)}

	gnome.viewBgColor = color.RGBA{
		R: uint8(colors.ViewBGcolor[0] * 255),
		G: uint8(colors.ViewBGcolor[1] * 255),
		B: uint8(colors.ViewBGcolor[2] * 255),
		A: uint8(colors.ViewBGcolor[3] * 255)}

	gnome.viewFgColor = color.RGBA{
		R: uint8(colors.ViewFGcolor[0] * 255),
		G: uint8(colors.ViewFGcolor[1] * 255),
		B: uint8(colors.ViewFGcolor[2] * 255),
		A: uint8(colors.ViewFGcolor[3] * 255)}

}

func (gnome *Gnome) fontScale() {

	// for any error below, we will use the default
	gnome.fontScaleFactor = 1

	// call gsettings get org.gnome.desktop.interface text-scaling-factor
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "text-scaling-factor")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}

	// get the text scaling factor
	ts := strings.TrimSpace(string(out))
	textScale, err := strconv.ParseFloat(ts, 32)
	if err != nil {
		return
	}

	// return the text scaling factor
	gnome.fontScaleFactor = float32(textScale)
}

func (gnome *Gnome) getFont() {

	gnome.font = theme.TextFont()
	// call gsettings get org.gnome.desktop.interface font-name
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "font-name")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		log.Println(string(out))
		return
	}
	// try to get the font as a TTF file
	fontFile := strings.TrimSpace(string(out))
	fontFile = strings.Trim(fontFile, "'")
	// the fontFile string is in the format: Name size, eg: "Sans Bold 12", so get the size
	parts := strings.Split(fontFile, " ")
	fontSize := parts[len(parts)-1]
	// convert the size to a float
	size, err := strconv.ParseFloat(fontSize, 32)
	// apply this to the fontScaleFactor
	gnome.fontSize = float32(size)
}

func (gnome *Gnome) setVariant() {
	// using the bgColor, detect if the theme is dark or light
	// if it is dark, set the variant to dark
	// if it is light, set the variant to light
	r, g, b, _ := gnome.bgColor.RGBA()

	brightness := (r/255*299 + g/255*587 + b/255*114) / 1000
	if brightness > 125 {
		gnome.variant = theme.VariantLight
	} else {
		gnome.variant = theme.VariantDark
	}
}

func (gnome *Gnome) getGTKVersion() int {
	// call gsettings get org.gnome.desktop.interface gtk-theme
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "gtk-theme")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		log.Println(string(out))
		return 3 // default to Gtk 3
	}
	themename := strings.TrimSpace(string(out))
	themename = strings.Trim(themename, "'")

	// ok so now, find if the theme is gtk4, either fallback to gtk3
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
		return 3 // default to Gtk 3
	}

	possiblePaths := []string{
		home + "/.local/share/themes/",
		home + "/.themes/",
		`/usr/local/share/themes/`,
		`/usr/share/themes/`,
	}

	for _, path := range possiblePaths {
		path = filepath.Join(path, themename)
		if _, err := os.Stat(path); err == nil {
			// found the theme directory
			// now check if it is gtk4
			if _, err := os.Stat(path + "gtk-4.0/gtk.css"); err == nil {
				// it is gtk4
				return 4
			} else {
				// it is gtk3
				return 3
			}
		}
	}
	return 3
}
