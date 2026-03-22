package watcher

import (
	"sync"
	"time"
)

type ConfigUpdate struct {
	ServiceID string
	Data      []byte
	Version   int
	Timestamp int64
}

type Subscriber struct {
	ID             string
	UpdateChan     chan *ConfigUpdate
	CurrentVersion int
	ConnectedAt    time.Time
}

type WatcherManager struct {
	mu          sync.RWMutex
	subscribers map[string][]*Subscriber
}

func NewWatcherManager() *WatcherManager {
	return &WatcherManager{
		subscribers: make(map[string][]*Subscriber),
	}
}
