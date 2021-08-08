# Technical Assesssment to build a message delivery system using Go - GRPC

Design and implement (with tests) a message delivery system using the Go programming
language. Communication between nodes in the system must occur at the network layer.

## Outline

The task is to create a message delivery system in Go using GPRC microservice, that behaves as a simple in-memory key-value storage.

Payload will be simple structure – a message value, indexed by user_ids (IDs) and a messageType value.

## Solution

The first approach was creating a mock grpc server to check if the proposed solution was possible :- This can be found in the `__test__` diectory. To run go test enter the command below in the application directory:

```
Run Go Test
make test
```

Next i created a `hub_test` directory to hold the hub server which the test client in the `client_test` directory calls in its init function before running the test functions in `client_test.go`. Also To run go test enter the command below in the application directory:

```
Run Go Test
make test
```

The main application to manually run can be found in the `app` directory. To run the server, cd into the `app/server` and enter the following command:

```
Run Server
go run .
```

to run the application and send an identity message, run the command below:

```
Run Server
go run . Who-Am-I identity
```

to run the application and send a list message, run the command below:

```
Run Server
go run . Where-am-i list
```

to run the application and send a relay message, run the command below:

```
Run Server
go run . foobar relay Userid -- Where Userid is id of user
go run . foobar relay 81
```

## Note

The appliaction can be extended on the client side to store the listIDs returned and then the userids that require the message to be relayed to them would be selected from the presaved list.
