package inventory

import (
	"sync/atomic"
)

// Use atomic to maintain simple incremental values.
// No need to RLock/Lock here.
type Statistics struct {
	checkouts atomic.Uint64
	checkins  atomic.Uint64
	search    atomic.Uint64
	requests  atomic.Uint64
}

func (s *Statistics) IncSearch() {
	s.search.Add(1)
}

func (s *Statistics) IncCheckouts() {
	s.checkouts.Add(1)
}

func (s *Statistics) IncCheckins() {
	s.checkins.Add(1)
}

func (s *Statistics) IncRequests() {
	s.requests.Add(1)
}

func (s *Statistics) DisplayStats() map[string]uint64 {
	m := make(map[string]uint64)
	m["requests"] = s.requests.Load()
	m["checkouts"] = s.checkouts.Load()
	m["checkins"] = s.checkins.Load()
	m["searches"] = s.search.Load()
	return m
}
