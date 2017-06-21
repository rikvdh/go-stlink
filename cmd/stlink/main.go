package main

import (
	"flag"
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/rikvdh/go-stlink"
)

var (
	serial = flag.String("serial", "", "ST-link serial, probe when empty")
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
		for _, d := range devs {
			fmt.Printf("STlink: %s\n", d.SerialNumber)
			dv, err := s.OpenDevice(d.SerialNumber)
			if err == nil {
				v, err := dv.GetTargetVoltage()
				if err == nil {
					fmt.Printf(" voltage: %.3f\n", v)
				} else {
					fmt.Printf(" voltage: %v\n", err)
				}
				m, err := dv.GetMode()
				if err == nil {
					fmt.Printf(" mode:    %s\n", m)
				} else {
					fmt.Printf(" mode:    %v\n", err)
				}
				dv.Close()
			}
		}
	} else {
		dev, err := s.OpenDevice(*serial)
		if err != nil {
			logrus.Fatalf("error opening device: %v", err)
		}

		mode, err := dev.GetMode()
		if err != nil {
			logrus.Fatalf("error getting mode: %v", err)
		}
		logrus.Infof("ST-link mode is: %s", mode)

		dev.Close()
	}
}
