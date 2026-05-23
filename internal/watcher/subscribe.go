package watcher

import "time"

func (wm *WatcherManager) Subscribe(serviceID, subscriberID string, currentVersion int) (*Subscriber, func()) {
	wm.mu.Lock()
	defer wm.mu.Unlock()

	sub := &Subscriber{
		ID:             subscriberID,
		UpdateChan:     make(chan *ConfigUpdate, 10),
		CurrentVersion: currentVersion,
		ConnectedAt:    time.Now(),
	}

	wm.subscribers[serviceID] = append(wm.subscribers[serviceID], sub)

	unsubscribe := func() {
		wm.mu.Lock()
		defer wm.mu.Unlock()

		subs := wm.subscribers[serviceID]
		for i, s := range subs {
			if s.ID == subscriberID {
				close(s.UpdateChan)
				wm.subscribers[serviceID] = append(subs[:i], subs[i+1:]...)
				break
			}
		}

		if len(wm.subscribers[serviceID]) == 0 {
			delete(wm.subscribers, serviceID)
		}
	}

	return sub, unsubscribe
}

func (wm *WatcherManager) GetSubscribersCount(serviceID string) int {
	wm.mu.RLock()
	defer wm.mu.RUnlock()
	return len(wm.subscribers[serviceID])
}

func (wm *WatcherManager) getSubscribersCopy(serviceID string) []*Subscriber {
	wm.mu.RLock()
	defer wm.mu.RUnlock()

	subs := wm.subscribers[serviceID]
	if len(subs) == 0 {
		return nil
	}

	c := make([]*Subscriber, len(subs))
	copy(c, subs)
	return c
}
