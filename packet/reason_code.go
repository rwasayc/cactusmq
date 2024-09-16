package packet

import "fmt"

// ToDo 完成修改
type RCode byte

func (r RCode) String() string {
	return fmt.Sprintf("RCode{code: %x, reason: %s}", byte(r), r.reason())
}

func (c RCode) Error() string {
	return c.String()
}

func (c RCode) Is(err error) bool {
	v, ok := err.(RCode)
	if !ok {
		return false
	}
	return c == v
}

var rcode2reason = map[RCode]string{
	RCSuccess:                             "Success",
	RCNormalDisconnection:                 "Normal disconnection",
	RCGrantedQoS1:                         "Granted QoS 1",
	RCGrantedQoS2:                         "Granted QoS 2",
	RCDisconnectWithWill:                  "Disconnect with Will Message",
	RCNoMatchingSubscribers:               "No matching subscribers",
	RCNoSubscriptionExisted:               "No subscription existed",
	RCContinueAuthentication:              "Continue authentication",
	RCReAuthenticate:                      "Re-authenticate",
	RCUnspecifiedError:                    "Unspecified error",
	RCMalformedPacket:                     "Malformed Packet",
	RCProtocolError:                       "Protocol Error",
	RCImplementationSpecific:              "Implementation specific error",
	RCUnsupportedProtocol:                 "Unsupported Protocol Version",
	RClientIDNotValid:                     "Client Identifier not valid",
	RCBadUsernameOrPassword:               "Bad User Name or Password",
	RCNotAuthorized:                       "Not authorized",
	RCServerUnavailable:                   "Server unavailable",
	RCServerBusy:                          "Server busy",
	RCBanned:                              "Banned",
	RCServerShuttingDown:                  "Server shutting down",
	RCBadAuthenticationMethod:             "Bad authentication method",
	RCKeepAliveTimeout:                    "Keep Alive timeout",
	RCSessionTakenOver:                    "Session taken over",
	RCTopicFilterInvalid:                  "Topic Filter invalid",
	RCTopicNameInvalid:                    "Topic Name invalid",
	RCPacketIDInUse:                       "Packet Identifier in use",
	RCPacketIDNotFound:                    "Packet Identifier not found",
	RCReceiveMaximumExceeded:              "Receive Maximum exceeded",
	RCTopicAliasInvalid:                   "Topic Alias invalid",
	RCPacketTooLarge:                      "Packet too large",
	RCMessageRateTooHigh:                  "Message rate too high",
	RCQuotaExceeded:                       "Quota exceeded",
	RCAdministrativeAction:                "Administrative action",
	RCPayloadFormatInvalid:                "Payload format invalid",
	RCRetainNotSupported:                  "Retain not supported",
	RCQoSNotSupported:                     "QoS not supported",
	RCUseAnotherServer:                    "Use another server",
	RCServerMoved:                         "Server moved",
	RCSharedSubscriptionsNotSupported:     "Shared Subscriptions not supported",
	RCConnectionRateExceeded:              "Connection rate exceeded",
	RCMaximumConnectTime:                  "Maximum connect time",
	RCSubscriptionIdentifiersNotSupported: "Subscription Identifiers not supported",
	RCWildcardSubscriptionsNotSupported:   "Wildcard Subscriptions not supported",
}

func (r RCode) reason() string {
	if reason, ok := rcode2reason[r]; ok {
		return reason
	}
	return "unknown reason code"
}

const (
	RCSuccess                             = RCode(0x00)
	RCNormalDisconnection                 = RCode(0x01)
	RCGrantedQoS1                         = RCode(0x02)
	RCGrantedQoS2                         = RCode(0x03)
	RCDisconnectWithWill                  = RCode(0x04)
	RCNoMatchingSubscribers               = RCode(0x10)
	RCNoSubscriptionExisted               = RCode(0x11)
	RCContinueAuthentication              = RCode(0x18)
	RCReAuthenticate                      = RCode(0x19)
	RCUnspecifiedError                    = RCode(0x80)
	RCMalformedPacket                     = RCode(0x81)
	RCMalformedPacketProtocolVersion      = RCode(0x81)
	RCMalformedPacketProtocolName         = RCode(0x81)
	RCProtocolError                       = RCode(0x82)
	RCImplementationSpecific              = RCode(0x83)
	RCUnsupportedProtocol                 = RCode(0x84)
	RClientIDNotValid                     = RCode(0x85)
	RCBadUsernameOrPassword               = RCode(0x86)
	RCNotAuthorized                       = RCode(0x87)
	RCServerUnavailable                   = RCode(0x88)
	RCServerBusy                          = RCode(0x89)
	RCBanned                              = RCode(0x8A)
	RCServerShuttingDown                  = RCode(0x8B)
	RCBadAuthenticationMethod             = RCode(0x8C)
	RCKeepAliveTimeout                    = RCode(0x8D)
	RCSessionTakenOver                    = RCode(0x8E)
	RCTopicFilterInvalid                  = RCode(0x8F)
	RCTopicNameInvalid                    = RCode(0x90)
	RCPacketIDInUse                       = RCode(0x91)
	RCPacketIDNotFound                    = RCode(0x92)
	RCReceiveMaximumExceeded              = RCode(0x93)
	RCTopicAliasInvalid                   = RCode(0x94)
	RCPacketTooLarge                      = RCode(0x95)
	RCMessageRateTooHigh                  = RCode(0x96)
	RCQuotaExceeded                       = RCode(0x97)
	RCAdministrativeAction                = RCode(0x98)
	RCPayloadFormatInvalid                = RCode(0x99)
	RCRetainNotSupported                  = RCode(0x9A)
	RCQoSNotSupported                     = RCode(0x9B)
	RCUseAnotherServer                    = RCode(0x9C)
	RCServerMoved                         = RCode(0x9D)
	RCSharedSubscriptionsNotSupported     = RCode(0x9E)
	RCConnectionRateExceeded              = RCode(0x9F)
	RCMaximumConnectTime                  = RCode(0xA0)
	RCSubscriptionIdentifiersNotSupported = RCode(0xA1)
	RCWildcardSubscriptionsNotSupported   = RCode(0xA2)
)
