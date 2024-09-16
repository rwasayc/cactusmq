package packet

import (
	"bytes"
	"log"
	"math"
)

// CONNECT – Connection Request
type ConnectionRequest struct {
	ProtocolName    ProtocolName       `json:"protocol_name"`
	ProtocolVersion ProtocolVersion    `json:"protocol_version"`
	Password        Password           `json:"password"`
	Username        FlagV[[]byte]      `json:"username"`
	ClientID        string             `json:"client_id"`
	Keepalive       uint16             `json:"keepalive"`
	CleanStart      FlagV[bool]        `json:"clean_start"`
	Reserved        bool               `json:"reserved"`
	Properties      *ConnectProperties `json:"properties"`
	Will            FlagV[ConnectWill] `json:"will"`
}

type ConnectFlags byte

func (flags ConnectFlags) CleanStart() bool {
	return 1&(flags>>1) > 0
}

func (flags ConnectFlags) WillFlag() bool {
	return 1&(flags>>2) > 0
}

func (flags ConnectFlags) WillQos() QoS {
	return QoS(3 & (flags >> 3))
}

func (flags ConnectFlags) WillRetain() bool {
	return 1&(flags>>5) > 0
}

func (flags ConnectFlags) PasswordFlag() bool {
	return 1&(flags>>6) > 0
}

func (flags ConnectFlags) UsernameFlag() bool {
	return 1&(flags>>7) > 0
}

func (flags ConnectFlags) Reserved() bool {
	return 1&flags > 0
}

type ConnectProperties struct {
	SessionExpiryInterval FlagV[uint32]   `json:"session_expiry_interval"`
	ReceiveMaximum        uint16          `json:"receive_maximum"`
	MaximumPacketSize     uint32          `json:"maximum_packet_size"`
	TopicAliasMaximum     uint16          `json:"topic_alias_maximum"`
	RequestResponseInfo   byte            `json:"request_response_info"`
	RequestProblemInfo    FlagV[byte]     `json:"request_problem_info"`
	UserProperty          []*UserProperty `json:"user_property"`
	AuthenticationMethod  string          `json:"authentication_method"`
	AuthenticationData    []byte          `json:"authentication_data"`
}

func (props *ConnectProperties) Decode(buf []byte) ([]byte, bool, error) {
	var err error
	var length uint32

	length, buf, err = decodeLength(buf)
	if err != nil {
		return buf, false, err
	}
	if length == 0 {
		return buf, false, nil
	}

	shouldRemain := len(buf) - int(length)
	for len(buf) > shouldRemain {
		var id Identifier
		id, buf, err = decodeIdentifier(buf)
		if err != nil {
			return buf, false, err
		}
		switch id {
		case IDSessionExpiryInterval:
			var sei uint32
			sei, buf, err = decodeUint32(buf)
			if err != nil {
				return buf, false, err
			}
			props.SessionExpiryInterval = NewFlagV(sei)
		case IDRequestProblemInformation:
			var rpi byte
			rpi, buf, err = decodeByte(buf)
			if err != nil {
				return buf, false, err
			}
			props.RequestProblemInfo = NewFlagV(rpi)
		case IDRequestResponseInformation:
			props.RequestResponseInfo, buf, err = decodeByte(buf)
			if err != nil {
				return buf, false, err
			}
		case IDReceiveMaximum:
			props.ReceiveMaximum, buf, err = decodeUint16(buf)
			if err != nil {
				return buf, false, err
			}
		case IDMaximumPacketSize:
			props.MaximumPacketSize, buf, err = decodeUint32(buf)
			if err != nil {
				return buf, false, err
			}
		case IDTopicAliasMaximum:
			props.TopicAliasMaximum, buf, err = decodeUint16(buf)
			if err != nil {
				return buf, false, err
			}
		case IDAuthenticationMethod:
			props.AuthenticationMethod, buf, err = decodeString(buf)
			if err != nil {
				return buf, false, err
			}
		case IDAuthenticationData:
			props.AuthenticationData, buf, err = decodeBytes(buf)
			if err != nil {
				return buf, false, err
			}
		case IDUserProperty:
			var k, v string
			k, v, buf, err = decodeStringPair(buf)
			if err != nil {
				return buf, false, err
			}
			if props.UserProperty == nil {
				props.UserProperty = []*UserProperty{}
			}
			props.UserProperty = append(props.UserProperty, &UserProperty{
				Key: k,
				Val: v,
			})
		default:
			return buf, false, RCMalformedPacket
		}
	}
	return buf, true, nil
}

