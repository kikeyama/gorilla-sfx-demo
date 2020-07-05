# Gorilla Mux Demo

Compatible with Go 1.14, [OpenTracing Go v1.1.0](https://github.com/opentracing/opentracing-go/tree/v1.1.0)  

## Known issue??

With opentracing-go v1.2.0, signalfx tracer causes panic as following.  

```
2020/07/05 14:44:45 http: panic serving 172.17.0.1:54674: interface conversion: *tracer.tracer is not opentracing.Tracer: missing method Extract
goroutine 34 [running]:
net/http.(*conn).serve.func1(0xc000208000)
    /usr/local/go/src/net/http/server.go:1772 +0x139
panic(0x8050e0, 0xc0002006f0)
    /usr/local/go/src/runtime/panic.go:975 +0x3e3
github.com/signalfx/signalfx-go-tracing/ddtrace/tracer.(*span).Tracer(0xc000230000, 0x872aff, 0xc000210940)
    /go/pkg/mod/github.com/signalfx/signalfx-go-tracing@v1.2.0/ddtrace/tracer/span.go:85 +0x45
github.com/opentracing/opentracing-go.ContextWithSpan(0x90c160, 0xc0002160c0, 0x9106e0, 0xc000230000, 0xc000230000, 0x910ae0)
    /go/pkg/mod/github.com/opentracing/opentracing-go@v1.2.0/gocontext.go:13 +0x9d
github.com/signalfx/signalfx-go-tracing/ddtrace/tracer.ContextWithSpan(...)
    /go/pkg/mod/github.com/signalfx/signalfx-go-tracing@v1.2.0/ddtrace/tracer/context.go:11
github.com/signalfx/signalfx-go-tracing/ddtrace/tracer.StartSpanFromContext(0x90c160, 0xc0002160c0, 0x872aff, 0xc, 0xc000200630, 0x5, 0x5, 0x0, 0x40f05d, 0x8535c0, ...)
    /go/pkg/mod/github.com/signalfx/signalfx-go-tracing@v1.2.0/ddtrace/tracer/context.go:36 +0xea
github.com/signalfx/signalfx-go-tracing/contrib/internal/httputil.TraceAndServe(0x9043e0, 0xc000160000, 0x90b220, 0xc000228000, 0xc00021e000, 0x8720d4, 0xa, 0x8701c4, 0x1, 0x0, ...)
    /go/pkg/mod/github.com/signalfx/signalfx-go-tracing@v1.2.0/contrib/internal/httputil/trace.go:39 +0x444
github.com/signalfx/signalfx-go-tracing/contrib/gorilla/mux.(*Router).ServeHTTP(0xc000020f30, 0x90b220, 0xc000228000, 0xc00021e000)
    /go/pkg/mod/github.com/signalfx/signalfx-go-tracing@v1.2.0/contrib/gorilla/mux/mux.go:104 +0x203
net/http.serverHandler.ServeHTTP(0xc000164000, 0x90b220, 0xc000228000, 0xc00021e000)
    /usr/local/go/src/net/http/server.go:2807 +0xa3
net/http.(*conn).serve(0xc000208000, 0x90c160, 0xc000216000)
    /usr/local/go/src/net/http/server.go:1895 +0x86c
created by net/http.(*Server).Serve
    /usr/local/go/src/net/http/server.go:2933 +0x35c
```

Use opentracing-go v1.1.0 to fix this issue.

```
go mod edit -require github.com/opentracing/opentracing-go@v1.1.0
```
