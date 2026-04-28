
# Pre-requisites

Install requirements

https://grpc.io/docs/languages/go/quickstart/

```shell
PS C:\Users\user\Desktop\Documents\GoLangProjects\grpc-go> go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
PS C:\Users\user\Desktop\Documents\GoLangProjects\grpc-go> go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Verify after successfull install:

```shell
PS C:\Users\user\Desktop\Documents\GoLangProjects\grpc-go> go version
go version go1.26.2 windows/amd64

PS C:\Users\user\Desktop\Documents\GoLangProjects\grpc-go> protoc --version
libprotoc 34.1
```

Create a symbolic link:

```shell
# Copy the existing file to the name that protoc is looking for
Copy-Item "C:\Users\user\go\bin\protoc-gen-go-grpc.exe" "C:\Users\user\go\bin\protoc-gen-go_grpc.exe"
```
# Setup our go project

```shell
PS C:\Users\user\Desktop\Documents\GoLangProjects\grpc-go> go mod init
go: cannot determine module path for source directory C:\Users\user\Desktop\Documents\GoLangProjects\grpc-go (outside GOPATH, module path must be specified)

Example usage:
        'go mod init example.com/m' to initialize a v0 or v1 module
        'go mod init example.com/m/v2' to initialize a v2 module

Run 'go help mod init' for more information.

PS C:\Users\user\Desktop\Documents\GoLangProjects\grpc-go> go mod init example.com/m
go: creating new go.mod: module example.com/m

PS C:\Users\user\Desktop\Documents\GoLangProjects\grpc-go> go mod tidy
go: warning: "all" matched no packages
```

Compile:

```shell
protoc --proto_path=greet/proto `
  --go_out=greet/proto --go_opt=paths=source_relative `
  --go-grpc_out=greet/proto --go-grpc_opt=paths=source_relative `
  greet/proto/dummy.proto
```

# Makefile

Clean:

```shell
user@DESKTOP-U6PJEA9:/mnt/c/Users/user/Desktop/Documents/GoLangProjects/grpc-go$ make clean
rm greet/proto/*.pb.go
rm proto-go
rm: cannot remove 'proto-go': No such file or directory
make: *** [Makefile:32: clean] Error 1
```

build:

```shell
user@DESKTOP-U6PJEA9:/mnt/c/Users/user/Desktop/Documents/GoLangProjects/grpc-go$ make build
protoc --proto_path=greet/proto --go_out=greet/proto --go_opt=paths=source_relative --go-grpc_out=greet/proto --go-grpc_opt=paths=source_relative greet/proto/dummy.proto
go build -o proto-go-course .
no Go files in /mnt/c/Users/user/Desktop/Documents/GoLangProjects/grpc-go
make: *** [Makefile:23: build] Error 1
```

make all:

```shell
user@DESKTOP-U6PJEA9:/mnt/c/Users/user/Desktop/Documents/GoLangProjects/grpc-go$ make all
protoc --proto_path=greet/proto --go_out=greet/proto --go_opt=paths=source_relative --go-grpc_out=greet/proto --go-grpc_opt=paths=source_relative greet/proto/*.proto
go build -o server ./greet/server/main.go
```

Verify the generated file:

```shell
user@DESKTOP-U6PJEA9:/mnt/c/Users/user/Desktop/Documents/GoLangProjects/grpc-go$ ls -rlht
total 16M
-rwxrwxrwx 1 user user  162 Apr 23 06:40 '~$Steps.docx'
drwxrwxrwx 1 user user 4.0K Apr 23 06:56  greet
-rwxrwxrwx 1 user user 3.3K Apr 23 07:34  go.sum
-rwxrwxrwx 1 user user  332 Apr 23 07:34  go.mod
-rwxrwxrwx 1 user user 1.6M Apr 23 08:49  Steps.docx
-rwxrwxrwx 1 user user 1.2K Apr 23 08:54  Makefile
-rwxrwxrwx 1 user user  15M Apr 23 08:54  server
-rwxrwxrwx 1 user user 2.8K Apr 23 08:55  README.md
```

and finally, run the server from Ubuntu:

```shell
user@DESKTOP-U6PJEA9:/mnt/c/Users/user/Desktop/Documents/GoLangProjects/grpc-go$ ./server
Listening on 0.0.0.0:5003

```

# Producer and consumer AWS Lambda functions communicating via Kafka events

Lambda doesn't allow to establish persistent TCP connections between two functions like a gRPC standard socket does.

The architecture instead is: Lambda A invokes Lambda B by means of API Gateway or some sort of mechanism like Kafka. In order to maintain the gRPC structure we require to define a shared contract between both the lambda functions.

Here is a simple example of how two lambda functions can use a protobuf serialization contract to communicate via kafka events. 

__The Protocol Buffer Definition (user.proto)__

This is the single source of truth for both functions.

```shell
syntax = "proto3";
package user.v1;

message UserRequest {
  int64 id = 1;
}

message UserResponse {
  string name = 1;
}
```


__The Producer (Lambda A)__

Lambda A is triggered by an event (e.g., an API Gateway request), serializes the data, and pushes it to Kafka.

```shell
package main

import (
"context"
"github.com/aws/aws-lambda-go/lambda"
"github.com/segmentio/kafka-go"
"google.golang.org/protobuf/proto"
"user-proto/pb" // Generated code from your .proto
)

func HandleRequest(ctx context.Context, event MyEvent) error {
// 1. Create your Protobuf message
user := &pb.User{Id: event.UserID, Name: "John"}

	// 2. Serialize to bytes
	payload, _ := proto.Marshal(user)

	// 3. Send to Kafka
	writer := &kafka.Writer{
		Addr:  kafka.TCP("your-kafka-broker:9092"),
		Topic: "user-updates",
	}
	return writer.WriteMessages(ctx, kafka.Message{Value: payload})
}

func main() {
  lambda.Start(HandleRequest)
}
```

__The Consumer (Lambda B)__

Lambda B is triggered by the Kafka topic. AWS handles the integration between MSK and Lambda, passing the message payload as anevent.

```shell
package main

import (
    "context"
    "github.com/aws/aws-lambda-go/events"
    "github.com/aws/aws-lambda-go/lambda"
    "google.golang.org/protobuf/proto"
    "user-proto/pb"
)

func HandleRequest(ctx context.Context, kafkaEvent events.KafkaEvent) error {
  for _, record := range kafkaEvent.Records["user-updates"] {
    // 1. Unmarshal the bytes back into the struct
    user := &pb.User{}
    proto.Unmarshal(record.Value, user)

	// 2. Process logic
	println("Processing user:", user.Name)
  }
  return nil
}

func main() {
  lambda.Start(HandleRequest)
}
```
