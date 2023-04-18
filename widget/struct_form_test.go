package widget

import (
	"reflect"
	"testing"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

func Test_parseFields(t *testing.T) {
	type Foo struct {
		Field1 string `fyne:"field"`
		Field2 int
		Field3 bool `fyne:"field"`
	}

	type Bar struct {
		Field1 float64 `fyne:"field"`
		Field2 time.Time
	}

	tests := []struct {
		name string
		args reflect.Type
		want []reflect.StructField
	}{
		{
			name: "parses fields with fyne tag",
			args: reflect.TypeOf(Foo{}),
			want: []reflect.StructField{
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
			},
		},
		{
			name: "returns empty slice when there are no fields with fyne tag",
			args: reflect.TypeOf(Bar{}),
			want: []reflect.StructField{
				{
					Name:   "Field1",
					Type:   reflect.TypeOf(float64(1.0)),
					Offset: 0,
					Index:  []int{0},
					Tag:    reflect.StructTag(`fyne:"field"`),
				}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseFields(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_createWidgets(t *testing.T) {
	type Foo struct {
		Field1 string `fyne:"field"`
		Field2 int    `fyne:"field"`
		Field3 bool
	}

	type Bar struct {
		Field1 float64 `fyne:"field"`
		Field2 time.Time
	}

	tests := []struct {
		name string
		args []reflect.StructField
		want []*widget.FormItem
	}{
		{
			name: "creates widgets for fields with fyne tag",
			args: []reflect.StructField{
				{
					Name: "Field1",
					Type: reflect.TypeOf(""),
					Tag:  reflect.StructTag(`fyne:"field"`),
				},
				{
					Name: "Field2",
					Type: reflect.TypeOf(0),
					Tag:  reflect.StructTag(`fyne:"field"`),
				},
			},
			want: []*widget.FormItem{
				{
					Text:   "Field1",
					Widget: widget.NewEntry(),
				},
				{
					Text:   "Field2",
					Widget: NewNumericalEntry(),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := createWidgets(tt.args)
			if len(got) != len(tt.want) {
				t.Errorf("got length %v, wanted %v", len(got), len(tt.want))
			}
			for i := 0; i < len(got) && i < len(tt.want); i += 1 {
				if reflect.TypeOf(got[i].Widget) != reflect.TypeOf(tt.want[i].Widget) {
					t.Errorf("createWidgets()[%v] = %v, want %v", i, reflect.TypeOf(got[i].Widget), reflect.TypeOf(tt.want[i].Widget))
				}
			}
		})
	}
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

func Test_buildValue(t *testing.T) {
	type args struct {
		t        reflect.Type
		fields   []reflect.StructField
		bindings []binding.DataItem
	}

	tests := []struct {
		name    string
		args    args
		want    reflect.Value
		wantErr bool
	}{
		{
			name: "Test Case 1 - Success",
			args: args{
				t:        reflect.TypeOf(testStruct{}),
				fields:   []reflect.StructField{{Name: "Name"}, {Name: "Age"}, {Name: "Address"}},
				bindings: []binding.DataItem{StringBind("John Doe"), IntBind(32), StringBind("1234 Main St.")},
			},
			want: reflect.ValueOf(&testStruct{Name: "John Doe", Age: 32, Address: "1234 Main St."}),
		},
		{
			name: "Test Case 2 - Invalid Binding",
			args: args{
				t:        reflect.TypeOf(testStruct{}),
				fields:   []reflect.StructField{{Name: "Name"}, {Name: "Age"}, {Name: "Address"}},
				bindings: []binding.DataItem{StringBind("John Doe"), StringBind("Invalid Age Value"), StringBind("1234 Main St.")},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := reflect.ValueOf(nil)
			err := false
			defer func() {
				if r := recover(); r != nil {
					err = true
				}
			}()
			got = buildValue(tt.args.t, tt.args.fields, tt.args.bindings)
			if err != tt.wantErr {
				t.Errorf("buildValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got.Interface(), tt.want.Interface()) {
				t.Errorf("buildValue() = %v, want %v", got, tt.want)
			}
		})
	}
}
