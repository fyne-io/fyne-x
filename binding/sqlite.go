package binding

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
)

// assert compliance with the Preferences interface
var _ fyne.Preferences = &SQLiteDatastore{}

// assert compliance with the error interface
var _ error = &ErrSQLiteDatastoreNoSuchKey{}

// ErrSQLiteDatastoreNoSuchKey indicates that the requested key either is
// not present within the datastore, or that it corresponds to a value of
// a type other than the one requested.
type ErrSQLiteDatastoreNoSuchKey struct {
	key string
	typ string
}

func (e *ErrSQLiteDatastoreNoSuchKey) Error() string {
	return fmt.Sprintf("No such key '%s' for value of requested type '%s'", e.key, e.typ)
}

// IsErrSQLiteDatastoreNoSuchKey return True if the given error is a
// ErrSQLiteDatastoreNoSuchKey.
func IsErrSQLiteDatastoreNoSuchKey(e error) bool {
	_, ok := e.(*ErrSQLiteDatastoreNoSuchKey)
	return ok
}

// SQLiteDatastore allows for storing data in an SQLite 3 database, mimicking
// the Fyne preferences API (in fact, SQLiteDatastore is a valid implementation
// of the Preferences API, though that is not its intended use case), including
// it's data-binding capabilities.  This can be useful for creating
// application-specific data storage for data that is not suitable for use with
// the Fyne preferences API - for example data that is relatively large. SQLite
// offers a number of advantages - correct handling of cross-platform file
// locking, good performance, ACID compliance, and so on.
//
// Although the preferences API is implemented, meaning that preferences
// bindings may be used for primitive types, this is discouraged for
// performance reasons - because the preferences listener API is not very
// granular, each time a listener must be triggered, the entire set of bound
// values must be scanned. Though this is fine for a small number of values, it
// means that performance will degrade linearly on the number of bindings.
// SQLiteDatastore should offer roughly O(L) performance assuming an average of
// L bindings per key.
//
// This API assumes that it is the only API that is used to access the given
// table - it may create, destroy, or modify any table in the database file as
// it wishes.
//
// A separate table is created for each binding type, using one column for the
// key, and one for the value. Keys must be unique across all binding type, and
// may not be empty.
//
// If a key is written with a particular type, and then with a second,
// different type, the original value is deleted and the key is now of the
// second type.
//
// The path may be either a file on-disk, or ":memory:" for an in-memory
// database.
//
// NOTE: this API assumes that overwriting a key with the same value should
// still trigger event listeners. e.g. X.SetString("foo", "bar");
// X.SetString("foo", "bar") would cause a data binding for key "foo" to update
// twice.
//
// NOTE: concurrent access to the same database by multiple processes is not
// supported and will result in undefined behavior. However, concurrent access
// from multiple goroutines in the same process is supported.
type SQLiteDatastore struct {
	db           *sqlite3.Conn
	metaLock     sync.Mutex // sqlite handles its own locking, this is only for metadata
	keytypes     map[string]string
	listeners    []func()
	keyListeners map[string][]binding.DataListener
}

func NewSQLiteDatastore(path string) (*SQLiteDatastore, error) {
	schema := `
	CREATE TABLE IF NOT EXISTS kvp_string (
		key TEXT PRIMARY KEY,
		value TEXT
	);

	CREATE TABLE IF NOT EXISTS kvp_bool (
		key TEXT PRIMARY KEY,
		value BOOLEAN
	);

	CREATE TABLE IF NOT EXISTS kvp_float (
		key TEXT PRIMARY KEY,
		value REAL
	);

	CREATE TABLE IF NOT EXISTS kvp_int(
		key TEXT PRIMARY KEY,
		value INTEGER
	);

	CREATE TABLE IF NOT EXISTS kvp_map(
		key TEXT PRIMARY KEY,
		value INTEGER
	);

	PRAGMA optimize;
	PRAGMA auto_vacuum=FULL;
	`

	db, err := sqlite3.Open(path)
	if err != nil {
		return nil, err
	}

	err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	ds := &SQLiteDatastore{
		db:           db,
		keytypes:     make(map[string]string),
		keyListeners: make(map[string][]binding.DataListener),
		listeners:    make([]func(), 0),
	}

	// Ensure that the database is closed automatically when the datastore
	// gets garbage collected, otherwise we might leak memory on the C
	// side.
	runtime.SetFinalizer(ds, func(ds *SQLiteDatastore) { ds.db.Close() })

	return ds, nil
}

