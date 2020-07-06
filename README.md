# Gorilla Mux Demo

Compatible with Go 1.14, [OpenTracing Go v1.1.0](https://github.com/opentracing/opentracing-go/tree/v1.1.0)  

## Endpoints

### `/` (root)

Return simple message  

### `/api/post`

Require post body with key "message"  

### `/api/trace/{id:[0-9a-z_-]+}`

Require query string with key "httpstatus" and int value (100-500) to generate mock http status code  

### `/api/grpc`

WIP  
