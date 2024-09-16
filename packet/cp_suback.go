package packet

import (
	"bytes"
	"fmt"
)

type SubscribeAcknowledgement struct {
	ClientID   string
	Properties SubscribeAcknowledgementProperties
	Payload    []*SubscribePayload
}

type SubscribeAcknowledgementProperties struct {
	UserProperty []*UserProperty
}

func (srp *SubscribeAcknowledgementProperties) Decode(buf []byte) ([]int, []byte, error) {
	var length uint32
	var err error
	var subscriptionIDs []int
	length, buf, err = decodeLength(buf)
	if err != nil {
		return subscriptionIDs, buf, err
	}
	if length == 0 {
		return subscriptionIDs, buf, nil
	}
	shouldRemain := len(buf) - int(length)
	for len(buf) > shouldRemain {
		var id Identifier
		id, buf, err = decodeIdentifier(buf)
		if err != nil {
			return subscriptionIDs, buf, err
		}
		switch id {
		case IDSubscriptionIdentifier:
			var subID int32
			subID, buf, err = decodeVarint(buf)
			if err != nil {
				return subscriptionIDs, buf, err
			}
			if subscriptionIDs == nil {
				subscriptionIDs = make([]int, 0)
			}
			subscriptionIDs = append(subscriptionIDs, int(subID))
		case IDUserProperty:
			var key string
			key, buf, err = decodeString(buf)
			if err == nil {
				var value string
				value, buf, err = decodeString(buf)
				if err == nil {
					if srp.UserProperty == nil {
						srp.UserProperty = make([]*UserProperty, 0, 1)
					}
					srp.UserProperty = append(srp.UserProperty, &UserProperty{Key: key, Val: value})
				}
			}
		default:
			err = fmt.Errorf("unknown identifier: %d", id)
		}
		if err != nil {
			return subscriptionIDs, buf, err
		}
	}
	return subscriptionIDs, buf, nil
}

func (srp *SubscribeAcknowledgementProperties) Encode(buf *bytes.Buffer, ids []int) error {
	var err error
	tmpBuf := bytes.NewBuffer(nil)
	if len(ids) > 0 {
		for _, id := range ids {
			err = tmpBuf.WriteByte(byte(IDSubscriptionIdentifier))
			if err != nil {
				return err
			}
			_, err = tmpBuf.Write(encodeVarint(int32(id)))
			if err != nil {
				return err
			}
		}
	}

	for _, userProperty := range srp.UserProperty {
		err = tmpBuf.WriteByte(byte(IDUserProperty))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeString(userProperty.Key))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeString(userProperty.Val))
		if err != nil {
			return err
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

func (sr *SubscribeAcknowledgement) Decode(buf []byte) error {
	var err error
	sr.ClientID, buf, err = decodeString(buf)
	if err != nil {
		return err
	}
	var subscriptionIDs []int
	subscriptionIDs, buf, err = sr.Properties.Decode(buf)
	if err != nil {
		return err
	}
	sr.Payload = make([]*SubscribePayload, 0)
	idx := 0
	for len(buf) > 0 {
		var payload = &SubscribePayload{}
		buf, err = payload.Decode(subscriptionIDs[idx], buf)
		if err != nil {
			return err
		}
		sr.Payload = append(sr.Payload, payload)
		idx++
	}
	return nil
}

func (sr *SubscribeAcknowledgement) Encode(ver ProtocolVersion, buf *bytes.Buffer) error {
	var err error
	_, err = buf.Write(encodeString(sr.ClientID))
	if err != nil {
		return err
	}
	ids := extractValues(sr.Payload, func(payload *SubscribePayload) int {
		return payload.SubscriptionID
	})
	err = sr.Properties.Encode(buf, ids)
	if err != nil {
		return err
	}
	for _, payload := range sr.Payload {
		err = payload.Encode(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sr *SubscribeAcknowledgement) Validate() RCode {
	return RCSuccess
}
