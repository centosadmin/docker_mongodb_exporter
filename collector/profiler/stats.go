package profiler

import (
	"sync"
)

type Stats struct {
	sync.Mutex
	Count          float64
	DocsExamined   float64
	KeysExamined   float64
	KeyUpdates     float64
	NReturned      float64
	NumYields      float64
	WriteConflicts float64
	Millis         float64
}

func (s *Stats) Get() *Stats {
	s.Lock()
	defer s.Unlock()
	return s
}
