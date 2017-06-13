package stlink

import (
	"encoding/hex"
	"fmt"

	"github.com/google/gousb"
)

type Stlink struct {
	usbctx *gousb.Context
}

func New() (*Stlink, error) {
	s := &Stlink{
		usbctx: gousb.NewContext(),
	}
	return s, nil
}

func (s *Stlink) Close() error {
	return s.usbctx.Close()
}

const (
	stVID        gousb.ID = 0x0483
	stlinkV1PID  gousb.ID = 0x3744
	stlinkV2PID  gousb.ID = 0x3748
	stlinkV21PID gousb.ID = 0x374b
)

func (s *Stlink) Probe() ([]Device, error) {
	var devlist []Device
	//s.usbctx.Debug(999)
	devs, err := s.usbctx.OpenDevices(func(desc *gousb.DeviceDesc) bool {
		if desc.Vendor == stVID && (desc.Product == stlinkV21PID || desc.Product == stlinkV2PID) {
			return true
		}
		return false
	})
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
			fmt.Printf("error retrieving sd: %v\n", err)
		}
		dev := Device{
			desc:         d.Desc,
			SerialNumber: sd,
		}
		// We check if the device ID is hex-encoded, otherwise do so
		if _, err := hex.DecodeString(sd); err != nil {
			dev.SerialNumber = hex.EncodeToString([]byte(sd))
		}
		devlist = append(devlist, dev)
	}
	return devlist, nil
}
