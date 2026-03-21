package models

type PortEntry struct {
	Port     uint16
	Protocol string
	PID      int
	Process  string
}
