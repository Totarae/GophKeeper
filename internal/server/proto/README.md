Генерация proto
````
protoc -I=internal/server/proto --go_out=internal/pkg/proto_gen --go-grpc_out=internal/pkg/proto_gen internal/server/proto/gohpkeeper.proto
````

Методы
````
grpcurl -plaintext localhost:50051 list
````
Ответ
````
gophkeeper.v1.AuthService
gophkeeper.v1.DataService
grpc.reflection.v1.ServerReflection
grpc.reflection.v1alpha.ServerReflection
````

