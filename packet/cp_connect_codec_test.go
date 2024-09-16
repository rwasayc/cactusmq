package packet

import (
	"bytes"
	"reflect"
	"testing"
)

type ConnectCodecTestcase struct {
	Name      string
	EncodeVer ProtocolVersion

	// data
	Request      *ConnectionRequest
	RequestBytes []byte
}

func TestConnect(t *testing.T) {
	encodeRunner := func(t *testing.T, tc ConnectCodecTestcase) bool {
		rcode := tc.Request.Validate()
		if rcode != RCSuccess {
			t.Errorf("expected \n%v\nbut got \n%v", RCSuccess, rcode)
			return false
		}
		buf := bytes.NewBuffer(nil)
		err := tc.Request.Encode(tc.EncodeVer, buf)
		if err != nil {
			t.Errorf("expected \n%v\nbut got \n%v", tc.RequestBytes, err)
			return false
		}
		if !bytes.Equal(buf.Bytes(), tc.RequestBytes) {
			t.Errorf("\nexpected \n%v\ngot \n%v", tc.RequestBytes, buf.Bytes())
			return false
		}
		return true
	}

	decodeRunner := func(t *testing.T, tc ConnectCodecTestcase) bool {
		request := &ConnectionRequest{}
		err := request.Decode(tc.RequestBytes)
		if err != nil {
			t.Errorf("decode expected \n%v\nbut got \n%v", JSON(tc.Request), err)
			return false
		}
		if !reflect.DeepEqual(tc.Request, request) {
			t.Errorf("decode expected \n%v\nbut got \n%v", JSON(tc.Request), JSON(request))
			return false
		}
		return true
	}

	for _, tc := range connectCodecTestcases {

		t.Run(tc.Name+" decode", func(t *testing.T) {
			if !decodeRunner(t, tc) {
				t.FailNow()
			}
		})
		t.Run(tc.Name+" encode", func(t *testing.T) {
			if !encodeRunner(t, tc) {
				t.FailNow()
			}
		})
	}
}

