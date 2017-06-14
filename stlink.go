package stlink

import (
	"encoding/hex"
	"errors"

	"github.com/google/gousb"
)

// Stlink context structure
type Stlink struct {
	usbctx *gousb.Context
}

// New creates a new Stlink context
func New() (*Stlink, error) {
	s := &Stlink{
		usbctx: gousb.NewContext(),
	}
	return s, nil
}

// Close closes the Slink context
func (s *Stlink) Close() error {
	return s.usbctx.Close()
}

const (
	stVID        gousb.ID = 0x0483
	stlinkV1PID  gousb.ID = 0x3744
	stlinkV2PID  gousb.ID = 0x3748
	stlinkV21PID gousb.ID = 0x374b
)

func (s *Stlink) probeAll() ([]*gousb.Device, error) {
	return s.usbctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		if desc.Vendor == stVID && (desc.Product == stlinkV21PID || desc.Product == stlinkV2PID) {
			return true
		}
		return false
	})
}

// Probe searches for a list of devices and returns them.
// returned devices don't need to be closed but can't be used either.
// Use OpenDevice to open a device by SerialNumber
func (s *Stlink) Probe() ([]Device, error) {
	var devlist []Device

	devs, err := s.probeAll()
	if err != nil {
		return nil, err
	}
	defer func() {
		for _, d := range devs {
			d.Close()
		}
	}()
	for _, d := range devs {
		// 3 is the iSerialNumber, no constant in gousb for this
		sd, err := d.GetStringDescriptor(3)
		if err != nil {
			return nil, err
		}
		dev := Device{
			desc:         d.Desc,
			SerialNumber: sd,
			opened:       false,
		}
		// We check if the device ID is hex-encoded, otherwise do so
		if _, err := hex.DecodeString(sd); err != nil {
			dev.SerialNumber = hex.EncodeToString([]byte(sd))
		}
		devlist = append(devlist, dev)
	}
	return devlist, nil
}

// OpenDevice opens a device by serial number. Giving serial
// as an empty string, OpenDevice takes the first ST-link
// it can find
func (s *Stlink) OpenDevice(serial string) (*Device, error) {
	devs, err := s.probeAll()
	if err != nil {
		return nil, err
	}
	for _, d := range devs {
		// 3 is the iSerialNumber, no constant in gousb for this
		sd, err := d.GetStringDescriptor(3)
		if err != nil {
			return nil, err
		}
		dev := &Device{
			desc:         d.Desc,
			dev:          d,
			SerialNumber: sd,
			opened:       true,
		}
		// We check if the device ID is hex-encoded, otherwise do so
		if _, err := hex.DecodeString(sd); err != nil {
			dev.SerialNumber = hex.EncodeToString([]byte(sd))
		}
		if dev.SerialNumber == serial || serial == "" {
			if err := dev.init(); err != nil {
				dev.Close()
				return nil, err
			}
			return dev, nil
		}
		dev.Close()
	}
	return nil, errors.New("device not found")
}
