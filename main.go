package main

import (
	"net/http"
	"log"
	"os"
	"strconv"
	"fmt"
	"encoding/json"
	"github.com/gorilla/mux"
//	"context"
//	"time"

	"github.com/golang/protobuf/jsonpb"

	"google.golang.org/grpc"
	pb "github.com/kikeyama/grpc-sfx-demo/pb"

	muxtrace "github.com/signalfx/signalfx-go-tracing/contrib/gorilla/mux"
	grpctrace "github.com/signalfx/signalfx-go-tracing/contrib/google.golang.org/grpc"
	"github.com/signalfx/signalfx-go-tracing/tracing"
)

//var logger log.Logger
var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

var c pb.AnimalServiceClient
var conn *grpc.ClientConn

const (
//	grpcPort     = ":50051"
//	defaultName = "world"
	grpcClientServiceName = "kikeyama_grpc_client"
)

type Message struct {
	Message string
}

type Healthz struct {
	Status string `json:"status"`
}

func RootHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("level=info message=\"Start handling root request\"")
	w.Write([]byte("Root Gorilla!\n"))
}

func HealthzHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("level=info message=\"Start handling healthz request\"")
	healthz := Healthz{"ok"}
	healthzJson, err := json.Marshal(healthz)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Printf("level=error message=\"unexpected error at /healthz\"")
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(healthzJson)
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
}

func IdHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("level=info message=\"Start handling ID request\"")
	// {id} in path
	vars := mux.Vars(r)

	httpstatus := r.FormValue("httpstatus")

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

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); !exists {
		return defaultVal
	} else {
		return value
	}
}

//func GrpcHandler(w http.ResponseWriter, r *http.Request) {
//	logger.Printf("level=info message=\"Start handling GRPC request\"")
//
//	grpcHost := getEnv("GRPC_HOST", "localhost")
////	grpcHost, exists := os.LookupEnv("GRPC_HOST")
////	if !exists {
////		grpcHost = "localhost"
////	}
//	grpcAddress := grpcHost + grpcPort
//
//	// enable signalfx trace
//	// Create the client interceptor using the grpc trace package.
//	si := grpctrace.StreamClientInterceptor(grpctrace.WithServiceName(grpcClientServiceName))
//	ui := grpctrace.UnaryClientInterceptor(grpctrace.WithServiceName(grpcClientServiceName))
//
//	// Set up a connection to the server.
//	conn, err_conn := grpc.Dial(grpcAddress, grpc.WithInsecure(), grpc.WithBlock(),
//		grpc.WithStreamInterceptor(si), grpc.WithUnaryInterceptor(ui))
//	if err_conn != nil {
//		log.Fatalf("did not connect: %v", err_conn)
//	}
//	defer conn.Close()
//	c := pb.NewDemoClient(conn)
//
//	// Contact the server and print out its response.
//	name := defaultName
//	if len(os.Args) > 1 {
//		name = os.Args[1]
//	}
//	ctx := r.Context()
////	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
////	defer cancel()
//	r2, err_r2 := c.GetMessageService(ctx, &pb.DemoRequest{Name: name})
//	if err_r2 != nil {
//		log.Fatalf("error: %v", err_r2)
//	}
//
//	w.Write([]byte(r2.GetMessage()))
//}

func openGrpcClient() error {

	logger.Printf("level=info message=\"Open gRPC client connection\"")

	grpcHost := getEnv("GRPC_HOST", "localhost")
	grpcPort := getEnv("GRPC_PORT", "50051")
	grpcAddress := grpcHost + ":" + grpcPort

	// enable signalfx trace
	// Create the client interceptor using the grpc trace package.
	si := grpctrace.StreamClientInterceptor(grpctrace.WithServiceName(grpcClientServiceName))
	ui := grpctrace.UnaryClientInterceptor(grpctrace.WithServiceName(grpcClientServiceName))

	// Set up a connection to the server.
	conn, err := grpc.Dial(grpcAddress, grpc.WithInsecure(), grpc.WithBlock(),
		grpc.WithStreamInterceptor(si), grpc.WithUnaryInterceptor(ui))
	if err != nil {
		log.Printf("level=error message=\"cannot connect grpc: %v\"", err)
		return err
	}
	c = pb.NewAnimalServiceClient(conn)

	return nil

}

func ListAnimalsHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("level=info message=\"List animals through gRPC\"")

//	grpcHost := getEnv("GRPC_HOST", "localhost")
//	grpcPort := getEnv("GRPC_PORT", "50051")
//	grpcAddress := grpcHost + ":" + grpcPort
//
//	// enable signalfx trace
//	// Create the client interceptor using the grpc trace package.
//	si := grpctrace.StreamClientInterceptor(grpctrace.WithServiceName(grpcClientServiceName))
//	ui := grpctrace.UnaryClientInterceptor(grpctrace.WithServiceName(grpcClientServiceName))
//
//	// Set up a connection to the server.
//	conn, err_conn := grpc.Dial(grpcAddress, grpc.WithInsecure(), grpc.WithBlock(),
//		grpc.WithStreamInterceptor(si), grpc.WithUnaryInterceptor(ui))
//	if err_conn != nil {
//		log.Fatalf("did not connect: %v", err_conn)
//	}
//	defer conn.Close()
//	c := pb.NewAnimalServiceClient(conn)

	// Contact the server and print out its response.
	ctx := r.Context()
	r2, err := c.ListAnimals(ctx, &pb.EmptyRequest{})
	if err != nil {
		logger.Printf("level=error message=\"failed to get response from grpc server: %v\"", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	m := jsonpb.Marshaler{EmitDefaults: true}
	m.Marshal(w, r2)

//	animalsJson, err := json.Marshal(r2.GetAnimals())
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		logger.Printf("level=error message=\"unable to marshall animals to json\"")
//		return
//	}
//	w.Write([]byte(animalsJson))
}

func GetAnimalHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	logger.Printf(fmt.Sprintf("level=info message=\"Get animal through gRPC with id: %s\"", id))

	ctx := r.Context()
	r2, err := c.GetAnimal(ctx, &pb.AnimalId{Id: id})
	if err != nil {
		logger.Printf("level=error message=\"failed to get response from grpc server: %v\"", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	m := jsonpb.Marshaler{EmitDefaults: true}
	m.Marshal(w, r2)
}

func main() {
	// Use signalfx tracing
	tracing.Start(tracing.WithGlobalTag("stage", "demo"), tracing.WithServiceName("kikeyama_gorilla"))
	defer tracing.Stop()

	r := muxtrace.NewRouter()

	err := openGrpcClient()
	if err != nil {
		logger.Fatalf("level=fatal message=\"failed to open grpc connection: %v\"", err)
	}
	defer conn.Close()

	//r := muxtrace.NewRouter(muxtrace.WithServiceName("kikeyama_gorilla"))	// service name doesn't work here
	// Routes consist of a path and a handler function.
	r.HandleFunc("/", RootHandler)
	r.HandleFunc("/healthz", HealthzHandler)
	r.HandleFunc("/api/post", PostHandler).Methods("POST")
	r.HandleFunc("/api/trace/{id:[0-9a-z_-]+}", IdHandler).Queries("httpstatus", "{httpstatus}")
//	r.HandleFunc("/api/grpc", GrpcHandler)
	r.HandleFunc("/api/grpc/animal", ListAnimalsHandler).Methods("GET")
	r.HandleFunc("/api/grpc/animal/{id:[0-9a-f-]+}", GetAnimalHandler).Methods("GET")

	// Bind to a port and pass our router in
	log.Fatal(http.ListenAndServe(":9090", r))
}
