package dialog

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// NewAbout creates a parallax about dialog using the app metadata along with the
// markdown content and links passed into this method.
// You should call Show on the returned dialog to display it.
func NewAbout(content string, links []*widget.Hyperlink, a fyne.App, w fyne.Window) dialog.Dialog {
	d := dialog.NewCustom("About", "OK", aboutContent(content, links, a), w)
	d.Resize(fyne.NewSize(400, 360))

	return d
}

// NewAboutWindow creates a parallax about window using the app metadata along with the
// markdown content and links passed into this method.
// You should call Show on the returned window to display it.
func NewAboutWindow(content string, links []*widget.Hyperlink, a fyne.App) fyne.Window {
	w := a.NewWindow("About")
	w.SetContent(aboutContent(content, links, a))
	w.Resize(fyne.NewSize(360, 300))

	return w
}

// ShowAbout opens a parallax about dialog using the app metadata along with the
// markdown content and links passed into this method.
func ShowAbout(content string, links []*widget.Hyperlink, a fyne.App, w fyne.Window) {
	d := NewAbout(content, links, a, w)
	d.Show()
}

// ShowAboutWindow opens a parallax about window using the app metadata along with the
// markdown content and links passed into this method.
func ShowAboutWindow(content string, links []*widget.Hyperlink, a fyne.App) {
	w := NewAboutWindow(content, links, a)
	w.Show()
}

func aboutContent(content string, links []*widget.Hyperlink, a fyne.App) fyne.CanvasObject {
	rich := widget.NewRichTextFromMarkdown(content)
	footer := aboutFooter(links)

	logo := canvas.NewImageFromResource(a.Metadata().Icon)
	logo.FillMode = canvas.ImageFillContain
	logo.SetMinSize(fyne.NewSize(128, 128))

	appData := widget.NewRichTextFromMarkdown(
		"## " + a.Metadata().Name + "\n**Version:** " + a.Metadata().Version)
	centerText(appData)
	space := canvas.NewRectangle(color.Transparent)
	space.SetMinSize(fyne.NewSquareSize(theme.Padding() * 4))
	body := container.NewVBox(
		space,
		logo,
		appData,
		container.NewCenter(rich))
	scroll := container.NewScroll(body)

	bgColor := withAlpha(theme.Color(theme.ColorNameBackground), 0xe0)
	shadowColor := withAlpha(theme.Color(theme.ColorNameBackground), 0x33)

	underlay := canvas.NewImageFromResource(a.Metadata().Icon)
	bg := canvas.NewRectangle(bgColor)
	underlayer := underLayout{}
	slideBG := container.New(underlayer, underlay)
	footerBG := canvas.NewRectangle(shadowColor)
	watchTheme(bg, footerBG)

	underlay.Resize(fyne.NewSize(512, 512))
	scroll.OnScrolled = func(p fyne.Position) {
		underlayer.offset = -p.Y / 3
		underlayer.Layout(slideBG.Objects, slideBG.Size())
	}

	bgClip := container.NewScroll(slideBG)
	bgClip.Direction = container.ScrollNone
	return container.NewStack(container.New(unpad{top: true}, bgClip, bg),
		container.NewBorder(nil,
			container.NewStack(footerBG, footer), nil, nil,
			container.New(unpad{top: true, bottom: true}, scroll)))
}

func aboutFooter(links []*widget.Hyperlink) fyne.CanvasObject {
	footer := container.NewHBox(layout.NewSpacer())
	for i, a := range links {
		footer.Add(a)
		if i < len(links)-1 {
			footer.Add(widget.NewLabel("-"))
		}
	}
	footer.Add(layout.NewSpacer())

	return footer
}

func centerText(rich *widget.RichText) {
	for _, s := range rich.Segments {
		if text, ok := s.(*widget.TextSegment); ok {
			text.Style.Alignment = fyne.TextAlignCenter
		}
	}
}

func watchTheme(bg, footer *canvas.Rectangle) {
	fyne.CurrentApp().Settings().AddListener(func(_ fyne.Settings) {
		bgColor := withAlpha(theme.Color(theme.ColorNameBackground), 0xe0)
		bg.FillColor = bgColor
		bg.Refresh()

		shadowColor := withAlpha(theme.Color(theme.ColorNameBackground), 0x33)
		footer.FillColor = shadowColor
		footer.Refresh()
	})
}

func withAlpha(c color.Color, alpha uint8) color.Color {
	r, g, b, _ := c.RGBA()
	return color.NRGBA{R: uint8(r >> 8), G: uint8(g >> 8), B: uint8(b >> 8), A: alpha}
}

type underLayout struct {
	offset float32
}

func (u underLayout) Layout(objs []fyne.CanvasObject, size fyne.Size) {
	under := objs[0]
	left := size.Width/2 - under.Size().Width/2
	under.Move(fyne.NewPos(left, u.offset-50))
}

func (u underLayout) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.Size{}
}

type unpad struct {
	top, bottom bool
}

func (u unpad) Layout(objs []fyne.CanvasObject, s fyne.Size) {
	pad := theme.Padding()
	var pos fyne.Position
	if u.top {
		pos.Y = -pad
	}
	size := s
	if u.top {
		size.Height += pad
	}
	if u.bottom {
		size.Height += pad
	}
	for _, o := range objs {
		o.Move(pos)
		o.Resize(size)
	}
}

func (u unpad) MinSize(_ []fyne.CanvasObject) fyne.Size {
	return fyne.NewSize(100, 100)
}