func (ds *SQLiteDatastore) addKeyListener(key string, listener binding.DataListener) {
	ds.metaLock.Lock()
	defer ds.metaLock.Unlock()

	_, ok := ds.keyListeners[key]
	if !ok {
		ds.keyListeners[key] = make([]binding.DataListener, 0)
	}

	ds.keyListeners[key] = append(ds.keyListeners[key], listener)
}

func (ds *SQLiteDatastore) removeKeyListener(key string, listener binding.DataListener) {
	ds.metaLock.Lock()
	defer ds.metaLock.Unlock()

	_, ok := ds.keyListeners[key]
	if !ok {
		return
	}

	newListeners := make([]binding.DataListener, 0)
	for _, l := range ds.keyListeners[key] {
		if l != listener {
			newListeners = append(newListeners, l)
		}
	}
	ds.keyListeners[key] = newListeners
}

func (ds *SQLiteDatastore) AddChangeListener(listener func()) {
	ds.metaLock.Lock()
	defer ds.metaLock.Unlock()

	ds.listeners = append(ds.listeners, listener)
}

func (ds *SQLiteDatastore) triggerListeners(key string) {
	ds.metaLock.Lock()
	defer ds.metaLock.Unlock()

	for _, l := range ds.listeners {
		l()
	}

	if key != "" {
		listeners, ok := ds.keyListeners[key]
		if !ok {
			return
		}

		for _, l := range listeners {
			l.DataChanged()
		}
	}
}

func (ds *SQLiteDatastore) set(key, typ string, value interface{}) error {
	defer ds.triggerListeners(key)

	query := fmt.Sprintf("INSERT INTO kvp_%s VALUES(?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value;", typ)

	var err error

	// Note that the type assertions here should never panic, because this
	// is a private method, and we're always calling it from a context
	// where we can guarantee the proper type.
	if typ == "string" {
		err = ds.db.Exec(query, key, value.(string))
	} else if typ == "bool" {
		err = ds.db.Exec(query, key, value.(bool))
	} else if typ == "float" {
		err = ds.db.Exec(query, key, value.(float64))
	} else if typ == "int" {
		err = ds.db.Exec(query, key, value.(int))
	} else {
		// This should never happen, since this is a private method.
		panic(typ)
	}

	if err != nil {
		return err
	}

	ds.metaLock.Lock()
	defer ds.metaLock.Unlock()
	ds.keytypes[key] = typ

	return nil
}

// Keys returns a list of all keys in the datastore irrespective of type.
func (ds *SQLiteDatastore) KeysTypes() ([]string, error) {
	keys, _, err := ds.KeysAndTypes()
	return keys, err
}

// KeysAndTypes returns a list of keys, as well as their corresponding data
// types.
//
// keys, types, err := X.KeysAndTypes()
func (ds *SQLiteDatastore) KeysAndTypes() ([]string, []string, error) {
	tables := []string{}
	query := "SELECT name FROM sqlite_master WHERE type='table';"
	stmt, err := ds.db.Prepare(query)
	if err != nil {
		return nil, nil, err
	}

	err = stmt.Exec()
	if err != nil {
		return nil, nil, err
	}

	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return nil, nil, err
		}
		if !hasRow {
			break
		}

		var name string
		err = stmt.Scan(&name)
		if err != nil {
			return nil, nil, err
		}

		if strings.HasPrefix(name, "kvp_") {
			tables = append(tables, name)
		}
	}

	keys := []string{}
	types := []string{}
	for _, table := range tables {
		query = fmt.Sprintf("SELECT key FROM %s;", table)
		stmt, err = ds.db.Prepare(query)
		if err != nil {
			return nil, nil, err
		}

		err = stmt.Exec()
		if err != nil {
			return nil, nil, err
		}

		for {
			hasRow, err := stmt.Step()
			if err != nil {
				return nil, nil, err
			}
			if !hasRow {
				break
			}

			var key string
			err = stmt.Scan(&key)
			if err != nil {
				return nil, nil, err
			}

			keys = append(keys, key)
			types = append(types, strings.TrimPrefix(table, "kvp_"))
		}

	}

	return keys, types, nil
}

