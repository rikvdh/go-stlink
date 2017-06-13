package main

import (
	"fmt"

	"github.com/Sirupsen/logrus"
	"github.com/rikvdh/go-stlink"
)

func main() {
	s, err := stlink.New()
	if err != nil {
		logrus.Fatalf("error getting Stlink context: %v\n", err)
	}

	devs, err := s.Probe()
	if err != nil {
		logrus.Fatalf("probe failed: %v\n", err)
	}
	for _, d := range devs {
		fmt.Printf("STlink: %s\n", d.SerialNumber)
	}
}
