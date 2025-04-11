package binding_test

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/test"
	xbinding "fyne.io/x/fyne/data/binding"

	"github.com/stretchr/testify/assert"
)

func waitOnChan(t *testing.T, propagated chan bool) {
	select {
	case <-propagated:
	case <-time.After(1 * time.Second):
		assert.Fail(t, "The test should not have timedout")
	}
}

func NewListener(data binding.DataItem) chan bool {
	propagated := make(chan bool)
	listener := binding.NewDataListener(func() {
		go func() {
			propagated <- true
		}()
	})
	data.AddListener(listener)
	<-propagated // ignore init callback

	return propagated
}

func TestJSONFromStringWithString(t *testing.T) {
	_ = test.NewTempApp(t)
	s := binding.NewString()

	assert.NotNil(t, s)

	json, err := xbinding.NewJSONFromString(s)

	assert.NoError(t, err)
	assert.NotNil(t, json)

	propagatedJSON := NewListener(json)

	childV, err := json.GetItemString("value")

	assert.NoError(t, err)
	assert.NotNil(t, childV)

	propagatedCV := NewListener(childV)

	childA, err := json.GetItemString("array", 0)

	assert.NoError(t, err)
	assert.NotNil(t, childA)

	propagatedCA := NewListener(childA)

	intChildV := binding.StringToInt(childV)

	assert.NotNil(t, intChildV)

	propagatedICV := NewListener(intChildV)

	assert.Equal(t, true, json.IsEmpty())

	err = s.Set(`{ "value": "7", "array": [ "test" ] }`)

	assert.NoError(t, err)

	waitOnChan(t, propagatedJSON)
	waitOnChan(t, propagatedCV)
	waitOnChan(t, propagatedCA)
	waitOnChan(t, propagatedICV)

	assert.Equal(t, false, json.IsEmpty())

	vs, err := childV.Get()

	assert.NoError(t, err)
	assert.NotNil(t, vs)
	assert.Equal(t, "7", vs)

	vs, err = childA.Get()

	assert.NoError(t, err)
	assert.NotNil(t, vs)
	assert.Equal(t, "test", vs)

	vi, err := intChildV.Get()

	assert.NoError(t, err)
	assert.NotNil(t, vi)
	assert.Equal(t, 7, vi)
}

func TestJSONFromStringWithFloat(t *testing.T) {
	s := binding.NewString()

	assert.NotNil(t, s)

	json, err := xbinding.NewJSONFromString(s)

	assert.NoError(t, err)
	assert.NotNil(t, json)

	propagatedJSON := NewListener(json)

	childV, err := json.GetItemFloat("value")

	assert.NoError(t, err)
	assert.NotNil(t, childV)

	propagatedCV := NewListener(childV)

	childA, err := json.GetItemFloat("array", 0)

	assert.NoError(t, err)
	assert.NotNil(t, childA)

	propagatedCA := NewListener(childA)

	err = s.Set(`{ "value": 7, "array": [ 42.8 ] }`)

	assert.NoError(t, err)

	waitOnChan(t, propagatedJSON)
	waitOnChan(t, propagatedCV)
	waitOnChan(t, propagatedCA)

	vs, err := childV.Get()

	assert.NoError(t, err)
	assert.Equal(t, 7.0, vs)

	vs, err = childA.Get()

	assert.NoError(t, err)
	assert.Equal(t, 42.8, vs)
}

func TestJSONFromStringWithInt(t *testing.T) {
	s := binding.NewString()

	assert.NotNil(t, s)

	json, err := xbinding.NewJSONFromString(s)

	assert.NoError(t, err)
	assert.NotNil(t, json)

	propagatedJSON := NewListener(json)

	childV, err := json.GetItemInt("value")

	assert.NoError(t, err)
	assert.NotNil(t, childV)

	propagatedCV := NewListener(childV)

	childA, err := json.GetItemInt("array", 0)

	assert.NoError(t, err)
	assert.NotNil(t, childA)

	propagatedCA := NewListener(childA)

	err = s.Set(`{ "value": 7, "array": [ 42 ] }`)

	assert.NoError(t, err)

	waitOnChan(t, propagatedJSON)
	waitOnChan(t, propagatedCV)
	waitOnChan(t, propagatedCA)

	vs, err := childV.Get()

	assert.NoError(t, err)
	assert.Equal(t, 7, vs)

	vs, err = childA.Get()

	assert.NoError(t, err)
	assert.Equal(t, 42, vs)
}

func TestJSONFromStringWithBool(t *testing.T) {
	s := binding.NewString()

	assert.NotNil(t, s)

	json, err := xbinding.NewJSONFromString(s)

	assert.NoError(t, err)
	assert.NotNil(t, json)

	propagatedJSON := NewListener(json)

	childV, err := json.GetItemBool("value")

	assert.NoError(t, err)
	assert.NotNil(t, childV)

	propagatedCV := NewListener(childV)

	childA, err := json.GetItemBool("array", 0)

	assert.NoError(t, err)
	assert.NotNil(t, childA)

	propagatedCA := NewListener(childA)

	err = s.Set(`{ "value": true, "array": [ true ] }`)

	assert.NoError(t, err)

	waitOnChan(t, propagatedJSON)
	waitOnChan(t, propagatedCV)
	waitOnChan(t, propagatedCA)

	vs, err := childV.Get()

	assert.NoError(t, err)
	assert.Equal(t, true, vs)

	vs, err = childA.Get()

	assert.NoError(t, err)
	assert.Equal(t, true, vs)
}
