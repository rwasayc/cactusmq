package packet

import (
	"bytes"
	"fmt"
)

type ConnectAcknowledgement struct {
	SessionPresent    bool                              `json:"session_present"`
	ConnectReasonCode RCode                             `json:"connect_reason_code"`
	Properties        *ConnectAcknowledgementProperties `json:"properties,omitempty"`
}

type ConnectAcknowledgementProperties struct {
	SessionExpiryInterval           uint32          `json:"session_expiry_interval,omitempty"`
	ReceiveMaximum                  uint16          `json:"receive_maximum,omitempty"`
	MaximumQoS                      QoS             `json:"maximum_qos,omitempty"`
	RetainAvailable                 uint8           `json:"retain_available,omitempty"`
	MaximumPacketSize               uint32          `json:"maximum_packet_size,omitempty"`
	AssignedClientIdentifier        string          `json:"assigned_client_identifier,omitempty"`
	TopicAliasMaximum               uint16          `json:"topic_alias_maximum,omitempty"`
	ReasonString                    string          `json:"reason_string,omitempty"`
	UserProperty                    []*UserProperty `json:"user_property,omitempty"`
	WildcardSubscriptionAvailable   uint8           `json:"wildcard_subscription_available,omitempty"`
	SubscriptionIdentifierAvailable uint8           `json:"subscription_identifier_available,omitempty"`
	SharedSubscriptionAvailable     uint8           `json:"shared_subscription_available,omitempty"`
	ServerKeepAlive                 uint16          `json:"server_keep_alive,omitempty"`
	ResponseInformation             string          `json:"response_information,omitempty"`
	ServerReference                 string          `json:"server_reference,omitempty"`
	AuthenticationMethod            string          `json:"authentication_method,omitempty"`
	AuthenticationData              []byte          `json:"authentication_data,omitempty"`
}

