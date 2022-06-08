package serial

import (
	"github.com/albenik/go-serial"
)

type SCtrlError int

const (
	ReadError SCtrlError = iota
	Timeout
	WriteError
	ZeroSize
	Okay
)

// FIXME: allow config
func SOpenPort() (serial.Port, error) {

	mode := &serial.Mode{9600, 8, serial.NoParity, serial.OneStopBit}

	pPort, err := serial.Open("COM1", mode)
	if err != nil {
		return nil, err
	}

	return pPort, nil
}

func SClosePort(port serial.Port) {
	port.Close()
}

func SRead(port serial.Port, buff *[]byte) SCtrlError {
	waiting, _ := port.ReadyToRead()
	if waiting == 0 {
		return ZeroSize
	}

	pBuff := make([]byte, waiting)

	read, _ := port.Read(pBuff)

	if read != int(waiting) {
		return ReadError
	}

	*buff = append(*buff, pBuff...)
	return Okay
}
