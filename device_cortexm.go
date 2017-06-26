package stlink

type CortexMPartNumber uint16

const (
	CortexMPartNumberUnknown CortexMPartNumber = 0x000
	CortexMPartNumberM0      CortexMPartNumber = 0xc20
	CortexMPartNumberM0Plus  CortexMPartNumber = 0xc60
	CortexMPartNumberM1      CortexMPartNumber = 0xc21
	CortexMPartNumberM3      CortexMPartNumber = 0xc23
	CortexMPartNumberM4      CortexMPartNumber = 0xc24
	CortexMPartNumberM7      CortexMPartNumber = 0xc27
)

func (c CortexMPartNumber) String() string {
	switch c {
	case CortexMPartNumberM0:
		return "ARM Cortex-M0"
	case CortexMPartNumberM0Plus:
		return "ARM Cortex-M0+"
	case CortexMPartNumberM1:
		return "ARM Cortex-M1"
	case CortexMPartNumberM3:
		return "ARM Cortex-M3"
	case CortexMPartNumberM4:
		return "ARM Cortex-M4"
	case CortexMPartNumberM7:
		return "ARM Cortex-M7"
	}
	return "unknown"
}

func (d *Device) CortexMPartNumber() (CortexMPartNumber, error) {
	if d.cpuID == 0 {
		_, err := d.CpuID()
		if err != nil {
			return CortexMPartNumberUnknown, err
		}
	}
	return CortexMPartNumber((d.cpuID >> 4) & 0xfff), nil
}
