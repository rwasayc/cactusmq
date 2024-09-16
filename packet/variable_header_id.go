package packet

type Identifier byte

const (
	IDPayloadFormatIndicator     Identifier = 0x01
	IDMessageExpiryInterval      Identifier = 0x02
	IDContentType                Identifier = 0x03
	IDResponseTopic              Identifier = 0x08
	IDCorrelationData            Identifier = 0x09
	IDSubscriptionIdentifier     Identifier = 0x0B
	IDSessionExpiryInterval      Identifier = 0x11
	IDAssignedClientID           Identifier = 0x12
	IDServerKeepAlive            Identifier = 0x13
	IDAuthenticationMethod       Identifier = 0x15
	IDAuthenticationData         Identifier = 0x16
	IDRequestProblemInformation  Identifier = 0x17
	IDWillDelayInterval          Identifier = 0x18
	IDRequestResponseInformation Identifier = 0x19
	IDResponseInformation        Identifier = 0x1A
	IDServerReference            Identifier = 0x1C
	IDReasonString               Identifier = 0x1F
	IDReceiveMaximum             Identifier = 0x21
	IDTopicAliasMaximum          Identifier = 0x22
	IDTopicAlias                 Identifier = 0x23
	IDMaximumQoS                 Identifier = 0x24
	IDRetainAvailable            Identifier = 0x25
	IDUserProperty               Identifier = 0x26
	IDMaximumPacketSize          Identifier = 0x27
	IDWildcardSubAvailable       Identifier = 0x28
	IDSubIDAvailable             Identifier = 0x29
	IDSharedSubAvailable         Identifier = 0x2A
)
