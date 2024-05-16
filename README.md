## Naive Bayes

# Install GRPC

go install google.golang.org/grpc@latest
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
go install golang.org/x/tools/cmd/goyacc@latest

## Generate the

protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative
protos/service.proto

# Creates the .y grammar file

goyacc -o calc.go calc.y

# Compiles everything .

go build -o calc main.go calc.go

# Run

./calc

# Go mod tidy

go mod tidy
