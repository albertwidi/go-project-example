package nsq

import (
	"errors"
	"fmt"
	"log"
	"sync"

	gonsq "github.com/nsqio/go-nsq"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	// _buffMultiplier is a multiplier for buffered channel
	// the multiplier is used with the number for concurrency
	// the result is _buffMultiplier * numberOfConcurrency
	_buffMultiplier = 10
	// ErrTopicWithChannelNotFound for error when channel and topic is not found
	ErrTopicWithChannelNotFound = errors.New("nsq: topic and channel not found")
	// prometheus metrics
	_nsqMessageRetrievedCount *prometheus.CounterVec
	_nsqHandleCount           *prometheus.CounterVec
	_nsqHandleDurationHist    *prometheus.HistogramVec
	_nsqWorkerCurrentGauge    *prometheus.GaugeVec
	_nsqThrottleGauge         *prometheus.GaugeVec
	_nsqMessageInBuffGauge    *prometheus.GaugeVec
)

// throwing fatal if prometheus metrics cannot be registered
// registration error should not happen if metrics name is different
// tested in nsq_test.go
func init() {
	_nsqMessageRetrievedCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "nsq_message_retrieved_total",
		Help: "total message being retrieved from nsq for certain topic and channel, retrieved doesn't mean it is been processed",
	}, []string{"topic", "channel"})
	if err := prometheus.Register(_nsqMessageRetrievedCount); err != nil {
		if !errors.As(err, &prometheus.AlreadyRegisteredError{}) {
			err = fmt.Errorf("error when registering nsqMessageRetrievedCount. err: %w", err)
			log.Fatal(err)
		}
	}
	_nsqHandleCount = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "nsq_handle_error_total",
		Help: "total of message being handled",
	}, []string{"topic", "channel", "error"})
	if err := prometheus.Register(_nsqHandleCount); err != nil {
		if !errors.As(err, &prometheus.AlreadyRegisteredError{}) {
			err = fmt.Errorf("error when registering nsqHandleCount. err: %w", err)
			log.Fatal(err)
		}
	}
	_nsqHandleDurationHist = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "nsq_message_handle_duration",
	}, []string{"topic", "channel"})
	if err := prometheus.Register(_nsqHandleDurationHist); err != nil {
		if !errors.As(err, &prometheus.AlreadyRegisteredError{}) {
			err = fmt.Errorf("error when registering nsqHandleDurationHist. err: %w", err)
			log.Fatal(err)
		}
	}
	_nsqWorkerCurrentGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "nsq_worker_count_current",
	}, []string{"topic", "channel"})
	if err := prometheus.Register(_nsqWorkerCurrentGauge); err != nil {
		if !errors.As(err, &prometheus.AlreadyRegisteredError{}) {
			err = fmt.Errorf("error when registering nsqWorkerCurrentGauge. err: %w", err)
			log.Fatal(err)
		}
	}
	_nsqThrottleGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "nsq_throttle_status",
	}, []string{"topic", "channel"})
	if err := prometheus.Register(_nsqThrottleGauge); err != nil {
		if !errors.As(err, &prometheus.AlreadyRegisteredError{}) {
			err = fmt.Errorf("error when registering nsqThrottleGauge. err: %w", err)
			log.Fatal(err)
		}
	}
	_nsqMessageInBuffGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "nsq_message_in_buffer",
	}, []string{"topic", "channel"})
	if err := prometheus.Register(_nsqMessageInBuffGauge); err != nil {
		if !errors.As(err, &prometheus.AlreadyRegisteredError{}) {
			err = fmt.Errorf("error when registering nsqThrottleGauge. err: %w", err)
			log.Fatal(err)
		}
	}
}

// ProducerBackend for NSQ
type ProducerBackend interface {
	Ping() error
	Publish(topic string, body []byte) error
	MultiPublish(topic string, body [][]byte) error
	Stop()
}

// ConsumerBackend for NSQ
type ConsumerBackend interface {
	Topic() string
	Channel() string
	Stop()
	AddHandler(handler gonsq.Handler)
	AddConcurrentHandlers(handler gonsq.Handler, concurrency int)
	ConnectToNSQLookupds(addresses []string) error
	Concurrency() int
	ChangeMaxInFlight(n int)
}

// Producer for nsq
type Producer struct {
	producer ProducerBackend
	topics   map[string]bool
}

// WrapProducer is a function to wrap the nsq producer
func WrapProducer(backend ProducerBackend, topics ...string) *Producer {
	p := Producer{
		producer: backend,
		topics:   make(map[string]bool),
	}
	for _, t := range topics {
		p.topics[t] = true
	}
	return &p
}

