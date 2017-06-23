package stlink

//go:generate go run cmd/getpartlist/main.go

type Target struct {
	Type       string
	Core       string
	Frequency  string
	FlashSize  string
	EepromSize string
	SramSize   string
}
