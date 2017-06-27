package stlink

import (
	"errors"
)

type FlashLoader interface {
	Init(voltage float32) error
}

func (d *Device) getFlashloader() (FlashLoader, error) {
	pn, err := d.DevID()
	if err != nil {
		return nil, err
	}
	/*v, err := d.TargetVoltage()
	if err != nil {
		return nil, err
	}*/

	switch pn {
	case ChipFamilySTM32F0, ChipFamilySTM32F09X, ChipFamilySTM32F0Small,
		ChipFamilySTM32F04, ChipFamilySTM32F0Can, ChipFamilySTM32F1Medium,
		ChipFamilySTM32F1Low, ChipFamilySTM32F1High, ChipFamilySTM32F1Connectivity,
		ChipFamilySTM32F1VLMedium, ChipFamilySTM32F3, ChipFamilySTM32F1VLHigh,
		ChipFamilySTM32F37x, ChipFamilySTM32F334, ChipFamilySTM32F3Small,
		ChipFamilySTM32F303High:
		// STM32FP
		return nil, nil
	case ChipFamilySTM32F2, ChipFamilySTM32F4, ChipFamilySTM32F4HD,
		ChipFamilySTM32F446, ChipFamilySTM32F4LP, ChipFamilySTM32F411RE,
		ChipFamilySTM32F4DE, ChipFamilySTM32F4DSI, ChipFamilySTM32F412,
		ChipFamilySTM32F410, ChipFamilySTM32F7, ChipFamilySTM32F7Advanced,
		ChipFamilySTM32F413, ChipFamilySTM32F7Foundation:
		// STM32FS
		/*l := stm32fs.STM32FS{}
		l.Init(v)
		return l, nil*/
	case ChipFamilySTM32F1XL:
		// STM32FPXL
		return nil, nil
	case ChipFamilySTM32L011, ChipFamilySTM32L0Cat2, ChipFamilySTM32L0,
		ChipFamilySTM32L0Cat5:
		//STM32L0
	case ChipFamilySTM32L1MediumLow, ChipFamilySTM32L1MediumHigh, ChipFamilySTM32L1Cat2,
		ChipFamilySTM32L1High, ChipFamilySTM32L152RE, ChipFamilySTM32L4,
		ChipFamilySTM32L434X, ChipFamilySTM32L4X6:
		// None
		return nil, nil
	}
	return nil, errors.New("unknown core")
}