func (props *ConnectProperties) Encode(buf *bytes.Buffer) error {
	if props == nil {
		buf.Write(encodeLength(0))
		return nil
	}
	emptyBuf := bytes.NewBuffer(nil)

	if props.SessionExpiryInterval.Flag() {
		emptyBuf.WriteByte(byte(IDSessionExpiryInterval))
		emptyBuf.Write(encodeUint32(props.SessionExpiryInterval.Value()))
	}
	if props.RequestProblemInfo.Flag() {
		emptyBuf.WriteByte(byte(IDRequestProblemInformation))
		emptyBuf.WriteByte(props.RequestProblemInfo.Value())
	}
	if props.RequestResponseInfo > 0 {
		emptyBuf.WriteByte(byte(IDRequestResponseInformation))
		emptyBuf.WriteByte(props.RequestResponseInfo)
	}
	if props.ReceiveMaximum > 0 {
		emptyBuf.WriteByte(byte(IDReceiveMaximum))
		emptyBuf.Write(encodeUint16(props.ReceiveMaximum))
	}
	if props.MaximumPacketSize > 0 {
		emptyBuf.WriteByte(byte(IDMaximumPacketSize))
		emptyBuf.Write(encodeUint32(props.MaximumPacketSize))
	}
	if props.TopicAliasMaximum > 0 {
		emptyBuf.WriteByte(byte(IDTopicAliasMaximum))
		emptyBuf.Write(encodeUint16(props.TopicAliasMaximum))
	}
	if props.AuthenticationMethod != "" {
		emptyBuf.WriteByte(byte(IDAuthenticationMethod))
		emptyBuf.Write(encodeString(props.AuthenticationMethod))
	}
	if props.AuthenticationData != nil {
		emptyBuf.WriteByte(byte(IDAuthenticationData))
		emptyBuf.Write(encodeBytes(props.AuthenticationData))
	}
	for _, prop := range props.UserProperty {
		emptyBuf.WriteByte(byte(IDUserProperty))
		emptyBuf.Write(encodeString(prop.Key))
		emptyBuf.Write(encodeString(prop.Val))
	}
	buf.Write(encodeLength(uint32(emptyBuf.Len())))
	buf.Write(emptyBuf.Bytes())
	return nil
}

type ConnectWill struct {
	Payload    []byte          `json:"payload"`    // -
	Topic      string          `json:"topic"`      // -
	Qos        FlagV[QoS]      `json:"qos"`        // -
	Retain     bool            `json:"retain"`     // -
	Properties *WillProperties `json:"properties"` // -
}

type WillProperties struct {
	MessageExpiryInterval uint32          `json:"message_expiry_interval"`
	PayloadFormat         FlagV[byte]     `json:"payload_format"`
	ContentType           string          `json:"content_type"`
	ResponseTopic         string          `json:"response_topic"`
	CorrelationData       []byte          `json:"correlation_data"`
	User                  []*UserProperty `json:"user_property"`
	WillDelayInterval     uint32          `json:"will_delay_interval"`
}

func (props *WillProperties) Decode(buf []byte) ([]byte, bool, error) {
	var err error
	var length uint32

	length, buf, err = decodeLength(buf)
	if err != nil {
		return buf, false, err
	}
	if length == 0 {
		return buf, false, nil
	}
	shouldRemain := len(buf) - int(length)

	for len(buf) > shouldRemain {
		var id byte
		id, buf, err = decodeByte(buf)
		if err != nil {
			return buf, false, err
		}
		switch Identifier(id) {
		case IDMessageExpiryInterval:
			props.MessageExpiryInterval, buf, err = decodeUint32(buf)
			if err != nil {
				return buf, false, err
			}
		case IDPayloadFormatIndicator:
			var pf byte
			pf, buf, err = decodeByte(buf)
			if err != nil {
				return buf, false, err
			}
			props.PayloadFormat = NewFlagV(pf)
		case IDContentType:
			props.ContentType, buf, err = decodeString(buf)
			if err != nil {
				return buf, false, err
			}
		case IDResponseTopic:
			props.ResponseTopic, buf, err = decodeString(buf)
			if err != nil {
				return buf, false, err
			}
		case IDCorrelationData:
			props.CorrelationData, buf, err = decodeBytes(buf)
			if err != nil {
				return buf, false, err
			}
		case IDUserProperty:
			var k, v string
			k, v, buf, err = decodeStringPair(buf)
			if err != nil {
				return buf, false, err
			}
			props.User = append(props.User, &UserProperty{
				Key: k,
				Val: v,
			})
		case IDWillDelayInterval:
			props.WillDelayInterval, buf, err = decodeUint32(buf)
			if err != nil {
				return buf, false, err
			}
		default:
			return buf, false, RCMalformedPacket
		}
	}

	return buf, true, nil
}

