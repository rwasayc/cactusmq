package packet

import (
	"encoding/binary"
	"unsafe"
)

// decodeBytes
func decodeBytes(buf []byte) ([]byte, []byte, error) {
	length, buf, err := decodeUint16(buf)
	if err != nil {
		return nil, buf, err
	}

	if int(length) > len(buf) {
		return nil, buf, RCMalformedPacket
	}

	return buf[:length], buf[length:], nil
}

func decodeIdentifier(buf []byte) (Identifier, []byte, error) {
	i, b, e := decodeByte(buf)
	if e != nil {
		return 0, b, e
	}
	return Identifier(i), b, nil
}

// decodeByte extracts the value of a byte from a byte array.
func decodeByte(buf []byte) (byte, []byte, error) {
	if len(buf) <= 0 {
		return 0, buf, RCMalformedPacket
	}
	return buf[0], buf[1:], nil
}

// decodeRCode extracts the value of a byte from a byte array.
func decodeRCode(buf []byte) (RCode, []byte, error) {
	if len(buf) <= 0 {
		return 0, buf, RCMalformedPacket
	}
	return RCode(buf[0]), buf[1:], nil
}

// decodeQos extracts the value of a byte from a byte array.
func decodeQos(buf []byte) (QoS, []byte, error) {
	if len(buf) <= 0 {
		return 0, buf, RCMalformedPacket
	}
	return QoS(buf[0]), buf[1:], nil
}

// decodeBool extracts the value of a byte from a byte array.
func decodeBool(buf []byte) (bool, []byte, error) {
	var b byte
	var err error
	b, buf, err = decodeByte(buf)
	if err != nil {
		return false, buf, err
	}
	return b > 0, buf, nil
}

// decodeUint16 extracts the value of two bytes from a byte array.
func decodeUint16(buf []byte) (uint16, []byte, error) {
	if len(buf) < 2 {
		return 0, buf, RCMalformedPacket
	}
	return binary.BigEndian.Uint16(buf[0:2]), buf[2:], nil
}

// decodeUint32 extracts the value of four bytes from a byte array.
func decodeUint32(buf []byte) (uint32, []byte, error) {
	if len(buf) < 4 {
		return 0, buf, RCMalformedPacket
	}
	return binary.BigEndian.Uint32(buf[:4]), buf[4:], nil
}

// encodeBool encodes a boolean value into a byte.
func encodeBool(b bool) byte {
	if b {
		return 1
	}
	return 0
}

func encodeString(val string) []byte {
	buf := make([]byte, 2, 32)
	binary.BigEndian.PutUint16(buf, uint16(len(val)))
	return append(buf, []byte(val)...)
}

// decodeString extracts a string from a byte array, beginning at an offset.
func decodeString(buf []byte) (string, []byte, error) {
	b, buf, err := decodeBytes(buf)
	if err != nil {
		return "", buf, err
	}

	if !validUTF8(b) { // [MQTT-1.5.4-1] [MQTT-3.1.3-5]
		return "", buf, RCMalformedPacket
	}

	return bytesToString(b), buf, nil
}

// decodeStringPair extracts a string from a byte array, beginning at an offset.
func decodeStringPair(buf []byte) (string, string, []byte, error) {
	var k, v []byte
	var err error
	k, buf, err = decodeBytes(buf)
	if err != nil {
		return "", "", buf, err
	}

	if !validUTF8(k) { // [MQTT-1.5.4-1] [MQTT-3.1.3-5]
		return "", "", buf, RCMalformedPacket
	}

	v, buf, err = decodeBytes(buf)
	if err != nil {
		return "", "", buf, err
	}

	if !validUTF8(v) { // [MQTT-1.5.4-1] [MQTT-3.1.3-5]
		return "", "", buf, RCMalformedPacket
	}

	return bytesToString(k), bytesToString(v), buf, nil
}

// encodeBytes encodes a byte array to a byte array. Used primarily for message payloads.
func encodeBytes(val []byte) []byte {
	buf := make([]byte, 2+len(val))
	binary.BigEndian.PutUint16(buf, uint16(len(val)))
	copy(buf[2:], val)
	return buf
}

// encodeUint16 encodes a uint16 to a byte array.
func encodeUint16(val uint16) []byte {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, val)
	return buf
}

// encodeUint32 encode a uint32 to a byte array.
func encodeUint32(val uint32) []byte {
	buf := make([]byte, 4)
	binary.BigEndian.PutUint32(buf, val)
	return buf
}

// bytesToString converts a byte slice to a string without allocating new memory.
func bytesToString(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

func decodeLength(buf []byte) (uint32, []byte, error) {
	return decodeVaruint(buf)
}

func decodeVaruint(buf []byte) (uint32, []byte, error) {
	var value uint32
	var multiplier uint32 = 1
	var encodedByte byte
	var i int

	for {
		if i >= len(buf) {
			return 0, buf, RCMalformedPacket
		}
		encodedByte = buf[i]
		i++
		value += uint32(encodedByte&127) * multiplier
		if multiplier > 128*128*128 {
			return 0, buf, RCMalformedPacket
		}
		multiplier *= 128
		if (encodedByte & 128) == 0 {
			break
		}
	}

	return value, buf[i:], nil
}

func decodeVarint(buf []byte) (int32, []byte, error) {
	var value int32
	var multiplier int32 = 1
	var encodedByte byte
	var i int

	for {
		if i >= len(buf) {
			return 0, buf, RCMalformedPacket
		}
		encodedByte = buf[i]
		i++
		value += int32(encodedByte&127) * multiplier
		if multiplier > 128*128*128 {
			return 0, buf, RCMalformedPacket
		}
		multiplier *= 128
		if (encodedByte & 128) == 0 {
			break
		}
	}

	return value, buf[i:], nil
}

func encodeLength(val uint32) []byte {
	return encodeVaruint(val)
}

func encodeVaruint(val uint32) []byte {
	var buf []byte
	for {
		encodedByte := byte(val % 128)
		val /= 128
		if val > 0 {
			encodedByte |= 128
		}
		buf = append(buf, encodedByte)
		if val == 0 {
			break
		}
	}
	return buf
}

func encodeVarint(val int32) []byte {
	var buf []byte
	for {
		encodedByte := byte(val % 128)
		val /= 128
		if val > 0 {
			encodedByte |= 128
		}
		buf = append(buf, encodedByte)
		if val == 0 {
			break
		}
	}
	return buf
}
