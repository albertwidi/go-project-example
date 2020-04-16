package nsq

import (
	"context"
	"time"
)

// Metrics middleware for nsq
// metrics that might be missleading:
// - throttled
// - message_in_buffer
// The metrics might be missleading because the message is not processed in ordered manner.
func Metrics(handler HandlerFunc) HandlerFunc {
	return func(ctx context.Context, message *Message) error {
		t := time.Now()
		e := "0"
		err := handler(ctx, message)
		if err != nil {
			e = "1"
		}
		_nsqHandleDurationHist.WithLabelValues(message.Topic, message.Channel).Observe(float64(time.Now().Sub(t).Milliseconds()))
		_nsqHandleCount.WithLabelValues(message.Topic, message.Channel, e).Add(1)
		_nsqWorkerCurrentGauge.WithLabelValues(message.Topic, message.Channel).Set(float64(message.Info.WorkerCurrent))
		_nsqThrottleGauge.WithLabelValues(message.Topic, message.Channel).Set(float64(message.Info.Throttled))
		_nsqMessageInBuffGauge.WithLabelValues(message.Topic, message.Channel).Set(float64(message.Info.MessageInBuffer))
		return err
	}
}

// ThrottleMiddleware implement MiddlewareFunc
type ThrottleMiddleware struct {
	// TimeDelay means the duration of time to pause message consumption
	TimeDelay time.Duration
}

// Throttle middleware for nsq.
// This middleware check whether there is some information about throttling in the message.
func (tm *ThrottleMiddleware) Throttle(handler HandlerFunc) HandlerFunc {
	return func(ctx context.Context, message *Message) error {
		// this means the worker is being throttled
		if message.Info.ThrottleFlag == 1 {
			time.Sleep(tm.TimeDelay)
			// set the status to be throttled because the middleware is active
			// this is needed for metrics
			message.Info.Throttled = 1
		}
		return handler(ctx, message)
	}
}
