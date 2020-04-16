package fakensq

import (
	"errors"
	"fmt"
	"sync"
	"time"

	nsqio "github.com/nsqio/go-nsq"
)

// FakeProducer struct
type FakeProducer struct {
	*FakeLookupd
}

// NewFakeProducer for publishing message to NSQ
func NewFakeProducer(consumer *FakeConsumer) *FakeProducer {
	p := FakeProducer{consumer.FakeLookupd}
	return &p
}

// Ping will always return nil
func (fp *FakeProducer) Ping() error {
	return nil
}

// Publish a message
// this function might block if the channel is full
func (fp *FakeProducer) Publish(topic string, message []byte) error {
	m := nsqio.Message{
		Body:     message,
		Delegate: &MessageDelegator{},
	}
	// publish to fakelookupd
	fp.publish(topic, &m)
	return nil
}

// MultiPublish message
func (fp *FakeProducer) MultiPublish(topic string, messages [][]byte) error {
	for _, message := range messages {
		if err := fp.Publish(topic, message); err != nil {
			return err
		}
	}
	return nil
}

// Stop fake producer
func (fp *FakeProducer) Stop() {
	return
}

// FakeConsumer struct
type FakeConsumer struct {
	*FakeLookupd
	config      ConsumerConfig
	data        []Message
	messageChan chan *nsqio.Message
	ErrChan     chan error
	handlers    []nsqHandler
	started     bool
	stopped     bool
	// nsq configuration
	MaxInFlight int
}

type nsqHandler struct {
	messageChan chan *nsqio.Message
	handler     nsqio.Handler
	stopChan    chan struct{}
}

func (h *nsqHandler) Stop() {
	h.stopChan <- struct{}{}
	return
}

// MessageDelegator implement Delegator of nsqio
type MessageDelegator struct {
}

func (mdm *MessageDelegator) OnFinish(message *nsqio.Message) {
	return
}

func (mdm *MessageDelegator) OnRequeue(m *nsqio.Message, t time.Duration, backoff bool) {
	return
}

func (mdm *MessageDelegator) OnTouch(m *nsqio.Message) {
	return
}

// Message mock
type Message struct {
	Name string
	Body []byte
}

// ConsumerConfig of fake nsq
type ConsumerConfig struct {
	Topic            string
	Channel          string
	Concurrency      int
	BufferMultiplier int
}

// Validate consumer configuration
func (cc *ConsumerConfig) Validate() error {
	if cc.Topic == "" || cc.Channel == "" {
		return errors.New("topic or channel cannot be empty")
	}
	if cc.Concurrency <= 0 {
		cc.Concurrency = 1
	}
	if cc.BufferMultiplier == 0 {
		cc.BufferMultiplier = 30
	}
	return nil
}

// NewFakeConsumer function
func NewFakeConsumer(config ConsumerConfig) (*FakeConsumer, error) {
	mock := FakeConsumer{
		FakeLookupd: &FakeLookupd{
			topicChannel: make(map[string]map[string]chan *nsqio.Message),
		},
		config:      config,
		messageChan: make(chan *nsqio.Message),
		ErrChan:     make(chan error),
	}
	return &mock, nil
}

// Topic return the consumer topic
func (cm *FakeConsumer) Topic() string {
	return cm.config.Topic
}

// Channel return the consumer channel
func (cm *FakeConsumer) Channel() string {
	return cm.config.Channel
}

// AddHandler for nsq
func (cm *FakeConsumer) AddHandler(handler nsqio.Handler) {
	cm.addHandlers(handler, 1)
}

// AddConcurrentHandlers for nsq
func (cm *FakeConsumer) AddConcurrentHandlers(handler nsqio.Handler, concurrency int) {
	cm.addHandlers(handler, concurrency)
}

