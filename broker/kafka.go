package broker

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	cluster "github.com/bsm/sarama-cluster"

	"github.com/Shopify/sarama"
)

type kBroker struct {
	p     sarama.SyncProducer
	addrs []string

	connected bool
	scMutex   sync.RWMutex
	Context   context.Context
}

type subscriber struct {
	t        string
	consumer *cluster.Consumer
	opts     SubscribeOptions
}

type kafkaMessage struct {
	t   string
	err error
	m   *Message
}

func (p *kafkaMessage) Topic() string {
	return p.t
}

func (p *kafkaMessage) Message() *Message {
	return p.m
}

func (p *kafkaMessage) Ack() error {
	return nil
}

func (p *kafkaMessage) Error() error {
	return p.err
}

func (s *subscriber) Topic() string {
	return s.t
}

func (s *subscriber) Unsubscribe() error {
	return s.consumer.Close()
}

func (k *kBroker) Connect() error {
	if k.isConnected() {
		return nil
	}

	conf := sarama.NewConfig()
	conf.Producer.Return.Successes = true
	conf.Producer.Timeout = 5 * time.Second
	p, err := sarama.NewSyncProducer(k.addrs, conf)

	if err != nil {
		return err
	}

	k.scMutex.Lock()
	k.p = p
	k.connected = true
	k.scMutex.Unlock()

	return nil
}

func (k *kBroker) Disconnect() error {
	if !k.isConnected() {
		return nil
	}
	k.scMutex.Lock()
	defer k.scMutex.Unlock()
	k.p.Close()
	k.connected = false
	return nil
}

func (k *kBroker) isConnected() bool {
	k.scMutex.RLock()
	defer k.scMutex.RUnlock()
	return k.connected
}

func (k *kBroker) Init(opts ...Option) {
}

func (k *kBroker) Publish(topic string, msg *Message, opts ...PublishOption) error {
	if !k.isConnected() {
		return errors.New("[kafka] broker not connected")
	}

	msgLog := &sarama.ProducerMessage{
		Topic: string(topic),
		Value: sarama.StringEncoder(msg.Data),
	}

	_, _, err := k.p.SendMessage(msgLog)
	return err
}

func (k *kBroker) Subscribe(topic string, handler Handler, opts ...SubscribeOption) (Subscriber, error) {
	opt := SubscribeOptions{
		AutoAck: true,
	}

	for _, o := range opts {
		o(&opt)
	}

	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true                 //是否接收返回的错误消息
	config.Consumer.Retry.Backoff = 1 * time.Second      //失败后再次尝试的间隔时间
	config.Consumer.MaxWaitTime = 250 * time.Millisecond //最大等待时间
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	consumer, err := cluster.NewConsumer(k.addrs, opt.GroupId, []string{topic}, config)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case msg, more := <-consumer.Messages():
				if more {

					data := &kafkaMessage{t: msg.Topic, m: &Message{Data: msg.Value}}
					err = handler(data)
					if opt.AutoAck {
						continue
					}
					consumer.MarkOffset(msg, "")
				}
			case ntf, more := <-consumer.Notifications():
				if more {
					fmt.Println("qqq", ntf)
				}
			case err, more := <-consumer.Errors():
				if more {
					fmt.Println("qqq", err)
				}
			}
		}
	}()

	return &subscriber{consumer: consumer, opts: opt, t: topic}, nil
}

func newKafkaBroker(opts ...Option) Broker {
	return nil
}

func NewBroker(opts ...Option) Broker {
	return newKafkaBroker(opts...)
}
