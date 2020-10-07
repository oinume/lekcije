package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"time"

	"go.opencensus.io/trace"

	"github.com/oinume/lekcije/server/config"
	"github.com/oinume/lekcije/server/open_census"
)

func main() {
	config.MustProcessDefault()
	exporter, flush, err := open_census.NewExporter(
		config.DefaultVars,
		"go-quickstart",
		true,
	)
	if err != nil {
		log.Fatalf("NewExporter failed: %v", err)
	}
	defer flush()
	trace.RegisterExporter(exporter)

	// 2. Configure 100% sample rate, otherwise, few traces will be sampled.
	trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})

	// 3. Create a span with the background context, making this the parent span.
	// A span must be closed.
	ctx, span := trace.StartSpan(context.Background(), "main")
	// 5b. Make the span close at the end of this function.
	defer span.End()

	for i := 0; i < 10; i++ {
		doWork(ctx)
	}
}

func doWork(ctx context.Context) {
	// 4. Start a child span. This will be a child span because we've passed
	// the parent span's ctx.
	_, span := trace.StartSpan(ctx, "doWork")
	// 5a. Make the span close at the end of this function.
	defer span.End()

	fmt.Println("doing busy work")
	time.Sleep(80 * time.Millisecond)
	buf := bytes.NewBuffer([]byte{0xFF, 0x00, 0x00, 0x00})
	num, err := binary.ReadVarint(buf)
	if err != nil {
		// 6. Set status upon error
		span.SetStatus(trace.Status{
			Code:    trace.StatusCodeUnknown,
			Message: err.Error(),
		})
	}

	// 7. Annotate our span to capture metadata about our operation
	span.Annotate([]trace.Attribute{
		trace.Int64Attribute("bytes to int", num),
	}, "Invoking doWork")
	time.Sleep(20 * time.Millisecond)
}
