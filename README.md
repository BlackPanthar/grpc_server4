# grpc_server4
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


