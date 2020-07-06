FROM golang:1.14
ENV GO111MODULE on
WORKDIR /go/src/work
ADD main.go /go/src/work
ADD go.mod /go/src/work
ADD go.sum /go/src/work
#RUN go mod init
#RUN go mod edit -require github.com/opentracing/opentracing-go@v1.1.0
RUN go build
CMD ["go", "run", "main.go"]