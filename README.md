# grpc-playground

## Install Taskfile

```bash
brew install go-task/tap/go-task
```

## Install gRPC protot gen

```bash
brew install protobuf
```

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

## Protoc gen

Example:
```bash
# cd error

protoc \
--plugin=protoc-gen-go=./bin/protoc-gen-go \
--plugin=protoc-gen-go-grpc=./bin/protoc-gen-go-grpc \
--go_out=. --go_opt=paths=source_relative \
--go-grpc_out=. --go-grpc_opt=paths=source_relative \
    api/api.proto
```
