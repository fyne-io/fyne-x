package widget

import (
	"errors"
	"fmt"
	"image"
	"image/png"
	"net/http"
)

var tileMap = make(map[string]image.Image)

func getTile(tileSource string, x, y, zoom int, cl *http.Client) (image.Image, error) {
	if tileSource == "" {
		return nil, errors.New("no tileSource provided")
	}

	u := fmt.Sprintf(tileSource, zoom, x, y)
	if tile, ok := tileMap[u]; ok {
		return tile, nil
	}

	req, err := http.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "Fyne-X Map Widget/0.1")
	res, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	img, err := png.Decode(res.Body)
	if err == nil {
		tileMap[u] = img
	}
	return img, err
}