func (props *WillProperties) Encode(buf *bytes.Buffer) error {
	tmpBuf := bytes.NewBuffer(nil)

	if props.MessageExpiryInterval > 0 {
		tmpBuf.WriteByte(byte(IDMessageExpiryInterval))
		tmpBuf.Write(encodeUint32(props.MessageExpiryInterval))
	}
	if props.PayloadFormat.Flag() {
		tmpBuf.WriteByte(byte(IDPayloadFormatIndicator))
		tmpBuf.WriteByte(props.PayloadFormat.Value())
	}
	if props.ContentType != "" {
		tmpBuf.WriteByte(byte(IDContentType))
		tmpBuf.Write(encodeString(props.ContentType))
	}
	if props.ResponseTopic != "" {
		tmpBuf.WriteByte(byte(IDResponseTopic))
		tmpBuf.Write(encodeString(props.ResponseTopic))
	}
	if props.CorrelationData != nil {
		tmpBuf.WriteByte(byte(IDCorrelationData))
		tmpBuf.Write(encodeBytes(props.CorrelationData))
	}
	if props.User != nil {
		for _, prop := range props.User {
			tmpBuf.WriteByte(byte(IDUserProperty))
			tmpBuf.Write(encodeString(prop.Key))
			tmpBuf.Write(encodeString(prop.Val))
		}
	}
	if props.WillDelayInterval > 0 {
		tmpBuf.WriteByte(byte(IDWillDelayInterval))
		tmpBuf.Write(encodeUint32(props.WillDelayInterval))
	}

	buf.Write(encodeLength(uint32(tmpBuf.Len())))
	buf.Write(tmpBuf.Bytes())
	return nil
}

func (cr *ConnectionRequest) Decode(buf []byte) (err error) {
	cr.ProtocolName, buf, err = decodeBytes(buf)
	if err != nil {
		return RCMalformedPacket
	}
	var pv byte
	pv, buf, err = decodeByte(buf)
	if err != nil {
		log.Println("pv, buf, err = decodeByte(buf)", err)
		return RCMalformedPacket
	}
	cr.ProtocolVersion = ProtocolVersion(pv)

	var flags byte
	flags, buf, err = decodeByte(buf)
	if err != nil {
		log.Println("flags, buf, err = decodeByte(buf)", err)
		return RCMalformedPacket
	}
	cflags := ConnectFlags(flags)

	if cflags.CleanStart() {
		cr.CleanStart = NewFlagV(true)
	}

	if cflags.Reserved() {
		cr.Reserved = true
	}

	cr.Keepalive, buf, err = decodeUint16(buf)
	if err != nil {
		log.Println("cr.Keepalive, buf, err = decodeUint16(buf)", err)
		return RCMalformedPacket
	}

	if cr.ProtocolVersion == ProtoVer5 {
		// todo
		prop := &ConnectProperties{}
		var has bool
		buf, has, err = prop.Decode(buf)
		if err != nil {
			log.Println("buf, err = cr.Properties.Decode(buf)", err)
			return RCMalformedPacket
		}
		if has {
			cr.Properties = prop
		}
	}

	cr.ClientID, buf, err = decodeString(buf)
	if err != nil {
		log.Println("cr.ClientID, buf, err = decodeString(buf)", err)
		return RCMalformedPacket
	}

	if cflags.WillFlag() {
		will := ConnectWill{}
		if cr.ProtocolVersion == ProtoVer5 {
			// todo
			var has bool
			var willProps *WillProperties = &WillProperties{}
			buf, has, err = willProps.Decode(buf)
			if err != nil {
				return RCMalformedPacket
			}
			if has {
				will.Properties = willProps
			}
		}
		will.Qos = NewFlagV(cflags.WillQos())
		will.Retain = cflags.WillRetain()

		will.Topic, buf, err = decodeString(buf)
		if err != nil {
			log.Println("will.Topic, buf, err = decodeString(buf)", err)
			return RCMalformedPacket
		}

		will.Payload, buf, err = decodeBytes(buf)
		if err != nil {
			log.Println("will.Payload, buf, err = decodeBytes(buf)", err)
			return RCMalformedPacket
		}
		cr.Will = NewFlagV(will)
	}

	if cflags.UsernameFlag() {
		var username []byte
		username, buf, err = decodeBytes(buf)
		if err != nil {
			log.Println("username, buf, err = decodeBytes(buf)", err)
			return RCMalformedPacket
		}
		cr.Username = NewFlagV(username)
	}

	if cflags.PasswordFlag() {
		var password []byte
		password, buf, err = decodeBytes(buf)
		if err != nil {
			log.Println("password, buf, err = decodeBytes(buf)", err)
			return RCMalformedPacket
		}
		cr.Password = NewPassword(password)
	}

	return nil
}

