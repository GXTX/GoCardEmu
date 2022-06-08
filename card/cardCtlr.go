package card

import "fmt"

func handlePacket() {
	if !runningCommand {
		if card.s != ILLEGAL_COMMAND {
			card.s = NO_JOB
		}

		card.p = NO_ERR
	}

	// do stuff / update status
	if *HasCard && card.r == NO_CARD {
		card.r = HAS_CARD_1

		if runningCommand && card.s == WAITING_FOR_CARD {
			card.s = RUNNING_COMMAND
		}
	}

	if runningCommand {
		switch currentCommand {
		case 0x10:
			Command10Initalize()
			break
		case 0x20:
			Command20ReadStatus()
			break
		case 0x33:
			break
		case 0x40:
			Command40Cancel()
			break
		case 0x53:
			break
		case 0x78:
			Command78PrintSettings()
			break
		case 0x7A:
			Command7ARegisterFont()
			break
		case 0x7B:
			Command7BPrintImage()
			break
		case 0x7C:
			break
		case 0x7D:
			Command7DErase()
			break
		case 0x7E:
			Command7EPrintBarcode()
			break
		case 0x80:
			Command80Eject()
			break
		case 0xA0:
			CommandA0Clean()
			break
		case 0xB0:
			break
		case 0xE1:
			break
		case 0xF0:
			break
		case 0xF1:
			break
		case 0xF5:
			break
		default:
			fmt.Printf("Unhandled command %x\n", currentCommand)
			// FIXME: set sill
			break
		}
		currentStep++
	}
}

func updateStatusInBuffer() {
	commandBuffer[1] = byte(card.r)
	commandBuffer[2] = byte(card.p)
	commandBuffer[3] = byte(card.s)

	fmt.Printf("%x, %x, %x, %x\n", currentCommand, card.r, card.p, card.s)
}

func ProcessPacket(b *[]byte) cStatus {
	sync := (*b)[0]

	if sync == 0x05 {
		handlePacket()
		fmt.Print("sync")
		*b = (*b)[1:]
		return ServerWaitingreply
	} else if sync != 0x02 {
		fmt.Print("missing stx")
		*b = (*b)[1:]
		return SyncError
	}

	if len(*b) < 8 {
		fmt.Print("not enough in the buffer to even bother")
		return SizeError
	}

	count := (*b)[1]
	if int(count) > len(*b)-2 { // does not count stx and sum
		// wrong size
		fmt.Print("wrong size")
		return SizeError
	}

	if (*b)[count] != 0x03 {
		// missing end
		fmt.Print("missing end")
		*b = (*b)[:count]
		return SyntaxError
	}

	actualSum := count

	// clear the current
	currentPacket = nil

	// loop through the rest of the buffer - 1
	for i := 2; i < int(count+1); i++ {
		val := (*b)[i]
		currentPacket = append(currentPacket, val)
		actualSum ^= val
	}

	packetSum := (*b)[count+1]

	// remove end_of_text
	currentPacket = currentPacket[:len(currentPacket)-1]

	// remove what we've handled
	*b = (*b)[(count + 2):]

	if packetSum != actualSum {
		// bad sum
		fmt.Print("bad sum")
		return ChecksumError
	}

	// MT2EXP "Transfer Card Data" interrupts Eject to do a CheckStatus, eject here otherwise the system will error
	if currentCommand == 0x80 {
		// FIXME: eject card
	}

	currentCommand = int(currentPacket[0])

	// remove host RPS and current command
	currentPacket = currentPacket[4:]

	card.softReset()
	card.s = RUNNING_COMMAND
	runningCommand = true
	currentStep = 0

	// pre format our resp buffer
	commandBuffer = make([]byte, 4)
	commandBuffer[0] = byte(currentCommand)
	commandBuffer[1] = byte(card.r)
	commandBuffer[2] = byte(card.p)
	commandBuffer[3] = byte(card.s)
	fmt.Printf("%x, %x, %x, %x\n", currentCommand, card.r, card.p, card.s)

	return Okay
}

func GeneratePacket(b *[]byte) {
	// Make sure the buffer we receive is clear
	*b = nil

	updateStatusInBuffer()

	count := len(commandBuffer) + 2

	buildPacket := make([]byte, count)

	buildPacket[0] = STX
	buildPacket[1] = byte(count)

	packetSum := byte(count)

	for pos, resp := range commandBuffer {
		buildPacket[pos+2] = resp
		packetSum ^= resp

	}

	buildPacket = append(buildPacket, []byte{ETX, packetSum ^ ETX}...)

	*b = append(*b, buildPacket...)
}
