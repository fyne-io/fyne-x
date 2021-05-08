package validation_test

import (
	"testing"

	"fyne.io/x/fyne/data/validation"

	"github.com/stretchr/testify/assert"
)

func TestPassword(t *testing.T) {
	pw := validation.NewPassword(100)

	assert.NoError(t, pw("5 horses Ran around"))
	assert.Error(t, pw("bad-password"))

	pw = validation.NewPassword(150)

	assert.NoError(t, pw("7-BreaD-Crumbs.^_SpeciaL"))
	assert.Error(t, pw("12345--12345"))
}
