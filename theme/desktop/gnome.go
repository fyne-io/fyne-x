//go:build linux
// +build linux

package desktop

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	"github.com/srwiley/oksvg"
)

// mapping to gnome/gtk icon names.
var gnomeIconMap = map[fyne.ThemeIconName]string{
	theme.IconNameInfo:     "dialog-information",
	theme.IconNameError:    "dialog-error",
	theme.IconNameQuestion: "dialog-question",

	theme.IconNameFolder:     "folder",
	theme.IconNameFolderNew:  "folder-new",
	theme.IconNameFolderOpen: "folder-open",
	theme.IconNameHome:       "go-home",
	theme.IconNameDownload:   "download",

	theme.IconNameDocument:        "document",
	theme.IconNameFileImage:       "image",
	theme.IconNameFileApplication: "binary",
	theme.IconNameFileText:        "text",
	theme.IconNameFileVideo:       "video",
	theme.IconNameFileAudio:       "audio",
	theme.IconNameComputer:        "computer",
	theme.IconNameMediaPhoto:      "photo",
	theme.IconNameMediaVideo:      "video",
	theme.IconNameMediaMusic:      "music",

	theme.IconNameConfirm: "dialog-apply",
	theme.IconNameCancel:  "cancel",

	theme.IconNameCheckButton:        "checkbox-symbolic",
	theme.IconNameCheckButtonChecked: "checkbox-checked-symbolic",
	theme.IconNameRadioButton:        "radio-symbolic",
	theme.IconNameRadioButtonChecked: "radio-checked-symbolic",

	theme.IconNameArrowDropDown:   "arrow-down",
	theme.IconNameArrowDropUp:     "arrow-up",
	theme.IconNameNavigateNext:    "go-right",
	theme.IconNameNavigateBack:    "go-left",
	theme.IconNameMoveDown:        "go-down",
	theme.IconNameMoveUp:          "go-up",
	theme.IconNameSettings:        "document-properties",
	theme.IconNameHistory:         "history-view",
	theme.IconNameList:            "view-list",
	theme.IconNameGrid:            "view-grid",
	theme.IconNameColorPalette:    "color-select",
	theme.IconNameColorChromatic:  "color-select",
	theme.IconNameColorAchromatic: "color-picker-grey",
}

// Map Fyne colorname to Adwaita/GTK color names
// See https://gnome.pages.gitlab.gnome.org/libadwaita/doc/main/named-colors.html
var gnomeColorMap = map[fyne.ThemeColorName]string{
	theme.ColorNameBackground:      "theme_bg_color,window_bg_color",
	theme.ColorNameForeground:      "theme_text_color,view_fg_color",
	theme.ColorNameButton:          "theme_base_color,view_bg_color",
	theme.ColorNameInputBackground: "theme_base_color,view_bg_color",
	theme.ColorNamePrimary:         "accent_color,success_color",
	theme.ColorNameError:           "error_color",
}

// Script to get the colors from the Gnome GTK/Adwaita theme.
const gjsColorScript = `
let gtkVersion = Number(ARGV[0] || 4);
imports.gi.versions.Gtk = gtkVersion + ".0";

const { Gtk, Gdk } = imports.gi;
if (gtkVersion === 3) {
  Gtk.init(null);
} else {
  Gtk.init();
}

const colors = {};
const win = new Gtk.Window();
const ctx = win.get_style_context();
const colorMap = %s;

for (let col in colorMap) {
  let [ok, bg] = [false, null];
  let found = false;
  colorMap[col].split(",").forEach((fetch) => {
    [ok, bg] = ctx.lookup_color(fetch);
    if (ok && !found) {
      found = true;
      colors[col] = [bg.red, bg.green, bg.blue, bg.alpha];
    }
  });
}

print(JSON.stringify(colors));
`

