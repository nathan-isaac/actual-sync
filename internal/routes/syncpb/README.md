# Protobuf generation

The protobuf files have already been generated to `sync.pb.go` from `sync.proto`.

## Dependencies

- [`protoc` compiler](https://github.com/protocolbuffers/protobuf)
- [`proto-gen-go` package](https://pkg.go.dev/google.golang.org/protobuf/cmd/protoc-gen-go).

For more information, refer [https://developers.google.com/protocol-buffers/docs/reference/go-generated](https://developers.google.com/protocol-buffers/docs/reference/go-generated)

## Regeneration

In case of change in `sync.proto` file, run the following command after setting up the dependencies:

```shell
protoc --go_out=internal/routes/syncpb/ --go_opt=paths=source_relative --proto_path=internal/routes/syncpb/ sync.proto
```
