package binding

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"

	"github.com/bvinc/go-sqlite-lite/sqlite3"
)

// SQLiteDatastoreMagic is stored in the special "metadata" table and is
// used to verify that this database really came from SQLiteDatatore.
const SQLiteDatastoreMagic = "9e1f63f7-a6b1-4d50-88e8-269ccca04d89"

// SQLiteDatastoreVersion identifies the file format version of
// SQLiteDatastore. The version history is as follows:
//
// 1 - initial version
const SQLiteDatastoreVersion = 1

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

// NewSQLiteDatastore instances a new SQLiteDatastore object.
//
// The path may be either a file on-disk, or ":memory:" for an in-memory
// database.
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

	CREATE TABLE IF NOT EXISTS sqlitedatastore_metadata(
		key TEXT PRIMARY KEY,
		value TEXT
	);

	PRAGMA optimize;
	PRAGMA auto_vacuum=FULL;
	`

	db, err := sqlite3.Open(path)
	if err != nil {
		return nil, err
	}

	ds := &SQLiteDatastore{
		db:           db,
		keytypes:     make(map[string]string),
		keyListeners: make(map[string][]binding.DataListener),
		listeners:    make([]func(), 0),
	}

	// Check if this is a new database - if not, make sure it matches our
	// SQLiteDataStore version number.
	tables, err := ds.getTables()
	if err != nil {
		return nil, err
	}
	if len(tables) > 0 {
		magic, err := ds.metadataGet("SQLiteDatastoreMagic")
		if err != nil {
			return nil, fmt.Errorf("failed to get magic string for non-empty database due to error: %v", err)
		}

		if magic != SQLiteDatastoreMagic {
			return nil, fmt.Errorf("magic '%s' did not match expected '%s'", magic, SQLiteDatastoreMagic)
		}

	}

	err = db.Exec(schema)
	if err != nil {
		return nil, err
	}

	err = ds.metadataSet("SQLiteDatastoreVersion", fmt.Sprintf("%d", SQLiteDatastoreVersion))
	if err != nil {
		return nil, err
	}
	err = ds.metadataSet("SQLiteDatastoreMagic", SQLiteDatastoreMagic)
	if err != nil {
		return nil, err
	}

	// Ensure that the database is closed automatically when the datastore
	// gets garbage collected, otherwise we might leak memory on the C
	// side.
	runtime.SetFinalizer(ds, func(ds *SQLiteDatastore) { ds.db.Close() })

	return ds, nil
}

func (ds *SQLiteDatastore) addKeyListener(key string, listener binding.DataListener) {

	// Ensure the listener list for this key is initialized.
	ds.metaLock.Lock()
	_, ok := ds.keyListeners[key]
	if !ok {
		ds.keyListeners[key] = make([]binding.DataListener, 0)
	}

	ds.keyListeners[key] = append(ds.keyListeners[key], listener)
	ds.metaLock.Unlock()
}

func (ds *SQLiteDatastore) removeKeyListener(key string, listener binding.DataListener) {
	ds.metaLock.Lock()
	defer ds.metaLock.Unlock()

	// If there is no listener list for this key, there is nothing to do.
	_, ok := ds.keyListeners[key]
	if !ok {
		return
	}

	// Iteratively create a new listener list excluding the one to be
	// removed. This trades simplicity for performance; it may cause
	// slowdowns linear on the number of listeners.
	newListeners := make([]binding.DataListener, 0)
	for _, l := range ds.keyListeners[key] {
		if l != listener {
			newListeners = append(newListeners, l)
		}
	}

	ds.keyListeners[key] = newListeners
}

// AddChangeListener implements fyne.Preferences
func (ds *SQLiteDatastore) AddChangeListener(listener func()) {
	ds.metaLock.Lock()
	ds.listeners = append(ds.listeners, listener)
	ds.metaLock.Unlock()
}

func (ds *SQLiteDatastore) triggerListeners(key string) {
	ds.metaLock.Lock()
	defer ds.metaLock.Unlock()

	// Trigger all listeners for the preferences listeners.
	for _, l := range ds.listeners {
		l()
	}

	// If a key is provided, trigger the listeners registered for that key
	// (e.g. the ones associated with our own data binding types - this is
	// where the performance comes from vs just using vanilla BindPrefernce
	// methods).
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

// GetVersion gets the format version of the SQLiteDatastore file. This
// may be useful for detecting when upgrades need to be performed, if the
// format version is ever updated.
func (ds *SQLiteDatastore) GetVersion() (int, error) {
	versionString, err := ds.metadataGet("SQLiteDatastoreVersion")
	if err != nil {
		return -1, err
	}

	version, err := strconv.Atoi(versionString)
	if err != nil {
		return -1, err
	}

	return version, nil
}

// GetAppName returns the application name if one has been set, and otherwise
// an empty string an error.
func (ds *SQLiteDatastore) GetAppName() (string, error) {
	appName, err := ds.metadataGet("SQLiteDatastoreAppName")
	if err != nil {
		return "", err
	}

	return appName, nil
}

// SetAppName sets the application name overwriting any which already exists.
//
// This functionality can be useful to allow applications to determine if
// a particular SQLiteDatastore file originated from an instance of that
// application.
func (ds *SQLiteDatastore) SetAppName(name string) error {
	return ds.metadataSet("SQLiteDatastoreAppName", name)
}

// GetAppVersion returns the application name if one has been set, and otherwise
// an empty string an error.
func (ds *SQLiteDatastore) GetAppVersion() (int, error) {
	versionString, err := ds.metadataGet("SQLiteDatastoreAppVersion")
	if err != nil {
		return -1, err
	}

	version, err := strconv.Atoi(versionString)
	if err != nil {
		return -1, err
	}

	return version, nil
}

// Close closes the underlying database connection, and removes all listeners
// from the datastore.
func (ds *SQLiteDatastore) Close() error {
	err := ds.db.Close()
	ds.keytypes = nil
	ds.listeners = nil
	ds.keyListeners = nil
	return err
}

// SetAppVersion sets the application name overwriting any which already
// exists.
//
// This functionality can be useful to allow applications to track changes to
// its own file format version, for example to determine if data from an
// older version needs to be migrated to a newer one.
func (ds *SQLiteDatastore) SetAppVersion(version int) error {
	return ds.metadataSet("SQLiteDatastoreAppVersion", fmt.Sprintf("%d", version))
}

// metadataGet is sued to read from the internal metadata table.
func (ds *SQLiteDatastore) metadataGet(key string) (string, error) {
	query := "SELECT value FROM sqlitedatastore_metadata WHERE key = ?;"
	stmt, err := ds.db.Prepare(query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	err = stmt.Exec(key)
	if err != nil {
		return "", err
	}

	haveRow, err := stmt.Step()
	if err != nil {
		return "", err
	}
	if !haveRow {
		return "", fmt.Errorf("missing key '%s' from metadata table", key)
	}

	var value string
	err = stmt.Scan(&value)
	if err != nil {
		return "", err
	}
	return value, nil

}

func (ds *SQLiteDatastore) metadataSet(key, value string) error {
	query := "INSERT INTO sqlitedatastore_metadata VALUES(?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value;"

	// Note that the type assertions here should never panic, because this
	// is a private method, and we're always calling it from a context
	// where we can guarantee the proper type.
	err := ds.db.Exec(query, key, value)
	if err != nil {
		return err
	}

	return nil
}

// set inserts the specified value into the appropriate table per it's typ[e]
// parameter. It will type-assert the value to the correct type, intentionally
// causing a panic if this fails.
func (ds *SQLiteDatastore) set(key, typ string, value interface{}) error {
	defer ds.triggerListeners(key)

	// Because we delete the old value, and we want set() to be atomic,
	// we will do this within a transaction.
	err := ds.db.Begin()
	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO kvp_%s VALUES(?, ?) ON CONFLICT(key) DO UPDATE SET value=excluded.value;", typ)

	// Note that we need to remove the value before we set it, since it may
	// have changed type and thus need to be removed from its previous
	// table.
	err = ds.remove(key)
	if (err != nil) && (!IsErrSQLiteDatastoreNoSuchKey(err)) {
		// It's OK if we get a no such key error, but anything else
		// is an SQLite error that we need to propagate.
		return err
	}

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

	// Commit the DELETE and INSERT together.
	err = ds.db.Commit()
	if err != nil {
		return err
	}

	ds.metaLock.Lock()
	ds.keytypes[key] = typ
	ds.metaLock.Unlock()

	return nil
}

// Keys returns a list of all keys in the datastore irrespective of type.
func (ds *SQLiteDatastore) Keys() ([]string, error) {
	keys, _, err := ds.KeysAndTypes()
	return keys, err
}

// getTables returns a list of all SQLite tables in the database.
func (ds *SQLiteDatastore) getTables() ([]string, error) {
	tables := []string{}
	query := "SELECT name FROM sqlite_master WHERE type='table';"
	stmt, err := ds.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.Exec()
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

		var name string
		err = stmt.Scan(&name)
		if err != nil {
			return nil, err
		}

		tables = append(tables, name)
	}

	return tables, nil
}

// KeysAndTypes returns a list of keys, as well as their corresponding data
// types.
//
// keys, types, err := X.KeysAndTypes()
func (ds *SQLiteDatastore) KeysAndTypes() ([]string, []string, error) {
	// NOTE: for internal use, it is faster to use ds.keytypes, which is
	// maintained to have the correct state, assuming that the database
	// file isn't tampered with under us (this is why access from multiple
	// processes isn't supported, among other things).
	//
	// There is no good technical reason not to also use ds.keytypes
	// to construct the keys and types lists, other than paranoia.

	// Select all tables and keep any that have a prefix kvp_; this is not
	// robust against malformed, corrupt, or tampered-with databases, but
	// this is safe so long as we are the only ones who have touched it.
	tables, err := ds.getTables()
	if err != nil {
		return nil, nil, err
	}
	prefixedTables := []string{}
	for _, name := range tables {
		if strings.HasPrefix(name, "kvp_") {
			prefixedTables = append(prefixedTables, name)
		}
	}
	tables = prefixedTables

	// Scan each table, then scan each key from it. We can infer the type
	// by just slicing the kvp_ off the table name.
	keys := []string{}
	types := []string{}
	for _, table := range tables {
		query := fmt.Sprintf("SELECT key FROM %s;", table)
		stmt, err := ds.db.Prepare(query)
		if err != nil {
			return nil, nil, err
		}
		defer stmt.Close()

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
	err := ds.remove(key)
	if err != nil {
		return err
	}

	// Deletion counts as a change.
	ds.triggerListeners(key)

	return nil
}

// remove removes the specified key from the database, but does not trigger
// any data bindings. This function should only be used internally within
// functions that will correctly trigger listeners.
//
// This exists because set() needs to be able to remove any old values of the
// key in cases where the type has changed without triggering listeners
// multiple times.
func (ds *SQLiteDatastore) remove(key string) error {
	ds.metaLock.Lock()
	defer ds.metaLock.Unlock()

	// Determine the key type so we know what table to remove it from.
	typ, ok := ds.keytypes[key]
	if !ok {
		return &ErrSQLiteDatastoreNoSuchKey{key: key, typ: "unknown"}
	}

	query := fmt.Sprintf("DELETE FROM kvp_%s WHERE key=?;", typ)
	err := ds.db.Exec(query, key)
	if err != nil {
		return err
	}

	// NOTE: the key is intentionally not removed from the keyListeners
	// map, as it may be re-added in the future.
	delete(ds.keytypes, key)

	return nil
}

// RemoveValue implements fyne.Preferences
func (ds *SQLiteDatastore) RemoveValue(key string) {
	err := ds.CheckedRemoveValue(key)
	if err != nil {
		fyne.LogError(fmt.Sprintf("SQLiteDatastore: failed to remove key '%s'", key), err)
	}
}

// get returns the value of the specified key, from the table corresponding
// to the specified type. Values read from the database are intentionally
// read into variables of the correct type in order to trigger type
// errors if they end up with the wrong type on disk somehow.
func (ds *SQLiteDatastore) get(key, typ string) (interface{}, error) {
	query := fmt.Sprintf("SELECT value FROM kvp_%s WHERE key = ?;", typ)
	stmt, err := ds.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

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

// StringWithFallback implements the fyne.Preferences interface.
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

// BoolWithFallback implements the fyne.Preferences interface.
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

// FloatWithFallback implements the fyne.Preferences interface.
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

// IntWithFallback implements the fyne.Preferences interface.
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

// BindBool creates a new bool binding for a key in the datastore.
func (ds *SQLiteDatastore) BindBool(key string) binding.Bool {
	bb := &sqlDsBoolBinding{sqlDsBinding: newSQLDsBinding(ds, key)}
	return bb
}

// Get implements binding.Bool
func (bb *sqlDsBoolBinding) Get() (bool, error) {
	return bb.ds.CheckedBool(bb.key)
}

// Set implements binding.Bool
func (bb *sqlDsBoolBinding) Set(value bool) error {
	return bb.ds.CheckedSetBool(bb.key, value)
}

// Declare conformance with binding.Float
var _ binding.Float = (*sqlDsFloatBinding)(nil)

type sqlDsFloatBinding struct {
	*sqlDsBinding
}

// BindFloat creates a new float binding for a key in the datastore.
func (ds *SQLiteDatastore) BindFloat(key string) binding.Float {
	bb := &sqlDsFloatBinding{sqlDsBinding: newSQLDsBinding(ds, key)}
	return bb
}

// Get implements binding.Float
func (bb *sqlDsFloatBinding) Get() (float64, error) {
	return bb.ds.CheckedFloat(bb.key)
}

// Set implements binding.Float
func (bb *sqlDsFloatBinding) Set(value float64) error {
	return bb.ds.CheckedSetFloat(bb.key, value)
}

// Declare conformance with binding.Int
var _ binding.Int = (*sqlDsIntBinding)(nil)

type sqlDsIntBinding struct {
	*sqlDsBinding
}

// BindInt creates a new int binding for a key in the datastore.
func (ds *SQLiteDatastore) BindInt(key string) binding.Int {
	bb := &sqlDsIntBinding{sqlDsBinding: newSQLDsBinding(ds, key)}
	return bb
}

// Get implements binding.Int
func (bb *sqlDsIntBinding) Get() (int, error) {
	return bb.ds.CheckedInt(bb.key)
}

// Set implements binding.Int
func (bb *sqlDsIntBinding) Set(value int) error {
	return bb.ds.CheckedSetInt(bb.key, value)
}

// Declare conformance with binding.String
var _ binding.String = (*sqlDsStringBinding)(nil)

type sqlDsStringBinding struct {
	*sqlDsBinding
}

// BindString creates a new string binding for a key in the datastore.
func (ds *SQLiteDatastore) BindString(key string) binding.String {
	bb := &sqlDsStringBinding{sqlDsBinding: newSQLDsBinding(ds, key)}
	return bb
}

// Get implements binding.String
func (bb *sqlDsStringBinding) Get() (string, error) {
	return bb.ds.CheckedString(bb.key)
}

// Set implements binding.String
func (bb *sqlDsStringBinding) Set(value string) error {
	return bb.ds.CheckedSetString(bb.key, value)
}

// sqlDsBinding is the underlying plumbing below all of our other bindings.
// This structure simply gets embedded into the corresponding structures for
// typed bindings.
type sqlDsBinding struct {
	ds  *SQLiteDatastore
	key string

	// listeners stores our own list of data listeners
	listeners []binding.DataListener

	// ownListner is the data listener we give back to the DataStore;
	// this is a wrapper around .triggerListeners(), which will in turn
	// call our own listeners.
	ownListener binding.DataListener
}

func newSQLDsBinding(ds *SQLiteDatastore, key string) *sqlDsBinding {
	b := &sqlDsBinding{
		ds:        ds,
		key:       key,
		listeners: make([]binding.DataListener, 0),
	}
	b.ownListener = binding.NewDataListener(b.triggerListeners)

	// Ensure that if this binding is garbage collected, we properly
	// remove its callback from the datastore. This avoids leaking
	// memory in the datastore's key listener map.
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
