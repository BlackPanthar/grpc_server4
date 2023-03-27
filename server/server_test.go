// This file contains test functions that test each of the corresponding RPC methods of the gRPC server.
//
// Each test function creates a context, initializes the global Ccontext object, creates a new instance of the server,
// calls the corresponding RPC method, and prints the response.
//
// If there is an error in calling the RPC method, the test function logs the error using t.Error and the test fails.
package main

import (
	"context"
	"fmt"
	types "grpc_server4/proto/generated"
	"testing"
)

// TestGetNodeInfo tests the GetNodeInfo RPC method of the gRPC server
func TestGetNodeInfo(t *testing.T) {
	// create a context
	ctx := context.Background()

	// initialize the global Ccontext object
	InitCcontext()

	// create a new instance of the server
	s := &server{}

	// call the GetNodeInfo RPC method
	ans, err := s.GetNodeInfo(ctx, &types.GetNodeInfoRequest{})
	if err != nil {
		t.Error(err)
	}

	// print the response
	fmt.Println(ans.DefaultNodeInfo)
}

// TestGetSyncing tests the GetSyncing RPC method of the gRPC server
func TestGetSyncing(t *testing.T) {
	// create a context
	ctx := context.Background()

	// initialize the global Ccontext object
	InitCcontext()

	// create a new instance of the server
	s := &server{}

	// call the GetSyncing RPC method
	ans, err := s.GetSyncing(ctx, &types.GetSyncingRequest{})
	if err != nil {
		t.Error(err)
	}

	// print the response
	fmt.Println(ans.Syncing)
}

// TestGetLatestBlock tests the GetLatestBlock RPC method of the gRPC server
func TestGetLatestBlock(t *testing.T) {
	// create a context
	ctx := context.Background()

	// initialize the global Ccontext object
	InitCcontext()

	// create a new instance of the server
	s := &server{}

	// call the GetLatestBlock RPC method
	ans, err := s.GetLatestBlock(ctx, &types.GetLatestBlockRequest{})
	if err != nil {
		t.Error(err)
	}

	// print the response
	fmt.Println(ans.BlockId)
}

// TestGetABCIInfo tests the GetABCIInfo RPC method of the gRPC server
func TestGetABCIInfo(t *testing.T) {
	// create a context
	ctx := context.Background()

	// initialize the global Ccontext object
	InitCcontext()

	// create a new instance of the server
	s := &server{}

	// call the GetABCIInfo RPC method
	ans, err := s.GetABCIInfo(ctx, &types.GetABCIInfoRequest{})
	if err != nil {
		t.Error(err)
	}

	// print the response
	fmt.Println(ans.Jsonrpc)
}

// TestGetStatusInfo tests the GetStatusInfo RPC method of the gRPC server
func TestGetStatusInfo(t *testing.T) {
	// create a context
	ctx := context.Background()

	// initialize the global Ccontext object
	InitCcontext()

	// create a new instance of the server
	s := &server{}

	// call the GetStatusInfo RPC method
	ans, err := s.GetStatusInfo(ctx, &types.GetStatusInfoRequest{})
	if err != nil {
		t.Error(err)
	}

	// print the response
	fmt.Println(ans.ResponseString)
}
