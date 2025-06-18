package store

import "sync"

// memoryCache holds the most recent GameStatus observed. It is safe for
// concurrent access by multiple goroutines.
var memoryCache struct {
	mu     sync.RWMutex
	status *GameStatus
}

// SaveGameStatus stores the latest GameStatus in the cache.
func SaveGameStatus(s *GameStatus) {
	memoryCache.mu.Lock()
	defer memoryCache.mu.Unlock()
	memoryCache.status = s
}

// GetGameStatus retrieves the previously stored GameStatus (if any).
// It returns nil if no GameStatus has been cached yet.
func GetGameStatus() *GameStatus {
	memoryCache.mu.RLock()
	defer memoryCache.mu.RUnlock()
	return memoryCache.status
}