func (cm *FakeConsumer) addHandlers(handler nsqio.Handler, concurrency int) {
	for i := 0; i < concurrency; i++ {
		h := nsqHandler{
			messageChan: cm.messageChan,
			handler:     handler,
			stopChan:    make(chan struct{}),
		}
		go func(h nsqHandler) {
			for {
				select {
				case <-h.stopChan:
					return
				case msg := <-cm.messageChan:
					err := h.handler.HandleMessage(msg)
					// send all error to the channel, and decide what to do
					if err != nil {
						ecm := ErrorConsumerFake{
							topic:   cm.config.Topic,
							channel: cm.config.Channel,
							message: msg.Body,
							err:     fmt.Errorf("nsq-mock-handler: error when handling message, with error: %w", err),
						}
						cm.ErrChan <- &ecm
					}
				}
			}
		}(h)
		cm.handlers = append(cm.handlers, h)
	}
}

// ConnectToNSQLookupds for nsq
func (cm *FakeConsumer) ConnectToNSQLookupds(addresses []string) error {
	cm.start()
	// register the topic, channel and the message channel to the fake lookupD
	// so the fake publisher able to publish the message into the right topic and channel
	cm.FakeLookupd.register(cm.Topic(), cm.Channel(), cm.messageChan)
	return nil
}

// ChangeMaxInFlight message in nsq consumer
func (cm *FakeConsumer) ChangeMaxInFlight(n int) {
	cm.MaxInFlight = n
}

// Concurrency return the number of conccurent worker
func (cm *FakeConsumer) Concurrency() int {
	return cm.config.Concurrency
}

// BufferMultiplier return the number of buffer multiplier
func (cm *FakeConsumer) BufferMultiplier() int {
	return cm.config.BufferMultiplier
}

// start will start the message sending
// the message sending mechanism will be concurrent
// based on how fast the consumer consumed the message
// using a single unbuffered channel
func (cm *FakeConsumer) start() {
	if cm.started {
		return
	}

	go func() {
		for _, d := range cm.data {
			if cm.stopped {
				return
			}

			delegate := MessageDelegator{}
			m := nsqio.Message{
				Body:     d.Body,
				Delegate: &delegate,
			}
			// send the message to the worker channel
			cm.messageChan <- &m
		}
	}()
	cm.started = true
}

// Stop consumer backend mock
func (cm *FakeConsumer) Stop() {
	for _, h := range cm.handlers {
		h.Stop()
	}
	cm.stopped = true

	return
}

// FakeLookupd for storing all information regarding topics and channel
type FakeLookupd struct {
	topicChannel map[string]map[string]chan *nsqio.Message
	mu           sync.Mutex
}

// register topic, channel and message channel to fake lookupd
func (fld *FakeLookupd) register(topic, channel string, messageChan chan *nsqio.Message) {
	_, ok := fld.topicChannel[topic]
	if !ok {
		fld.topicChannel[topic] = make(map[string]chan *nsqio.Message)
	}
	fld.topicChannel[topic][channel] = messageChan
}

// publish the message to desired topic and multiplex via channel
func (fld *FakeLookupd) publish(topic string, message *nsqio.Message) error {
	// block because we might see blocking in publishing message to channel
	// no point of accessing this function concurrently
	fld.mu.Lock()
	defer fld.mu.Unlock()

	t, ok := fld.topicChannel[topic]
	if !ok {
		return nil
	}

	// publish to all possible channel
	// expect this to be blocking
	for _, messageChan := range t {
		messageChan <- message
	}
	return nil
}

// ErrorConsumerFake for throwing error from the mock consumer
type ErrorConsumerFake struct {
	topic   string
	channel string
	message []byte
	err     error
}

// Error return the error string from error consumer mock
func (ecm *ErrorConsumerFake) Error() string {
	return ecm.err.Error()
}

// Is implementation of error
func (ecm *ErrorConsumerFake) Is(err, target error) bool {
	return errors.Is(err, target)
}

// As implementation of error
func (ecm *ErrorConsumerFake) As(err error, target interface{}) bool {
	return errors.As(err, target)
}

// Unwrap implementation of error
func (ecm *ErrorConsumerFake) Unwrap(err error) error {
	return errors.Unwrap(err)
}

// Topic return the topic of error
func (ecm *ErrorConsumerFake) Topic() string {
	return ecm.topic
}

// Channel return the channel of error
func (ecm *ErrorConsumerFake) Channel() string {
	return ecm.channel
}

// Message return the message that return error
func (ecm *ErrorConsumerFake) Message() []byte {
	return ecm.message
}
