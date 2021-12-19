package tracerx

import (
	"net/http"
	"net/url"

	"log"
	"os"

	"github.com/PolarPanda611/trinity-micro/core/logx"
	"github.com/uber/jaeger-client-go"
	jaegerConfig "github.com/uber/jaeger-client-go/config"
	jaegerLog "github.com/uber/jaeger-client-go/log"
	"github.com/uber/jaeger-lib/metrics"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

const (
	defaultComponentName = "net/http"
	TraceIDKey           = "trace-id"
)

type Config struct {
	ServiceName     string
	JaegerAgentHost string
}

func Init(c Config) {
	os.Setenv("JAEGER_SERVICE_NAME", c.ServiceName)
	os.Setenv("JAEGER_AGENT_HOST", c.JaegerAgentHost)
	cfg, err := jaegerConfig.FromEnv()
	if err != nil {
		// parsing errors might happen here, such as when we get a string where we expect a number
		log.Printf("Could not parse Jaeger env vars: %s", err.Error())
		return
	}

	cfg.ServiceName = c.ServiceName
	cfg.Sampler.Type = "const"
	cfg.Sampler.Param = 1
	jLogger := jaegerLog.StdLogger
	jMetricsFactory := metrics.NullFactory
	if _, err := cfg.InitGlobalTracer(
		c.ServiceName,
		jaegerConfig.Logger(jLogger),
		jaegerConfig.Metrics(jMetricsFactory),
	); err != nil {
		log.Printf("Could not initialize jaeger tracer: %s", err.Error())
	}
}

type mwOptions struct {
	opNameFunc    func(r *http.Request) string
	spanObserver  func(span opentracing.Span, r *http.Request)
	urlTagFunc    func(u *url.URL) string
	componentName string
}

// MWOption controls the behavior of the Middleware.
type MWOption func(*mwOptions)

// OperationNameFunc returns a MWOption that uses given function f
// to generate operation name for each server-side span.
func OperationNameFunc(f func(r *http.Request) string) MWOption {
	return func(options *mwOptions) {
		options.opNameFunc = f
	}
}

func ChiOpenTracer(options ...MWOption) func(next http.Handler) http.Handler {
	opts := mwOptions{
		opNameFunc: func(r *http.Request) string {
			return "HTTP " + r.Method
		},
		spanObserver: func(span opentracing.Span, r *http.Request) {},
		urlTagFunc: func(u *url.URL) string {
			return u.String()
		},
	}
	for _, opt := range options {
		opt(&opts)
	}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tr := opentracing.GlobalTracer()
			carrier := opentracing.HTTPHeadersCarrier(r.Header)
			ctx, _ := tr.Extract(opentracing.HTTPHeaders, carrier)
			op := opts.opNameFunc(r)
			sp := tr.StartSpan(op, ext.RPCServerOption(ctx))
			ext.HTTPMethod.Set(sp, r.Method)
			ext.HTTPUrl.Set(sp, opts.urlTagFunc(r.URL))
			opts.spanObserver(sp, r)

			// set component name, use "net/http" if caller does not specify
			componentName := opts.componentName
			if componentName == "" {
				componentName = defaultComponentName
			}
			ext.Component.Set(sp, componentName)
			sc, ok := sp.Context().(jaeger.SpanContext)
			var traceID string
			if ok {
				traceID = sc.TraceID().String()
			}
			newCtx := logx.InjectCtx(r.Context(), logx.FromCtx(r.Context()).WithField(TraceIDKey, traceID))
			r = r.WithContext(opentracing.ContextWithSpan(newCtx, sp))
			ww := &copyStatusWriter{
				w: w,
			}
			next.ServeHTTP(ww, r)
			ext.HTTPStatusCode.Set(sp, uint16(ww.status))
			sp.Finish()
		})
	}
}

type copyStatusWriter struct {
	w      http.ResponseWriter
	status int
}

func (w *copyStatusWriter) Header() http.Header {
	return w.w.Header()
}
func (w *copyStatusWriter) Write(content []byte) (int, error) {
	return w.w.Write(content)
}

func (w *copyStatusWriter) WriteHeader(statusCode int) {
	w.status = statusCode
	w.w.WriteHeader(statusCode)
}
