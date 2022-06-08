package main

import (
	"GoCardEmu/card"
	"GoCardEmu/serial"
	"os"
	"time"
)

var insertedCard bool

func main() {
	pPort, _ := serial.SOpenPort()

	var readBuffer []byte
	var writeBuffer []byte

	//card.SetCardPtr(&insertedCard)
	card.HasCard = &insertedCard

	running := true
	for running {
		if serial.SRead(pPort, &readBuffer) != serial.SCtrlError(card.Okay) && len(readBuffer) == 0 {
			time.Sleep(500 / 2 * time.Microsecond)
			continue
		}

		cardStatus := card.ProcessPacket(&readBuffer)

		// we need to ack as soon as we receive a valid command, then wait for a ENQ
		if cardStatus == card.Okay {
			pPort.Write([]byte{card.ACK})
		} else if cardStatus == card.ServerWaitingreply {
			writeBuffer = nil
			card.GeneratePacket(&writeBuffer)
			pPort.Write(writeBuffer)
		}

		time.Sleep(500 * time.Microsecond)
	}

	os.Exit(0)
}
