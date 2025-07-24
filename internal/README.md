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
````