// Script to get icons from theme.
const gjsIconsScript = `
let gtkVersion = Number(ARGV[0] || 4);
imports.gi.versions.Gtk = gtkVersion + ".0";
const iconSize = 96; // can be 8, 16, 24, 32, 48, 64, 96

const { Gtk, Gdk } = imports.gi;
if (gtkVersion === 3) {
  Gtk.init(null);
} else {
  Gtk.init();
}

let iconTheme = null;
const icons = %s; // the icon list to get
const iconset = {};

if (gtkVersion === 3) {
  iconTheme = Gtk.IconTheme.get_default();
} else {
  iconTheme = Gtk.IconTheme.get_for_display(Gdk.Display.get_default());
}

icons.forEach((name) => {
  try {
    if (gtkVersion === 3) {
      const icon = iconTheme.lookup_icon(name, iconSize, 0);
      iconset[name] = icon.get_filename();
    } else {
      const icon = iconTheme.lookup_icon(name, null, null, iconSize, null, 0);
      iconset[name] = icon.file.get_path();
    }
  } catch (e) {
    iconset[name] = null;
  }
});

print(JSON.stringify(iconset));
`

// GnomeTheme theme, based on the Gnome desktop manager. This theme uses GJS and gsettings to get
// the colors and font from the Gnome desktop.
type GnomeTheme struct {
	colors map[fyne.ThemeColorName]color.Color
	icons  map[string]string

	scaleFactor float32
	font        fyne.Resource
	fontSize    float32
	iconCache   map[string]fyne.Resource

	versionNumber int
	themeName     string
}

// Color returns the color for the given color name
//
// Implements: fyne.Theme
func (gnome *GnomeTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {

	if col, ok := gnome.colors[name]; ok {
		return col
	}

	return theme.DefaultTheme().Color(name, variant)
}

// Font returns the font for the given name.
//
// Implements: fyne.Theme
func (gnome *GnomeTheme) Font(s fyne.TextStyle) fyne.Resource {
	if gnome.font == nil {
		return theme.DefaultTheme().Font(s)
	}
	return gnome.font
}

// Icon returns the icon for the given name.
//
// Implements: fyne.Theme
func (gnome *GnomeTheme) Icon(i fyne.ThemeIconName) fyne.Resource {
	if icon, found := gnomeIconMap[i]; found {
		if resource := gnome.loadIcon(icon); resource != nil {
			return resource
		}
	}
	return theme.DefaultTheme().Icon(i)
}

// Invert is a specific Gnome/GTK option to invert the theme color for background of window and some input
// widget. This to help to imitate some GTK application with "views" inside the window.
func (gnome *GnomeTheme) Invert() {

	gnome.colors[theme.ColorNameBackground],
		gnome.colors[theme.ColorNameInputBackground],
		gnome.colors[theme.ColorNameButton] =
		gnome.colors[theme.ColorNameButton],
		gnome.colors[theme.ColorNameBackground],
		gnome.colors[theme.ColorNameBackground]
}

// Size returns the size for the given name. It will scale the detected Gnome font size
// by the Gnome font factor.
//
// Implements: fyne.Theme
func (gnome *GnomeTheme) Size(s fyne.ThemeSizeName) float32 {
	switch s {
	case theme.SizeNameText:
		return gnome.scaleFactor * gnome.fontSize
	}
	return theme.DefaultTheme().Size(s) * gnome.scaleFactor
}

// applyColors sets the colors for the Gnome theme. Colors are defined by a GJS script.
func (gnome *GnomeTheme) applyColors(gtkVersion int, wg *sync.WaitGroup) {

	if wg != nil {
		defer wg.Done()
	}
	// we will call gjs to get the colors
	gjs, err := exec.LookPath("gjs")
	if err != nil {
		log.Println("To activate the theme, please install gjs", err)
		return
	}

	// create a temp file to store the colors
	f, err := ioutil.TempFile("", "fyne-theme-gnome-*.js")
	if err != nil {
		log.Println(err)
		return
	}
	defer os.Remove(f.Name())

	// generate the js object from gnomeColorMap
	colormap := "{\n"
	for col, fetch := range gnomeColorMap {
		colormap += fmt.Sprintf(`    "%s": "%s",`+"\n", col, fetch)
	}
	colormap += "}"

	// write the script to the temp file
	script := fmt.Sprintf(gjsColorScript, colormap)
	_, err = f.WriteString(script)
	if err != nil {
		log.Println(err)
		return
	}

	// run the script
	cmd := exec.Command(gjs,
		f.Name(), strconv.Itoa(gtkVersion),
		fmt.Sprintf("%0.2f", 1.0),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("gjs error:", err, string(out))
		return
	}

	// decode json
	var colors map[fyne.ThemeColorName][]float32
	err = json.Unmarshal(out, &colors)
	if err != nil {
		log.Println("gjs error:", err, string(out))
		return
	}
	for name, rgba := range colors {
		// convert string arry to colors
		gnome.colors[name] = gnome.parseColor(rgba)
	}

}

// applyFont gets the font name from gsettings and set the font size. This also calls
// setFont() to set the font.
func (gnome *GnomeTheme) applyFont(wg *sync.WaitGroup) {

	if wg != nil {
		defer wg.Done()
	}

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
	if err != nil {
		log.Println(err)
		return
	}
	// apply this to the fontScaleFactor
	gnome.fontSize = float32(size)

	// try to get the font as a TTF file
	gnome.setFont(strings.Join(parts[:len(parts)-1], " "))
}

// applyFontScale find the font scaling factor in settings.
func (gnome *GnomeTheme) applyFontScale(wg *sync.WaitGroup) {
	if wg != nil {
		defer wg.Done()
	}
	// for any error below, we will use the default
	gnome.scaleFactor = 1

	// call gsettings get org.gnome.desktop.interface text-scaling-factor
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "text-scaling-factor")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return
	}

	// get the text scaling factor
	ts := strings.TrimSpace(string(out))
	scaleValue, err := strconv.ParseFloat(ts, 32)
	if err != nil {
		return
	}

	// return the text scaling factor
	gnome.scaleFactor = float32(scaleValue)
}

