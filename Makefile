gen-proto:
	@echo "Generating Go files"
	protoc --go_out=plugins=grpc:./ ./proto/message.proto

test:
	@echo "Running test"
	cd __test__ && go test ./...