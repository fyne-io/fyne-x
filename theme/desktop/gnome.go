package desktop

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
	ft "fyne.io/fyne/v2/theme"
)

// Script to get the colors from the Gnome GTK/Adwaita theme.
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

// GnomeTheme theme, based on the Gnome desktop manager. This theme uses GJS and gsettings to get
// the colors and font from the Gnome desktop.
type GnomeTheme struct {
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
	gnome := &GnomeTheme{}
	if gtkVersion == -1 {
		gtkVersion = gnome.getGTKVersion()
	}
	gnome.decodeTheme(gtkVersion, ft.VariantDark)

	return gnome
}

// Color returns the color for the given color name
//
// Implements: fyne.Theme
func (gnome *GnomeTheme) Color(name fyne.ThemeColorName, _ fyne.ThemeVariant) color.Color {

	switch name {
	case ft.ColorNameBackground:
		return gnome.bgColor
	case ft.ColorNameForeground:
		return gnome.fgColor
	case ft.ColorNameButton, ft.ColorNameInputBackground:
		return gnome.viewBgColor
	default:
		return ft.DefaultTheme().Color(name, gnome.variant)
	}
}

// Icon returns the icon for the given name.
//
// Implements: fyne.Theme
func (gnome *GnomeTheme) Icon(i fyne.ThemeIconName) fyne.Resource {
	return ft.DefaultTheme().Icon(i)
}

// Font returns the font for the given name.
//
// Implements: fyne.Theme
func (gnome *GnomeTheme) Font(s fyne.TextStyle) fyne.Resource {
	if gnome.font == nil {
		return ft.DefaultTheme().Font(s)
	}
	return gnome.font
}

// Size returns the size for the given name. It will scale the detected Gnome font size
// by the Gnome font factor.
//
// Implements: fyne.Theme
func (g *GnomeTheme) Size(s fyne.ThemeSizeName) float32 {
	switch s {
	case ft.SizeNameText:
		return g.fontScaleFactor * g.fontSize
	}
	return ft.DefaultTheme().Size(s)
}

func (gnome *GnomeTheme) decodeTheme(gtkVersion int, variant fyne.ThemeVariant) {
	// default
	gnome.bgColor = ft.DefaultTheme().Color(ft.ColorNameBackground, variant)
	gnome.fgColor = ft.DefaultTheme().Color(ft.ColorNameForeground, variant)
	gnome.fontSize = ft.DefaultTheme().Size(ft.SizeNameText)
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

func (gnome *GnomeTheme) getColors(gtkVersion int) {

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

func (gnome *GnomeTheme) fontScale() {

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

func (gnome *GnomeTheme) getFont() {

	gnome.font = ft.TextFont()
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

	// try to get the font as a TTF file
	gnome.setFont(strings.Join(parts[:len(parts)-1], " "))
}

func (gnome *GnomeTheme) setVariant() {
	// using the bgColor, detect if the theme is dark or light
	// if it is dark, set the variant to dark
	// if it is light, set the variant to light
	r, g, b, _ := gnome.bgColor.RGBA()

	brightness := (r/255*299 + g/255*587 + b/255*114) / 1000
	if brightness > 125 {
		gnome.variant = ft.VariantLight
	} else {
		gnome.variant = ft.VariantDark
	}
}

func (gnome *GnomeTheme) getGTKVersion() int {
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
			// now check if it is gtk4 compatible
			if _, err := os.Stat(path + "gtk-4.0/gtk.css"); err == nil {
				// it is gtk4
				return 4
			}
			if _, err := os.Stat(path + "gtk-3.0/gtk.css"); err == nil {
				return 3
			}
		}
	}
	return 3 // default, but that may be a false positive now
}

func (gnome *GnomeTheme) setFont(fontname string) {

	fontpath, err := getFontPath(fontname)
	if err != nil {
		log.Println(err)
		return
	}

	ext := filepath.Ext(fontpath)
	if ext != ".ttf" {
		font, err := converToTTF(fontpath)
		if err != nil {
			log.Println(err)
			return
		}
		gnome.font = fyne.NewStaticResource(fontpath, font)
	} else {
		font, err := ioutil.ReadFile(fontpath)
		if err != nil {
			log.Println(err)
			return
		}
		gnome.font = fyne.NewStaticResource(fontpath, font)
	}

}