// applyIcons gets the icon theme from gsettings and call GJS script to get the icon set.
func (gnome *GnomeTheme) applyIcons(gtkVersion int, wg *sync.WaitGroup) {

	if wg != nil {
		defer wg.Done()
	}

	gjs, err := exec.LookPath("gjs")
	if err != nil {
		log.Println("To activate the theme, please install gjs", err)
		return
	}
	// create the list of icon to get
	var icons []string
	for _, icon := range gnomeIconMap {
		icons = append(icons, icon)
	}
	iconSet := "[\n"
	for _, icon := range icons {
		iconSet += fmt.Sprintf(`    "%s",`+"\n", icon)
	}
	iconSet += "]"

	gjsIconList := fmt.Sprintf(gjsIconsScript, iconSet)

	// write the script to a temp file
	f, err := ioutil.TempFile("", "fyne-theme-gnome-*.js")
	if err != nil {
		log.Println(err)
		return
	}
	defer os.Remove(f.Name())

	// write the script to the temp file
	_, err = f.WriteString(gjsIconList)
	if err != nil {
		log.Println(err)
		return
	}

	// Call gjs with 2 version, 3 and 4 to complete the icon, this because
	// gtk version is sometimes not available or icon is not fully completed...
	// It's a bit tricky but it works.
	for _, gtkVersion := range []string{"3", "4"} {
		// run the script
		cmd := exec.Command(gjs,
			f.Name(), gtkVersion,
		)
		out, err := cmd.CombinedOutput()
		if err != nil {
			log.Println("gjs error:", err, string(out))
			return
		}

		tmpicons := map[string]*string{}
		// decode json to apply to the gnome theme
		err = json.Unmarshal(out, &tmpicons)
		if err != nil {
			log.Println(err)
			return
		}
		for k, v := range tmpicons {
			if _, ok := gnome.icons[k]; !ok {
				if v != nil && *v != "" {
					gnome.icons[k] = *v
				}
			}
		}
	}
}

