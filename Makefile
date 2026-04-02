.PHONY: all proto fmt test

PROTO_SRC := proto/lxdr/v1/lxdr.proto
PROTO_OUT := lxdr/lxdr.pb.go

all: proto fmt test

proto: $(PROTO_OUT)

$(PROTO_OUT): $(PROTO_SRC)
	protoc --go_out=. --proto_path=. $(PROTO_SRC)
	mv lxdr/lxdr/lxdr.pb.go $(PROTO_OUT)
	rmdir lxdr/lxdr

fmt:
	gofmt -w lxdr/*.go

test:
	env GOCACHE=$(CURDIR)/.gocache go test ./...
