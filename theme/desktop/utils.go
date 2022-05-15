package desktop

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// getFontPath will detect the font path from the font name taken from gsettings. As the font is not exactly
// the one that fc-match can find, we need to do some extra work to rebuild the name with style.
func getFontPath(fontname string) (string, error) {

	// This to transoform CamelCase to Camel-Case
	camelRegExp := regexp.MustCompile(`([a-z\-])([A-Z])`)

	// get all possible styles in fc-list
	allstyles := []string{}
	cmd := exec.Command("fc-list", "--format", "%{style}\n")
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	styles := strings.Split(string(out), "\n")
	for _, style := range styles {
		if style != "" {
			split := strings.Split(style, ",")
			for _, s := range split {
				allstyles = append(allstyles, s)
				// we also need to add a "-" for camel cases
				s = camelRegExp.ReplaceAllString(s, "$1-$2")
				allstyles = append(allstyles, s)
			}
		}
	}

	// Find the styles, remove it from the nmae, this make a correct fc-match query
	fontstyle := []string{}
	for _, style := range allstyles {
		// remove the style, we add a "space" to avoid this case: Semi-Condensed contains Condensed
		// and there is always a space before the style (because the font name prefixes the string)
		if strings.Contains(fontname, " "+style) {
			fontstyle = append(fontstyle, style)
			fontname = strings.ReplaceAll(fontname, style, "")
		}
	}

	// we can now search
	// fc-math ... "Font Name:Font Style
	var fontpath string
	cmd = exec.Command("fc-match", "-f", "%{file}", fontname+":"+strings.Join(fontstyle, " "))
	out, err = cmd.CombinedOutput()
	log.Println(string(out), fontname)
	if err != nil {
		log.Println(err)
		log.Println(string(out))
		return "", err
	}

	// get the font path with fc-list command
	fontpath = string(out)
	fontpath = strings.TrimSpace(fontpath)
	return fontpath, nil
}

// converToTTF will convert a font to a ttf file. This requires the fontconfig package.
func converToTTF(fontpath string) ([]byte, error) {

	// convert the font to a ttf file
	basename := filepath.Base(fontpath)
	tempTTF := filepath.Join(os.TempDir(), "fyne-"+basename+".ttf")

	// Convert to TTF
	ffScript := `Open("%s");Generate("%s")`
	script := fmt.Sprintf(ffScript, fontpath, tempTTF)
	cmd := exec.Command("fontforge", "-c", script)
	cmd.Env = append(cmd.Env, "FONTFORGE_LANGUAGE=ff")

	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println(err)
		log.Println(string(out))
		return nil, err
	}
	defer os.Remove(tempTTF)
	log.Println("TTF font generated: ", tempTTF)

	// read the temporary ttf file
	return ioutil.ReadFile(tempTTF)
}
