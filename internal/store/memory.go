package store

import (
	"sync"
)

// memoryCache holds the most recent GameStatus observed for each TargetGame.
// It is safe for concurrent access by multiple goroutines.
var memoryCache struct {
	mu       sync.RWMutex
	statuses map[string]*GameStatus
}

func init() {
	memoryCache.statuses = make(map[string]*GameStatus)
}

// SaveGameStatus stores/updates the latest GameStatus for its corresponding
// TargetGame in the cache.
func SaveGameStatus(s *GameStatus) {
	memoryCache.mu.Lock()
	defer memoryCache.mu.Unlock()
	key := s.Target.GameKey
	memoryCache.statuses[key] = s
}

// GetGameStatus retrieves the previously stored GameStatus for the provided
// TargetGame. It returns nil if no GameStatus has been cached yet.
func GetGameStatus(key string) *GameStatus {
	memoryCache.mu.RLock()
	defer memoryCache.mu.RUnlock()
	return memoryCache.statuses[key]
}
