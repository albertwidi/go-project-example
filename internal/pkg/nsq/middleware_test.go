package nsq

import (
	"context"
	"errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	fakensq "github.com/albertwidi/go-project-example/internal/pkg/nsq/fakensq"
)

func TestThrottleMiddleware(t *testing.T) {
	t.Parallel()

	var (
		topic             = "test_topic"
		channel           = "test_channel"
		errChan           = make(chan error)
		messageNum        int32
		currentMessageNum int32
		messageThrottled  int32

		// testing expect
		messageExpect = "test middleware throttling"
		// to make sure that error is being sent back
		errNil = errors.New("error should be nil")
	)

	// we are using fake consumer, this means the concurrency is always 1
	// and the number of message buffer is 1 * _bufferMultiplier
	consumer, err := fakensq.NewFakeConsumer(topic, channel, nil)
	if err != nil {
		t.Error(err)
		return
	}
	producer := fakensq.NewFakeProducer(consumer)

	wc, err := WrapConsumers([]string{"test"}, consumer)
	if err != nil {
		t.Error(err)
		return
	}

	tm := ThrottleMiddleware{TimeDelay: time.Second}
	// use throttle middleware
	wc.Use(
		tm.Throttle,
	)

	wc.Handle(topic, channel, func(ctx context.Context, message *Message) error {
		atomic.AddInt32(&currentMessageNum, 1)

		if string(message.Message.Body) != messageExpect {
			err := fmt.Errorf("epecting message %s but got %s", messageExpect, string(message.Message.Body))
			errChan <- err
			return err
		}

		// this means this is the first message, sleep to make other message to wait
		// because if the handler is not finished, the worker is not back to consume state
		// to make sure the buffer is filled first before consuming more message
		if currentMessageNum == 1 {
			time.Sleep(time.Millisecond * 100)
		}

		// check whether a throttled message is exists
		// this message should exists because throttle middleware is used
		if message.Info.Throttled == 1 {
			atomic.AddInt32(&messageThrottled, 1)
			errChan <- errNil
			return nil
		}

		// this means the test have reach the end of message
		if currentMessageNum == messageNum {
			if messageThrottled < 1 {
				err := errors.New("message is never throttled")
				errChan <- err
				return err
			}
			errChan <- errNil
			return nil
		}

		errChan <- errNil
		return err
	})

	if err := wc.Start(); err != nil {
		t.Error(err)
		return
	}

	// send message as much as (_buffMultiplier/2) + 3 to tirgger the throttle mechanism
	// why we need at least +3 message?
	//
	// c = consumed
	// d = done
	// m = message in buffer
	// <nil> = no message, buffer is empty
	//
	// _buffMultiplier/2 + 3 = 8 message
	// | m | m | m | m | m | m | m | m | <nil> | <nil> |
	//   1   2   3   4   5   6   7   8     9      10
	// message_length: 8
	//
	//
	// when the program start, the message will be consumed into the worker right away
	// then the worker will pause themself for a while, so message is not consumed
	// at this point, this is the look in the buffer:
	// | m | m | m | m | m | m | m | m | <nil> | <nil> | <nil > |
	//   c   1   2   3   4   5   6   7     8       9      10
	// message_length: 7
	//
	// at this point when consuming more message which is message number 1, the buffer will become:
	// | m | m | m | m | m | m | m | m | <nil> | <nil> | <nil > | <nil> |
	//   d   c   1   2   3   4   5   6     7       8       9       10
	// message_length: 6
	//
	// when consuming the message, evaluation of the buffer length will kick in
	// this is where the evaluator for thorttle know that the number of message
	// is more than half of the buffer size, then throttle mechanism will be invoked
	// this is why, with lower number of message the test won't pass,
	// because it depends on message number in the buffer
	for i := 1; i <= (_buffMultiplier/2)+3; i++ {
		if err := producer.Publish(topic, []byte(messageExpect)); err != nil {
			t.Error(err)
			return
		}
		messageNum++
	}

	for i := 1; i <= int(messageNum); i++ {
		err = <-errChan
		if err != errNil {
			t.Error(err)
			return
		}
	}
	close(errChan)

	if err := wc.Stop(); err != nil {
		t.Error(err)
		return
	}
}