// Publish message to nsqd
func (p *Producer) Publish(topic string, body []byte) error {
	if ok := p.topics[topic]; !ok {
		return errors.New("nsq: topic is not allowed to be published by this producer")
	}
	return p.producer.Publish(topic, body)
}

// MultiPublish message to nsqd
func (p *Producer) MultiPublish(topic string, body [][]byte) error {
	if ok := p.topics[topic]; !ok {
		return errors.New("nsq: topic is not allowed to be published by this producer")
	}
	return p.producer.MultiPublish(topic, body)
}

// Consumer for nsq
type Consumer struct {
	backends        map[string]map[string]ConsumerBackend
	handlers        []*nsqHandler
	middlewares     []MiddlewareFunc
	lookupdsAddress []string
	mu              sync.Mutex
}

// WrapConsumers of gonsq
func WrapConsumers(lookupdsAddr []string, backends ...ConsumerBackend) (*Consumer, error) {
	if lookupdsAddr == nil || len(lookupdsAddr) == 0 {
		return nil, errors.New("nsq: lookupd address cannot be empty")
	}

	b := make(map[string]map[string]ConsumerBackend)
	for _, c := range backends {
		topic := c.Topic()
		channel := c.Channel()

		if b[topic] == nil {
			b[topic] = make(map[string]ConsumerBackend)
		}
		b[topic][channel] = c
	}

	c := Consumer{
		backends:        b,
		lookupdsAddress: lookupdsAddr,
	}
	return &c, nil
}

// Backends return information regarding topic and channel that avaialbe
func (c *Consumer) Backends() map[string]map[string]bool {
	m := map[string]map[string]bool{}
	for topic, channels := range c.backends {
		for channel := range channels {
			if m[topic] == nil {
				m[topic] = map[string]bool{}
			}
			m[topic][channel] = true
		}
	}
	return m
}

// Use the middleware
// use should be called before handle function
// this function will avoid to add the same middleware twice
// if the same middleware is used, it will skip the addition
func (c *Consumer) Use(middleware ...MiddlewareFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	// check whether the middleware is already exits
	// if middleware already exist, avoid adding the middleware
	for _, m := range middleware {
		found := false
		for _, im := range c.middlewares {
			if &im == &m {
				found = true
				break
			}
		}
		if !found {
			c.middlewares = append(c.middlewares, m)
		}
	}
}

// Handle the consumer
func (c *Consumer) Handle(topic, channel string, handler HandlerFunc) {
	c.mu.Lock()
	defer c.mu.Unlock()

	for i := range c.middlewares {
		handler = c.middlewares[len(c.middlewares)-i-1](handler)
	}
	h := &nsqHandler{
		topic:    topic,
		channel:  channel,
		handler:  handler,
		stopChan: make(chan struct{}),
	}
	c.handlers = append(c.handlers, h)
}

// Start the consumer
func (c *Consumer) Start() error {
	for _, handler := range c.handlers {
		backend, ok := c.backends[handler.topic][handler.channel]
		if !ok {
			return fmt.Errorf("nsq: backend with topoc %s and channel %s not found. error: %w", handler.topic, handler.channel, ErrTopicWithChannelNotFound)
		}
		// set concurrency for handler
		handler.SetConcurrency(backend.Concurrency())

		// create a default handler for handling nsq handler
		dh := defaultHandler{
			nsqHandler:      handler,
			consumerBackend: backend,
		}
		// consumerConcurrency for consuming message from NSQ
		// most of the time we don't need consumerConcurrency because consuming message from NSQ is very fast
		// the handler or the true consumer might need time to handle the message
		// so we need to keep the number of message consumer low, for now it is 1:30
		consumerConcurrency := handler.concurrency / 30
		if consumerConcurrency > 1 {
			backend.AddConcurrentHandlers(&dh, consumerConcurrency)
		} else {
			backend.AddHandler(&dh)
		}
		// change the MaxInFlight to buffLength as the number of message won't exceed the buffLength
		backend.ChangeMaxInFlight(dh.buffLength)

		if err := backend.ConnectToNSQLookupds(c.lookupdsAddress); err != nil {
			return err
		}
		// invoke all handler to work
		// depends on the concurrency that is initiated
		for i := 0; i < handler.concurrency; i++ {
			go handler.Work()
		}
	}
	return nil
}

// Stop all the nsq consumer
func (c *Consumer) Stop() error {
	for _, channels := range c.backends {
		for _, backend := range channels {
			backend.Stop()
		}
	}
	for _, handler := range c.handlers {
		// stop all the handler worker based on concurrency number
		// this step is expected to be blocking
		// wait until all worker is exited
		for i := 0; i < handler.concurrency; i++ {
			handler.Stop()
		}
	}
	return nil
}
