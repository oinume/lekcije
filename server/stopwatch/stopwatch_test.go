package stopwatch

import (
	"testing"
	"time"
)

func TestSyncStopwatch_Mark(t *testing.T) {
	s := NewSync()
	s.Start()
	time.Sleep(time.Millisecond * 100)
	s.Mark("sleep1")
	time.Sleep(time.Millisecond * 100)
	s.Mark("sleep2")
	s.Stop()
	println(s.Report())
}
