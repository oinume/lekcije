package stopwatch

import (
	"bytes"
	"fmt"
	"sync"
	"time"
)

//https://github.com/oinume/lampetty-commons/blob/master/src/main/java/net/lampetty/commons/time/MarkingStopwatch.java

type Stopwatch interface {
	Start() Stopwatch
	Stop() Stopwatch
	Mark(name string) Stopwatch
	GetTotalTime() time.Duration
	Report() string
}

type SyncStopwatch struct {
	startedAt time.Time
	stoppedAt time.Time
	marks     []*mark
	sync.Mutex
}

type mark struct {
	name string
	at   time.Time
}

func NewSync() *SyncStopwatch {
	return &SyncStopwatch{}
}

func (s *SyncStopwatch) Start() Stopwatch {
	s.Lock()
	s.startedAt = time.Now()
	s.Unlock()
	return s
}

func (s *SyncStopwatch) Stop() Stopwatch {
	s.Lock()
	if s.stoppedAt.IsZero() {
		s.stoppedAt = time.Now()
	}
	s.Unlock()
	return s
}

func (s *SyncStopwatch) Mark(name string) Stopwatch {
	//if (stoppedTime != -1) {
	//	throw new IllegalStateException("Already stopped.");
	//}
	s.Lock()
	s.marks = append(s.marks, &mark{
		name: name,
		at:   time.Now(),
	})
	s.Unlock()
	return s
}

func (s *SyncStopwatch) GetTotalTime() time.Duration {
	return s.stoppedAt.Sub(s.startedAt)
}

func (s *SyncStopwatch) Report() string {
	s.Stop()
	s.Lock()
	s.marks = append(s.marks, &mark{
		name: "__stop__",
		at:   s.stoppedAt,
	})
	s.Unlock()

	var b bytes.Buffer
	fmt.Fprintf(&b, "%-41.40s %-11s %-15s %s\n", "NAME", "TIME(ms)", "CUMULATIVE(ms)", "PERCENTAGE")

	previousTime := s.startedAt
	totalTime := s.GetTotalTime()
	for _, mark := range s.marks {
		duration := mark.at.Sub(previousTime)
		cumulative := mark.at.Sub(s.startedAt)
		percentage := float64(duration) / float64(totalTime) * 100
		fmt.Fprintf(
			&b, " %-41.40s %-11d %-15d %.2f%%\n",
			mark.name, duration/time.Millisecond,
			cumulative/time.Millisecond,
			float64(percentage),
		)

		previousTime = mark.at
	}
	return b.String()
}
