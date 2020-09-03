package broker

type Broker interface {
	Connect() error
	Init(...Option) error
	Options() Options
	Disconnect() error
	Publish(topic string, m *Message, opts ...PublishOption) error
	Subscribe(topic string, h Handler, opt ...SubscribeOption) (Subscriber, error)
	String() string
}

type Handler func(Event) error

type Message struct {
	Data []byte
}

type Event interface {
	Topic() string
	Message() *Message
	Ack() error
	Error() error
}

type Subscriber interface {
	Topic() string
	Unsubscribe() error
}

var (
	DefaultBroker Broker = NewBroker()
)

func Init(opts ...Option) error {
	return DefaultBroker.Init(opts...)
}

func Connect() error {
	return DefaultBroker.Connect()
}

func Disconnect() error {
	return DefaultBroker.Disconnect()
}

func Publish(topic string, msg *Message, opts ...PublishOption) error {
	return DefaultBroker.Publish(topic, msg, opts...)
}

func Subscribe(topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	return DefaultBroker.Subscribe(topic, handler, opts...)
}

func String() string {
	return DefaultBroker.String()
}
