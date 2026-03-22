package watcher

func (wm *WatcherManager) Notify(serviceID string, update *ConfigUpdate) {
	wm.mu.RLock()
	subs := wm.subscribers[serviceID]
	wm.mu.RUnlock()

	for _, sub := range subs {
		if update.Version > sub.CurrentVersion {
			select {
			case sub.UpdateChan <- update:
				sub.CurrentVersion = update.Version
			default:
			}
		}
	}
}
