package main

import (
	"flag"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/rikvdh/go-stlink"
)

var (
	serial = flag.String("serial", "", "ST-link serial, probe when empty")
	flash  = flag.Bool("f", false, "flash or no..")
	halt   = flag.Bool("h", false, "halt the core")
	run    = flag.Bool("r", false, "run")
	reset  = flag.Bool("re", false, "reset")
)

func main() {
	flag.Parse()
	s, err := stlink.New()
	if err != nil {
		logrus.Fatalf("error getting Stlink context: %v\n", err)
	}

	if *serial == "" {
		devs, err := s.Probe()
		if err != nil {
			logrus.Fatalf("probe failed: %v\n", err)
		}
		logrus.Infof("found %d devices", len(devs))
		for _, d := range devs {
			probeDevice(s, d.SerialNumber)
		}
	} else {
		if *flash {
			runFlash(s, *serial)
		} else if *halt {
			logrus.Infof("stlink: %s", *serial)
			dv, err := s.OpenDevice(*serial)
			if err != nil {
				panic(err)
			}
			defer dv.Close()

			if err := dv.ForceDebug(); err != nil {
				panic(err)
			}
			fmt.Printf("%s", dv)
		} else if *run {
			logrus.Infof("stlink: %s", *serial)
			dv, err := s.OpenDevice(*serial)
			if err != nil {
				panic(err)
			}
			defer dv.Close()
			panic(dv.Run())
		} else if *reset {
			logrus.Infof("stlink: %s", *serial)
			dv, err := s.OpenDevice(*serial)
			if err != nil {
				panic(err)
			}
			defer panic(dv.HardReset())
			dv.Close()
		} else {
			probeDevice(s, *serial)
		}
	}
}

func runFlash(s *stlink.Stlink, serial string) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.Debugf("stlink: %s", serial)
	dv, err := s.OpenDevice(serial)
	if err != nil {
		panic(err)
	}
	defer dv.Close()

	if err := dv.ForceDebug(); err != nil {
		panic(err)
	}

	if err := dv.EnterSWDMode(); err != nil {
		panic(err)
	}

	f, err := dv.GetFlashloader()
	if err != nil {
		panic(err)
	}
	if err := f.Unlock(); err != nil {
		panic(err)
	}
	if err := f.Lock(); err != nil {
		panic(err)
	}
}

func probeDevice(s *stlink.Stlink, serial string) {
	fmt.Printf("STlink: %s\n", serial)
	dv, err := s.OpenDevice(serial)
	if err != nil {
		return
	}
	fmt.Printf("%s", dv)
	dv.Close()
}
