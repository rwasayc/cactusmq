package packet

import "bytes"

type PublishReceived struct {
	PacketID   uint16
	ReasonCode RCode
	Properties BaseProperties
}

func (pa *PublishReceived) Decode(buf []byte) error {
	var err error
	pa.PacketID, buf, err = decodeUint16(buf)
	if err != nil {
		return err
	}
	var code byte
	code, buf, err = decodeByte(buf)
	if err != nil {
		return err
	}
	pa.ReasonCode = RCode(code)

	buf, err = pa.Properties.Decode(buf)
	if err != nil {
		return err
	}
	return nil
}

func (pa *PublishReceived) Encode(ver ProtocolVersion, buf *bytes.Buffer) error {
	var err error
	_, err = buf.Write(encodeUint16(pa.PacketID))
	if err != nil {
		return err
	}
	err = buf.WriteByte(byte(pa.ReasonCode))
	if err != nil {
		return err
	}
	err = pa.Properties.Encode(buf)
	if err != nil {
		return err
	}
	return nil
}

func (pa *PublishReceived) Validate() RCode {
	return RCSuccess
}
