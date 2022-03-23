package logx

import (
	"bytes"
	"context"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

type ContextKey string

const (
	LogxContext ContextKey = "logx-context"
	UrlPath     string     = "url.path"
	UrlMethod   string     = "url.method"
	UrlPattern  string     = "url.pattern"
)

type Config struct {
	ServiceName  string
	Env          string
	MinimalLevel logrus.Level
	Output       io.Writer
}

func Init(c Config) logrus.FieldLogger {
	if c.ServiceName == "" {
		log.Fatal("init logger failed, service name cannot be empty")
	}
	if c.Env == "" {
		log.Fatal("init logger failed, env cannot be empty")
	}
	if c.MinimalLevel == 0 {
		c.MinimalLevel = logrus.WarnLevel
	}
	if c.Output == nil {
		c.Output = os.Stdout
	}
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(c.Output)
	return logger.WithFields(logrus.Fields{
		"service": c.ServiceName,
		"env":     c.Env,
	})
}

var _ http.ResponseWriter = new(recordResponseWriter)
var bufferPool = sync.Pool{
	New: func() interface{} {
		return new(bytes.Buffer)
	},
}

type recordResponseWriter struct {
	w      http.ResponseWriter
	buf    *bytes.Buffer
	status int
}

func (w *recordResponseWriter) Header() http.Header {
	return w.w.Header()
}
func (w *recordResponseWriter) Write(content []byte) (int, error) {
	rsp := io.MultiWriter(w.w, w.buf)
	return rsp.Write(content)
}

func (w *recordResponseWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.w.WriteHeader(statusCode)
}

func SessionLogger(logger logrus.FieldLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), LogxContext, logger)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}
func InjectCtx(ctx context.Context, logger logrus.FieldLogger) context.Context {
	return context.WithValue(ctx, LogxContext, logger)
}

// LoggerFromCtx
// if not exist will panic
func FromCtx(ctx context.Context) logrus.FieldLogger {
	logger, ok := ctx.Value(LogxContext).(logrus.FieldLogger)
	if !ok {
		panic("please use logx.InitLogger to init logger ")
	}
	return logger
}

func ChiLoggerRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()
		chiCtx := chi.RouteContext(r.Context())
		ctxLogger := FromCtx(r.Context()).WithFields(
			map[string]interface{}{
				UrlPath:    r.URL.Path,
				UrlMethod:  r.Method,
				UrlPattern: chiCtx.RoutePattern(),
			},
		)
		r = r.WithContext(InjectCtx(r.Context(), ctxLogger))
		ctxLogger.Infof("request start at %v", now.String())
		dump, err := httputil.DumpRequest(r, false)
		if err != nil {
			ctxLogger.Warnf("failed to DumpRequest, err: %v", err)
		} else {
			ctxLogger.Infof("request dump: %s", string(dump))
		}
		buf := bufferPool.Get().(*bytes.Buffer)
		ww := recordResponseWriter{buf: buf, w: w}
		defer func() {
			buf.Reset()
			bufferPool.Put(buf)
			ctxLogger.Infof("request end, Code: %v, Response: %s, latency: %v ", ww.status, ww.buf.String(), time.Since(now))
		}()
		next.ServeHTTP(&ww, r)
	})
}
