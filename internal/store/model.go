package store

import "time"

// Which game are *you* watching?
type TargetGame struct {
	Gym      string    // "Seattle Central College #1"
	Type     string    // "Saturday Morning Hoops"
	Level    string    // "Recreational-CoEd"
	DateTime time.Time // 2025-06-21 10:00 America/Los_Angeles
}

// What did we see on the page?
type GameStatus struct {
	Target     TargetGame
	IsFull     bool
	Found      bool // false => game not yet posted
	Players    map[string][]string
	ObservedAt time.Time
}
