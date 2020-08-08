package main

import (
	"net/http"
	"log"
	"os"
	"strconv"
	"encoding/json"
	"github.com/gorilla/mux"

	muxtrace "github.com/signalfx/signalfx-go-tracing/contrib/gorilla/mux"
	//httptrace "github.com/signalfx/signalfx-go-tracing/contrib/net/http"
	"github.com/signalfx/signalfx-go-tracing/tracing"
	//"github.com/signalfx/signalfx-go-tracing/ddtrace/tracer"
	//"github.com/opentracing/opentracing-go"
	//spanlog "github.com/opentracing/opentracing-go/log"
)

//var logger log.Logger
var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

//logger := log.New(os.Stdout, "[Test]", log.LstdFlags|log.Lmicroseconds|log.Llongfile)

type Message struct {
	Message string
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
//	span := tracer.StartSpan("RootHandler", tracer.ResourceName("/"))
//	defer span.Finish()

	// Retrieve Trace ID and Span ID
//	traceID := span.Context().TraceID()
//	spanID := span.Context().SpanID()
	// span.Context().TraceID undefined (type opentracing.SpanContext has no field or method TraceID)

//	logger.Printf("Start handling root request signalfx.trace_id=%d signalfx.span_id=%d", traceId, spanID)

//	span, _ := opentracing.StartSpanFromContext(r.Context(), "RootHandler")
//	defer span.Finish()
//	span.LogFields(spanlog.String("message", "root handler span"))
//
//	// span.Context().TraceID undefined (type opentracing.SpanContext has no field or method TraceID)
//	traceID := span.Context().TraceID()
//	spanID := span.Context().SpanID()

	logger.Printf("level=info message=\"Start handling root request\"")
	w.Write([]byte("Root Gorilla!\n"))
}

func PostHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("level=info message=\"Start handling post request\"")
	var message Message
	err := json.NewDecoder(r.Body).Decode(&message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusTeapot)
		logger.Printf("level=error message=\"ボディがガラ空きやで\"")
		return
	}
	logger.Printf("level=info message=\"%+v\"", message)
	w.Write([]byte(message.Message))

	//body, err := json.Marshal(r.Body)
	//if err != nil {
	//	w.WriteHeader(418)
	//	logger.Printf("ボディがガラ空きやで")
	//}
	////resp := "POST %+v Gorilla!\n", r.PostForm
	//logger.Printf(string(body))
	//w.Write(body)
}

func IdHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("level=info message=\"Start handling ID request\"")
	// {id} in path
	vars := mux.Vars(r)
	//resp := "ID %v Gorilla!\n", vars["id"]

	httpstatus := r.FormValue("httpstatus")
	//err := decoder.Decode(&httpstatus, r.URL.Query())
	//if err != nil {
	//	ResponseBadRequest(w, err.Error())
	//	return
	//}

	intHttpstatus, err := strconv.Atoi(httpstatus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if intHttpstatus < 400 {
		logger.Printf("level=info message=\"HTTP Status: " + httpstatus + "\"")
	} else {
		logger.Printf("level=error message=\"エラーです HTTP Status: " + httpstatus + "\"")
		http.Error(w, "HTTPステータスコードが" + httpstatus + "なのでエラーですYO", intHttpstatus)
	}

	w.Write([]byte(vars["id"]))
}

func GrpcHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("level=info message=\"Start handling GRPC request\"")
	w.Write([]byte("GRPC Gorilla!\n"))
}

func main() {
	// Use signalfx tracing
	tracing.Start(tracing.WithGlobalTag("stage", "demo"), tracing.WithServiceName("kikeyama_gorilla"))
	defer tracing.Stop()

	// Use ddtrace
	//tracer.Start(tracer.WithGlobalTag("stage", "demo"))
	//defer tracer.Stop()

	//r := mux.NewRouter()
	r := muxtrace.NewRouter()
	//r := muxtrace.NewRouter(muxtrace.WithServiceName("kikeyama_gorilla"))	// service name doesn't work here
	//r := httptrace.NewServeMux(httptrace.WithServiceName("kikeyama_gorilla"))
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/api/post", PostHandler).Methods("POST")
	r.HandleFunc("/api/trace/{id:[0-9a-z_-]+}", IdHandler).Queries("httpstatus", "{httpstatus}")
	r.HandleFunc("/api/grpc", GrpcHandler)

	// Bind to a port and pass our router in
	//log.Fatal(http.ListenAndServe(":9090", r))
	log.Fatal(http.ListenAndServe(":9090", r))
}
