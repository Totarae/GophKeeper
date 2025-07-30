zap logger for logging
bcrypt 
jwt
grpc

testify
data-dog/go-sqlmock

Сделал отдельно репу для пользователя, отдельно для информации

Накрытие
````
go test -coverprofile="coverage.out" ./...
go tool cover -html="coverage.out"
> go tool cover -func coverage.out
go test -covermode=atomic -coverpkg=./... -coverprofile=coverage.out ./...

````

Собрать
````
docker compose build --no-cache
````