// CheckedRemoveValue works identically to RemoveValue(), but returns an error
// if one was encountered.
func (ds *SQLiteDatastore) CheckedRemoveValue(key string) error {

	ds.metaLock.Lock()
	typ, ok := ds.keytypes[key]
	if !ok {
		return &ErrSQLiteDatastoreNoSuchKey{key: key, typ: "unknown"}
	}
	ds.metaLock.Unlock()

	ds.triggerListeners(key)

	query := fmt.Sprintf("DELETE FROM kvp_%s WHERE key=?;", typ)
	err := ds.db.Exec(query, key)
	if err != nil {
		return err
	}

	ds.metaLock.Lock()
	delete(ds.keytypes, key)
	ds.metaLock.Unlock()

	return nil
}

// RemoveValue implements fyne.Preferences
func (ds *SQLiteDatastore) RemoveValue(key string) {
	err := ds.CheckedRemoveValue(key)
	if err != nil {
		fyne.LogError(fmt.Sprintf("SQLiteDatastore: failed to remove key '%s'", key), err)
	}
}

func (ds *SQLiteDatastore) get(key, typ string) (interface{}, error) {
	query := fmt.Sprintf("SELECT value FROM kvp_%s WHERE key = ?;", typ)
	stmt, err := ds.db.Prepare(query)

	if err != nil {
		return nil, err
	}

	err = stmt.Exec(key)
	if err != nil {
		return nil, err
	}

	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return nil, err
		}
		if !hasRow {
			break
		}

		// NOTE: I know that sqlite3 can accept an interface{} type via
		// Scan(), however I'm breaking it out in this particular way
		// to base the final type on the desired one, rather than what
		// is stored in the database (ideally triggering an error if
		// the type in the database is the wrong one).

		if typ == "string" {
			var value string
			err = stmt.Scan(&value)
			if err != nil {
				return nil, err
			}
			return value, err

		} else if typ == "bool" {
			var value bool
			err = stmt.Scan(&value)
			if err != nil {
				return nil, err
			}
			return value, err

		} else if typ == "float" {
			var value float64
			err = stmt.Scan(&value)
			if err != nil {
				return nil, err
			}
			return value, err

		} else if typ == "int" {
			var value int
			err = stmt.Scan(&value)
			if err != nil {
				return nil, err
			}
			return value, err

		} else {
			// This should never happen, since get() is an internal
			// method.
			panic(typ)
		}

	}

	return nil, &ErrSQLiteDatastoreNoSuchKey{key: key, typ: typ}
}

// String returns the string associated with the given key, or the empty string
// if it does not exist or is not of type string.
//
// String implements the fyne.Preferences interface.
func (ds *SQLiteDatastore) String(key string) string {
	return ds.StringWithFallback(key, "")
}

// String implements the fyne.Preferences interface.
func (ds *SQLiteDatastore) StringWithFallback(key, fallback string) string {
	value, err := ds.CheckedString(key)
	if err != nil {
		if !IsErrSQLiteDatastoreNoSuchKey(err) {
			// no need to log errors if it's just a missing key
			fyne.LogError(fmt.Sprintf("SQLiteDatastore: failed to get string '%s'", key), err)
		}
		return fallback
	}
	return value
}

// CheckedString works similarly to String(), but also returns an error if
// one was encountered.
func (ds *SQLiteDatastore) CheckedString(key string) (string, error) {
	value, err := ds.get(key, "string")
	if err != nil {
		return "", err
	}
	return value.(string), nil
}

// SetString implements the fyne.Preferences interface.
func (ds *SQLiteDatastore) SetString(key, value string) {
	err := ds.CheckedSetString(key, value)
	if err != nil {
		fyne.LogError(fmt.Sprintf("SQLiteDatastore: failed to set string '%s' to '%s'", key, value), err)
	}
}

// CheckedSetString works identically to SetString(), but also returns an error
// if one occurred.
func (ds *SQLiteDatastore) CheckedSetString(key, value string) error {
	return ds.set(key, "string", value)
}

// Bool returns the bool associated with the given key, or false if it does not
// exist or is not of type bool.
//
// Bool implements the fyne.Preferences interface.
func (ds *SQLiteDatastore) Bool(key string) bool {
	return ds.BoolWithFallback(key, false)
}

