package card

func Command10Initalize() {
	switch currentStep {
	case 1:
		if card.r == HAS_CARD_1 {
			//eject
		}
		break
	default:
		break
	}

	if currentStep > 1 {
		runningCommand = false
	}
}

func Command20ReadStatus() {
	switch currentStep {
	default:
		card.softReset()
		runningCommand = false
		break
	}
}

func Command40Cancel() {
	switch currentStep {
	default:
		card.softReset()
		runningCommand = false
		break
	}
}

func Command78PrintSettings() {
	switch currentStep {
	default:
		card.softReset()
		runningCommand = false
		break
	}
}

func Command7ARegisterFont() {
	switch currentStep {
	default:
		card.softReset()
		runningCommand = false
		break
	}
}

func Command7BPrintImage() {
	switch currentStep {
	case 1:
		if !*HasCard {
			// FIXME: setperr
		}
		break
	default:
		break
	}

	if currentStep > 1 {
		runningCommand = false
	}
}

func Command7DErase() {
	switch currentStep {
	case 1:
		if !*HasCard {
			// FIXME: setperr
		}
		break
	default:
		break
	}

	if currentStep > 1 {
		runningCommand = false
	}
}

func Command7EPrintBarcode() {
	switch currentStep {
	case 1:
		if !*HasCard {
			// FIXME: setperr
		}
		break
	default:
		break
	}

	if currentStep > 1 {
		runningCommand = false
	}
}

func Command80Eject() {
	switch currentStep {
	case 1: // Special for "Transfer Card Data" in MT2EXP, we need 2 RUNNING_COMMAND replies
		break
	case 2:
		if !*HasCard {
			// FIXME: eject
		} else {
			// FIXME: s ille
		}
		break
	default:
		break
	}

	if currentStep > 2 {
		runningCommand = false
	}
}

func CommandA0Clean() {
	switch currentStep {
	case 1:
		if !*HasCard {
			card.s = WAITING_FOR_CARD
			currentStep--
		}
		break
	case 2:
		break
	case 3:
		// FIXME: eject
		break
	default:
		break
	}

	if currentStep > 3 {
		runningCommand = false
	}
}
