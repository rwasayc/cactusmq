package packet

import (
	"bytes"
	"fmt"
)

type PublishMessage struct {
	// fixed header
	DUP      bool
	QoSLevel QoS
	Retain   bool
	// variable header
	TopicName  string
	PacketID   uint16
	Properties PublishMessageProperties
	// payload
	Payload []byte
}

type PublishMessageProperties struct {
	PayloadFormatIndicator PayloadFormatIndicator `json:"payload_format_indicator"`
	MessageExpiryInterval  uint32                 `json:"message_expiry_interval"`
	TopicAlias             uint16                 `json:"topic_alias"`
	ResponseTopic          string                 `json:"response_topic"`
	CorrelationData        []byte                 `json:"correlation_data"`
	UserProperty           []*UserProperty        `json:"user_property"`
	SubscriptionIdentifier []byte                 `json:"subscription_identifier"`
	ContentType            string                 `json:"content_type"`
}

func (pmp *PublishMessageProperties) Encode(buf *bytes.Buffer) error {
	var err error
	tmpBuf := bytes.NewBuffer(nil)

	if pmp.PayloadFormatIndicator != 0 {
		err = tmpBuf.WriteByte(byte(IDPayloadFormatIndicator))
		if err != nil {
			return err
		}
		err = tmpBuf.WriteByte(byte(pmp.PayloadFormatIndicator))
		if err != nil {
			return err
		}
	}
	if pmp.MessageExpiryInterval != 0 {
		err = tmpBuf.WriteByte(byte(IDMessageExpiryInterval))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeUint32(pmp.MessageExpiryInterval))
		if err != nil {
			return err
		}
	}
	if pmp.TopicAlias != 0 {
		err = tmpBuf.WriteByte(byte(IDTopicAlias))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeUint16(pmp.TopicAlias))
		if err != nil {
			return err
		}
	}
	if pmp.ResponseTopic != "" {
		err = tmpBuf.WriteByte(byte(IDResponseTopic))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeString(pmp.ResponseTopic))
		if err != nil {
			return err
		}
	}
	if pmp.CorrelationData != nil {
		err = tmpBuf.WriteByte(byte(IDCorrelationData))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeBytes(pmp.CorrelationData))
		if err != nil {
			return err
		}
	}
	if pmp.UserProperty != nil {
		// pmp.UserProperty
		for _, up := range pmp.UserProperty {
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
	if pmp.SubscriptionIdentifier != nil {
		err = tmpBuf.WriteByte(byte(IDSubscriptionIdentifier))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeBytes(pmp.SubscriptionIdentifier))
		if err != nil {
			return err
		}
	}
	if pmp.ContentType != "" {
		err = tmpBuf.WriteByte(byte(IDContentType))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeString(pmp.ContentType))
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

func (pmp *PublishMessageProperties) Decode(buf []byte) ([]byte, error) {
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
		case IDPayloadFormatIndicator:
			var b byte
			b, buf, err = decodeByte(buf)
			if err == nil {
				pmp.PayloadFormatIndicator = PayloadFormatIndicator(b)
			}
		case IDMessageExpiryInterval:
			pmp.MessageExpiryInterval, buf, err = decodeUint32(buf)
		case IDTopicAlias:
			pmp.TopicAlias, buf, err = decodeUint16(buf)
		case IDResponseTopic:
			pmp.ResponseTopic, buf, err = decodeString(buf)
		case IDCorrelationData:
			pmp.CorrelationData, buf, err = decodeBytes(buf)
		case IDUserProperty:
			var key string
			key, buf, err = decodeString(buf)
			if err == nil {
				var value string
				value, buf, err = decodeString(buf)
				if err == nil {
					pmp.UserProperty = append(pmp.UserProperty, &UserProperty{Key: key, Val: value})
				}
			}
		case IDSubscriptionIdentifier:
			pmp.SubscriptionIdentifier, buf, err = decodeBytes(buf)
		case IDContentType:
			pmp.ContentType, buf, err = decodeString(buf)
		default:
			err = fmt.Errorf("unknown property: %d", id)
		}
		if err != nil {
			return buf, err
		}
	}

	return buf, nil
}

func (pm *PublishMessage) Decode(buf []byte) error {
	var err error
	pm.TopicName, buf, err = decodeString(buf)
	if err != nil {
		return err
	}
	pm.PacketID, buf, err = decodeUint16(buf)
	if err != nil {
		return err
	}
	buf, err = pm.Properties.Decode(buf)
	if err != nil {
		return err
	}
	pm.Payload, buf, err = decodeBytes(buf)
	if err != nil {
		return err
	}
	return nil
}

func (pm *PublishMessage) Encode(ver ProtocolVersion, buf *bytes.Buffer) error {
	var err error
	_, err = buf.Write(encodeString(pm.TopicName))
	if err != nil {
		return err
	}
	_, err = buf.Write(encodeUint16(pm.PacketID))
	if err != nil {
		return err
	}
	err = pm.Properties.Encode(buf)
	if err != nil {
		return err
	}
	_, err = buf.Write(encodeBytes(pm.Payload))
	if err != nil {
		return err
	}
	return nil
}

func (pm *PublishMessage) Validate() RCode {
	return RCSuccess
}
