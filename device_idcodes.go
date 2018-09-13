package stlink

import (
	"encoding/binary"
	"errors"
)

func (d *Device) CoreID() (uint32, error) {
	tx := make([]byte, cmdSize, cmdSize)
	tx[0] = byte(stlinkCmdDebug)
	tx[1] = byte(stlinkCmdDebugReadCoreid)
	err := d.write(tx)
	if err != nil {
		return 0, err
	}

	var v uint32
	if err := binary.Read(d.inEp, binary.LittleEndian, &v); err != nil {
		return 0, err
	}
	return v, nil
}

const (
	cortexMIDCodeAddress   uint32 = 0xE0042000
	cortexMIDCodeM0Address uint32 = 0x40015800
)

func (d *Device) ChipID() (uint32, error) {
	pn, err := d.CortexMPartNumber()
	if err != nil {
		return 0, err
	}
	if pn == CortexMPartNumberM0 || pn == CortexMPartNumberM0Plus {
		return d.Read32(cortexMIDCodeM0Address)
	}
	return d.Read32(cortexMIDCodeAddress)
}

const cortexMCpuIDRegisterAddress uint32 = 0xE000ED00

func (d *Device) CpuID() (uint32, error) {
	id, err := d.Read32(cortexMCpuIDRegisterAddress)
	if err == nil {
		d.cpuID = id
	}
	return id, err
}

type ChipFamily uint16

const (
	ChipFamilySTM32F0             ChipFamily = 0x440
	ChipFamilySTM32F0Can          ChipFamily = 0x448
	ChipFamilySTM32F0Small        ChipFamily = 0x444
	ChipFamilySTM32F04            ChipFamily = 0x445
	ChipFamilySTM32F09X           ChipFamily = 0x442
	ChipFamilySTM32F1Connectivity ChipFamily = 0x418
	ChipFamilySTM32F1High         ChipFamily = 0x414
	ChipFamilySTM32F1Low          ChipFamily = 0x412
	ChipFamilySTM32F1Medium       ChipFamily = 0x410
	ChipFamilySTM32F1VLHigh       ChipFamily = 0x428
	ChipFamilySTM32F1VLMedium     ChipFamily = 0x420
	ChipFamilySTM32F1XL           ChipFamily = 0x430
	ChipFamilySTM32F2             ChipFamily = 0x411
	ChipFamilySTM32F3             ChipFamily = 0x422
	ChipFamilySTM32F3Small        ChipFamily = 0x439
	ChipFamilySTM32F303High       ChipFamily = 0x446
	ChipFamilySTM32F334           ChipFamily = 0x438
	ChipFamilySTM32F37x           ChipFamily = 0x432
	ChipFamilySTM32F4             ChipFamily = 0x413
	ChipFamilySTM32F4DE           ChipFamily = 0x433
	ChipFamilySTM32F4DSI          ChipFamily = 0x434
	ChipFamilySTM32F4HD           ChipFamily = 0x419
	ChipFamilySTM32F4LP           ChipFamily = 0x423
	ChipFamilySTM32F410           ChipFamily = 0x458
	ChipFamilySTM32F411RE         ChipFamily = 0x431
	ChipFamilySTM32F412           ChipFamily = 0x441
	ChipFamilySTM32F413           ChipFamily = 0x463
	ChipFamilySTM32F446           ChipFamily = 0x421
	ChipFamilySTM32F7             ChipFamily = 0x449
	ChipFamilySTM32F7Advanced     ChipFamily = 0x451
	ChipFamilySTM32F7Foundation   ChipFamily = 0x452
	ChipFamilySTM32L0             ChipFamily = 0x417
	ChipFamilySTM32L0Cat2         ChipFamily = 0x425
	ChipFamilySTM32L0Cat5         ChipFamily = 0x447
	ChipFamilySTM32L011           ChipFamily = 0x457
	ChipFamilySTM32L1Cat2         ChipFamily = 0x429
	ChipFamilySTM32L1High         ChipFamily = 0x436
	ChipFamilySTM32L1MediumLow    ChipFamily = 0x416
	ChipFamilySTM32L1MediumHigh   ChipFamily = 0x427
	ChipFamilySTM32L152RE         ChipFamily = 0x437
	ChipFamilySTM32L4             ChipFamily = 0x415
	ChipFamilySTM32L434X          ChipFamily = 0x435
	ChipFamilySTM32L4X6           ChipFamily = 0x461
)

type ChipFamilyGroup string

const (
	ChipFamilyGroupSTM32L0 ChipFamilyGroup = "stm32l0"
	ChipFamilyGroupSTM32L1 ChipFamilyGroup = "stm32l1"
	ChipFamilyGroupSTM32L4 ChipFamilyGroup = "stm32l4"
	ChipFamilyGroupSTM32F0 ChipFamilyGroup = "stm32f0"
	ChipFamilyGroupSTM32F1 ChipFamilyGroup = "stm32f1"
	ChipFamilyGroupSTM32F2 ChipFamilyGroup = "stm32f2"
	ChipFamilyGroupSTM32F3 ChipFamilyGroup = "stm32f3"
	ChipFamilyGroupSTM32F4 ChipFamilyGroup = "stm32f4"
	ChipFamilyGroupSTM32F7 ChipFamilyGroup = "stm32f7"
)

