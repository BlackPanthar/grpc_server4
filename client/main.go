// This file implements a client for interacting with a Tendermint gRPC server.
//
// It supports various commands, such as GetNodeInfo, GetSyncing, GetLatestBlock, GetBlockByHeight, GetLatestValidatorSet,
// GetValidatorSetByHeight, GetABCIInfo, and GetStatusInfo.
//
// When executed, the CLI parses the command-line arguments to determine which command to run and calls the corresponding
// gRPC method on the Tendermint server. It then prints the response to standard output in JSON format. If an error occurs,
// the CLI logs the error and exits.

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	types "grpc_server4/proto/generated"
	"log"
	"os"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/v12/app"
	"google.golang.org/grpc"
)

// The GRPC_SERVER_ADDRESS constant represents the address of the gRPC server to be used by the client.

const (
	GRPC_SERVER_ADDRESS = "localhost:9090"
	CHAIN_ID            = "osmosis-1"
)

// The Ccontext variable is of type client.Context and is used for client initialization.

var Ccontext = client.Context{}.WithChainID(CHAIN_ID)

// The InitCcontext function initializes the Ccontext

func InitCcontext() {
	encodingConfig := app.MakeEncodingConfig()
	Ccontext = Ccontext.WithCodec(encodingConfig.Marshaler).
		WithInterfaceRegistry(encodingConfig.InterfaceRegistry).
		WithTxConfig(encodingConfig.TxConfig).
		WithLegacyAmino(encodingConfig.Amino).
		WithInput(os.Stdin).
		WithBroadcastMode(flags.BroadcastSync).
		WithViper("OSMOSIS").
		WithSignModeStr(flags.SignModeDirect)
	conf := sdk.GetConfig()
	conf.SetBech32PrefixForAccount("osmo", "osmopub")
}

func main() {
	// Initialize the global Ccontext object
	InitCcontext()

	// Parse command line arguments
	flag.Parse()
	args := flag.Args()

	// Ensure that the command is specified in the arguments
	if len(args) == 0 {
		fmt.Println("command should be oneof: GetNodeInfo,GetSyncing,GetLatestBlock,GetBlockByHeight,GetLatestValidatorSet,GetValidatorSetByHeight, GetABCIInfo, GetStatusInfo")
		return
	}

	// Establish a gRPC connection to the server
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

	// Create a new instance of the gRPC client
	c := types.NewGrpcQueryServiceClient(grpcConn)

	// Execute the requested command
	switch args[0] {
	case "GetNodeInfo":
		// Call the GetNodeInfo RPC method and print the response
		ctx := context.Background()
		r, err := c.GetNodeInfo(ctx, &types.GetNodeInfoRequest{})
		if err != nil {
			log.Fatalf("GetNodeInfo err: %v", err)
			return
		}
		out, err := json.Marshal(r)
		if err != nil {
			log.Fatalf("GetNodeInfo err: %v", err)
			return
		}
		fmt.Println(string(out))

	case "GetSyncing":
		// Call the GetSyncing RPC method and print the response
		ctx := context.Background()
		r, err := c.GetSyncing(ctx, &types.GetSyncingRequest{})
		if err != nil {
			log.Fatalf("GetSyncing err: %v", err)
			return
		}
		fmt.Println("syncing is:", r.Syncing)

	case "GetLatestBlock":
		// Call the GetLatestBlock RPC method and print the response
		ctx := context.Background()
		r, err := c.GetLatestBlock(ctx, &types.GetLatestBlockRequest{})
		if err != nil {
			log.Fatalf("GetLatestBlock err: %v", err)
		}
		out, err := json.Marshal(r)
		if err != nil {
			log.Fatalf("GetLatestBlock err: %v", err)
			return
		}
		fmt.Println(string(out))

	case "GetBlockByHeight":
		// Call the GetBlockByHeight RPC method with the specified block height and print the response
		ctx := context.Background()
		if len(args) < 2 {
			fmt.Println("this command need the height!!!")
			return
		}
		height, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			log.Fatalf("GetSyncing err: %v", err)
		}

		r, err := c.GetBlockByHeight(ctx, &types.GetBlockByHeightRequest{Height: height})
		if err != nil {
			log.Fatalf("GetBlockByHeight err: %v", err)
			return
		}
		out, err := json.Marshal(r)
		if err != nil {
			log.Fatalf("GetBlockByHeight err: %v", err)
			return
		}
		// Print the response
		fmt.Println(string(out))
	case "GetLatestValidatorSet":
		// Call the GetLatestValidatorSet RPC method and print the response
		ctx := context.Background()
		r, err := c.GetLatestValidatorSet(ctx, &types.GetLatestValidatorSetRequest{})
		if err != nil {
			log.Fatalf("GetLatestValidatorSet err: %v", err)
			return
		}
		out, err := json.Marshal(r)
		if err != nil {
			log.Fatalf("GetLatestValidatorSet err: %v", err)
			return
		}
		fmt.Println(string(out))
	case "GetValidatorSetByHeight":
		// Call the GetValidatorSetByHeight RPC method with the specified block height and print the response
		ctx := context.Background()
		if len(args) < 2 {
			fmt.Println("this command need the height!!!")
			return
		}
		height, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			log.Fatalf("GetSyncing err: %v", err)
		}
		r, err := c.GetValidatorSetByHeight(ctx, &types.GetValidatorSetByHeightRequest{Height: height})
		if err != nil {
			log.Fatalf("GetValidatorSetByHeight err: %v", err)
		}
		out, err := json.Marshal(r)
		if err != nil {
			log.Fatalf("GetValidatorSetByHeight err: %v", err)
			return
		}
		fmt.Println(string(out))

	case "GetABCIInfo":
		// Call the GetABCIInfo RPC method and print the response
		ctx := context.Background()
		conn, err := grpc.Dial(GRPC_SERVER_ADDRESS, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("dial err: %v", err)
		}
		defer conn.Close()
		r, err := c.GetABCIInfo(ctx, &types.GetABCIInfoRequest{})
		if err != nil {
			log.Fatalf("GetABCIInfo err: %v", err)
		}
		out, err := json.Marshal(r)
		if err != nil {
			log.Fatalf("GetABCIInfo err: %v", err)
			return
		}
		fmt.Println(string(out))

	case "GetStatusInfo":
		// Call the GetStatusInfo RPC method and print the response
		ctx := context.Background()
		conn, err := grpc.Dial(GRPC_SERVER_ADDRESS, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("dial err: %v", err)
		}
		defer conn.Close()
		r, err := c.GetStatusInfo(ctx, &types.GetStatusInfoRequest{})
		if err != nil {
			log.Fatalf("GetStatusInfo err: %v", err)
		}
		out, err := json.Marshal(r)
		if err != nil {
			log.Fatalf("GetStatusInfo err: %v", err)
			return
		}
		fmt.Println(string(out))
	default:
		// If the command is not recognized, print the available commands to the user
		fmt.Println("command should be oneof: GetNodeInfo,GetSyncing,GetLatestBlock,GetBlockByHeight,GetLatestValidatorSet,GetValidatorSetByHeight, GetABCIInfo, GetStatusInfo")
	}
}
