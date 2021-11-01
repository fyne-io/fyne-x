package binding

import (
	"fmt"
	"testing"

	"fyne.io/fyne/v2/data/binding"
	"github.com/stretchr/testify/assert"
)

func TestSQLiteDatastore_String(t *testing.T) {
	ds, err := NewSQLiteDatastore(":memory:")
	assert.Nil(t, err)

	count := 0
	b := ds.BindString("foo")
	b.AddListener(binding.NewDataListener(func() { count++ }))

	err = ds.CheckedSetString("foo", "bar")
	assert.Nil(t, err)

	assert.Equal(t, ds.String("foo"), "bar")
	assert.Equal(t, ds.String("baz"), "")
	assert.Equal(t, ds.StringWithFallback("baz", "quux"), "quux")
	assert.Equal(t, ds.StringWithFallback("foo", "quux"), "bar")

	ds.SetString("foo", "bar")
	ds.SetString("foo", "baz")
	assert.Equal(t, 3, count)

}

func TestSQLiteDatastore_Bool(t *testing.T) {
	ds, err := NewSQLiteDatastore(":memory:")
	assert.Nil(t, err)

	count := 0
	b := ds.BindBool("foo")
	b.AddListener(binding.NewDataListener(func() { count++ }))

	err = ds.CheckedSetBool("foo", true)
	assert.Nil(t, err)

	assert.Equal(t, ds.Bool("foo"), true)
	assert.Equal(t, ds.Bool("baz"), false)
	assert.Equal(t, ds.BoolWithFallback("baz", true), true)
	assert.Equal(t, ds.BoolWithFallback("foo", false), true)

	assert.Equal(t, 1, count)
	err = b.Set(false)
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
	assert.Equal(t, false, ds.Bool("foo"))
}

func TestSQLiteDatastore_Float(t *testing.T) {
	ds, err := NewSQLiteDatastore(":memory:")
	assert.Nil(t, err)

	count := 0
	b := ds.BindFloat("foo")
	b.AddListener(binding.NewDataListener(func() { count++ }))

	err = ds.CheckedSetFloat("foo", 1.7)
	assert.Nil(t, err)

	assert.Equal(t, ds.Float("foo"), 1.7)
	assert.Equal(t, ds.Float("baz"), 0.0)
	assert.Equal(t, ds.FloatWithFallback("baz", 2.3), 2.3)
	assert.Equal(t, ds.FloatWithFallback("foo", 4.5), 1.7)

	assert.Equal(t, 1, count)
	err = b.Set(9.8)
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
	assert.Equal(t, 9.8, ds.Float("foo"))
}

func TestSQLiteDatastore_Int(t *testing.T) {
	ds, err := NewSQLiteDatastore(":memory:")
	assert.Nil(t, err)

	count := 0
	b := ds.BindInt("foo")
	b.AddListener(binding.NewDataListener(func() { count++ }))

	err = ds.CheckedSetInt("foo", 1)
	assert.Nil(t, err)

	assert.Equal(t, ds.Int("foo"), 1)
	assert.Equal(t, ds.Int("baz"), 0)
	assert.Equal(t, ds.IntWithFallback("baz", 2), 2)
	assert.Equal(t, ds.IntWithFallback("foo", 3), 1)
	assert.Equal(t, 1, count)

	err = b.Set(9)
	assert.Nil(t, err)
	assert.Equal(t, 2, count)
	assert.Equal(t, 9, ds.Int("foo"))
}

func TestSQLiteDatastore_KeysAndTypes(t *testing.T) {
	ds, err := NewSQLiteDatastore(":memory:")
	assert.Nil(t, err)

	err = ds.CheckedSetInt("foo", 1)
	assert.Nil(t, err)
	err = ds.CheckedSetFloat("bar", 3.7)
	assert.Nil(t, err)
	err = ds.CheckedSetBool("baz", true)
	assert.Nil(t, err)
	err = ds.CheckedSetString("quux", "spam")
	assert.Nil(t, err)

	keys, types, err := ds.KeysAndTypes()
	assert.Nil(t, err)
	assert.Contains(t, keys, "foo")
	assert.Contains(t, keys, "bar")
	assert.Contains(t, keys, "baz")
	assert.Contains(t, keys, "quux")
	assert.Equal(t, len(keys), 4)

	assert.Contains(t, types, "int")
	assert.Contains(t, types, "float")
	assert.Contains(t, types, "bool")
	assert.Contains(t, types, "string")
	assert.Equal(t, len(types), 4)
}

