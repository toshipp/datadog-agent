// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License Version 2.0.
// This product includes software developed at Datadog (https://www.datadoghq.com/).
// Copyright 2016-2019 Datadog, Inc.

package sender

import (
	"context"
	"time"

	"github.com/DataDog/datadog-agent/pkg/util/log"

	"github.com/DataDog/datadog-agent/pkg/logs/message"
)

const (
	batchTimeout   = 5 * time.Second
	maxBatchSize   = 20
	maxContentSize = 1000000
)

// batchStrategy contains all the logic to send logs in batch.
type batchStrategy struct {
	buffer       *MessageBuffer
	formatter    Formatter
	batchTimeout time.Duration
}

// NewBatchStrategy returns a new batchStrategy.
func NewBatchStrategy(formatter Formatter) Strategy {
	return &batchStrategy{
		buffer:       NewMessageBuffer(maxBatchSize, maxContentSize),
		formatter:    formatter,
		batchTimeout: batchTimeout,
	}
}

// Send accumulates messages to a buffer and sends them when the buffer is full or outdated.
func (s *batchStrategy) Send(inputChan chan *message.Message, outputChan chan *message.Message, send func([]byte) error) {
	flushTimer := time.NewTimer(s.batchTimeout)
	defer func() {
		flushTimer.Stop()
	}()

	for {
		select {
		case message, isOpen := <-inputChan:
			if !isOpen {
				// inputChan has been closed, no more payload are expected
				s.sendBuffer(outputChan, send)
				return
			}
			added := s.buffer.AddMessage(message)
			if !added || s.buffer.IsFull() {
				// message buffer is full, either reaching max batch size or max content size,
				// send the payload now
				if !flushTimer.Stop() {
					// make sure the timer won't tick concurrently
					select {
					case <-flushTimer.C:
					default:
					}
				}
				s.sendBuffer(outputChan, send)
				flushTimer.Reset(s.batchTimeout)
			}
			if !added {
				// it's possible that the message could not be added because the buffer was full
				// so we need to retry once again
				s.buffer.AddMessage(message)
			}
		case <-flushTimer.C:
			// the first message that was added to the buffer has been here for too long,
			// send the payload now
			s.sendBuffer(outputChan, send)
			flushTimer.Reset(s.batchTimeout)
		}
	}
}

// sendBuffer sends all the messages that are stored in the buffer and forwards them
// to the next stage of the pipeline.
func (s *batchStrategy) sendBuffer(outputChan chan *message.Message, send func([]byte) error) {
	if s.buffer.IsEmpty() {
		return
	}

	messages := s.buffer.GetMessages()
	defer s.buffer.Clear()

	err := send(s.formatter.Format(messages))
	if err != nil {
		if err == context.Canceled {
			return
		}
		log.Warnf("Could not send payload: %v", err)
	}

	for _, message := range messages {
		outputChan <- message
	}
}
