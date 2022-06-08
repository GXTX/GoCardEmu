package card

//
// Command structure
// [stx][len][command][r][p][s][command data][etx][sum]
//

var currentPacket []byte
var currentCommand int
var runningCommand bool
var currentStep int
var commandBuffer []byte // resp buff
var HasCard *bool

type cStatus int

const (
	SizeError cStatus = iota
	SyncError
	SyntaxError
	ChecksumError
	EmptyResponseError
	ServerWaitingreply
	Okay
)

const STX = 0x02
const ETX = 0x03
const REQ = 0x05
const ACK = 0x06

type rStatus byte

const (
	NO_CARD rStatus = 0x30
	HAS_CARD_1
	CARD_STATUS_ERR
	HAS_CARD_2
	EJECTING_CARD
)

type pStatus byte

const (
	NO_ERR pStatus = 0x30
	READ_ERR
	WRITE_ERR
	CARD_JAM
	MOTOR_ERROR
	PRINT_ERR
	ILLEGAL_ERR = 0x38
	BATTERY_ERR = 0x40
	SYSTEM_ERR
	TRACK_1_READ_ERR = 0x51
	TRACK_2_READ_ERR
	TRACK_3_READ_ERR
	TRACK_1_AND_2_READ_ERR
	TRACK_1_AND_3_READ_ERR
	TRACK_2_AND_3_READ_ERR
)

type sStatus byte

const (
	NO_JOB          sStatus = 0x30
	ILLEGAL_COMMAND         = 0x32
	RUNNING_COMMAND         = 0x33
	WAITING_FOR_CARD
	DISPENSER_EMPTY
	NO_DISPENSER
	CARD_FULL
)

type cardStatus struct {
	r rStatus
	p pStatus
	s sStatus
}

var card = cardStatus{r: NO_CARD}

func (c *cardStatus) hardReset() {
	c.r = NO_CARD
	c.softReset()
}

func (c *cardStatus) softReset() {
	c.p = NO_ERR
	c.s = NO_JOB
}
