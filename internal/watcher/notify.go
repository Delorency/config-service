package watcher

import (
	"log"
	"time"
)

func (wm *WatcherManager) Notify(serviceID string, update *ConfigUpdate) {
	subs := wm.getSubscribersCopy(serviceID)
	if len(subs) == 0 {
		return
	}

	for _, sub := range subs {
		func() {
			sub.mu.Lock()
			defer sub.mu.Unlock()

			if update.Version <= sub.CurrentVersion {
				return
			}
			sub.CurrentVersion = update.Version
		}()

		select {
		case sub.UpdateChan <- update:

		case <-time.After(1 * time.Second):
			log.Printf("Subscriber %s for service %s is slow, dropping update",
				sub.ID, serviceID)
		}
	}
}
