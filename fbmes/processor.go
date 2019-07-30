package fbmes

type Messaging struct {
}

type messagingProcessor interface {
	ProcessMessage(m Messaging)
}