// Bool implements the fyne.Preferences interface.
func (ds *SQLiteDatastore) BoolWithFallback(key string, fallback bool) bool {
	value, err := ds.CheckedBool(key)
	if err != nil {
		if !IsErrSQLiteDatastoreNoSuchKey(err) {
			// no need to log errors if it's just a missing key
			fyne.LogError(fmt.Sprintf("SQLiteDatastore: failed to get bool '%s'", key), err)
		}
		return fallback
	}
	return value
}

// CheckedBool works similarly to Bool(), but also returns an error if
// one was encountered.
func (ds *SQLiteDatastore) CheckedBool(key string) (bool, error) {
	value, err := ds.get(key, "bool")
	if err != nil {
		return false, err
	}
	return value.(bool), nil
}

// SetBool implements the fyne.Preferences interface.
func (ds *SQLiteDatastore) SetBool(key string, value bool) {
	err := ds.CheckedSetBool(key, value)
	if err != nil {
		fyne.LogError(fmt.Sprintf("SQLiteDatastore: failed to set bool '%s' to '%v'", key, value), err)
	}
}

// CheckedSetBool works identically to SetBool(), but also returns an error
// if one occurred.
func (ds *SQLiteDatastore) CheckedSetBool(key string, value bool) error {
	return ds.set(key, "bool", value)
}

// Float returns the float associated with the given key, or 0.0 if it does not
// exist or is not of type float.
//
// Float implements the fyne.Preferences interface.
func (ds *SQLiteDatastore) Float(key string) float64 {
	return ds.FloatWithFallback(key, 0.0)
}

// Float implements the fyne.Preferences interface.
func (ds *SQLiteDatastore) FloatWithFallback(key string, fallback float64) float64 {
	value, err := ds.CheckedFloat(key)
	if err != nil {
		if !IsErrSQLiteDatastoreNoSuchKey(err) {
			// no need to log errors if it's just a missing key
			fyne.LogError(fmt.Sprintf("SQLiteDatastore: failed to get float '%s'", key), err)
		}
		return fallback
	}
	return value
}

// CheckedFloat works similarly to Float(), but also returns an error if
// one was encountered.
func (ds *SQLiteDatastore) CheckedFloat(key string) (float64, error) {
	value, err := ds.get(key, "float")
	if err != nil {
		return 0.0, err
	}
	return value.(float64), nil
}

// SetFloat implements the fyne.Preferences interface.
func (ds *SQLiteDatastore) SetFloat(key string, value float64) {
	err := ds.CheckedSetFloat(key, value)
	if err != nil {
		fyne.LogError(fmt.Sprintf("SQLiteDatastore: failed to set float '%s' to '%f'", key, value), err)
	}
}

// CheckedSetFloat works identically to SetFloat(), but also returns an error
// if one occurred.
func (ds *SQLiteDatastore) CheckedSetFloat(key string, value float64) error {
	return ds.set(key, "float", value)
}

// Int returns the int associated with the given key, or 0 if it does not
// exist or is not of type int.
//
// Int implements the fyne.Preferences interface.
func (ds *SQLiteDatastore) Int(key string) int {
	return ds.IntWithFallback(key, 0)
}

// Int implements the fyne.Preferences interface.
func (ds *SQLiteDatastore) IntWithFallback(key string, fallback int) int {
	value, err := ds.CheckedInt(key)
	if err != nil {
		if !IsErrSQLiteDatastoreNoSuchKey(err) {
			// no need to log errors if it's just a missing key
			fyne.LogError(fmt.Sprintf("SQLiteDatastore: failed to get int '%s'", key), err)
		}
		return fallback
	}
	return value
}

// CheckedInt works similarly to Int(), but also returns an error if
// one was encountered.
func (ds *SQLiteDatastore) CheckedInt(key string) (int, error) {
	value, err := ds.get(key, "int")
	if err != nil {
		return 0, err
	}
	return value.(int), nil
}

// SetInt implements the fyne.Preferences interface.
func (ds *SQLiteDatastore) SetInt(key string, value int) {
	err := ds.CheckedSetInt(key, value)
	if err != nil {
		fyne.LogError(fmt.Sprintf("SQLiteDatastore: failed to get int '%s'", key), err)
	}
}

