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
