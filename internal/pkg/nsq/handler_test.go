package nsq

import (
	"context"
	"testing"
	"time"

	fakensq "github.com/kosanapp/backend/internal/pkg/nsq/fakensq"
	gonsq "github.com/nsqio/go-nsq"
)

func TestNSQHandlerSetConcurrency(t *testing.T) {
	t.Parallel()

	concurrencies := []int{0, 1, -1, -10, 10, 20}

	for _, c := range concurrencies {
		t.Logf("concurrency: %d", c)
		buffLen := c * _buffMultiplier
		handler := nsqHandler{}
		handler.SetConcurrency(c)

		// set to default if value is less than 0
		if c <= 0 {
			c = 1
			buffLen = _buffMultiplier
		}

		if handler.concurrency != c {
			t.Errorf("expecting concurrency %d but got %d", c, handler.concurrency)
			return
		}

		if handler.buffLength != buffLen {
			t.Errorf("expecting buffer length of %d but got %d", buffLen, handler.buffLength)
			return
		}
	}
}

func TestNSQHandlerConcurrencyControl(t *testing.T) {
	t.Parallel()

	concurrency := 5

	handler := nsqHandler{}
	handler.SetConcurrency(concurrency)

	for i := 1; i <= 5; i++ {
		go handler.Work()
		// wait until the goroutines scheduled
		// this might be too long, but its ok
		time.Sleep(time.Millisecond * 100)
		if handler.workerNumber != i {
			t.Errorf("start: expecting number worker number of %d but got %d", i, handler.workerNumber)
			return
		}
	}

	for i := handler.workerNumber; i < 0; i-- {
		handler.Stop()
		if handler.workerNumber != i-1 {
			t.Errorf("stop: expecting worker number of %d but got %d", i-1, handler.workerNumber)
			return
		}
	}
}

func TestDefaultHandlerHandleMessage(t *testing.T) {
	t.Parallel()

	var (
		topic       = "test_topic"
		channel     = "test_channel"
		messageBody = []byte("test message")
	)

	backend, err := fakensq.NewFakeConsumer("random", "random", nil)
	if err != nil {
		t.Error(err)
		return
	}

	df := defaultHandler{
		nsqHandler: &nsqHandler{
			topic:       topic,
			channel:     channel,
			messageBuff: make(chan *Message, 1),
		},
		// using buffered channel with length 1
		// because in this test we don't listen the message using a worker
		// and sending the message to this channel will block
		consumerBackend: backend,
	}

	nsqMessage := &gonsq.Message{
		Body:     messageBody,
		Delegate: &fakensq.MessageDelegator{},
	}
	if err := df.HandleMessage(nsqMessage); err != nil {
		t.Error(err)
		return
	}

	message := <-df.messageBuff
	if message.Topic != topic {
		t.Errorf("expecting topic %s but got %s", topic, message.Topic)
		return
	}
	if message.Channel != channel {
		t.Errorf("expecting channel %s but got %s", channel, message.Channel)
		return
	}
	if string(message.Message.Body) != string(messageBody) {
		t.Errorf("expecting message body %s but got %s", string(message.Message.Body), string(messageBody))
		return
	}
}

func TestDefaultHandlerThrottle(t *testing.T) {
	t.Parallel()

	backend, err := fakensq.NewFakeConsumer("random", "random", nil)
	if err != nil {
		t.Error(err)
		return
	}

	df := defaultHandler{
		nsqHandler: &nsqHandler{
			messageBuff: make(chan *Message, 3),
		},
		consumerBackend: backend,
	}
	doneChan := make(chan struct{})

	// TODO: change the static number for limiter to dynamic
	for i := 0; i < 3; i++ {
		m := &gonsq.Message{}
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()

		go func() {
			if err := df.HandleMessage(m); err != nil {
				t.Error(err)
				return
			}
			doneChan <- struct{}{}
		}()

		select {
		case <-doneChan:
			continue
		case <-ctx.Done():
			if i != 1 {
				t.Error("error while current buffer is less than half")
				return
			}
			if backend.MaxInFlight != 0 {
				t.Error("error: max in flight is not being set to 0")
				return
			}
			return
		}
	}
}