// findThemeInformation decodes the theme from the gsettings and Gtk API.
func (gnome *GnomeTheme) findThemeInformation(gtkVersion int, variant fyne.ThemeVariant) {
	// make things faster in concurrent mode

	themename := gnome.getThemeName()
	if themename == "" {
		return
	}
	gnome.themeName = themename
	wg := sync.WaitGroup{}
	wg.Add(4)
	go gnome.applyColors(gtkVersion, &wg)
	go gnome.applyIcons(gtkVersion, &wg)
	go gnome.applyFont(&wg)
	go gnome.applyFontScale(&wg)
	wg.Wait()
}

// getGTKVersion gets the available GTK version for the given theme. If the version cannot be
// determine, it will return 3 wich is the most common used version.
func (gnome *GnomeTheme) getGTKVersion() (version int) {

	version = 3

	// ok so now, find if the theme is gtk4, either fallback to gtk3
	home, err := os.UserHomeDir()
	if err != nil {
		log.Println(err)
		return
	}

	possiblePaths := []string{
		home + "/.local/share/themes/",
		home + "/.themes/",
		`/usr/local/share/themes/`,
		`/usr/share/themes/`,
	}

	for _, path := range possiblePaths {
		path = filepath.Join(path, gnome.themeName)
		if _, err := os.Stat(path); err == nil {
			// now check if it is gtk4 compatible
			if _, err := os.Stat(path + "gtk-4.0/gtk.css"); err == nil {
				// it is gtk4
				version = 3
				return
			}
			if _, err := os.Stat(path + "gtk-3.0/gtk.css"); err == nil {
				version = 3
				return
			}
		}
	}
	return // default, but that may be a false positive now
}

// getThemeName gets the current theme name.
func (gnome *GnomeTheme) getThemeName() string {
	// call gsettings get org.gnome.desktop.interface gtk-theme
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "gtk-theme")
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		log.Println(string(out))
		return ""
	}
	themename := strings.TrimSpace(string(out))
	themename = strings.Trim(themename, "'")
	return themename
}

// loadIcon loads the icon from gnome theme, if the icon was already loaded, so the cached version is returned.
func (gnome *GnomeTheme) loadIcon(name string) (resource fyne.Resource) {
	var ok bool

	if resource, ok = gnome.iconCache[name]; ok {
		return
	}

	defer func() {
		// whatever the result is, cache it
		// even if it is nil
		gnome.iconCache[name] = resource
	}()

	if filename, ok := gnome.icons[name]; ok {
		content, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Println("Error while loading icon", err)
			return
		}
		if strings.HasSuffix(filename, ".svg") {
			// we need to ensure that the svg can be opened by Fyne
			buff := bytes.NewBuffer(content)
			_, err := oksvg.ReadIconStream(buff)
			if err != nil {
				// try to convert it to png with a converter
				log.Println("Cannot load file", filename, err, ", try to convert with tools")
				resource, err = convertSVGtoPNG(filename)
				if err != nil {
					log.Println("Cannot convert file", filename, ":", err)
				}
				return
			}
		}
		resource = fyne.NewStaticResource(filename, content)
		return
	}
	return
}

// parseColor converts a float32 array to color.Color.
func (*GnomeTheme) parseColor(col []float32) color.Color {
	return color.RGBA{
		R: uint8(col[0] * 255),
		G: uint8(col[1] * 255),
		B: uint8(col[2] * 255),
		A: uint8(col[3] * 255),
	}
}

// setFont sets the font for the theme - this method calls getFontPath() and converToTTF
// if needed.
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

// NewGnomeTheme returns a new Gnome theme based on the given gtk version. If gtkVersion is <= 0,
// the theme will try to determine the higher Gtk version available for the current GtkTheme.
func NewGnomeTheme(gtkVersion int) fyne.Theme {
	gnome := &GnomeTheme{
		fontSize:    theme.DefaultTheme().Size(theme.SizeNameText),
		iconCache:   map[string]fyne.Resource{},
		icons:       map[string]string{},
		colors:      map[fyne.ThemeColorName]color.Color{},
		font:        theme.DefaultTextFont(),
		scaleFactor: 1.0,
	}

	if gtkVersion <= 0 {
		// detect gtkVersion
		gtkVersion = gnome.getGTKVersion()
	}

	gnome.findThemeInformation(gtkVersion, theme.VariantDark)
	return gnome
}
