package binding

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"fyne.io/fyne/v2/data/binding"
	"github.com/stretchr/testify/assert"
)

func TestSQLiteDatastore_Version(t *testing.T) {
	ds, err := NewSQLiteDatastore(":memory:")
	assert.Nil(t, err)

	version, err := ds.GetVersion()
	assert.Nil(t, err)
	assert.Equal(t, SQLiteDatastoreVersion, version)
}

func TestSqliteDatastore_Metadata(t *testing.T) {
	// This test uses a temporary file, to ensure metadata correctly
	// persists across closing/opening the database.

	tempFile, err := ioutil.TempFile(os.TempDir(), "TestSqliteDatastore_Metadata")
	if err != nil {
		t.Fatal("failed to create temporary file", err)
	}

	defer os.Remove(tempFile.Name())

	ds, err := NewSQLiteDatastore(tempFile.Name())
	assert.Nil(t, err)

	err = ds.SetAppName("test app")
	assert.Nil(t, err)
	err = ds.SetAppVersion(123)
	assert.Nil(t, err)

	name, err := ds.GetAppName()
	assert.Nil(t, err)
	assert.Equal(t, "test app", name)

	version, err := ds.GetAppVersion()
	assert.Nil(t, err)
	assert.Equal(t, 123, version)

	ds.Close()

	ds2, err := NewSQLiteDatastore(tempFile.Name())
	assert.Nil(t, err)

	name, err = ds2.GetAppName()
	assert.Nil(t, err)
	assert.Equal(t, "test app", name)

	version, err = ds2.GetAppVersion()
	assert.Nil(t, err)
	assert.Equal(t, 123, version)
}

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
	t.Logf("keys=%v", keys)
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

	b := newSQLDsBinding(ds, "foo")
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

func TestSQLiteDatastore_TypeChange(t *testing.T) {
	// Make sure that changing the type of a key works correctly, and does
	// not lead to creating any duplicate keys.
	ds, err := NewSQLiteDatastore(":memory:")
	assert.Nil(t, err)

	ds.SetString("abc", "xyz")
	keys, types, err := ds.KeysAndTypes()
	t.Logf("1 keys=%v types=%v", keys, types)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(keys))
	assert.Equal(t, 1, len(types))

	ds.SetFloat("abc", 1.234)
	keys, types, err = ds.KeysAndTypes()
	t.Logf("2 keys=%v types=%v", keys, types)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(keys))
	assert.Equal(t, 1, len(types))

	ds.SetInt("abc", 5)
	keys, types, err = ds.KeysAndTypes()
	t.Logf("3 keys=%v types=%v", keys, types)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(keys))
	assert.Equal(t, 1, len(types))

	ds.SetBool("abc", false)
	keys, types, err = ds.KeysAndTypes()
	t.Logf("4 keys=%v types=%v", keys, types)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(keys))
	assert.Equal(t, 1, len(types))

	ds.SetString("abc", "xyz")
	keys, types, err = ds.KeysAndTypes()
	t.Logf("5 keys=%v types=%v", keys, types)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(keys))
	assert.Equal(t, 1, len(types))
}

func TestSQLiteDatastore_StringList(t *testing.T) {
	fmt.Printf("----\n")
	ds, err := NewSQLiteDatastore(":memory:")
	assert.Nil(t, err)

	err = ds.CheckedSetStringList("foo", []string{"one", "two", "three"})
	assert.Nil(t, err)

	list, err := ds.CheckedStringList("foo")
	assert.Nil(t, err)
	assert.Equal(t, []string{"one", "two", "three"}, list)

	list = ds.StringListWithFallback("foo", []string{"four"})
	assert.Nil(t, err)
	assert.Equal(t, []string{"one", "two", "three"}, list)

	list = ds.StringListWithFallback("bar", []string{"four"})
	assert.Nil(t, err)
	assert.Equal(t, []string{"four"}, list)

	list = ds.StringList("bar")
	assert.Equal(t, ([]string)(nil), list)

	sb := ds.BindStringList("foo")
	list, err = sb.Get()
	assert.Nil(t, err)
	assert.Equal(t, []string{"one", "two", "three"}, list)

	err = sb.SetValue(2, "frob")
	assert.Nil(t, err)
	list, err = sb.Get()
	assert.Nil(t, err)
	assert.Equal(t, []string{"one", "two", "frob"}, list)

	err = sb.SetValue(-1, "frob")
	assert.NotNil(t, err)

	err = sb.SetValue(3, "frob")
	assert.NotNil(t, err)

	two, err := sb.GetValue(1)
	assert.Nil(t, err)
	assert.Equal(t, "two", two)

	assert.Equal(t, 3, sb.Length())

	firstItem, err := sb.GetItem(0)
	assert.Nil(t, err)

	value, err := firstItem.(binding.String).Get()
	assert.Nil(t, err)
	assert.Equal(t, "one", value)

	count := 0
	firstItem.AddListener(binding.NewDataListener(func() { count++ }))

	err = firstItem.(binding.String).Set("spam")
	assert.Nil(t, err)
	assert.Equal(t, 1, count)

	err = ds.CheckedSetStringList("foo", []string{"eggs"})
	assert.Nil(t, err)
	assert.Equal(t, 2, count)

	err = sb.SetValue(0, "quux")
	assert.Nil(t, err)
	assert.Equal(t, 3, count)

	firstItem.RemoveListener(firstItem.(*sqlDsStringListItem).listeners[0])

	err = firstItem.(binding.String).Set("abc")
	assert.Nil(t, err)
	assert.Equal(t, 3, count)
}
