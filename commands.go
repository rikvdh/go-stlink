package stlink

type stlinkCmd uint8

const (
	stlinkCmdGetVersion               stlinkCmd = 0xf1
	stlinkCmdDebug                              = 0xf2
	stlinkCmdDfu                                = 0xf3
	stlinkCmdDfuExit                            = 0x07
	stlinkCmdDfuGetVersion                      = 0x08
	stlinkCmdGetCurrentMode                     = 0xf5
	stlinkCmdGetTargetVoltage                   = 0xf7
	stlinkCmdDebugGetStatus                     = 0x01
	stlinkCmdDebugForce                         = 0x02
	stlinkCmdDebugReadMem32                     = 0x07
	stlinkCmdDebugWriteMem32                    = 0x08
	stlinkCmdDebugWriteMem8                     = 0x0d
	stlinkCmdDebugEnterMode                     = 0x20
	stlinkCmdDebugEnterSwd                      = 0xa3
	stlinkCmdDebugEnterJtag                     = 0x00
	stlinkCmdDebugExit                          = 0x21
	stlinkCmdDebugReadCoreid                    = 0x22
	stlinkCmdDebugResetsys                      = 0x03
	stlinkCmdDebugReadallregs                   = 0x04
	stlinkCmdDebugReadReg                       = 0x33
	stlinkCmdDebugWriteReg                      = 0x34
	stlinkCmdDebugRunCore                       = 0x09
	stlinkCmdDebugStepCore                      = 0x0a
	stlinkCmdDebugWriteRegpc                    = 0x34
	stlinkCmdDebugHardReset                     = 0x3c
	stlinkCmdDebugReadcoreregs                  = 0x3a
	stlinkCmdDebugSetfp                         = 0x0b
	stlinkCmdDebugJtagWritedebug32bit           = 0x35
	stlinkCmdDebugJtagReaddebug32bit            = 0x36
	stlinkCmdDebugSwdSetFreq                    = 0x43
)
