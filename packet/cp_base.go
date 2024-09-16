package packet

import (
	"bytes"
	"fmt"
)

type BaseProperties struct {
	ReasonString string
	UserProperty []*UserProperty
}

func (pa *BaseProperties) Encode(buf *bytes.Buffer) error {
	var err error
	tmpBuf := bytes.NewBuffer(nil)
	if pa.ReasonString != "" {
		err = tmpBuf.WriteByte(byte(IDReasonString))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeString(pa.ReasonString))
		if err != nil {
			return err
		}
	}
	if len(pa.UserProperty) > 0 {
		for _, up := range pa.UserProperty {
			err = tmpBuf.WriteByte(byte(IDUserProperty))
			if err != nil {
				return err
			}
			_, err = tmpBuf.Write(encodeString(up.Key))
			if err != nil {
				return err
			}
			_, err = tmpBuf.Write(encodeString(up.Val))
			if err != nil {
				return err
			}
		}
	}

	_, err = buf.Write(encodeLength(uint32(tmpBuf.Len())))
	if err != nil {
		return err
	}
	if tmpBuf.Len() > 0 {
		_, err = buf.Write(tmpBuf.Bytes())
		if err != nil {
			return err
		}
	}
	return nil

}
func (pa *BaseProperties) Decode(buf []byte) ([]byte, error) {
	var length uint32
	var err error
	length, buf, err = decodeLength(buf)
	if err != nil {
		return buf, err
	}
	if length == 0 {
		return buf, nil
	}
	shouldRemain := len(buf) - int(length)
	for len(buf) > shouldRemain {
		var id Identifier
		id, buf, err = decodeIdentifier(buf)
		if err != nil {
			return buf, err
		}
		switch id {
		case IDReasonString:
			pa.ReasonString, buf, err = decodeString(buf)
		case IDUserProperty:
			var key string
			key, buf, err = decodeString(buf)
			if err == nil {
				var value string
				value, buf, err = decodeString(buf)
				if err == nil {
					if pa.UserProperty == nil {
						pa.UserProperty = make([]*UserProperty, 0, 1)
					}
					pa.UserProperty = append(pa.UserProperty, &UserProperty{Key: key, Val: value})
				}
			}
		default:
			err = fmt.Errorf("unknown identifier: %d", id)
		}
		if err != nil {
			return buf, err
		}
	}

	return buf, nil
}
