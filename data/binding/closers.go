// Package binding provides extended sources of data binding.
package binding

import (
	"io"

	"fyne.io/fyne/v2/data/binding"
)

// StringCloser is an extension of the String interface that allows resources to be freed
// using the standard `Close()` method.
type StringCloser interface {
	binding.String
	io.Closer
}
