package binding

import (
	"errors"
	"sync"

	"fyne.io/fyne/v2/data/binding"

	jsonvalue "github.com/Andrew-M-C/go.jsonvalue"
)

// JSONValue supports binding a JSON object
type JSONValue interface {
	binding.DataItem

	GetItemString(firstParam interface{}, params ...interface{}) (binding.String, error)
	GetItemFloat(firstParam interface{}, params ...interface{}) (binding.Float, error)
	GetItemInt(firstParam interface{}, params ...interface{}) (binding.Int, error)
	GetItemBool(firstParam interface{}, params ...interface{}) (binding.Bool, error)

	IsEmpty() bool
}

type databoundJSON struct {
	JSONValue

	self   binding.Untyped
	rwlock sync.RWMutex
	source binding.String
	last   string
	err    error
}

// Internal type for the children accessor
type childJSON struct {
	source *databoundJSON
	err    error

	target []interface{}
	first  interface{}
}

type childJSONString struct {
	binding.String
	generic childJSON
}

type childJSONFloat struct {
	binding.Float
	generic childJSON
}

type childJSONInt struct {
	binding.Int
	generic childJSON
}

type childJSONBool struct {
	binding.Bool
	generic childJSON
}

var (
	errWrongType = errors.New("wrong type provided")
)

// NewJSONFromString return a data binding to a JSON object synchronized with the `String`
// binding used to create the new binding. The JSON object is not exposed. You have to get children
// data binding targeting a specific path to actually get a data out of that JSON.
func NewJSONFromString(data binding.String) (JSONValue, error) {
	ret := &databoundJSON{self: binding.NewUntyped(), rwlock: sync.RWMutex{}, source: data, last: "{}"}
	data.AddListener(binding.NewDataListener(ret.changed))

	return ret, nil
}

// IsEmpty report true only if the data binding has already received a fully valid JSON object
func (json *databoundJSON) IsEmpty() bool {
	v, err := json.get()
	if err != nil {
		return true
	}

	if !v.IsObject() {
		return true
	}

	return v.Len() == 0
}

func (json *databoundJSON) lock() {
	json.rwlock.Lock()
}

func (json *databoundJSON) unlock() {
	json.rwlock.Unlock()
}

func (json *databoundJSON) rlock() {
	json.rwlock.RLock()
}

func (json *databoundJSON) runlock() {
	json.rwlock.RUnlock()
}

func (json *databoundJSON) get() (*jsonvalue.V, error) {
	if json.err != nil {
		return nil, json.err
	}

	i, err := json.self.Get()
	if err != nil {
		return nil, err
	}

	structured, ok := i.(*jsonvalue.V)
	if !ok {
		return nil, errWrongType
	}

	return structured, nil
}

func (json *databoundJSON) set(value *jsonvalue.V) error {
	s, err := value.MarshalString()
	if err != nil {
		return err
	}

	return json.source.Set(s)
}

func (json *databoundJSON) changed() {
	s, err := json.source.Get()
	if err != nil {
		json.err = err
		return
	}

	var structured *jsonvalue.V

	if s == "" {
		structured = &jsonvalue.V{}
	} else {
		structured, err = jsonvalue.UnmarshalString(s)
	}

	json.err = err
	if err != nil {
		structured = &jsonvalue.V{}
	}

	json.last = s
	json.self.Set(structured)
}

func (json *databoundJSON) AddListener(listener binding.DataListener) {
	json.self.AddListener(listener)
}

func (json *databoundJSON) RemoveListener(listener binding.DataListener) {
	json.self.RemoveListener(listener)
}

// Return a `String` binding linked with the specificed path to the JSON object provided by this data binding.
//
// The parameters follow the jsonvalue.GetString logic and only a String value can be fetched by this binding from
// the JSON object.
func (json *databoundJSON) GetItemString(firstParam interface{}, params ...interface{}) (binding.String, error) {
	ret := &childJSONString{String: binding.NewString(), generic: childJSON{source: json, target: make([]interface{}, 0)}}

	ret.generic.first = firstParam
	ret.generic.target = append(ret.generic.target, params...)

	json.AddListener(binding.NewDataListener(ret.changed))

	return ret, nil
}

func (child *childJSONString) changed() {
	child.generic.source.rlock()
	defer child.generic.source.runlock()

	structured, err := child.generic.source.get()
	child.generic.err = err
	if err != nil {
		return
	}

	var s string = ""

	if structured.IsObject() {
		s, err = structured.GetString(child.generic.first, child.generic.target...)
		child.generic.err = err
		if err != nil {
			return
		}
	}

	child.String.Set(s)
}

func (child *childJSONString) Get() (string, error) {
	if child.generic.err != nil {
		return "", child.generic.err
	}
	return child.String.Get()
}

