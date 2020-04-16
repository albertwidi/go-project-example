package nsq

import (
	"context"
	"testing"
	"time"

	fakensq "github.com/kyolabs/mono/gopkg/nsq/fakensq"
	gonsq "github.com/nsqio/go-nsq"
)

func TestNSQHandlerSetConcurrency(t *testing.T) {
	t.Parallel()

	cases := []struct {
		concurrency            int
		buffMultiplier         int
		expectConcurreny       int
		expectBufferMultiplier int
	}{
		{
			concurrency:            1,
			buffMultiplier:         10,
			expectConcurreny:       1,
			expectBufferMultiplier: 10,
		},
		{
			concurrency:            -1,
			buffMultiplier:         -1,
			expectConcurreny:       1,
			expectBufferMultiplier: 20,
		},
		{
			concurrency:            1,
			buffMultiplier:         1,
			expectConcurreny:       1,
			expectBufferMultiplier: 1,
		},
		{
			concurrency:            1,
			buffMultiplier:         -1,
			expectConcurreny:       1,
			expectBufferMultiplier: 20,
		},
	}

	var (
		topic   = "test_concurrency"
		channel = "test_concurrency"
	)

	for _, c := range cases {
		t.Logf("concurrency: %d", c.concurrency)
		t.Logf("buff_multiplier: %d", c.buffMultiplier)

		backend, err := fakensq.NewFakeConsumer(fakensq.ConsumerConfig{Topic: topic, Channel: channel, Concurrency: c.concurrency, BufferMultiplier: c.buffMultiplier})
		if err != nil {
			t.Error(err)
			return
		}

		wc, err := WrapConsumers(ConsumerConfig{
			LookupdsAddr: []string{"testing"},
		}, backend)
		if err != nil {
			t.Error(err)
			return
		}
		// Trigger the creation of handler.
		wc.Handle(topic, channel, nil)

		handler := wc.handlers[0]
		if handler == nil {
			t.Error("handler should not be nil, as handle is triggered")
		}

		if handler.concurrency != c.expectConcurreny {
			t.Errorf("expecting concurrency %d but got %d", c.expectConcurreny, handler.concurrency)
			return
		}
		if handler.buffMultiplier != c.expectBufferMultiplier {
			t.Errorf("expecting buffer multiplier of %d but got %d", c.expectBufferMultiplier, handler.buffMultiplier)
			return
		}
		if handler.buffLength != c.expectConcurreny*c.expectBufferMultiplier {
			t.Errorf("expecting buffer length of %d but got %d", c.expectConcurreny*c.expectBufferMultiplier, handler.buffLength)
			return
		}

	}
}

func TestNSQHandlerConcurrencyControl(t *testing.T) {
	t.Parallel()

	var (
		topic       = "test_concurrency"
		channel     = "test_concurrency"
		concurrency = 5
	)

	backend, err := fakensq.NewFakeConsumer(fakensq.ConsumerConfig{Topic: topic, Channel: channel})
	if err != nil {
		t.Error(err)
		return
	}

	wc, err := WrapConsumers(ConsumerConfig{
		LookupdsAddr: []string{"testing"},
		Concurrency:  concurrency,
	}, backend)
	if err != nil {
		t.Error(err)
		return
	}
	// Trigger the creation of handler.
	wc.Handle(topic, channel, nil)

	handler := wc.handlers[0]
	if handler == nil {
		t.Error("handler should not be nil, as handle is triggered")
	}

	for i := 1; i <= 5; i++ {
		go handler.Work()
		// Wait until the goroutines scheduled
		// this might be too long, but its ok.
		time.Sleep(time.Millisecond * 10)
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

	backend, err := fakensq.NewFakeConsumer(fakensq.ConsumerConfig{Topic: topic, Channel: channel})
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
		// Using buffered channel with length 1,
		// because in this test we don't listen the message using a worker,
		// and sending the message to this channel will block.
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

	topic := "random"
	channel := "random"
	backend, err := fakensq.NewFakeConsumer(fakensq.ConsumerConfig{Topic: topic, Channel: channel})
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
