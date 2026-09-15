module github.com/mimokpl/miwin-toolkit/protoc-gen-go-http

go 1.25.0

require (
	github.com/mimokpl/miwin-toolkit/protoc-gen-common v1.1.1
	google.golang.org/protobuf v1.36.12
)

require google.golang.org/genproto/googleapis/api v0.0.0-20260911204522-f61a6ca850bd

replace github.com/mimokpl/miwin-toolkit/protoc-gen-common => ../protoc-gen-common
