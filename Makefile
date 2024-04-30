compile:
	protoc --go_out=protos/gen/ --go_opt=paths=source_relative --go-grpc_out=protos/gen/ --go-grpc_opt=paths=source_relative protos/leases.proto

clean:
	rm -rf protos/gen/*

build:
	go build