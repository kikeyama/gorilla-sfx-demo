FROM golang:1.14 as build
ENV GO111MODULE on
WORKDIR /go/src/work
ADD main.go /go/src/work
ADD go.mod /go/src/work
ADD go.sum /go/src/work
#RUN go mod init
#RUN go mod edit -require github.com/opentracing/opentracing-go@v1.1.0
RUN CGO_ENABLED=0 go build -o /bin/gorilla-sfx-demo

FROM scratch
COPY --from=build /bin/gorilla-sfx-demo /bin/gorilla-sfx-demo
CMD ["/bin/gorilla-sfx-demo"]