func (cr *ConnectionRequest) Encode(ver ProtocolVersion, buf *bytes.Buffer) error {
	var err error
	// encode protocol name
	_, err = buf.Write(encodeBytes(cr.ProtocolName))
	if err != nil {
		return err
	}
	// encode protocol version
	err = buf.WriteByte(byte(cr.ProtocolVersion))
	if err != nil {
		return err
	}
	// encode connect flags
	{
		flags := byte(0)
		if cr.Will.Flag() {
			flags |= 0b00000100
			flags |= byte(cr.Will.Value().Qos.Value()) << 3
			if cr.Will.Value().Retain {
				flags |= 0b00100000
			}
		}
		if cr.Username.Flag() {
			flags |= 0b10000000
		}
		if cr.Password.Flag() {
			flags |= 0b01000000
		}
		if cr.Reserved {
			flags |= 0b00000001
		}
		if cr.CleanStart.Flag() {
			flags |= 0b00000010
		}

		err = buf.WriteByte(flags)
		if err != nil {
			return err
		}
	}
	// encode keepalive
	_, err = buf.Write(encodeUint16(cr.Keepalive))
	if err != nil {
		return err
	}
	// encode properties
	if cr.ProtocolVersion == ProtoVer5 {
		// todo
		err = cr.Properties.Encode(buf)
		if err != nil {
			return err
		}
	}

	// encode client id
	_, err = buf.Write(encodeString(cr.ClientID))
	if err != nil {
		return err
	}

	// encode will
	if cr.Will.Flag() {
		will := cr.Will.Value()
		if cr.ProtocolVersion == ProtoVer5 {
			// todo
			err = will.Properties.Encode(buf)
			if err != nil {
				return err
			}
		}
		buf.Write(encodeString(will.Topic))
		buf.Write(encodeBytes(will.Payload))
	}

	// encode username
	if cr.Username.Flag() {
		// todo
		_, err = buf.Write(encodeBytes(cr.Username.Value()))
		if err != nil {
			return err
		}
	}

	// encode password
	if cr.Password.Flag() {
		// todo
		_, err = buf.Write(encodeBytes(cr.Password.Value()))
		if err != nil {
			return err
		}
	}
	return nil
}

func (cr *ConnectionRequest) Validate() RCode {
	// check protocol version
	if !cr.ProtocolVersion.IsValid() {
		return RCUnsupportedProtocol
	}
	// check protocol name
	if !cr.ProtocolName.IsValid() {
		return RCUnsupportedProtocol
	}
	// check reserved
	if cr.Reserved {
		return RCMalformedPacket
	}
	// check password
	{
		if cr.Password.Flag() && len(cr.Password.Value()) == 0 {
			return RCMalformedPacket
		}
		if !cr.Password.Flag() && len(cr.Password.Value()) > 0 {
			return RCMalformedPacket
		}
	}
	// check username
	{
		if cr.Username.Flag() && len(cr.Username.Value()) == 0 {
			return RCMalformedPacket
		}
		if !cr.Username.Flag() && len(cr.Username.Value()) > 0 {
			return RCMalformedPacket
		}
	}
	// check will
	{
		if cr.Will.Flag() {
			if len(cr.Will.Value().Topic) == 0 {
				return RCMalformedPacket
			}
			if cr.Will.Value().Qos.Flag() {
				if cr.Will.Value().Qos.Value() > QoS2 {
					return RCMalformedPacket
				}
			}
		}
	}
	// check client id
	if len(cr.ClientID) > math.MaxUint16 {
		return RCMalformedPacket
	}
	return RCSuccess
}
