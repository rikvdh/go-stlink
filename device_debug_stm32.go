package stlink

const (
	AIRCRReg uint32 = 0xe000ed0c
	DHCSRReg uint32 = 0xe000edf0
	DEMCRReg uint32 = 0xe000edfc

	AIRCRKey            uint32 = 0x05fa0000
	AIRCRSysResetReqBit uint32 = 0x00000004
	AIRCRSysResetReq    uint32 = AIRCRKey | AIRCRSysResetReqBit

	DHCSRKey           uint32 = 0xa05f0000
	DHCSRDebugEnBit    uint32 = 0x00000001
	DHCSRHaltBit       uint32 = 0x00000002
	DHCSRStepBit       uint32 = 0x00000004
	DHCSRStatusHaltBit uint32 = 0x00020000
	DHCSRDebugDis      uint32 = DHCSRKey
	DHCSRDebugEn       uint32 = DHCSRKey | DHCSRDebugEnBit
	DHCSRHalt          uint32 = DHCSRKey | DHCSRDebugEnBit | DHCSRHaltBit
	DHCSRStep          uint32 = DHCSRKey | DHCSRDebugEnBit | DHCSRStepBit

	DEMCRRunAfterReset  uint32 = 0x00000000
	DEMCRHaltAfterReset uint32 = 0x00000001
)

func (d *Device) CoreResetHalt() error {
	if err := d.Write32(DHCSRReg, DHCSRHalt); err != nil {
		return err
	}
	if err := d.Write32(DEMCRReg, DEMCRHaltAfterReset); err != nil {
		return err
	}
	if err := d.Write32(AIRCRReg, AIRCRSysResetReq); err != nil {
		return err
	}
	_, err := d.Read32(AIRCRReg)
	return err
}
