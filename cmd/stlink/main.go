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
				name, err := dv.GetName()
				if err == nil {
					fmt.Printf(" name:    %s\n", name)
				} else {
					fmt.Printf(" name:    %v\n", err)
				}
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
				ver, err := dv.GetVersion()
				if err == nil {
					fmt.Printf(" version: %s\n", ver)
				} else {
					fmt.Printf(" version: %v\n", err)
				}
				status, err := dv.GetStatus()
				if err == nil {
					fmt.Printf(" status:  %s\n", status)
				} else {
					fmt.Printf(" status:  %v\n", err)
				}
				cid, err := dv.GetCoreID()
				if err == nil {
					fmt.Printf(" coreid:  %08x\n", cid)
				} else {
					fmt.Printf(" coreid:  %v\n", err)
				}
				//print(string.format("cpuid: 0x%08x", dev:cpuid()))
				//print(string.format("chipid: 0x%08x", dev:chipid()))
				//print(string.format("devid: 0x%03x", dev:devid()))
				//print(string.format("flashSize: %d", dev:flashSize()))
				//print(string.format("r0: 0x%08x", dev:readReg()))
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
