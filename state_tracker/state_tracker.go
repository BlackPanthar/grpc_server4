package state_tracker

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"grpc_server4/types"
	"log"
	"time"
)

type StateTrackerStruct struct {
	TestResult []TestResult `json:"test_result"`
}
type TestResult struct {
	Height int64  `json:"height"`
	Hash   string `json:"hash"`
}

const GRPC_SERVER_ADDRESS = "localhost:9090"

func StateTracker() {
	ctx := context.Background()
	grpcConn, err := grpc.Dial(
		GRPC_SERVER_ADDRESS, // your gRPC server address.
		grpc.WithInsecure(), // The SDK doesn't support any transport security mechanism.
	)
	if err != nil {
		log.Fatalf("dial err: %v", err)
	}
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			fmt.Println("grpcConn.Close() err:", err)
		}
	}(grpcConn)
	c := types.NewGrpcQueryServiceClient(grpcConn)
	resp, err := c.GetLatestBlock(ctx, &types.GetLatestBlockRequest{})
	if err != nil {
		log.Fatalf("get latest block: %v", err)
	}
	ans := &StateTrackerStruct{TestResult: make([]TestResult, 0, 0)}
	height := resp.Block.Header.Height + 1
	ans.TestResult = append(ans.TestResult, TestResult{
		Height: resp.Block.Header.Height,
		Hash:   hex.EncodeToString(resp.BlockId.Hash),
	})
	time.Sleep(30 * time.Second)
	for i := 0; i < 4; i++ {
		block, err := c.GetBlockByHeight(ctx, &types.GetBlockByHeightRequest{Height: height})
		if err != nil {
			log.Fatalf("GetBlockByHeight: %v", err)
		}
		height++
		ans.TestResult = append(ans.TestResult, TestResult{
			Height: block.Block.Header.Height,
			Hash:   hex.EncodeToString(resp.BlockId.Hash),
		})
	}
	jsonStr, err := json.Marshal(ans)
	fmt.Println(string(jsonStr))
}