func (child *childJSONString) Set(val string) error {
	child.generic.source.lock()
	defer child.generic.source.unlock()

	structured, err := child.generic.source.get()
	if err != nil {
		return err
	}

	structured.SetString(val).At(child.generic.target)
	return child.generic.source.set(structured)
}

// Return a `Float` binding linked with the specificed path to the JSON object provided by this data binding.
//
// The parameters follow the jsonvalue.GetFloat64 logic and only a Numeric value can be fetched by this binding
// from the JSON object.
func (json *databoundJSON) GetItemFloat(firstParam interface{}, params ...interface{}) (binding.Float, error) {
	ret := &childJSONFloat{Float: binding.NewFloat(), generic: childJSON{source: json, target: make([]interface{}, 0)}}

	ret.generic.first = firstParam
	ret.generic.target = append(ret.generic.target, params...)

	json.AddListener(binding.NewDataListener(ret.changed))

	return ret, nil
}

func (child *childJSONFloat) changed() {
	child.generic.source.rlock()
	defer child.generic.source.runlock()

	structured, err := child.generic.source.get()
	child.generic.err = err
	if err != nil {
		return
	}

	var f float64

	if structured.IsObject() {
		f, err = structured.GetFloat64(child.generic.first, child.generic.target...)
		child.generic.err = err
		if err != nil {
			return
		}
	}

	child.Float.Set(f)
}

func (child *childJSONFloat) Get() (float64, error) {
	if child.generic.err != nil {
		return 0, child.generic.err
	}
	return child.Float.Get()
}

func (child *childJSONFloat) Set(val float64) error {
	child.generic.source.lock()
	defer child.generic.source.unlock()

	structured, err := child.generic.source.get()
	if err != nil {
		return err
	}

	structured.SetFloat64(val).At(child.generic.target)
	return child.generic.source.set(structured)
}

// Return a `Int` binding linked with the specificed path to the JSON object provided by this data binding.
//
// The parameters follow the jsonvalue.GetInt logic and only a Numeric value can be fetched by this binding
// from the JSON object.
func (json *databoundJSON) GetItemInt(firstParam interface{}, params ...interface{}) (binding.Int, error) {
	ret := &childJSONInt{Int: binding.NewInt(), generic: childJSON{source: json, target: make([]interface{}, 0)}}

	ret.generic.first = firstParam
	ret.generic.target = append(ret.generic.target, params...)

	json.AddListener(binding.NewDataListener(ret.changed))

	return ret, nil
}

func (child *childJSONInt) changed() {
	child.generic.source.rlock()
	defer child.generic.source.runlock()

	structured, err := child.generic.source.get()
	child.generic.err = err
	if err != nil {
		return
	}

	var f int

	if structured.IsObject() {
		f, err = structured.GetInt(child.generic.first, child.generic.target...)
		child.generic.err = err
		if err != nil {
			return
		}
	}

	child.Int.Set(f)
}

func (child *childJSONInt) Get() (int, error) {
	if child.generic.err != nil {
		return 0, child.generic.err
	}
	return child.Int.Get()
}

func (child *childJSONInt) Set(val int) error {
	child.generic.source.lock()
	defer child.generic.source.unlock()

	structured, err := child.generic.source.get()
	if err != nil {
		return err
	}

	structured.SetInt(val).At(child.generic.target)
	return child.generic.source.set(structured)
}

// Return a `Bool` binding linked with the specificed path to the JSON object provided by this data binding.
//
// The parameters follow the jsonvalue.GetBool logic and only a boolean value can be fetched by this binding
// from the JSON object.
func (json *databoundJSON) GetItemBool(firstParam interface{}, params ...interface{}) (binding.Bool, error) {
	ret := &childJSONBool{Bool: binding.NewBool(), generic: childJSON{source: json, target: make([]interface{}, 0)}}

	ret.generic.first = firstParam
	ret.generic.target = append(ret.generic.target, params...)

	json.AddListener(binding.NewDataListener(ret.changed))

	return ret, nil
}

func (child *childJSONBool) changed() {
	child.generic.source.rlock()
	defer child.generic.source.runlock()

	structured, err := child.generic.source.get()
	child.generic.err = err
	if err != nil {
		return
	}

	var b bool

	if structured.IsObject() {
		b, err = structured.GetBool(child.generic.first, child.generic.target...)
		child.generic.err = err
		if err != nil {
			return
		}
	}

	child.Bool.Set(b)
}

func (child *childJSONBool) Get() (bool, error) {
	if child.generic.err != nil {
		return false, child.generic.err
	}
	return child.Bool.Get()
}

func (child *childJSONBool) Set(val bool) error {
	child.generic.source.lock()
	defer child.generic.source.unlock()

	structured, err := child.generic.source.get()
	if err != nil {
		return err
	}

	structured.SetBool(val).At(child.generic.target)
	return child.generic.source.set(structured)
}
