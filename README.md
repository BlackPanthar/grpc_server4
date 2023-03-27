## grpc_server
grpc server

### first download packages
use```go mod tidy``` 

### start the grpc server
```
cd server
go build
./server
```

### function check
client
```
cd client
go build
./client   [Option] 

where [Option] can be GetNodeInfo, GetSyncing, GetLatestBlock, GetBlockByHeight [Height] e.g height: 8700000, GetLatestValidatorSet,GetValidatorSetByHeight [Height] eg height: 8658239
```
test grpc file
```
cd test
go test
```
test state_tracker
```
cd state_tracker
go test
```

Few Notes:

My goal is simply to build a GRpc client/server application in go that can retrieve the same result as a URL like this https://rpc.osmosis.zone/abci_ and I followed the following steps:

1. Define the protocol buffer messages and services for the Tendermint RPC API. The Tendermint RPC API specification can be found in the rpc.proto file of the Tendermint repository. The messages and services in this file are used to interact with the Tendermint node via gRPC.

2. Generate the gRPC client and server code in Go using the protoc compiler and the protoc-gen-go-grpc plugin, i specifically used liprotoc 3.14.0 and protoc-gen-go v1.25.0-devel

protoc -I . \
-I /home/emeka/code/goinstalls/protobuf-1.3.3-alpha.regen.1 \
-I /home/emeka/code/goinstalls/cosmos-sdk-0.45.1 \
-I /home/emeka/code/goinstalls/cosmos-proto-1.0.0-beta.2 \
-I /home/emeka/code/goinstalls/tendermint-0.34.13 \
-I /home/emeka/code/googleapis \
--proto_path=/home/emeka/code/goinstalls/cosmos-sdk-0.45.1/proto \
--proto_path=/home/emeka/code/goinstalls/tendermint-0.34.13/proto \
--proto_path=/home/emeka/code/googleapis \
--go_out=. \
--go-grpc_out=../types \
rpc.proto

3. Implement the gRPC server in Go. Each method should handle the incoming request, perform the appropriate action against the Tendermint node, and return the response to the client.

4. Implement the gRPC client in Go. The client should connect to the server using the appropriate gRPC connection options, and then use the generated client code to call the methods defined in the Tendermint RPC API services. Each method call should send the appropriate request to the server and wait for the response.

5. Run the server in one terminal window, and then run the client in another terminal window. The client should be able to connect to the server and invoke the Tendermint RPC API methods, which should return the same results as the example URL for example https://rpc.osmosis.zone/abci_info (changing formats may result in an interim check that both calls are not empty with no errors, as the test, we'll see)


PS: info.json is the state tracker json and is generated in the state_tracker folder
