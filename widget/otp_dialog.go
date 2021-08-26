package widget

import (
	"bytes"
	_ "embed"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/x/fyne/widget/loaders"
)

// OTPDialog represents an implementation of the GriddedInput in a dialog
// expecting a multi-factor authentication code.
type OTPDialog struct {
	dialog.Dialog

	dismissed bool
}

// ShowOTPDialog creates and shows an OTP dialog. Title wil be set to title and the dismiss button set to dismiss.
// If onComplete is not nil, it will be called when the grid is full. If it returns an error, the error will be
// shown and the dialog will remain shown. If it returns nil, the dialog will close. If onDismiss is not nil,
// it will be called if the dismiss button is pressed on the dialog.
func ShowOTPDialog(title, dismiss string, onComplete fyne.StringValidator, onDismiss func(), parent fyne.Window) {
	NewOTPDialog(title, dismiss, onComplete, onDismiss, parent).Show()
}

// NewOTPDialog creates a new OTP Dialog. Title wil be set to title and the dismiss button set to dismiss.
// If onComplete is not nil, it will be called when the grid is full. If it returns an error, the error will be
// shown and the dialog will remain shown. If it returns nil, the dialog will close. If onDismiss is not nil,
// it will be called if the dismiss button is pressed on the dialog.
func NewOTPDialog(title, dismiss string, onComplete fyne.StringValidator, onDismiss func(), parent fyne.Window) *OTPDialog {
	// Initialize a dialog instance
	otpDialog := &OTPDialog{dismissed: true}

	// Create components

	input := NewGriddedEntry(Digits, 6)
	input.Separator = canvas.NewText("-", theme.FocusColor())
	errField := NewWarningLabel()
	progress := loaders.NewGridLoader()
	progress.Hide()

	// Create the dialog

	otpDialog.Dialog = dialog.NewCustom(
		title, dismiss,
		container.New(otpDialog,
			canvas.NewImageFromReader(bytes.NewReader(MFAIcon), "mfa.png"),
			input,
			errField,
			progress,
		),
		parent,
	)

	// Set the callbacks

	otpDialog.SetOnClosed(func() {
		if otpDialog.dismissed && onDismiss != nil {
			onDismiss()
		}
	})
	input.OnCompletion = func(s string) {
		errField.Reset()

		input.Disable()
		defer input.Enable()

		progress.Show()
		defer progress.Hide()

		fyne.CurrentApp().Driver().CanvasForObject(input).Unfocus()

		if onComplete != nil {
			if err := onComplete(s); err != nil {
				errField.Warn(err.Error())
				fyne.CurrentApp().Driver().CanvasForObject(input).Focus(input.entries[input.selected])
				return
			}
		}

		otpDialog.dismissed = false
		otpDialog.Hide()
	}

	return otpDialog
}

func (o *OTPDialog) Layout(obj []fyne.CanvasObject, size fyne.Size) {
	// icon
	iconHeight := size.Height * .66
	obj[0].Resize(fyne.NewSize(size.Width/1.25, iconHeight))
	obj[0].Move(fyne.NewPos(size.Width/8-theme.Padding()*2, -theme.Padding()*4))

	// the input
	obj[1].Move(fyne.NewPos(size.Width/8, size.Height/2+theme.Padding()))
	obj[1].Resize(fyne.NewSize(size.Width*.8, size.Height/3))

	// errors
	obj[2].Move(obj[1].Position().Add(fyne.NewDelta(-obj[1].Size().Width/3+theme.Padding()*4, -theme.Padding()*2)))
	obj[2].Resize(fyne.NewSize(size.Width-theme.Padding(), size.Height-theme.Padding()))

	// progress
	obj[3].Resize(fyne.NewSize(size.Width/2, size.Height/2+theme.Padding()*6))
	obj[3].Move(fyne.NewPos(size.Width/2-theme.Padding()*6, size.Height/3))
}

func (o *OTPDialog) MinSize(obj []fyne.CanvasObject) fyne.Size {
	contentMin := obj[1].MinSize()

	width := contentMin.Width * 3
	height := contentMin.Height * 3

	return fyne.NewSize(width, height)
}

//go:embed mfa.png
var MFAIcon []byte
