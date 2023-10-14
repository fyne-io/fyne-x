package binding

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
)

// Since it's a callback we're testing, we let it process for some time before checking for result.
const callbackWait = 10 * time.Millisecond

func TestDisableableBinding(t *testing.T) {
	test.NewApp()

	bound := NewDisableableBinding()
	widget1 := widget.NewEntry()
	widget2 := widget.NewSelectEntry([]string{"test1", "test2"})

	if widget1.Disabled() || widget2.Disabled() {
		t.Errorf("test ended early because of wrong initial values on widgets.")
		return
	}

	// Adding widgets
	bound.AddWidgets(widget1, widget2)

	boundCheck := widget.NewCheckWithData("My bound checkbox", bound)

	// Checking not inverted
	boundCheck.SetChecked(true)

	// Since it's a callback we're testing, we let it process for some time before checking for result.
	<-time.After(callbackWait)
	if widget1.Disabled() || widget2.Disabled() {
		t.Errorf("Widget1 or 2 was not enabled.")
		return
	}

	boundCheck.SetChecked(false)
	<-time.After(callbackWait)
	if !widget1.Disabled() || !widget2.Disabled() {
		t.Errorf("Widget1 or 2 was not disabled.")
		return
	}

	// Checking inverted
	bound.SetInverted(true)

	boundCheck.SetChecked(true)
	<-time.After(callbackWait)
	if !widget1.Disabled() || !widget2.Disabled() {
		t.Errorf("Widget1 or 2 was not disabled.")
		return
	}

	boundCheck.SetChecked(false)
	<-time.After(callbackWait)
	if widget1.Disabled() || widget2.Disabled() {
		t.Errorf("Widget1 or 2 was not enabled.")
		return
	}
}
