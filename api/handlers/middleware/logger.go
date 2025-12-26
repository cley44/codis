package middleware

import (
	"bytes"
	"io"
	"log/slog"
	"runtime/debug"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/samber/oops"
)

const maxBodySize = 2000

func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Recovery middleware won't catch errors from this middleware.
		// For safety reason, i catch panics again.
		defer func() {
			// recover
			if err := recover(); err != nil {
				slog.Error("Panic recovered", slog.Any("error", err), slog.Any("stacktrace", string(debug.Stack())))
			}
		}()

		pathIsAuth := strings.HasPrefix(ctx.Request.URL.Path, "/auth")

		// dump request body
		br := newBodyReader(ctx.Request.Body, maxBodySize, !pathIsAuth)
		ctx.Request.Body = br

		// dump response body
		bw := newBodyWriter(ctx.Writer, maxBodySize, !pathIsAuth)
		ctx.Writer = bw

		// some evil middlewares modify this values
		path := ctx.Request.URL.Path

		if ctx.Request.URL.RawQuery != "" {
			path += "?" + ctx.Request.URL.RawQuery
		}

		if ctx.Request.URL.RawFragment != "" {
			path += "#" + ctx.Request.URL.RawFragment
		}

		ctx.Next()

		statusCode := ctx.Writer.Status()

		end := time.Now()

		entry := slog.With(slog.Group("http",
			slog.Int("status", statusCode),
			// slog.Int("latency_ms", int(latency.Milliseconds())),
			slog.String("time", end.Format(time.RFC3339)),
			slog.String("method", ctx.Request.Method),
			slog.String("path", path),
			// slog.String("ip", c.ClientIP()),
			slog.String("user-agent", ctx.Request.UserAgent()),
		))

		if len(ctx.Errors) > 0 {
			entry.Error(
				ctx.Errors.String(),
				slog.Any("error", ctx.Errors[0]),
				slog.String("req_body", br.String()),
				slog.String("res_body", bw.String()),
			)
			test, _ := oops.AsOops(ctx.Errors[0])
			println(test.Stacktrace())
		}
	}
}

type bodyWriter struct {
	gin.ResponseWriter
	body    *bytes.Buffer
	maxSize int
	bytes   int
}

// implements gin.ResponseWriter.
func (w *bodyWriter) Write(b []byte) (int, error) {
	length := len(b)

	if w.body != nil {
		if w.body.Len()+length > w.maxSize {
			w.body.Truncate(min(w.maxSize, length, w.body.Len()))
			w.body.Write(b[:min(w.maxSize-w.body.Len(), length)])
		} else {
			w.body.Write(b)
		}
	}

	w.bytes += length //nolint:staticcheck
	return w.ResponseWriter.Write(b)
}

func (w *bodyWriter) String() string {
	if w.body != nil {
		return w.body.String()
	}
	return ""
}

func newBodyWriter(writer gin.ResponseWriter, maxSize int, recordBody bool) *bodyWriter {
	var body *bytes.Buffer
	if recordBody {
		body = bytes.NewBufferString("")
	}

	return &bodyWriter{
		ResponseWriter: writer,
		body:           body,
		maxSize:        maxSize,
		bytes:          0,
	}
}

type bodyReader struct {
	io.ReadCloser
	body    *bytes.Buffer
	maxSize int
	bytes   int
}

// implements io.Reader.
func (r *bodyReader) Read(b []byte) (int, error) {
	n, err := r.ReadCloser.Read(b)
	if r.body != nil && r.body.Len() < r.maxSize {
		if r.body.Len()+n > r.maxSize {
			r.body.Write(b[:min(r.maxSize-r.body.Len(), n)])
		} else {
			r.body.Write(b[:n])
		}
	}
	r.bytes += n
	return n, err
}

func (r *bodyReader) String() string {
	if r.body != nil {
		return r.body.String()
	}
	return ""
}

func newBodyReader(reader io.ReadCloser, maxSize int, recordBody bool) *bodyReader {
	var body *bytes.Buffer
	if recordBody {
		body = bytes.NewBufferString("")
	}

	return &bodyReader{
		ReadCloser: reader,
		body:       body,
		maxSize:    maxSize,
		bytes:      0,
	}
}
