package widget

import (
	"reflect"
	"regexp"
	"strings"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type StructForm struct {
	widget.BaseWidget
	structType reflect.Type
	form       *widget.Form
	canvas     fyne.Canvas
	fields     []reflect.StructField
	widgets    []*widget.FormItem
	bindings   []binding.DataItem
	submit     func(s interface{})
}

func NewStructForm(s interface{}, submit func(s interface{}), cancel func()) *StructForm {
	sf := &StructForm{}
	sf.ExtendBaseWidget(sf)

	sf.structType = reflect.TypeOf(s)
	sf.fields = parseFields(sf.structType)
	sf.widgets = createWidgets(sf.fields)
	sf.submit = submit

	sf.form = &widget.Form{
		OnCancel: cancel,
		OnSubmit: sf.Submit,
		Items:    sf.widgets,
	}

	return sf
}

func buildValue(t reflect.Type, fields []reflect.StructField, bindings []binding.DataItem) reflect.Value {
	val := reflect.New(t)

	for i, f := range fields {
		var v interface{}
		var err error

		switch bindings[i].(type) {
		case binding.Int:
			v, err = bindings[i].(binding.Int).Get()
		case binding.String:
			v, err = bindings[i].(binding.String).Get()
		case binding.Float:
			v, err = bindings[i].(binding.Float).Get()
		}

		if err != nil {
			panic(err)
		}
		val.FieldByName(f.Name).Set(
			reflect.ValueOf(v),
		)
	}

	return val
}

func (sf *StructForm) Submit() {
	sf.submit(buildValue(
		sf.structType, sf.fields, sf.bindings,
	).Interface())
}

func parseFields(t reflect.Type) []reflect.StructField {
	fields := make([]reflect.StructField, 0)
	for i := 0; i < t.NumField(); i += 1 {
		f := t.Field(i)
		if strings.Contains(f.Tag.Get("fyne"), "field") {
			fields = append(fields, f)
		}
	}
	return fields
}

func createWidgets(fs []reflect.StructField) []*widget.FormItem {
	widgets := make([]*widget.FormItem, len(fs))
	binding := make([]*binding.DataItem, len(fs))
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	for i, f := range fs {
		w, b := fieldToWidget(f)
		widgets[i] = &widget.FormItem{
			Text:   matchFirstCap.ReplaceAllString(f.Name, "${1} ${2}"),
			Widget: w,
		}
		binding[i] = &b
	}
	return widgets
}

func fieldToWidget(f reflect.StructField) (fyne.CanvasObject, binding.DataItem) {
	switch f.Type.Kind() {
	case reflect.Float32, reflect.Float64:
		w := NewNumericalEntry()
		w.AllowFloat = true
		b := binding.NewString()
		w.Bind(b)
		return w, binding.StringToFloat(b)
	case reflect.Int:
		w := NewNumericalEntry()
		b := binding.NewString()
		w.Bind(b)
		return w, binding.StringToInt(b)
	case reflect.TypeOf(time.Time{}).Kind():
		b := binding.NewUntyped()
		w := NewCalendar(time.Now(), func(t time.Time) {
			b.Set(&t)
		})
		return w, b
	default:
		w := widget.NewEntry()
		b := binding.NewString()
		w.Bind(b)
		return w, b
	}
}

func (sf *StructForm) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(sf.form)
}
