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
		logrus.Infof("found %d devices", len(devs))
		for _, d := range devs {
			probeDevice(s, d.SerialNumber)
		}
	} else {
		probeDevice(s, *serial)
	}
}

func probeDevice(s *stlink.Stlink, serial string) {
	fmt.Printf("STlink: %s\n", serial)
	dv, err := s.OpenDevice(serial)
	if err != nil {
		return
	}
	name, err := dv.Name()
	if err == nil {
		fmt.Printf(" name:    %s\n", name)
	} else {
		fmt.Printf(" name:    %v\n", err)
	}
	m, err := dv.Mode()
	if err == nil {
		fmt.Printf(" mode:    %s\n", m)
	} else {
		fmt.Printf(" mode:    %v\n", err)
	}
	ver, err := dv.Version()
	if err == nil {
		fmt.Printf(" version: %s\n", ver)
	} else {
		fmt.Printf(" version: %v\n", err)
	}
	v, err := dv.TargetVoltage()
	if err == nil {
		fmt.Printf(" voltage: %.3f\n", v)
	} else {
		fmt.Printf(" voltage: %v\n", err)
	}
	if v < 1 {
		fmt.Printf(" target voltage too low!\n")
		return
	}
	status, err := dv.Status()
	if err == nil {
		fmt.Printf(" status:  %s\n", status)
	} else {
		fmt.Printf(" status:  %v\n", err)
	}
	cid, err := dv.CoreID()
	if err == nil {
		fmt.Printf(" coreid:  %08x\n", cid)
	} else {
		fmt.Printf(" coreid:  %v\n", err)
	}
	cpu, err := dv.CpuID()
	if err == nil {
		fmt.Printf(" cpu:     %08x\n", cpu)
	} else {
		fmt.Printf(" cpu:     %v\n", err)
	}
	chip, err := dv.ChipID()
	if err == nil {
		fmt.Printf(" chip:    %08x\n", chip)
	} else {
		fmt.Printf(" chip:    %v\n", err)
	}
	dev, err := dv.DevID()
	if err == nil {
		fmt.Printf(" dev:     %03x\n", dev)
	} else {
		fmt.Printf(" dev:     %v\n", err)
	}
	pn, err := dv.CortexMPartNumber()
	if err == nil {
		fmt.Printf(" part-no: %s\n", pn)
	} else {
		fmt.Printf(" part-no: %v\n", err)
	}
	dv.Close()
}
