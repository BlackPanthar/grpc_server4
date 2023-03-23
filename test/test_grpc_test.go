package test

import (
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"grpc_server4/types"
	"log"
	"testing"
)

const (
	GRPC_SERVER_ADDRESS = "localhost:9090"
	CHAIN_ID            = "osmosis-1"
)

func TestGrpc(t *testing.T) {
	grpcConn, err := grpc.Dial(
		GRPC_SERVER_ADDRESS, // your gRPC server address.
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
	c := types.NewGrpcQueryServiceClient(grpcConn)
	ctx := context.Background()

	nodeinfo, err := c.GetNodeInfo(ctx, &types.GetNodeInfoRequest{})
	if err != nil {
		log.Fatalf("GetNodeInfo err: %v", err)
		return
	}
	out, err := json.Marshal(nodeinfo)
	if err != nil {
		log.Fatalf("GetNodeInfo err: %v", err)
		return
	}
	fmt.Println(string(out))

	syncingInfo, err := c.GetSyncing(ctx, &types.GetSyncingRequest{})
	if err != nil {
		log.Fatalf("GetSyncing err: %v", err)
		return
	}
	fmt.Println("Syncing:", syncingInfo.Syncing)

	latestBlock, err := c.GetLatestBlock(ctx, &types.GetLatestBlockRequest{})
	if err != nil {
		log.Fatalf("GetLatestBlock err: %v", err)
	}
	out, err = json.Marshal(latestBlock)
	if err != nil {
		log.Fatalf("GetLatestBlock err: %v", err)
		return
	}
	fmt.Println(string(out))

	heightBlock, err := c.GetBlockByHeight(ctx, &types.GetBlockByHeightRequest{Height: 8658239})
	if err != nil {
		log.Fatalf("GetBlockByHeight err: %v", err)
		return
	}
	out, err = json.Marshal(heightBlock)
	if err != nil {
		log.Fatalf("GetBlockByHeight err: %v", err)
		return
	}
	fmt.Println(string(out))

	latestValidators, err := c.GetLatestValidatorSet(ctx, &types.GetLatestValidatorSetRequest{})
	if err != nil {
		log.Fatalf("GetLatestValidatorSet err: %v", err)
		return
	}
	out, err = json.Marshal(latestValidators)
	if err != nil {
		log.Fatalf("GetLatestValidatorSet err: %v", err)
		return
	}
	fmt.Println(string(out))

	heightValidators, err := c.GetValidatorSetByHeight(ctx, &types.GetValidatorSetByHeightRequest{Height: 8658239})
	if err != nil {
		log.Fatalf("GetValidatorSetByHeight err: %v", err)
	}
	out, err = json.Marshal(heightValidators)
	if err != nil {
		log.Fatalf("GetValidatorSetByHeight err: %v", err)
		return
	}
	fmt.Println(string(out))

}
