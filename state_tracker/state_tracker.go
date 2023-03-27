// This package contains a function StateTracker()
//
// It tracks the blockchain state of a Tendermint node via gRPC by calling the GetLatestBlock()
// and GetBlockByHeight() methods of the GrpcQueryServiceClient provided by the gRPC server.
//
//It stores the results of these calls in a StateTrackerStruct struct and
//outputs them as JSON to a file named "info.json".

package state_tracker

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	types "grpc_server4/proto/generated"
	"log"
	"os"
	"time"

	"google.golang.org/grpc"
)

// Struct that holds the state tracker test result.
type StateTrackerStruct struct {
	TestResult []TestResult `json:"test_result"`
}

// Struct that holds a single test result of a specific height and hash.
type TestResult struct {
	Height int64  `json:"height"`
	Hash   string `json:"hash"`
}

// Address of the gRPC server.
const GRPC_SERVER_ADDRESS = "localhost:9090"

// Function that tracks the state of the Tendermint node.
func StateTracker() {
	ctx := context.Background()
	// Create a gRPC client connection.
	grpcConn, err := grpc.Dial(
		GRPC_SERVER_ADDRESS,
		grpc.WithInsecure(),
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
	// Create a gRPC client instance.
	c := types.NewGrpcQueryServiceClient(grpcConn)
	// Get the latest block of the Tendermint node.
	resp, err := c.GetLatestBlock(ctx, &types.GetLatestBlockRequest{})
	if err != nil {
		log.Fatalf("get latest block: %v", err)
	}
	// Initialize the state tracker test result with the latest block.
	ans := &StateTrackerStruct{TestResult: make([]TestResult, 0, 0)}
	height := resp.Block.Header.Height + 1
	ans.TestResult = append(ans.TestResult, TestResult{
		Height: resp.Block.Header.Height,
		Hash:   hex.EncodeToString(resp.BlockId.Hash),
	})
	// Wait for 30 seconds to let the node create new blocks.
	time.Sleep(30 * time.Second)
	// Get the next five blocks and add them to the state tracker test result.
	for i := 0; i < 5; i++ {
		block, err := c.GetBlockByHeight(ctx, &types.GetBlockByHeightRequest{Height: height})
		if err != nil {
			log.Fatalf("GetBlockByHeight: %v", err)
		}
		height++
		ans.TestResult = append(ans.TestResult, TestResult{
			Height: block.Block.Header.Height,
			Hash:   hex.EncodeToString(block.BlockId.Hash),
		})
	}
	// Convert the state tracker test result to JSON and print it.
	jsonStr, err := json.Marshal(ans)
	fmt.Println(string(jsonStr))
	// Write the state tracker test result to a file.
	filePtr, err := os.Create("info.json")
	if err != nil {
		fmt.Println("create file failed! err:", err.Error())
		return
	}
	defer filePtr.Close()
	encoder := json.NewEncoder(filePtr)
	err = encoder.Encode(ans)
	if err != nil {
		fmt.Println("encode err:", err.Error())
	}
}
