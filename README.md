# how to use 
### first you need to down load the packages
use```go mod tidy``` 

### second you need to build the grpc server and start it
```
cd server
go build
./server
```

### and then you can start your test
to use the client
```
cd client
go build
./client   [Option] 

where [Option] can be GetNodeInfo, GetSyncing, GetLatestBlock, GetBlockByHeight [Height] e.g height: 8700000, GetLatestValidatorSet,GetValidatorSetByHeight [Height] eg height: 8658239
```
to use test grpc file
```
cd test
go test
```
to use test state_tracker
```
cd state_tracker
go test
```

