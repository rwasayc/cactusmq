package packet

import (
	"bytes"
	"reflect"
	"testing"
)

type PubAckCodecTestcase struct {
	Name      string
	EncodeVer ProtocolVersion

	// data
	Request      *PublishAcknowledgement
	RequestBytes []byte
}

func TestPubAck(t *testing.T) {
	encodeRunner := func(t *testing.T, tc PubAckCodecTestcase) bool {
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

	decodeRunner := func(t *testing.T, tc PubAckCodecTestcase) bool {
		request := &PublishAcknowledgement{}
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

	for _, tc := range PubAckCodecTestcases {
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

var PubAckCodecTestcases = []PubAckCodecTestcase{
	{
		Name:      "basic",
		EncodeVer: ProtoVer5,
		Request: &PublishAcknowledgement{
			PacketID:   1,
			ReasonCode: RCNormalDisconnection,
			Properties: BaseProperties{
				ReasonString: "reason string",
				UserProperty: []*UserProperty{
					{Key: "user1", Val: "value1"},
					{Key: "user2", Val: "value2"},
				},
			},
		},
		RequestBytes: []byte{
			0, 1, // Packet ID
			byte(RCNormalDisconnection),                                            // Reason Code
			48,                                                                     // Properties Length
			byte(IDReasonString),                                                   // Reason String ID
			0, 13, 'r', 'e', 'a', 's', 'o', 'n', ' ', 's', 't', 'r', 'i', 'n', 'g', // Reason String Value
			byte(IDUserProperty),
			0, 5, 'u', 's', 'e', 'r', '1', 0, 6, 'v', 'a', 'l', 'u', 'e', '1', // User Property 1
			byte(IDUserProperty),
			0, 5, 'u', 's', 'e', 'r', '2', 0, 6, 'v', 'a', 'l', 'u', 'e', '2', // User Property 2
		},
	},
}
