package fetcher

import (
	"context"
	"crypto/tls"
	"net/http/httptrace"

	"go.opencensus.io/trace"
)

type HTTPClientTracer struct {
	ctx                 context.Context
	clientTrace         *httptrace.ClientTrace
	spanPrefix          string
	getConnSpan         *trace.Span
	dnsSpan             *trace.Span
	connectSpan         *trace.Span
	tlsHandshakeSpan    *trace.Span
	waitForResponseSpan *trace.Span
}

func NewHTTPClientTracer(ctx context.Context, spanPrefix string) *HTTPClientTracer {
	// TODO: Add annotation for all spans
	tracer := &HTTPClientTracer{
		ctx:        ctx,
		spanPrefix: spanPrefix,
	}
	clientTrace := &httptrace.ClientTrace{
		GetConn:              tracer.getConn,
		GotConn:              tracer.gotConn,
		PutIdleConn:          nil,
		GotFirstResponseByte: tracer.gotFirstResponseByte,
		Got100Continue:       nil,
		Got1xxResponse:       nil,
		DNSStart:             tracer.dnsStart,
		DNSDone:              tracer.dnsDone,
		ConnectStart:         tracer.connectStart,
		ConnectDone:          tracer.connectDone,
		TLSHandshakeStart:    tracer.tlsHandshakeStart,
		TLSHandshakeDone:     tracer.tlsHandshakeDone,
		WroteHeaderField:     nil,
		WroteHeaders:         nil,
		Wait100Continue:      nil,
		WroteRequest:         tracer.wroteRequest,
	}
	tracer.clientTrace = clientTrace
	return tracer
}

func (t *HTTPClientTracer) Trace() *httptrace.ClientTrace {
	return t.clientTrace
}

func (t *HTTPClientTracer) FinishSpans() {
	t.finishSpan(t.getConnSpan)
	t.finishSpan(t.dnsSpan)
	t.finishSpan(t.connectSpan)
	t.finishSpan(t.tlsHandshakeSpan)
	t.finishSpan(t.waitForResponseSpan)
}

func (t *HTTPClientTracer) startSpan(name string) (context.Context, *trace.Span) {
	return trace.StartSpan(t.ctx, t.spanPrefix+name)
}

func (t *HTTPClientTracer) finishSpan(span *trace.Span) {
	if span != nil {
		span.End()
	}
}

func (t *HTTPClientTracer) getConn(hostPort string) {
	_, t.getConnSpan = t.startSpan("getConn")
}

func (t *HTTPClientTracer) gotConn(connInfo httptrace.GotConnInfo) {
	t.finishSpan(t.getConnSpan)
}

func (t *HTTPClientTracer) dnsStart(info httptrace.DNSStartInfo) {
	_, t.dnsSpan = t.startSpan("dns")
}

func (t *HTTPClientTracer) dnsDone(dnsInfo httptrace.DNSDoneInfo) {
	t.finishSpan(t.dnsSpan)
}

func (t *HTTPClientTracer) connectStart(network, addr string) {
	_, t.connectSpan = t.startSpan("connect")
}

func (t *HTTPClientTracer) connectDone(network, addr string, err error) {
	// TODO: handle err
	t.finishSpan(t.connectSpan)
}

func (t *HTTPClientTracer) tlsHandshakeStart() {
	_, t.tlsHandshakeSpan = t.startSpan("tlsHandshake")
}

func (t *HTTPClientTracer) tlsHandshakeDone(state tls.ConnectionState, err error) {
	// TODO: handle err
	t.finishSpan(t.tlsHandshakeSpan)
}

func (t *HTTPClientTracer) wroteRequest(info httptrace.WroteRequestInfo) {
	_, t.waitForResponseSpan = t.startSpan("waitForResponse")
}

func (t *HTTPClientTracer) gotFirstResponseByte() {
	t.finishSpan(t.waitForResponseSpan)
}
