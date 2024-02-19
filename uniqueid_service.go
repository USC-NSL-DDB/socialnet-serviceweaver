package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ServiceWeaver/weaver"
)

type IUniqueIdService interface {
	ComposeUniqueId(context.Context, PostType) int64
}

type UniqueIdService struct {
	weaver.Implements[IUniqueIdService]

	mu               sync.Mutex
	currentTimestamp int64
	counter          int
	customEpoch      int64
	machineId       string
}

func (s *UniqueIdService) Init(context.Context) error {
	s.currentTimestamp = -1
	s.counter = 0
	// Custom Epoch (January 1, 2018 Midnight GMT = 2018-01-01T00:00:00Z)
	s.customEpoch = 1514764800000
	s.machineId = "0000000000000000"
	return nil
}

func (s *UniqueIdService) ComposeUniqueId(_ context.Context, postType PostType) int64 {
	timestamp := time.Now().UnixNano()/int64(time.Millisecond) - s.customEpoch
	idx := s.GetCounter(timestamp)

	// Converting timestamp and counter to hex strings and composing the unique ID
	timestampHex := fmt.Sprintf("%x", timestamp)
	if len(timestampHex) > 10 {
		timestampHex = timestampHex[len(timestampHex)-10:]
	} else if len(timestampHex) < 10 {
		timestampHex = fmt.Sprintf("%0*s", 10, timestampHex)
	}

	counterHex := fmt.Sprintf("%x", idx)
	if len(counterHex) > 3 {
		counterHex = counterHex[len(counterHex)-3:]
	} else if len(counterHex) < 3 {
		counterHex = fmt.Sprintf("%0*s", 3, counterHex)
	}

	postIDStr := s.machineId + timestampHex + counterHex
	var postID int64
	fmt.Sscanf(postIDStr, "%x", &postID)
	postID = postID & 0x7FFFFFFFFFFFFFFF

	return postID
}

// GetCounter - Manages the incrementation of the counter within the same millisecond timestamp
func (s *UniqueIdService) GetCounter(timestamp int64) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.currentTimestamp > timestamp {
		panic("Timestamps are not incremental.")
	}
	if s.currentTimestamp == timestamp {
		s.counter++
		return s.counter
	} else {
		s.currentTimestamp = timestamp
		s.counter = 0
		return s.counter
	}
}
