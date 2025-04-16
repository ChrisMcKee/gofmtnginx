package stats

import (
	"fmt"
	"sync"
	"time"
)

// Stats holds formatting statistics
type Stats struct {
	FilesProcessed int
	FilesSkipped   int
	FilesFailed    int
	StartTime      time.Time
	mu             sync.Mutex
}

// New creates a new Stats instance
func New() *Stats {
	return &Stats{
		StartTime: time.Now(),
	}
}

// IncrementProcessed increments the processed files counter
func (s *Stats) IncrementProcessed() {
	s.mu.Lock()
	s.FilesProcessed++
	s.mu.Unlock()
}

// IncrementSkipped increments the skipped files counter
func (s *Stats) IncrementSkipped() {
	s.mu.Lock()
	s.FilesSkipped++
	s.mu.Unlock()
}

// IncrementFailed increments the failed files counter
func (s *Stats) IncrementFailed() {
	s.mu.Lock()
	s.FilesFailed++
	s.mu.Unlock()
}

// Duration returns the time elapsed since start
func (s *Stats) Duration() time.Duration {
	return time.Since(s.StartTime)
}

// String returns a formatted string representation of the stats
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
