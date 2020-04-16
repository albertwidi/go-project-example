package nsq

import (
	"context"
	"sync"
	"time"

	gonsq "github.com/nsqio/go-nsq"
)

// HandlerFunc for nsq
type HandlerFunc func(ctx context.Context, message *Message) error

// MiddlewareFunc for nsq middleware
type MiddlewareFunc func(handler HandlerFunc) HandlerFunc

// Message for nsq
type Message struct {
	Topic   string
	Channel string
	Message *gonsq.Message
	Info    *Info
}

// ID return message id from gonsq message.ID
func (m *Message) ID() gonsq.MessageID {
	return m.Message.ID
}

// Finish call gonsq message finish
func (m *Message) Finish() {
	m.Message.Finish()
}

// Requeue call gonsq message requeue
func (m *Message) Requeue(delay time.Duration) {
	m.Message.Requeue(delay)
}

// RequeueWithoutBackoff call gonsq message requeueWithoutBackoff
func (m *Message) RequeueWithoutBackoff(delay time.Duration) {
	m.Message.RequeueWithoutBackoff(delay)
}

// Info for message
type Info struct {
	WorkerTotal     int
	WorkerCurrent   int
	MessageInBuffer int
	ThrottleFlag    int
	Throttled       int
}

type nsqHandler struct {
	handler        HandlerFunc
	concurrency    int
	buffMultiplier int
	workerNumber   int
	topic          string
	channel        string
	// messageBuff is a buffered channel to buffer messages
	messageBuff chan *Message
	buffLength  int
	stopChan    chan struct{}
	mu          sync.Mutex
	throttle    bool
}

// SetThrottle to set the handler status if throttled or not
func (nh *nsqHandler) SetThrottle(throttle bool) {
	nh.mu.Lock()
	nh.throttle = throttle
	nh.mu.Unlock()
}

// Work to handle nsq message
func (nh *nsqHandler) Work() {
	nh.mu.Lock()
	// Guard with lock,
	// don't let worker number goes more than concurrency number.
	if nh.workerNumber == nh.concurrency {
		return
	}
	nh.workerNumber++
	nh.mu.Unlock()

	for {
		select {
		case <-nh.stopChan:
			return
		case message := <-nh.messageBuff:
			// Add information about worker to message
			// this will add additional allocation and memory
			// but essential to monitor the number of worker.
			message.Info.WorkerTotal = nh.concurrency
			message.Info.WorkerCurrent = nh.workerNumber
			message.Info.MessageInBuffer = len(nh.messageBuff)
			// Set the next message throttle flag to 1
			// because the handler is set to throttle from the default handler.
			if nh.throttle {
				message.Info.ThrottleFlag = 1
			}
			nh.handler(context.Background(), message)
		}
	}
}

// Stop the work of nsq handler
func (nh *nsqHandler) Stop() {
	nh.stopChan <- struct{}{}
	nh.workerNumber--
}

type defaultHandler struct {
	*nsqHandler
	// backend
	consumerBackend ConsumerBackend
}

// HandleMessage of nsq
func (dfh *defaultHandler) HandleMessage(message *gonsq.Message) error {
	_nsqMessageRetrievedCount.WithLabelValues(dfh.topic, dfh.channel).Add(1)
	// Message in the buffer should always less than bufferLength/2
	// if its already more than half of the buffer size, we should pause the consumption
	// and wait for the buffer to be consumed first.
	if len(dfh.messageBuff) > (dfh.buffLength / 2) {
		// set the handler throttle to true, so all message will be throttled right away
		dfh.SetThrottle(true)
		// pause the message consumption to NSQD by set the MaxInFlight to 0
		dfh.consumerBackend.ChangeMaxInFlight(0)
		for {
			// Sleep every one second to check whether the message number is already decreased in the buffer,
			// it might be better to have a lower evaluation interval, but need some metrics first.
			// The default throttling here won't affect the message consumer because messages already buffered
			// but will have some effect for the nsqd itself because we pause the message consumption from nsqd.
			time.Sleep(time.Second * 1)
			if len(dfh.messageBuff) < (dfh.buffLength / 2) {
				// resume the message consumption to NSQD by set the MaxInFlight to buffer size
				dfh.consumerBackend.ChangeMaxInFlight(dfh.buffLength)
				dfh.SetThrottle(false)
				break
			}
		}
	}
	dfh.messageBuff <- &Message{
		Topic:   dfh.topic,
		Channel: dfh.channel,
		Message: message,
		// allocate for info
		Info: &Info{},
	}
	return nil
}
