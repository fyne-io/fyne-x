package widget

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const outlineInset = float32(1.45)

type TracedButton struct {
	widget.Button

	s *canvas.Shader
}

func NewTracedButton(text string, tapped func()) fyne.CanvasObject {
	t := &TracedButton{}
	t.Text = text
	t.OnTapped = tapped
	t.ExtendBaseWidget(t)

	radius := theme.Current().Size(theme.SizeNameButtonRadius)
	shader := canvas.NewShader("roundedButtonTrace", []byte(traceShaderSource), []byte(traceShaderSourceES))
	uniforms := map[string]float32{
		"inset":  outlineInset,
		"radius": radius,
		"stroke": 1.15,
		"speed":  0.34,
	}
	putShaderColor(uniforms, "primary", theme.ColorForWidget(theme.ColorNamePrimary, t))
	putShaderColor(uniforms, "button", theme.ColorForWidget(theme.ColorNameButton, t))
	shader.Uniforms = uniforms

	anim := canvas.NewShaderAnimation(shader)
	anim.Start()

	t.s = shader
	return t
}

func (t *TracedButton) CreateRenderer() fyne.WidgetRenderer {
	t.ExtendBaseWidget(t)
	br := t.Button.CreateRenderer()
	r := &tracedButtonRender{WidgetRenderer: br}

	objs := br.Objects()
	r.objs = append(objs, t.s)
	return r
}

func (t *TracedButton) Refresh() {
	uniforms := t.s.Uniforms
	putShaderColor(uniforms, "primary", theme.ColorForWidget(theme.ColorNamePrimary, t))
	putShaderColor(uniforms, "button", theme.ColorForWidget(theme.ColorNameButton, t))
	t.s.Uniforms = uniforms

	// TODO apply the shader color / radius
	t.Button.Refresh()
}

type tracedButtonRender struct {
	fyne.WidgetRenderer

	objs []fyne.CanvasObject
}

func (r *tracedButtonRender) Objects() []fyne.CanvasObject {
	return r.objs
}

func (r *tracedButtonRender) Layout(size fyne.Size) {
	r.WidgetRenderer.Layout(size)
	r.objs[len(r.objs)-1].Resize(size)
}

func putShaderColor(uniforms map[string]float32, name string, c color.Color) {
	r, g, b, _ := c.RGBA()
	uniforms[name+"R"] = float32(r) / 65535
	uniforms[name+"G"] = float32(g) / 65535
	uniforms[name+"B"] = float32(b) / 65535
}
