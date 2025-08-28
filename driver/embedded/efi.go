//go:build tamago

package fynex

import (
	"image"
	"io"
	"log"
	"time"

	"github.com/usbarmory/go-boot/uefi"
	"github.com/usbarmory/go-boot/uefi/x64"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/driver/embedded"
)

type efi struct {
	buf   []byte
	queue chan embedded.Event
}

func NewUEFIDriver() embedded.Driver {
	return &efi{queue: make(chan embedded.Event)}
}

func (u *efi) Queue() chan embedded.Event {
	return u.queue
}

func (u *efi) Run() {
	err := u.runLoop(fyne.CurrentApp(), u.queue)
	if err != nil {
		u.handleError(err)
	} else {
		_ = x64.UEFI.Runtime.ResetSystem(uefi.EfiResetShutdown)
	}
}

func (u *efi) handleError(err error) {
	log.Println("> ", err)
	time.Sleep(time.Second * 3)

	_ = x64.UEFI.Runtime.ResetSystem(uefi.EfiResetWarm)
}

func (u *efi) Render(img image.Image) {
	gop, _ := x64.UEFI.Boot.GetGraphicsOutput()
	i := 0

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pix := img.At(x, y)

			r, g, b, _ := pix.RGBA()
			u.buf[i+0] = byte(b)
			u.buf[i+1] = byte(g)
			u.buf[i+2] = byte(r)
			u.buf[i+3] = 0

			i += 4
		}
	}

	_ = gop.Blt(u.buf, uefi.EfiBltBufferToVideo, 0, 0, 0, 0, uint64(width), uint64(height), 0)
}

func (u *efi) runLoop(a fyne.App, queue chan embedded.Event) error {
	// we have to Run on a goroutine because the UEFI.Console is used on main in other code...
	wait := make(chan struct{})
	go func() {
		a.Run()
		wait <- struct{}{}
	}()

	defer func() {
		a.Quit()
		<-wait
	}()

	data := make([]byte, 4)
	for {
		n, err := x64.UEFI.Console.Read(data)
		if err == io.EOF {
			return err
		}
		if err != nil {
			fyne.LogError("failed to read", err)
			continue
		}

		for i := 0; i < n && i < len(data); i++ {
			switch data[i] {
			case 0:
				continue
			case 17, 23: // ctrl+Q or Esc
				return nil
			}

			triggerKey(uint16(data[i]), embedded.KeyPressed, queue)
			triggerKey(uint16(data[i]), embedded.KeyReleased, queue)
		}
	}
}

func (u *efi) ScreenSize() fyne.Size {
	gop, err := x64.UEFI.Boot.GetGraphicsOutput()
	if err != nil {
		u.handleError(err)
		return fyne.NewSize(0, 0)
	}

	mode, _ := gop.GetMode()
	info, _ := mode.GetInfo()
	ww, hh := uint64(info.HorizontalResolution),
		uint64(info.VerticalResolution)

	u.buf = make([]byte, ww*hh*4)
	return fyne.NewSize(float32(ww), float32(hh))
}
