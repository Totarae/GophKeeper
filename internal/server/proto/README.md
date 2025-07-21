Генерация proto
````
protoc -I=internal/server/proto --go_out=internal/pkg/proto_gen --go-grpc_out=internal/pkg/proto_gen internal/server/proto/gohpkeeper.proto
````