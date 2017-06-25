package stlink

//go:generate go run cmd/getpartlist/main.go

type Target struct {
	Type       string
	Core       string
	Frequency  uint
	FlashSize  uint
	EepromSize uint
	SramSize   uint
}
