# DOCUMENTATION

This README.md document is not completed, please go through the complete documentation in the file included in this project:

* [event-driven-kafka-communication-using-grpc-protobuf-schema.docx](event-driven-kafka-communication-using-grpc-protobuf-schema.docx)
* [AWS-Lambda-functions-communication-via-event-driven-kafka-events-validated-via-grpc-protobuf-schema.docx](AWS-Lambda-functions-communication-via-event-driven-kafka-events-validated-via-grpc-protobuf-schema.docx)
* [grpc-go.docx](grpc-go.docx)
* [Typescript.docx](Typescript.docx)

The below documentation is not yet completed, and will contain extracts of the above  documents.

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
"user-proto/pb" // Generated code from user.proto
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

To ensure this implementation remains robust for production environments, the following best practices need to be considered:

1. Schema Registry

Instead of manually distributing the user.proto file, consider using the AWS Glue Schema Registry. This allows your producers and consumers to interact with a centralized repository, enabling schema validation and compatibility checks (e.g., Forward/Backward compatibility) at runtime rather than compile time.

2. Error Handling and DLQs

Since Lambda consumes Kafka events in batches, a failure in processing one message can affect the entire batch.

Checkpointing: Ensure you are handling errors within the loop.

Dead Letter Queues (DLQs): Configure your Lambda-Kafka event source mapping to send failed batches to an SQS queue or a separate Kafka topic after the retry policy is exhausted.

3. Kafka Producer Persistence

In your producer code, you are creating a new kafka.Writer inside the HandleRequest function.

Optimization: This creates a new TCP connection on every invocation, which is expensive and inefficient.

Refactor: Initialize the kafka.Writer in the init() function or as a global variable outside of the HandleRequest function. This allows the connection to be reused across warm Lambda invocations.

4. Protobuf Best Practices

Versioning: Always keep your .proto files in a dedicated repository or a shared library module.

Compatibility: Never change the field numbers in existing messages. When adding new fields, always assign a new, unique tag number to prevent breaking changes.

# gRPC backward compatibility

Maintaining backward compatibility in gRPC is primarily managed through careful stewardship of your Protocol Buffers (.proto) files, as these define the contract between your clients and servers.

Because Protobuf uses field tags (the numbers assigned to fields) rather than names for binary serialization, you have significant flexibility to evolve your API without breaking existing communication.

## Key Principles for Backward Compatibility


1. Additive Changes are SafeYou can always add new fields to a message.

Existing clients will simply ignore these new fields because they are not part of their generated code, and the new server will handle the absence of these fields in requests from older clients (usually by defaulting to the zero-value).

2. Field Tags are Permanent

Never change a field tag number. The tag number is the identifier on the wire. Changing it is equivalent to changing the field entirely, which will break serialization for any client expecting the old tag.  Never reuse a tag number. Even if you delete a field, do not use its number for a new field.

3. How to Safely Remove Fields

To remove a field without breaking compatibility:

* Use the `reserved` keyword: This prevents future developers from accidentally reusing the field name or tag number.Protocol Buffers

* Do not delete the field entirely: Simply commenting it out is not enough because someone might accidentally reuse the tag. Marking it reserved ensures the protocol remains stable.

4. Avoid Changing Data Types

Changing the type of an existing field (e.g., from int32 to string) is a breaking change. If you need a different type, it is standard practice to add a new field with a new tag number and deprecate the old one.

* Managing Breaking Changes

When you absolutely must make a breaking change that cannot be handled via additive evolution, you should version your service.

* Package Versioning: Include a version identifier in the package name (e.g., package my.service.v1;).Side-by-Side Deployment: When you move to v2, keep the v1 service running. This allows you to migrate clients gradually over time. You can host both v1 and v2 on the same server instance.

__Esquema V1 (Original)__

En la primera versión, el mensaje era plano y sencillo.

```shell

syntax = "proto3";

package mi_empresa.usuario.v1;

message Usuario {
  int64 id = 1;
  string email = 2;
  string nombre = 3;
  string apellido = 4;
}
```

__Esquema V2 (Evolucionado)__

Aquí aplicamos las reglas de compatibilidad: reservamos los campos antiguos y añadimos la nueva estructura, garantizando que el servicio siga funcionando con clientes antiguos.

```shell
syntax = "proto3";

package mi_empresa.usuario.v2;

message Usuario {
  // Campos originales mantenidos
  int64 id = 1;
  string email = 2;

  // Campos obsoletos marcados como reservados
  // Nunca deben reutilizarse estos números (3, 4) ni nombres
  reserved 3, 4;
  reserved "nombre", "apellido";

  // Nueva estructura añadida
  NombreCompleto nombre_completo = 5;

  // Nuevo campo añadido (compatible)
  bool es_verificado = 6;
}

message NombreCompleto {
  string primer_nombre = 1;
  string segundo_nombre = 2;
  string apellidos = 3;
}
```

