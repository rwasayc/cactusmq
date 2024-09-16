package packet

import (
	"bytes"
	"reflect"
	"testing"
)

type ConnackCodecTestcase struct {
	Name      string
	EncodeVer ProtocolVersion

	// data
	Request      *ConnectAcknowledgement
	RequestBytes []byte
}

func TestConnack(t *testing.T) {
	encodeRunner := func(t *testing.T, tc ConnackCodecTestcase) bool {
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

	decodeRunner := func(t *testing.T, tc ConnackCodecTestcase) bool {
		request := &ConnectAcknowledgement{}
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

	for _, tc := range ConnackCodecTestcases {
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

var ConnackCodecTestcases = []ConnackCodecTestcase{
	{
		Name:      "basic",
		EncodeVer: ProtoVer5,
		Request: &ConnectAcknowledgement{
			SessionPresent:    true,
			ConnectReasonCode: RCNormalDisconnection,
			Properties: &ConnectAcknowledgementProperties{
				SessionExpiryInterval:    10,
				ReceiveMaximum:           100,
				MaximumQoS:               QoS1,
				RetainAvailable:          1,
				MaximumPacketSize:        1024,
				AssignedClientIdentifier: "id1",
				TopicAliasMaximum:        10,
				ReasonString:             "success",
				UserProperty: []*UserProperty{
					{
						Key: "user1",
						Val: "value1",
					},
					{
						Key: "user2",
						Val: "value2",
					},
				},
				WildcardSubscriptionAvailable:   1,
				SubscriptionIdentifierAvailable: 2,
				SharedSubscriptionAvailable:     3,
				ServerKeepAlive:                 10,
				ResponseInformation:             "resp1",
				ServerReference:                 "ref1",
				AuthenticationMethod:            "auth1",
				AuthenticationData:              []byte("data1"),
			},
		},
		RequestBytes: []byte{
			1,                             // Session Present
			1,                             // Connect Reason Code
			108,                           // Properties Length
			byte(IDSessionExpiryInterval), // Session Expiry Interval ID
			0, 0, 0, 10,                   // Session Expiry Interval Value
			byte(IDReceiveMaximum), // Receive Maximum ID
			0, 100,                 // Receive Maximum Value
			byte(IDMaximumQoS),        // Maximum QoS ID
			byte(QoS1),                // Maximum QoS Value
			byte(IDRetainAvailable),   // Retain Available ID
			1,                         // Retain Available Value
			byte(IDMaximumPacketSize), // Maximum Packet Size ID
			0, 0, 0b00000100, 0,       // Maximum Packet Size Value
			byte(IDAssignedClientID), // Assigned Client ID ID
			0, 3, 'i', 'd', '1',      // Assigned Client ID Value
			byte(IDTopicAliasMaximum), // Topic Alias Maximum ID
			0, 10,                     // Topic Alias Maximum Value
			byte(IDReasonString),                    // Reason String ID
			0, 7, 's', 'u', 'c', 'c', 'e', 's', 's', // Reason String Value
			byte(IDUserProperty),                                              // User Property ID
			0, 5, 'u', 's', 'e', 'r', '1', 0, 6, 'v', 'a', 'l', 'u', 'e', '1', // User Property Value
			byte(IDUserProperty),                                              // User Property ID
			0, 5, 'u', 's', 'e', 'r', '2', 0, 6, 'v', 'a', 'l', 'u', 'e', '2', // User Property Value
			byte(IDWildcardSubAvailable), // Wildcard Subscription Available ID
			1,                            // Wildcard Subscription Available Value
			byte(IDSubIDAvailable),       // Subscription Identifier Available ID
			2,                            // Subscription Identifier Available Value
			byte(IDSharedSubAvailable),   // Shared Subscription Available ID
			3,                            // Shared Subscription Available Value
			byte(IDServerKeepAlive),      // Server Keep Alive ID
			0, 10,                        // Server Keep Alive Value
			byte(IDResponseInformation),   // Response Information ID
			0, 5, 'r', 'e', 's', 'p', '1', // Response Information Value
			byte(IDServerReference),  // Server Reference ID
			0, 4, 'r', 'e', 'f', '1', // Server Reference Value
			byte(IDAuthenticationMethod),  // Authentication Method ID
			0, 5, 'a', 'u', 't', 'h', '1', // Authentication Method Value
			byte(IDAuthenticationData),    // Authentication Data ID
			0, 5, 'd', 'a', 't', 'a', '1', // Authentication Data Value
		},
	},
}
