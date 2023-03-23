package main

import (
	"context"
	"fmt"
	"grpc_server4/types"
	"testing"
)

func TestGetNodeInfo(t *testing.T) {
	ctx := context.Background()
	InitCcontext()
	s := &server{}
	ans, err := s.GetNodeInfo(ctx, &types.GetNodeInfoRequest{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ans.DefaultNodeInfo)
}

func TestGetSyncing(t *testing.T) {
	ctx := context.Background()
	InitCcontext()
	s := &server{}
	ans, err := s.GetSyncing(ctx, &types.GetSyncingRequest{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ans.Syncing)
}

func TestGetLatestBlock(t *testing.T) {
	ctx := context.Background()
	InitCcontext()
	s := &server{}
	ans, err := s.GetLatestBlock(ctx, &types.GetLatestBlockRequest{})
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ans.BlockId)
}
