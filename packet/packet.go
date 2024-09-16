package packet

import "bytes"

type Codec interface {
	Encode(ProtocolVersion, *bytes.Buffer) error
	Decode([]byte) error
}

const MaxRemainingLength = 268435455

type PayloadFormatIndicator byte

const (
	PFI_BYTE PayloadFormatIndicator = 0 // Payload is unspecified bytes
	PFI_UTF8 PayloadFormatIndicator = 1 // Payload is UTF-8 Encoded Character Data
)

var protoVer2ProtocolName = map[ProtocolVersion]ProtocolName{
	ProtoVer31:  FixedProtocolNameV31,
	ProtoVer311: FixedProtocolNameV311,
	ProtoVer5:   FixedProtocolNameV5,
}

// MQTT QoS 0,1,2
type QoS byte

const (
	QoS0 QoS = iota // At most once delivery
	QoS1            // At least once delivery
	QoS2            // Exactly once delivery
)

func (q QoS) IsValid() bool {
	return q >= QoS0 && q <= QoS2
}

type ProtocolVersion byte

const (
	ProtoVer31  ProtocolVersion = 3
	ProtoVer311 ProtocolVersion = 4
	ProtoVer5   ProtocolVersion = 5
)

func (p ProtocolVersion) IsValid() bool {
	_, ok := protoVer2ProtocolName[p]
	return ok
}

type ProtocolName []byte

var FixedProtocolNameV5 ProtocolName = []byte{'M', 'Q', 'T', 'T'}
var FixedProtocolNameV311 ProtocolName = []byte{'M', 'Q', 'T', 'T'}
var FixedProtocolNameV31 ProtocolName = []byte{'M', 'Q', 'I', 's', 'd', 'p'}

func (p ProtocolName) IsValid() bool {
	return bytes.Equal(p, FixedProtocolNameV5) || bytes.Equal(p, FixedProtocolNameV311) || bytes.Equal(p, FixedProtocolNameV31)
}
