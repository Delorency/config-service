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
	mu             sync.RWMutex
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

func (wm *WatcherManager) Stop() {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	for serviceID, subs := range wm.subscribers {
		for _, sub := range subs {
			close(sub.UpdateChan)
		}
		delete(wm.subscribers, serviceID)
	}
}