func (cap *ConnectAcknowledgementProperties) Encode(buf *bytes.Buffer) error {
	var err error
	tmpBuf := bytes.NewBuffer(nil)

	if cap.SessionExpiryInterval != 0 {
		err = tmpBuf.WriteByte(byte(IDSessionExpiryInterval))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeUint32(cap.SessionExpiryInterval))
		if err != nil {
			return err
		}
	}
	if cap.ReceiveMaximum != 0 {
		err = tmpBuf.WriteByte(byte(IDReceiveMaximum))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeUint16(cap.ReceiveMaximum))
		if err != nil {
			return err
		}
	}
	if cap.MaximumQoS != 0 {
		err = tmpBuf.WriteByte(byte(IDMaximumQoS))
		if err != nil {
			return err
		}
		err = tmpBuf.WriteByte(byte(cap.MaximumQoS))
		if err != nil {
			return err
		}
	}
	if cap.RetainAvailable != 0 {
		err = tmpBuf.WriteByte(byte(IDRetainAvailable))
		if err != nil {
			return err
		}
		err = tmpBuf.WriteByte(cap.RetainAvailable)
		if err != nil {
			return err
		}
	}
	if cap.MaximumPacketSize != 0 {
		err = tmpBuf.WriteByte(byte(IDMaximumPacketSize))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeUint32(cap.MaximumPacketSize))
		if err != nil {
			return err
		}
	}
	if cap.AssignedClientIdentifier != "" {
		err = tmpBuf.WriteByte(byte(IDAssignedClientID))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeString(cap.AssignedClientIdentifier))
		if err != nil {
			return err
		}
	}
	if cap.TopicAliasMaximum != 0 {
		err = tmpBuf.WriteByte(byte(IDTopicAliasMaximum))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeUint16(cap.TopicAliasMaximum))
		if err != nil {
			return err
		}
	}
	if cap.ReasonString != "" {
		err = tmpBuf.WriteByte(byte(IDReasonString))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeString(cap.ReasonString))
		if err != nil {
			return err
		}
	}
	if cap.UserProperty != nil {
		for _, prop := range cap.UserProperty {
			err = tmpBuf.WriteByte(byte(IDUserProperty))
			if err != nil {
				return err
			}
			_, err = tmpBuf.Write(encodeString(prop.Key))
			if err != nil {
				return err
			}
			_, err = tmpBuf.Write(encodeString(prop.Val))
			if err != nil {
				return err
			}
		}
	}
	if cap.WildcardSubscriptionAvailable != 0 {
		err = tmpBuf.WriteByte(byte(IDWildcardSubAvailable))
		if err != nil {
			return err
		}
		err = tmpBuf.WriteByte(cap.WildcardSubscriptionAvailable)
		if err != nil {
			return err
		}
	}
	if cap.SubscriptionIdentifierAvailable != 0 {
		err = tmpBuf.WriteByte(byte(IDSubIDAvailable))
		if err != nil {
			return err
		}
		err = tmpBuf.WriteByte(cap.SubscriptionIdentifierAvailable)
		if err != nil {
			return err
		}
	}
	if cap.SharedSubscriptionAvailable != 0 {
		err = tmpBuf.WriteByte(byte(IDSharedSubAvailable))
		if err != nil {
			return err
		}
		err = tmpBuf.WriteByte(cap.SharedSubscriptionAvailable)
		if err != nil {
			return err
		}
	}
	if cap.ServerKeepAlive != 0 {
		err = tmpBuf.WriteByte(byte(IDServerKeepAlive))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeUint16(cap.ServerKeepAlive))
		if err != nil {
			return err
		}
	}
	if cap.ResponseInformation != "" {
		err = tmpBuf.WriteByte(byte(IDResponseInformation))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeString(cap.ResponseInformation))
		if err != nil {
			return err
		}
	}
	if cap.ServerReference != "" {
		err = tmpBuf.WriteByte(byte(IDServerReference))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeString(cap.ServerReference))
		if err != nil {
			return err
		}
	}
	if cap.AuthenticationMethod != "" {
		err = tmpBuf.WriteByte(byte(IDAuthenticationMethod))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeString(cap.AuthenticationMethod))
		if err != nil {
			return err
		}
	}
	if cap.AuthenticationData != nil {
		err = tmpBuf.WriteByte(byte(IDAuthenticationData))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(encodeUint16(uint16(len(cap.AuthenticationData))))
		if err != nil {
			return err
		}
		_, err = tmpBuf.Write(cap.AuthenticationData)
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

func (cap *ConnectAcknowledgementProperties) Decode(buf []byte) ([]byte, error) {
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
		case IDSessionExpiryInterval:
			cap.SessionExpiryInterval, buf, err = decodeUint32(buf)
		case IDReceiveMaximum:
			cap.ReceiveMaximum, buf, err = decodeUint16(buf)
		case IDMaximumQoS:
			cap.MaximumQoS, buf, err = decodeQos(buf)
		case IDRetainAvailable:
			cap.RetainAvailable, buf, err = decodeByte(buf)
		case IDMaximumPacketSize:
			cap.MaximumPacketSize, buf, err = decodeUint32(buf)
		case IDAssignedClientID:
			cap.AssignedClientIdentifier, buf, err = decodeString(buf)
		case IDTopicAliasMaximum:
			cap.TopicAliasMaximum, buf, err = decodeUint16(buf)
		case IDReasonString:
			cap.ReasonString, buf, err = decodeString(buf)
		case IDUserProperty:
			var key string
			key, buf, err = decodeString(buf)
			if err == nil {
				var value string
				value, buf, err = decodeString(buf)
				if err == nil {
					if cap.UserProperty == nil {
						cap.UserProperty = make([]*UserProperty, 0, 1)
					}
					cap.UserProperty = append(cap.UserProperty, &UserProperty{Key: key, Val: value})
				}
			}
		case IDWildcardSubAvailable:
			cap.WildcardSubscriptionAvailable, buf, err = decodeByte(buf)
		case IDSubIDAvailable:
			cap.SubscriptionIdentifierAvailable, buf, err = decodeByte(buf)
		case IDSharedSubAvailable:
			cap.SharedSubscriptionAvailable, buf, err = decodeByte(buf)
		case IDServerKeepAlive:
			cap.ServerKeepAlive, buf, err = decodeUint16(buf)
		case IDResponseInformation:
			cap.ResponseInformation, buf, err = decodeString(buf)
		case IDServerReference:
			cap.ServerReference, buf, err = decodeString(buf)
		case IDAuthenticationMethod:
			cap.AuthenticationMethod, buf, err = decodeString(buf)
		case IDAuthenticationData:
			cap.AuthenticationData, buf, err = decodeBytes(buf)
		default:
			err = fmt.Errorf("unknown property id: %d", id)
		}
		if err != nil {
			return buf, err
		}
	}

	return buf, nil
}

func (ca *ConnectAcknowledgement) Decode(buf []byte) error {
	var err error
	ca.SessionPresent, buf, err = decodeBool(buf)
	if err != nil {
		return err
	}
	ca.ConnectReasonCode, buf, err = decodeRCode(buf)
	if err != nil {
		return err
	}
	ca.Properties = &ConnectAcknowledgementProperties{}
	buf, err = ca.Properties.Decode(buf)
	if err != nil {
		return err
	}
	return nil
}

func (ca *ConnectAcknowledgement) Encode(ver ProtocolVersion, buf *bytes.Buffer) error {
	var err error
	if ca.SessionPresent {
		err = buf.WriteByte(1)
	} else {
		err = buf.WriteByte(0)
	}
	if err != nil {
		return err
	}
	err = buf.WriteByte(byte(ca.ConnectReasonCode))
	if err != nil {
		return err
	}
	if ca.Properties != nil {
		err = ca.Properties.Encode(buf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ca *ConnectAcknowledgement) Validate() RCode {
	return RCSuccess
}
