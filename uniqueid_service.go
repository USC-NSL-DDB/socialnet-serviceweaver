package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
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
	// TODO: How to obtain the machine id?
	machineId        string
}

// Custom Epoch (January 1, 2018 Midnight GMT = 2018-01-01T00:00:00Z)
const CustomEpoch = 1514764800000

func (s *UniqueIdService) Init(context.Context) error {
	s.currentTimestamp = -1
	s.counter = 0
	s.machineId = "0000000000000000"
	return nil
}

func (s *UniqueIdService) ComposeUniqueId(_ context.Context, postType PostType) int64 {
	timestamp := time.Now().UnixNano()/int64(time.Millisecond) - CustomEpoch
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

func GetMachineId(netif string) string {
	macAddrFilename := "/sys/class/net/" + netif + "/address"

	macAddrFile, err := os.Open(macAddrFilename)
	if err != nil {
		log.Fatalf("Cannot read MAC address from net interface %s: %v", netif, err)
		return ""
	}
	defer macAddrFile.Close()

	scanner := bufio.NewScanner(macAddrFile)
	scanner.Scan()
	mac := scanner.Text()
	if mac == "" {
		log.Fatalf("Cannot read MAC address from net interface %s", netif)
		return ""
	}

	log.Printf("MAC address = %s", mac)

	macHash := fmt.Sprintf("%x", HashMacAddressPid(mac))

	if len(macHash) > 3 {
		macHash = macHash[len(macHash)-3:]
	} else if len(macHash) < 3 {
		macHash = strings.Repeat("0", 3-len(macHash)) + macHash
	}

	return macHash
}

func HashMacAddressPid(mac string) uint16 {
	var hash uint16 = 0
	pid := os.Getpid() // Get the current process ID
	macPid := mac + strconv.Itoa(pid)

	for i, char := range macPid {
		if i < len(mac) { // Ensure we only consider the MAC address length for the hash calculation
			hash += uint16(char) << ((i & 1) * 8)
		}
	}
	return hash
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
