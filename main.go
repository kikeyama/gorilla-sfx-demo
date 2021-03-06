package main

import (
	"net/http"
	"log"
	"os"
	"strconv"
	"fmt"
	"encoding/json"
	"github.com/gorilla/mux"

	"github.com/golang/protobuf/jsonpb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "github.com/kikeyama/grpc-sfx-demo/pb"

	muxtrace "github.com/signalfx/signalfx-go-tracing/contrib/gorilla/mux"
	grpctrace "github.com/signalfx/signalfx-go-tracing/contrib/google.golang.org/grpc"
	"github.com/signalfx/signalfx-go-tracing/tracing"
)

var logger = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds|log.Lshortfile)

var c pb.AnimalServiceClient
var conn *grpc.ClientConn

const (
	grpcClientServiceName = "kikeyama_grpc_client"
)

type Message struct {
	Message string
}

type Healthz struct {
	Status string `json:"status"`
}

type HTTPStatus struct {
	Code    int    `json:"code,omitempty"`
	Status  string `json:"status,omitempty"`
	Message string `json:"message,omitempty"`
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

func openGrpcClient() error {

	logger.Printf("level=info message=\"Start open gRPC client connection\"")

	grpcHost := getEnv("GRPC_HOST", "localhost")
	grpcPort := getEnv("GRPC_PORT", "50051")
	grpcAddress := grpcHost + ":" + grpcPort

	// enable signalfx trace
	// Create the client interceptor using the grpc trace package.
	si := grpctrace.StreamClientInterceptor(grpctrace.WithServiceName(grpcClientServiceName))
	ui := grpctrace.UnaryClientInterceptor(grpctrace.WithServiceName(grpcClientServiceName))

	// Set up a connection to the server.
	conn, err := grpc.Dial(
		grpcAddress, 
		grpc.WithInsecure(), 
//		grpc.WithBlock(),
		grpc.WithStreamInterceptor(si), 
		grpc.WithUnaryInterceptor(ui),
	)
	if err != nil {
		log.Printf("level=error message=\"cannot connect grpc: %v\"", err)
		return err
	}
	c = pb.NewAnimalServiceClient(conn)

	logger.Printf("level=info message=\"Finish open gRPC client connection\"")
	return nil

}

func ListAnimalsHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("level=info message=\"List animals through gRPC\"")

	// Contact the server and print out its response.
	ctx := r.Context()
	r2, err := c.ListAnimals(ctx, &pb.Empty{})
	if err != nil {
		logger.Printf("level=error message=\"failed to get response from grpc server: %v\"", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	m := jsonpb.Marshaler{EmitDefaults: true}
	m.Marshal(w, r2)
}

func GetAnimalHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	logger.Printf(fmt.Sprintf("level=info message=\"Get animal through gRPC with id: %s\"", id))

	ctx := r.Context()
	r2, err := c.GetAnimal(ctx, &pb.AnimalId{Id: id})
	if err != nil {
//		r2Json, _ := json.Marshal(r2)
		if status.Code(err) == codes.NotFound {
			logger.Printf("level=error message=\"document not found: %v\"", err)
			w.Header().Set("Content-Type", "application/json")
//			http.Error(w, string(r2Json), http.StatusNotFound)
//			http.Error(w, "{}", http.StatusNotFound)
			httpStatus := HTTPStatus{
				Code:    http.StatusNotFound,
				Status:  "Error",
				Message: status.Convert(err).Message(),
			}
			httpStatusJson, _ := json.Marshal(httpStatus)
			w.WriteHeader(http.StatusNotFound)
			w.Write(httpStatusJson)
			return
		}
		logger.Printf("level=error message=\"failed to get response from grpc server: %v\"", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	m := jsonpb.Marshaler{EmitDefaults: true}
	m.Marshal(w, r2)

//	animalJson, err := json.Marshal(r2)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		logger.Printf("level=error message=\"unable to marshall animal to json\"")
//		return
//	}
//	w.Write([]byte(animalJson))
}

func CreateAnimalHandler(w http.ResponseWriter, r *http.Request) {
	logger.Printf("level=info message=\"Create animal through gRPC\"")

	var pbAnimal pb.Animal
	err := json.NewDecoder(r.Body).Decode(&pbAnimal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		logger.Printf("level=error message=\"error in parse request body json\"")
		return
	}

	// Contact the server and print out its response.
	ctx := r.Context()
	r2, err := c.CreateAnimal(ctx, &pbAnimal)
	if err != nil {
		logger.Printf("level=error message=\"failed to get response from grpc server: %v\"", err)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusCreated)
	m := jsonpb.Marshaler{EmitDefaults: true}
	m.Marshal(w, r2)
}

func DeleteAnimalHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	logger.Printf(fmt.Sprintf("level=info message=\"Delete animal through gRPC with id: %s\"", id))

	ctx := r.Context()
	_, err := c.DeleteAnimal(ctx, &pb.AnimalId{Id: id})
	if err != nil {
//		r2Json, _ := json.Marshal(r2)
		if status.Code(err) == codes.NotFound {
			logger.Printf("level=error message=\"document not found: %v\"", err)
			w.Header().Set("Content-Type", "application/json")
//			http.Error(w, string(r2Json), http.StatusNotFound)
//			http.Error(w, "{}", http.StatusNotFound)
			httpStatus := HTTPStatus{
				Code:    http.StatusNotFound,
				Status:  "Error",
				Message: status.Convert(err).Message(),
			}
			httpStatusJson, _ := json.Marshal(httpStatus)
			w.WriteHeader(http.StatusNotFound)
			w.Write(httpStatusJson)
			return
		}
		logger.Printf("level=error message=\"failed to get response from grpc server: %v\"", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	httpStatus := HTTPStatus{
		Code:    http.StatusOK,
		Status:  "Success",
		Message: http.StatusText(http.StatusOK),
	}
	httpStatusJson, err := json.Marshal(httpStatus)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		logger.Printf("level=error message=\"unexpected error at not found\"")
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
//	m := jsonpb.Marshaler{EmitDefaults: true}
//	m.Marshal(w, r2)

	w.Write(httpStatusJson)
}

func NotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger.Printf(fmt.Sprintf("level=error message=\"page not found at %s\"", r.RequestURI))
		httperror := HTTPStatus{
			Code:    http.StatusNotFound,
			Message: http.StatusText(http.StatusNotFound),
		}
		httperrorJson, err := json.Marshal(httperror)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			logger.Printf("level=error message=\"unexpected error at not found\"")
			return
		}
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.WriteHeader(http.StatusNotFound)
		w.Write(httperrorJson)
	})
}

func main() {
	// Use signalfx tracing
	tracing.Start(tracing.WithGlobalTag("stage", "demo"), tracing.WithServiceName("kikeyama_gorilla"))
	defer tracing.Stop()

	r := muxtrace.NewRouter()

	// open gRPC connection	as client
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
	r.HandleFunc("/api/grpc/animal", ListAnimalsHandler).Methods("GET")
	r.HandleFunc("/api/grpc/animal/{id:[0-9a-f-]+}", GetAnimalHandler).Methods("GET")
	r.HandleFunc("/api/grpc/animal", CreateAnimalHandler).Methods("POST")
	r.HandleFunc("/api/grpc/animal/{id:[0-9a-f-]+}", DeleteAnimalHandler).Methods("DELETE")

	r.NotFoundHandler = NotFoundHandler()

	// Bind to a port and pass our router in
	logger.Printf("level=info message=\"Open http connection\"")
	log.Fatal(http.ListenAndServe(":9090", r))
}
