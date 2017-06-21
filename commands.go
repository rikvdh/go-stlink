package stlink

type stlinkCmd uint8

const (
	stlinkCmdGetVersion               stlinkCmd = 0xf1
	stlinkCmdDebug                    stlinkCmd = 0xf2
	stlinkCmdDfu                      stlinkCmd = 0xf3
	stlinkCmdDfuExit                  stlinkCmd = 0x07
	stlinkCmdDfuGetVersion            stlinkCmd = 0x08
	stlinkCmdGetCurrentMode           stlinkCmd = 0xf5
	stlinkCmdGetTargetVoltage         stlinkCmd = 0xf7
	stlinkCmdDebugGetStatus           stlinkCmd = 0x01
	stlinkCmdDebugForce               stlinkCmd = 0x02
	stlinkCmdDebugReadMem32           stlinkCmd = 0x07
	stlinkCmdDebugWriteMem32          stlinkCmd = 0x08
	stlinkCmdDebugWriteMem8           stlinkCmd = 0x0d
	stlinkCmdDebugEnterMode           stlinkCmd = 0x20
	stlinkCmdDebugEnterSwd            stlinkCmd = 0xa3
	stlinkCmdDebugEnterJtag           stlinkCmd = 0x00
	stlinkCmdDebugExit                stlinkCmd = 0x21
	stlinkCmdDebugReadCoreid          stlinkCmd = 0x22
	stlinkCmdDebugResetsys            stlinkCmd = 0x03
	stlinkCmdDebugReadallregs         stlinkCmd = 0x04
	stlinkCmdDebugReadReg             stlinkCmd = 0x33
	stlinkCmdDebugWriteReg            stlinkCmd = 0x34
	stlinkCmdDebugRunCore             stlinkCmd = 0x09
	stlinkCmdDebugStepCore            stlinkCmd = 0x0a
	stlinkCmdDebugWriteRegpc          stlinkCmd = 0x34
	stlinkCmdDebugHardReset           stlinkCmd = 0x3c
	stlinkCmdDebugReadcoreregs        stlinkCmd = 0x3a
	stlinkCmdDebugSetfp               stlinkCmd = 0x0b
	stlinkCmdDebugJtagWritedebug32bit stlinkCmd = 0x35
	stlinkCmdDebugJtagReaddebug32bit  stlinkCmd = 0x36
	stlinkCmdDebugSwdSetFreq          stlinkCmd = 0x43
)
