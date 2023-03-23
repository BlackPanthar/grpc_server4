package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"grpc_server4/types"
	"log"
	"os"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/v12/app"
	"google.golang.org/grpc"
)

// hello_client

const (
	GRPC_SERVER_ADDRESS = "localhost:9090"
	CHAIN_ID            = "osmosis-1"
)

var Ccontext = client.Context{}.WithChainID(CHAIN_ID)

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
	InitCcontext()
	flag.Parse()
	args := flag.Args()
	if len(args) == 0 {
		fmt.Println("command should be oneof: GetNodeInfo,GetSyncing,GetLatestBlock,GetBlockByHeight,GetLatestValidatorSet,GetValidatorSetByHeight")
		return

	}

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
	switch args[0] {
	case "GetNodeInfo":
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
		ctx := context.Background()
		r, err := c.GetSyncing(ctx, &types.GetSyncingRequest{})
		if err != nil {
			log.Fatalf("GetSyncing err: %v", err)
			return
		}
		fmt.Println("syncing is:", r.Syncing)
	case "GetLatestBlock":
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
		fmt.Println(string(out))
	case "GetLatestValidatorSet":
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
	default:
		fmt.Println("command should be oneof: GetNodeInfo,GetSyncing,GetLatestBlock,GetBlockByHeight,GetLatestValidatorSet,GetValidatorSetByHeight")
	}

}
