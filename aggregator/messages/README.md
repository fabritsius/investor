# Messages

This directory contains protocol buffers for communication between the aggregator and plugins.

## Generating from `.proto` files

1. Make sure you have protocol buffers installed:<br>
  `go install google.golang.org/protobuf/cmd/protoc-gen-go`
2. From the project root directory run:<br>
  `protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative  ./aggregator/messages/messages.proto`