var connectCodecTestcases = []ConnectCodecTestcase{
	{
		Name:      "basic v3.1",
		EncodeVer: ProtoVer31,
		Request: &ConnectionRequest{
			ProtocolName:    FixedProtocolNameV31,
			ProtocolVersion: ProtoVer31,
			ClientID:        "id1",
			Keepalive:       10,
		},
		RequestBytes: []byte{
			0, 6, // Protocol Name - MSB+LSB
			'M', 'Q', 'I', 's', 'd', 'p', // Protocol Name
			3,     // Protocol Version
			0,     // Packet Flags
			0, 10, // Keepalive
			0, 3, // Client ID - MSB+LSB
			'i', 'd', '1', // Client ID "id1"
		},
	},
	{
		Name:      "basic v3.1.1",
		EncodeVer: ProtoVer311,
		Request: &ConnectionRequest{
			ProtocolName:    FixedProtocolNameV311,
			ProtocolVersion: ProtoVer311,
			ClientID:        "id1",
			Keepalive:       10,
		},
		RequestBytes: []byte{
			0, 4, // Protocol Name - MSB+LSB
			'M', 'Q', 'T', 'T', // Protocol Name
			4,     // Protocol Version
			0,     // Packet Flags
			0, 10, // Keepalive
			0, 3, // Client ID - MSB+LSB
			'i', 'd', '1', // Client ID "id1"
		},
	},
	{
		Name:      "basic v5",
		EncodeVer: ProtoVer5,
		Request: &ConnectionRequest{
			ProtocolName:    FixedProtocolNameV5,
			ProtocolVersion: ProtoVer5,
			ClientID:        "id1",
			Keepalive:       10,
		},
		RequestBytes: []byte{
			0, 4, // Protocol Name - MSB+LSB
			'M', 'Q', 'T', 'T', // Protocol Name
			5,     // Protocol Version
			0,     // Packet Flags
			0, 10, // Keepalive
			0,    // Properties Length - MSB+LSB
			0, 3, // Client ID - MSB+LSB
			'i', 'd', '1', // Client ID "id1"
		},
	},
	{
		Name:      "basic v5 with username&password",
		EncodeVer: ProtoVer5,
		Request: &ConnectionRequest{
			ProtocolName:    FixedProtocolNameV5,
			ProtocolVersion: ProtoVer5,
			ClientID:        "id1",
			Keepalive:       10,
			Password:        NewSPassword("password"),
			Username:        NewFlagV([]byte("username")),
		},
		RequestBytes: []byte{
			0, 4, // Protocol Name - MSB+LSB
			'M', 'Q', 'T', 'T', // Protocol Name
			5,          // Protocol Version
			0b11000000, // Packet Flags
			0, 10,      // Keepalive
			0,    // Properties Length - MSB+LSB
			0, 3, // Client ID - MSB+LSB
			'i', 'd', '1', // Client ID "id1"
			0, 8, // UserName - MSB+LSB
			'u', 's', 'e', 'r', 'n', 'a', 'm', 'e', // UserName "username"
			0, 8, // Password - MSB+LSB
			'p', 'a', 's', 's', 'w', 'o', 'r', 'd', // Password "password"
		},
	},
	{
		Name:      "basic v5 with full properties",
		EncodeVer: ProtoVer5,
		Request: &ConnectionRequest{
			ProtocolName:    FixedProtocolNameV5,
			ProtocolVersion: ProtoVer5,
			ClientID:        "id1",
			Keepalive:       10,
			Properties: &ConnectProperties{
				SessionExpiryInterval: NewFlagV[uint32](10),
				RequestProblemInfo:    NewFlagV[byte](1),
				ReceiveMaximum:        1024,
				RequestResponseInfo:   3,
				MaximumPacketSize:     1024,
				TopicAliasMaximum:     10,
				AuthenticationMethod:  "scma    ",
				AuthenticationData:    []byte("authdata"),
				UserProperty: []*UserProperty{
					{Key: "key", Val: "value"},
				},
			},
		},
		RequestBytes: []byte{
			0, 4, // Protocol Name - MSB+LSB
			'M', 'Q', 'T', 'T', // Protocol Name
			5,          // Protocol Version
			0b00000000, // Packet Flags
			0, 10,      // Keepalive

			55,                            // Properties Length - MSB+LSB
			byte(IDSessionExpiryInterval), // Session Expiry Interval ID
			0, 0, 0, 10,                   // Session Expiry Interval 10
			byte(IDRequestProblemInformation),  // Request Problem Information ID
			1,                                  // Request Problem Information 10
			byte(IDRequestResponseInformation), // Request Response Information ID
			3,                                  // Request Response Information 3
			byte(IDReceiveMaximum),             // Receive Maximum ID
			4, 0,                               // Receive Maximum 1024
			byte(IDMaximumPacketSize), // Maximum Packet Size ID
			0, 0, 4, 0,                // Maximum Packet Size 1024
			byte(IDTopicAliasMaximum), // Topic Alias Maximum ID
			0, 10,                     // Topic Alias Maximum 10
			byte(IDAuthenticationMethod),                 // Authentication Method ID
			0, 8, 's', 'c', 'm', 'a', ' ', ' ', ' ', ' ', // Authentication Method "scma"
			byte(IDAuthenticationData),                   // Authentication Data ID
			0, 8, 'a', 'u', 't', 'h', 'd', 'a', 't', 'a', // Authentication Data "authdata"
			byte(IDUserProperty), // User Property ID
			0, 3, 'k', 'e', 'y',  // User Property "key"
			0, 5, 'v', 'a', 'l', 'u', 'e', // User Property "value"
			0, 3, // Client ID - MSB+LSB
			'i', 'd', '1', // Client ID "id1"
		},
	},
	{
		Name:      "basic v5 with will",
		EncodeVer: ProtoVer5,
		Request: &ConnectionRequest{
			ProtocolName:    FixedProtocolNameV5,
			ProtocolVersion: ProtoVer5,
			ClientID:        "id1",
			Keepalive:       10,
			Will: FlagV[ConnectWill]{
				V: &ConnectWill{
					Topic:   "topic",
					Qos:     NewFlagV(QoS1),
					Retain:  true,
					Payload: []byte("payload"),
					Properties: &WillProperties{
						MessageExpiryInterval: 10,
						PayloadFormat:         NewFlagV[byte](1),
						ContentType:           "content_type",
						ResponseTopic:         "response_topic",
						CorrelationData:       []byte("correlation_data"),
						User: []*UserProperty{
							{Key: "key", Val: "value"},
						},
						WillDelayInterval: 10,
					},
				},
			},
		},
		RequestBytes: []byte{
			0, 4, // Protocol Name - MSB+LSB
			'M', 'Q', 'T', 'T', // Protocol Name
			5,     // Protocol Version
			44,    // Packet Flags
			0, 10, // Keepalive
			0,    // Properties Length - MSB+LSB
			0, 3, // Client ID - MSB+LSB
			'i', 'd', '1', // Client ID "id1"
			76,          // Will Properties Length - MSB+LSB
			2,           // IDMessageExpiryInterval ID
			0, 0, 0, 10, // IDMessageExpiryInterval Value
			1,     // IDPayloadFormatIndicator ID
			1,     // IDPayloadFormatIndicator Value
			3,     // IDContentType ID
			0, 12, // IDContentType Length
			'c', 'o', 'n', 't', 'e', 'n', 't', '_', 't', 'y', 'p', 'e', // IDContentType Value
			8,     // IDResponseTopic ID
			0, 14, // IDResponseTopic Length
			'r', 'e', 's', 'p', 'o', 'n', 's', 'e', '_', 't', 'o', 'p', 'i', 'c', // IDResponseTopic Value
			9,     // IDCorrelationData ID
			0, 16, // IDCorrelationData Length
			'c', 'o', 'r', 'r', 'e', 'l', 'a', 't', 'i', 'o', 'n', '_', 'd', 'a', 't', 'a', // IDCorrelationData Value
			38,   // IDUserProperty ID
			0, 3, // IDUserProperty Key Length
			'k', 'e', 'y', // IDUserProperty Key
			0, 5, // IDUserProperty Value Length
			'v', 'a', 'l', 'u', 'e', // IDUserProperty Value
			24,          // IDWillDelayInterval ID
			0, 0, 0, 10, // IDWillDelayInterval Value
			0, 5, // Will Topic Length
			116, 111, 112, 105, 99, // Will Topic "topic"
			0, 7, // Will Payload Length
			112, 97, 121, 108, 111, 97, 100, // Will Payload "payload"
		},
	},
}
