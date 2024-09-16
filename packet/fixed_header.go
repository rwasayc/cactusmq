package packet

// FixedHeader is a struct that represents the fixed header of a MQTT packet.
type FixedHeader struct {
	flags byte
	typ   CPType
	rlen  uint32 // Remaining Length
}

// GetType returns the type of the control packet.
func (fh *FixedHeader) GetType() CPType {
	if fh == nil {
		return 0
	}
	return fh.typ
}

// GetRemainLength returns the remaining length of the control packet.
func (fh *FixedHeader) GetRemainingLength() uint32 {
	if fh == nil {
		return 0
	}
	return fh.rlen
}

func (fh *FixedHeader) GetFlags0() bool {
	if fh == nil {
		return false
	}
	return fh.flags&0x01 == 0x01
}
func (fh *FixedHeader) GetFlags1() bool {
	if fh == nil {
		return false
	}
	return fh.flags&0x02 == 0x02
}

func (fh *FixedHeader) GetFlags2() bool {
	if fh == nil {
		return false
	}
	return fh.flags&0x04 == 0x04
}

func (fh *FixedHeader) GetFlags3() bool {
	if fh == nil {
		return false
	}
	return fh.flags&0x08 == 0x08
}

// CPType is a byte that represents the type of the control packet.
type CPType byte

// cpType2String is a map that maps CPType to its string representation.
var cpType2String = map[CPType]string{
	Reserved:    "Reserved",
	CONNECT:     "CONNECT",
	CONNACK:     "CONNACK",
	PUBLISH:     "PUBLISH",
	PUBACK:      "PUBACK",
	PUBREC:      "PUBREC",
	PUBREL:      "PUBREL",
	PUBCOMP:     "PUBCOMP",
	SUBSCRIBE:   "SUBSCRIBE",
	SUBACK:      "SUBACK",
	UNSUBSCRIBE: "UNSUBSCRIBE",
	UNSUBACK:    "UNSUBACK",
	PINGREQ:     "PINGREQ",
	PINGRESP:    "PINGRESP",
	DISCONNECT:  "DISCONNECT",
	AUTH:        "AUTH",
}

// String returns the string representation of the control packet type.
func (t CPType) String() string {
	if str, ok := cpType2String[t]; ok {
		return str
	}
	return "Invalid Control Packet Type"
}

// Constants for the control packet types.
const (
	Reserved    CPType = iota // Value: 0, Flow: Forbidden, Note: Reserved for future use
	CONNECT                   // Value: 1, Flow: Client to Server, Note: Connection request
	CONNACK                   // Value: 2, Flow: Server to Client, Note: Connect acknowledgment
	PUBLISH                   // Value: 3, Flow: Client to Server or Server to Client, Note: Publish message
	PUBACK                    // Value: 4, Flow: Client to Server or Server to Client, Note: Publish acknowledgment (QoS 1)
	PUBREC                    // Value: 5, Flow: Client to Server or Server to Client, Note: Publish received (QoS 2 delivery part 1)
	PUBREL                    // Value: 6, Flow: Client to Server or Server to Client, Note: Publish release (QoS 2 delivery part 2)
	PUBCOMP                   // Value: 7, Flow: Client to Server or Server to Client, Note: Publish complete (QoS 2 delivery part 3)
	SUBSCRIBE                 // Value: 8, Flow: Client to Server, Note: Subscribe request
	SUBACK                    // Value: 9, Flow: Server to Client, Note: Subscribe acknowledgment
	UNSUBSCRIBE               // Value: 10, Flow: Client to Server, Note: Unsubscribe request
	UNSUBACK                  // Value: 11, Flow: Server to Client, Note: Unsubscribe Acknowledgment
	PINGREQ                   // Value: 12, Flow: Client to Server, Note: PING Request
	PINGRESP                  // Value: 13, Flow: Server to Client, Note: PING Response
	DISCONNECT                // Value: 14, Flow: Client to Server or Server to Client, Note: Disconnect notification
	AUTH                      // Value: 15, Flow: Client to Server or Server to Client, Note: Authentication exchange
)
