gen-proto:
	@echo "Generating Go files"
	protoc --go_out=plugins=grpc:./ ./proto/message.proto

test:
	@echo "Running Mock test"
	cd __test__ && go test ./...
	@echo "Running client test"
	go test -v ./... -cover

vet:
	@echo "Running Go Vet"
	go vet -v ./...