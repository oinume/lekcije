package logger

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
	"time"

	"github.com/oinume/lekcije/server/config"
	"github.com/uber-go/zap"
)

var (
	AccessLogger = zap.New(zap.NewJSONEncoder(zap.RFC3339Formatter("ts")), zap.Output(os.Stdout))
	AppLogger    = zap.New(zap.NewJSONEncoder(zap.RFC3339Formatter("ts")), zap.Output(os.Stderr))
)

func init() {
	if !config.IsProductionEnv() {
		AppLogger.SetLevel(zap.DebugLevel)
	}
}

func InitializeAccessLogger(writer io.Writer) {
	AccessLogger = zap.New(
		zap.NewJSONEncoder(zap.RFC3339Formatter("ts")),
		zap.Output(zap.AddSync(writer)),
	)
}

func InitializeAppLogger(writer io.Writer) {
	AppLogger = zap.New(
		zap.NewJSONEncoder(zap.RFC3339Formatter("ts")),
		zap.Output(zap.AddSync(writer)),
	)
	if !config.IsProductionEnv() {
		AppLogger.SetLevel(zap.DebugLevel)
	}
}

type LoggingHTTPTransport struct {
	DumpHeaderBody bool
}

func (t *LoggingHTTPTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	var reqDump bytes.Buffer
	chunked := false
	if t.DumpHeaderBody {
		if len(req.Header) > 0 {
			fmt.Fprintln(&reqDump, "--- Request Header ---")
			for k, v := range req.Header {
				fmt.Fprintf(&reqDump, "%s: %s\n", k, strings.Join(v, ","))
			}
			chunked = len(req.TransferEncoding) > 0 && req.TransferEncoding[0] == "chunked"
			if len(req.TransferEncoding) > 0 {
				fmt.Fprintf(&reqDump, "Transfer-Encoding: %s\r\n", strings.Join(req.TransferEncoding, ","))
			}
			if req.Close {
				fmt.Fprintf(&reqDump, "Connection: close\r\n")
			}
		}
		if req.Body != nil {
			fmt.Fprintln(&reqDump, "--- Request Body ---")
			dump, _ := dumpRequestBody(req, chunked)
			fmt.Fprint(&reqDump, string(dump))
		}
	}

	start := time.Now()
	resp, err := http.DefaultTransport.RoundTrip(req)
	if err != nil {
		return nil, err
	}
	end := time.Now()

	var out bytes.Buffer
	fmt.Fprintf(
		&out, "%s %s %d %d\n",
		req.Method, req.URL,
		resp.StatusCode, time.Duration(end.Sub(start).Nanoseconds())/time.Millisecond,
	)

	var respDump bytes.Buffer
	if t.DumpHeaderBody {
		if len(resp.Header) > 0 {
			fmt.Fprintln(&respDump, "--- Response Header ---")
			for k, v := range resp.Header {
				fmt.Fprintf(&respDump, "%s: %s\n", k, strings.Join(v, ","))
			}
		}
		save := resp.Body
		savecl := resp.ContentLength
		if resp.Body != nil {
			save, resp.Body, _ = drainBody(resp.Body)
			fmt.Fprintln(&respDump, "--- Response Body ---")
			_, _ = io.Copy(&respDump, resp.Body)
		}
		resp.Body = save
		resp.ContentLength = savecl
	}

	if s := reqDump.String(); s != "" {
		fmt.Fprintln(&out, s)
	}
	if s := respDump.String(); s != "" {
		fmt.Fprintln(&out, s)
	}

	if !config.IsProductionEnv() {
		fmt.Println(out.String())
	}

	return resp, err
}

func (t *LoggingHTTPTransport) CancelRequest(r *http.Request) {
}

func dumpRequestBody(req *http.Request, chunked bool) ([]byte, error) {
	// https://github.com/golang/go/blob/master/src/net/http/httputil/dump.go#L187 のDumpRequestを参考にしている
	if req.Body == nil {
		return []byte{}, nil
	}

	var err error
	var save io.ReadCloser
	save, req.Body, err = drainBody(req.Body)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	var dest io.Writer = &b
	if chunked {
		dest = httputil.NewChunkedWriter(dest)
	}
	_, err = io.Copy(dest, req.Body)
	if chunked {
		dest.(io.Closer).Close()
		io.WriteString(&b, "\r\n")
	}

	req.Body = save
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func drainBody(b io.ReadCloser) (r1, r2 io.ReadCloser, err error) {
	var buf bytes.Buffer
	if _, err = buf.ReadFrom(b); err != nil {
		return nil, nil, err
	}
	if err = b.Close(); err != nil {
		return nil, nil, err
	}
	return ioutil.NopCloser(&buf), ioutil.NopCloser(bytes.NewReader(buf.Bytes())), nil
}