func TestSQLiteDatastore_WritePrecedence(t *testing.T) {

	ds, err := NewSQLiteDatastore(":memory:")
	assert.Nil(t, err)

	err = ds.CheckedSetInt("foo", 1)
	assert.Nil(t, err)

	err = ds.CheckedSetFloat("foo", 1.7)
	assert.Nil(t, err)

	err = ds.CheckedSetString("foo", "baz")
	assert.Nil(t, err)

	assert.Equal(t, ds.String("foo"), "baz")
}

func TestSQLiteDatastore_RemoveValue(t *testing.T) {
	ds, err := NewSQLiteDatastore(":memory:")
	assert.Nil(t, err)

	err = ds.CheckedSetInt("foo", 1)
	assert.Nil(t, err)
	err = ds.CheckedSetFloat("bar", 3.7)
	assert.Nil(t, err)
	err = ds.CheckedSetBool("baz", true)
	assert.Nil(t, err)
	err = ds.CheckedSetString("quux", "spam")
	assert.Nil(t, err)

	keys, types, err := ds.KeysAndTypes()
	assert.Nil(t, err)
	assert.Contains(t, keys, "foo")
	assert.Contains(t, keys, "bar")
	assert.Contains(t, keys, "baz")
	assert.Contains(t, keys, "quux")
	assert.Equal(t, len(keys), 4)

	assert.Contains(t, types, "int")
	assert.Contains(t, types, "float")
	assert.Contains(t, types, "bool")
	assert.Contains(t, types, "string")
	assert.Equal(t, len(types), 4)

	err = ds.CheckedRemoveValue("foo")
	assert.Nil(t, err)

	keys, types, err = ds.KeysAndTypes()
	assert.Nil(t, err)
	assert.NotContains(t, keys, "foo")
	assert.Contains(t, keys, "bar")
	assert.Contains(t, keys, "baz")
	assert.Contains(t, keys, "quux")
	assert.Equal(t, len(keys), 3)

	assert.NotContains(t, types, "int")
	assert.Contains(t, types, "float")
	assert.Contains(t, types, "bool")
	assert.Contains(t, types, "string")
	assert.Equal(t, len(types), 3)

	ds.RemoveValue("quux")

	keys, types, err = ds.KeysAndTypes()
	assert.Nil(t, err)
	assert.Contains(t, keys, "bar")
	assert.Contains(t, keys, "baz")
	assert.NotContains(t, keys, "quux")
	assert.Equal(t, len(keys), 2)

	assert.Contains(t, types, "float")
	assert.Contains(t, types, "bool")
	assert.NotContains(t, types, "string")
	assert.Equal(t, len(types), 2)
}

func TestSQLiteDatastore_Listeners(t *testing.T) {
	ds, err := NewSQLiteDatastore(":memory:")
	assert.Nil(t, err)

	count := 0

	ds.AddChangeListener(func() { count++ })

	ds.SetInt("foo", 1)
	ds.SetInt("bar", 2)
	ds.SetInt("baz", 2)

	assert.Equal(t, 3, count)
}

func TestSQLiteDatastore_binding(t *testing.T) {
	ds, err := NewSQLiteDatastore(":memory:")
	assert.Nil(t, err)

	count := 0

	fmt.Printf("----\n")
	b := newSqlDsBinding(ds, "foo")
	b.AddListener(binding.NewDataListener(func() { count++ }))

	ds.SetInt("foo", 7)
	ds.SetString("foo", "frob")
	ds.RemoveValue("foo")

	// Deleting the value should count.
	assert.Equal(t, 3, count)

	// Removing the value should not remove the corresponding listeners.
	ds.SetInt("foo", 7)
	ds.RemoveValue("foo")

	assert.Equal(t, 5, count)
}