// CheckedSetInt works identically to SetInt(), but also returns an error
// if one occurred.
func (ds *SQLiteDatastore) CheckedSetInt(key string, value int) error {
	return ds.set(key, "int", value)
}

// Declare conformance with binding.Bool
var _ binding.Bool = (*sqlDsBoolBinding)(nil)

type sqlDsBoolBinding struct {
	*sqlDsBinding
}

func (ds *SQLiteDatastore) BindBool(key string) binding.Bool {
	bb := &sqlDsBoolBinding{sqlDsBinding: newSqlDsBinding(ds, key)}
	return bb
}

func (bb *sqlDsBoolBinding) Get() (bool, error) {
	return bb.ds.CheckedBool(bb.key)
}

func (bb *sqlDsBoolBinding) Set(value bool) error {
	return bb.ds.CheckedSetBool(bb.key, value)
}

// Declare conformance with binding.Float
var _ binding.Float = (*sqlDsFloatBinding)(nil)

type sqlDsFloatBinding struct {
	*sqlDsBinding
}

func (ds *SQLiteDatastore) BindFloat(key string) binding.Float {
	bb := &sqlDsFloatBinding{sqlDsBinding: newSqlDsBinding(ds, key)}
	return bb
}

func (bb *sqlDsFloatBinding) Get() (float64, error) {
	return bb.ds.CheckedFloat(bb.key)
}

func (bb *sqlDsFloatBinding) Set(value float64) error {
	return bb.ds.CheckedSetFloat(bb.key, value)
}

// Declare conformance with binding.Int
var _ binding.Int = (*sqlDsIntBinding)(nil)

type sqlDsIntBinding struct {
	*sqlDsBinding
}

func (ds *SQLiteDatastore) BindInt(key string) binding.Int {
	bb := &sqlDsIntBinding{sqlDsBinding: newSqlDsBinding(ds, key)}
	return bb
}

func (bb *sqlDsIntBinding) Get() (int, error) {
	return bb.ds.CheckedInt(bb.key)
}

func (bb *sqlDsIntBinding) Set(value int) error {
	return bb.ds.CheckedSetInt(bb.key, value)
}

// Declare conformance with binding.String
var _ binding.String = (*sqlDsStringBinding)(nil)

type sqlDsStringBinding struct {
	*sqlDsBinding
}

func (ds *SQLiteDatastore) BindString(key string) binding.String {
	bb := &sqlDsStringBinding{sqlDsBinding: newSqlDsBinding(ds, key)}
	return bb
}

func (bb *sqlDsStringBinding) Get() (string, error) {
	return bb.ds.CheckedString(bb.key)
}

func (bb *sqlDsStringBinding) Set(value string) error {
	return bb.ds.CheckedSetString(bb.key, value)
}

type sqlDsBinding struct {
	ds          *SQLiteDatastore
	key         string
	listeners   []binding.DataListener
	ownListener binding.DataListener
}

func newSqlDsBinding(ds *SQLiteDatastore, key string) *sqlDsBinding {
	b := &sqlDsBinding{
		ds:        ds,
		key:       key,
		listeners: make([]binding.DataListener, 0),
	}
	b.ownListener = binding.NewDataListener(b.triggerListeners)

	// Ensure that if this binding is garbage collected, we properly
	// remove its callback from the datastore.
	runtime.SetFinalizer(b, func(b *sqlDsBinding) {
		b.ds.removeKeyListener(b.key, b.ownListener)
	})

	ds.addKeyListener(key, b.ownListener)

	return b
}

func (b *sqlDsBinding) triggerListeners() {
	for _, l := range b.listeners {
		l.DataChanged()
	}
}

// AddListener implements binding.DataItem
func (b *sqlDsBinding) AddListener(listener binding.DataListener) {
	b.listeners = append(b.listeners, listener)
}

// RemoveListener implements binding.DataItem
func (b *sqlDsBinding) RemoveListener(listener binding.DataListener) {
	newListeners := make([]binding.DataListener, 0)
	for _, l := range b.listeners {
		if l != listener {
			newListeners = append(newListeners, l)
		}
	}
	b.listeners = newListeners
}

func (b *sqlDsBinding) get(key, typ string) (interface{}, error) {
	return b.ds.get(key, typ)
}

func (b *sqlDsBinding) set(key, typ string, value interface{}) error {
	return b.ds.set(key, typ, value)
}
