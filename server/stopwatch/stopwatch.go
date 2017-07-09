package stopwatch

//https://github.com/oinume/lampetty-commons/blob/master/src/main/java/net/lampetty/commons/time/MarkingStopwatch.java

type Stopwatch struct{}

func New() *Stopwatch {
	return &Stopwatch{}
}

func (s *Stopwatch) Start() *Stopwatch {
	return s
}
