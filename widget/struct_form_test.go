package widget

import (
	"reflect"
	"testing"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

type testFoo struct {
	Field1 string `fyne:"field"`
	Field2 int
	Field3 bool      `fyne:"field"`
	Field4 time.Time `fyne:"field"`
}

type testBar struct {
	Field1 float64 `fyne:"field"`
	Field2 time.Time
}

type testStruct struct {
	Name    string `fyne:"field"`
	Age     int    `fyne:"field"`
	Address string `fyne:"field"`
}

func StringBind(s string) binding.DataItem {
	return binding.BindString(&s).(binding.DataItem)
}

func IntBind(i int) binding.DataItem {
	return binding.BindInt(&i).(binding.DataItem)
}

func FloatBind(i float64) binding.DataItem {
	return binding.BindFloat(&i).(binding.DataItem)
}

func TestNewStructForm(t *testing.T) {
	type args struct {
		s      interface{}
		submit func(s interface{}, err error)
		cancel func()
	}
	type res struct {
		fields   []reflect.StructField
		widgets  []*widget.FormItem
		bindings []binding.DataItem
	}
	tests := []struct {
		name string
		args args
		want res
	}{
		{
			"testFoo",
			args{
				testFoo{},
				func(s interface{}, err error) {
					if err != nil {
						t.Error(err)
					}
					val, ok := s.(*testFoo)
					if !ok {
						t.Errorf("testFoo struct not created %v from %v", val, s)
					}
				},
				func() {},
			},
			res{
				fields: []reflect.StructField{
					{
						Name:   "Field1",
						Type:   reflect.TypeOf(""),
						Offset: 0,
						Index:  []int{0},
						Tag:    reflect.StructTag(`fyne:"field"`),
					},
					{
						Name:   "Field3",
						Type:   reflect.TypeOf(true),
						Offset: 24,
						Index:  []int{2},
						Tag:    reflect.StructTag(`fyne:"field"`),
					},
					{
						Name:   "Field4",
						Type:   reflect.TypeOf(time.Time{}),
						Offset: 32,
						Index:  []int{3},
						Tag:    reflect.StructTag(`fyne:"field"`),
					},
				},
				widgets: []*widget.FormItem{
					{
						Text:   "Field1",
						Widget: widget.NewEntry(),
					},
					{
						Text:   "Field3",
						Widget: widget.NewCheck("", func(b bool) {}),
					},
					{
						Text:   "Field4",
						Widget: NewCalendar(time.Now(), func(t time.Time) {}),
					},
				},
			},
		},
		{
			"testBar",
			args{
				testBar{},
				func(s interface{}, err error) {
					if err != nil {
						t.Error(err)
					}
					val, ok := s.(*testBar)
					if !ok {
						t.Errorf("testFoo struct not created %v from %v", val, s)
					}
				},
				func() {},
			},
			res{
				fields: []reflect.StructField{
					{
						Name:   "Field1",
						Type:   reflect.TypeOf(float64(1.0)),
						Offset: 0,
						Index:  []int{0},
						Tag:    reflect.StructTag(`fyne:"field"`),
					},
				},
				widgets: []*widget.FormItem{
					{
						Text:   "Field1",
						Widget: NewNumericalEntry(),
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewStructForm(tt.args.s, tt.args.submit, tt.args.cancel)
			if len(got.fields) != len(got.widgets) {
				t.Errorf("NewStructForm() len(fields) = %v != len(widgets) = %v", len(got.fields), len(got.widgets))
			}
			if len(got.fields) != len(got.bindings) {
				t.Errorf("NewStructForm() len(fields) = %v != len(bindings) = %v", len(got.fields), len(got.bindings))
			}

			if len(got.fields) != len(tt.want.fields) {
				t.Errorf("NewStructForm() len(fields) = %v, want %v", len(got.fields), len(tt.want.fields))
			}
			if !reflect.DeepEqual(got.fields, tt.want.fields) {
				t.Errorf("NewStructForm().fields got %v want %v", got.fields, tt.want.fields)
			}

			if len(got.widgets) != len(tt.want.widgets) {
				t.Errorf("NewStructForm().widgets got length %v, wanted %v", len(got.widgets), len(tt.want.widgets))
			}
			for i := 0; i < len(got.widgets) && i < len(tt.want.widgets); i += 1 {
				if reflect.TypeOf(got.widgets[i].Widget) != reflect.TypeOf(tt.want.widgets[i].Widget) {
					t.Errorf(
						"NewStructForm().widgets[%v] = %v, want %v", i,
						reflect.TypeOf(got.widgets[i].Widget), reflect.TypeOf(tt.want.widgets[i].Widget),
					)
				}
			}
			got.Submit()
		})
	}
}
