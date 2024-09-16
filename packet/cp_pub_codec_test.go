package packet

import (
	"bytes"
	"reflect"
	"testing"
)

type PubCodecTestcase struct {
	Name      string
	EncodeVer ProtocolVersion

	// data
	Request      *PublishMessage
	RequestBytes []byte
}

func TestPub(t *testing.T) {
	encodeRunner := func(t *testing.T, tc PubCodecTestcase) bool {
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

	decodeRunner := func(t *testing.T, tc PubCodecTestcase) bool {
		request := &PublishMessage{}
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

	for _, tc := range PubCodecTestcases {
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

var PubCodecTestcases = []PubCodecTestcase{
	{
		Name:      "basic",
		EncodeVer: ProtoVer5,
		Request: &PublishMessage{
			TopicName: "topic",
			PacketID:  1,
			Payload:   []byte("payload"),
			Properties: PublishMessageProperties{
				PayloadFormatIndicator: 1,
				MessageExpiryInterval:  10,
				TopicAlias:             2,
				ContentType:            "text/plain",
				ResponseTopic:          "response",
				CorrelationData:        []byte{'c', 'o', 'r', 'r', 'e', 'l', 'a', 't', 'i', 'o', 'n', 'a', 't', 'i', 'o', 'n'},
				UserProperty: []*UserProperty{
					{Key: "user1", Val: "value1"},
					{Key: "user2", Val: "value2"},
				},
				SubscriptionIdentifier: []byte("subscription"),
			},
		},
		RequestBytes: []byte{
			0, 5, 't', 'o', 'p', 'i', 'c', // Topic Name
			0, 1, // Packet ID
			100, // Properties Length
			byte(IDPayloadFormatIndicator),
			1, // Payload Format Indicator
			byte(IDMessageExpiryInterval),
			0, 0, 0, 10, // Message Expiry Interval
			byte(IDTopicAlias),
			0, 2, // Topic Alias
			byte(IDResponseTopic),
			0, 8, 'r', 'e', 's', 'p', 'o', 'n', 's', 'e', // Response Topic
			byte(IDCorrelationData),
			0, 16, 'c', 'o', 'r', 'r', 'e', 'l', 'a', 't', 'i', 'o', 'n', 'a', 't', 'i', 'o', 'n', // Correlation Data
			byte(IDUserProperty),
			0, 5, 'u', 's', 'e', 'r', '1', 0, 6, 'v', 'a', 'l', 'u', 'e', '1', // User Property
			byte(IDUserProperty),
			0, 5, 'u', 's', 'e', 'r', '2', 0, 6, 'v', 'a', 'l', 'u', 'e', '2', // User Property
			byte(IDSubscriptionIdentifier),
			0, 12, 's', 'u', 'b', 's', 'c', 'r', 'i', 'p', 't', 'i', 'o', 'n', // Subscription Identifier
			byte(IDContentType),
			0, 10, 't', 'e', 'x', 't', '/', 'p', 'l', 'a', 'i', 'n', // Content Type
			0, 7, 'p', 'a', 'y', 'l', 'o', 'a', 'd', // Payload
		},
	},
}
