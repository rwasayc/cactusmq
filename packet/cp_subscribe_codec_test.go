package packet

import (
	"bytes"
	"reflect"
	"testing"
)

type SubscribeCodecTestcase struct {
	Name      string
	EncodeVer ProtocolVersion

	// data
	Request      *SubscribeRequest
	RequestBytes []byte
}

func TestSubscribe(t *testing.T) {
	encodeRunner := func(t *testing.T, tc SubscribeCodecTestcase) bool {
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

	decodeRunner := func(t *testing.T, tc SubscribeCodecTestcase) bool {
		request := &SubscribeRequest{}
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

	for _, tc := range SubscribeCodecTestcases {
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

var SubscribeCodecTestcases = []SubscribeCodecTestcase{
	{
		Name:      "basic",
		EncodeVer: ProtoVer5,
		Request: &SubscribeRequest{
			ClientID: "clientID",
			Properties: SubscribeRequestProperties{
				UserProperty: []*UserProperty{
					{Key: "user1", Val: "value1"},
					{Key: "user2", Val: "value2"},
				},
			},
			Payload: []*SubscribePayload{
				{
					TopicFilter:       "topic1",
					QoS:               QoS1,
					NoLocal:           true,
					RetainAsPublished: true,
					RetainHandling:    RetainHandling(1),
					SubscriptionID:    1,
				},
			},
		},
		RequestBytes: []byte{
			0, 8, // Client ID
			'c', 'l', 'i', 'e', 'n', 't', 'I', 'D',
			34,                                                                // properties length
			byte(IDSubscriptionIdentifier),                                    // Subscription Identifier ID
			1,                                                                 // Subscription Identifier Value
			byte(IDUserProperty),                                              // User Property ID
			0, 5, 'u', 's', 'e', 'r', '1', 0, 6, 'v', 'a', 'l', 'u', 'e', '1', // User Property 1
			byte(IDUserProperty),                                              // User Property ID
			0, 5, 'u', 's', 'e', 'r', '2', 0, 6, 'v', 'a', 'l', 'u', 'e', '2', // User Property 2
			0, 6, //  Topic Filter Length
			't', 'o', 'p', 'i', 'c', '1',
			29, // flags
		},
	},
}