// GroupName gets the chip family group name (e.g "stm32f0", "stm32l0" ...)
func (cf ChipFamily) Group() ChipFamilyGroup {
	switch cf {
	case ChipFamilySTM32F0, ChipFamilySTM32F0Can, ChipFamilySTM32F0Small,
		ChipFamilySTM32F04, ChipFamilySTM32F09X:
		return ChipFamilyGroupSTM32F0
	case ChipFamilySTM32F1Connectivity, ChipFamilySTM32F1High, ChipFamilySTM32F1Low,
	ChipFamilySTM32F1Medium, ChipFamilySTM32F1VLHigh, ChipFamilySTM32F1VLMedium, ChipFamilySTM32F1XL:
		return ChipFamilyGroupSTM32F1
	case ChipFamilySTM32F2:
		return ChipFamilyGroupSTM32F2
	case ChipFamilySTM32F3, ChipFamilySTM32F3Small, ChipFamilySTM32F303High,
		ChipFamilySTM32F334, ChipFamilySTM32F37x:
		return ChipFamilyGroupSTM32F3
	case ChipFamilySTM32F4, ChipFamilySTM32F4DE, ChipFamilySTM32F4DSI, ChipFamilySTM32F4HD,
		ChipFamilySTM32F4LP, ChipFamilySTM32F410, ChipFamilySTM32F411RE, ChipFamilySTM32F412,
		ChipFamilySTM32F413, ChipFamilySTM32F446:
		return ChipFamilyGroupSTM32F4
	case ChipFamilySTM32F7, ChipFamilySTM32F7Advanced, ChipFamilySTM32F7Foundation:
		return ChipFamilyGroupSTM32F7
	case ChipFamilySTM32L0, ChipFamilySTM32L0Cat2, ChipFamilySTM32L0Cat5, ChipFamilySTM32L011:
		return ChipFamilyGroupSTM32L0
	case ChipFamilySTM32L1Cat2, ChipFamilySTM32L1High, ChipFamilySTM32L1MediumLow,
		ChipFamilySTM32L1MediumHigh, ChipFamilySTM32L152RE:
		return ChipFamilyGroupSTM32L1
	case ChipFamilySTM32L4, ChipFamilySTM32L434X, ChipFamilySTM32L4X6:
		return ChipFamilyGroupSTM32L4
	}
	return "unknown"
}

func (d *Device) DevID() (ChipFamily, error) {
	id, err := d.ChipID()
	return ChipFamily(id & 0xfff), err
}

func (d *Device) FlashSize() (uint16, error) {
	pn, err := d.DevID()
	if err != nil {
		return 0, err
	}
	switch pn {
	case ChipFamilySTM32F0, ChipFamilySTM32F09X, ChipFamilySTM32F0Can,
		ChipFamilySTM32F0Small, ChipFamilySTM32F04, ChipFamilySTM32F3,
		ChipFamilySTM32F37x, ChipFamilySTM32F334, ChipFamilySTM32F3Small,
		ChipFamilySTM32F303High:
		return d.Read16(0x1ffff7cc)

	case ChipFamilySTM32F1Medium, ChipFamilySTM32F1High, ChipFamilySTM32F1Low,
		ChipFamilySTM32F1Connectivity, ChipFamilySTM32F1VLMedium, ChipFamilySTM32F1VLHigh,
		ChipFamilySTM32F1XL:
		return d.Read16(0x1ffff7e0)

	case ChipFamilySTM32F2, ChipFamilySTM32F4, ChipFamilySTM32F4HD,
		ChipFamilySTM32F446, ChipFamilySTM32F411RE, ChipFamilySTM32F4DE,
		ChipFamilySTM32F4DSI, ChipFamilySTM32F412, ChipFamilySTM32F410,
		ChipFamilySTM32F413, ChipFamilySTM32F7Foundation, ChipFamilySTM32F4LP:
		return d.Read16(0x1fff7a22)

	case ChipFamilySTM32L4, ChipFamilySTM32L434X, ChipFamilySTM32L4X6:
		return d.Read16(0x1fff75e0)

	case ChipFamilySTM32L011, ChipFamilySTM32L0Cat2, ChipFamilySTM32L0,
		ChipFamilySTM32L0Cat5:
		return d.Read16(0x1ff8007c)

	case ChipFamilySTM32L1MediumLow, ChipFamilySTM32L1Cat2:
		return d.Read16(0x1ff8004c)

	case ChipFamilySTM32L1MediumHigh, ChipFamilySTM32L1High, ChipFamilySTM32L152RE:
		return d.Read16(0x1ff800cc)

	case ChipFamilySTM32F7, ChipFamilySTM32F7Advanced:
		return d.Read16(0x1ff0f442)

	}
	return 0, errors.New("unknown core")
}
