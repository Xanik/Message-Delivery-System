# Technical Assesssment to build a message delivery system using Go - GRPC

Design and implement (with tests) a message delivery system using the Go programming
language. Communication between nodes in the system must occur at the network layer.

## Outline

The task is to create a message delivery system in Go using GPRC microservice, that behaves as a simple in-memory key-value storage.

Payload will be simple structure – a message value, indexed by user_ids (IDs).

Go is an errors first language – this approach would be used to write a safe program.

```
Run Go Test
*make test

Run Go Vet
*make vet
```
