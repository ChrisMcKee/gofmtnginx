package stats

import (
	"fmt"
	"sync"
	"time"
)

type Stats struct {
	FilesProcessed int
	FilesSkipped   int
	FilesFailed    int
	StartTime      time.Time
	mu             sync.Mutex
}

func New() *Stats {
	return &Stats{
		StartTime: time.Now(),
	}
}

func (s *Stats) IncrementProcessed() {
	s.mu.Lock()
	s.FilesProcessed++
	s.mu.Unlock()
}

func (s *Stats) IncrementSkipped() {
	s.mu.Lock()
	s.FilesSkipped++
	s.mu.Unlock()
}

func (s *Stats) IncrementFailed() {
	s.mu.Lock()
	s.FilesFailed++
	s.mu.Unlock()
}

func (s *Stats) Duration() time.Duration {
	return time.Since(s.StartTime)
}

func (s *Stats) String() string {
	return fmt.Sprintf("\nFormatting Statistics:\n"+
		"✅ Files processed: %d\n"+
		"💨 Files skipped: %d\n"+
		"❌ Files failed: %d\n"+
		"⏱️ Total time: %v\n",
		s.FilesProcessed,
		s.FilesSkipped,
		s.FilesFailed,
		s.Duration())
}
