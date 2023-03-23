package main

import (
	"context"
	"fmt"
	"grpc_server4/types"
	"log"
	"net"
	"os"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/grpc/tmservice"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/v12/app"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	anypb "google.golang.org/protobuf/types/known/anypb"
)

const GRPC_SERVER_ADDRESS = "grpc.osmosis.zone:9090"
const CHAIN_ID = "osmosis-1"
const NODE_URL = "https://osmosis-mainnet-rpc.allthatnode.com:26657"

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
		WithSignModeStr(flags.SignModeDirect).
		WithNodeURI(NODE_URL)
	conf := sdk.GetConfig()
	conf.SetBech32PrefixForAccount("osmo", "osmopub")
}

type server struct {
	ctx context.Context
	types.UnimplementedGrpcQueryServiceServer
}

func (s *server) GetNodeInfo(ctx context.Context, req *types.GetNodeInfoRequest) (*types.GetNodeInfoResponse, error) {
	grpcConn, _ := grpc.Dial(
		GRPC_SERVER_ADDRESS, // your gRPC server address.
		grpc.WithInsecure(), // The SDK doesn't support any transport security mechanism.
	)
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			fmt.Println("grpcConn.Close() err:", err)
		}
	}(grpcConn)
	client := tmservice.NewServiceClient(grpcConn)
	resp, err := client.GetNodeInfo(ctx, &tmservice.GetNodeInfoRequest{})
	if err != nil {
		return nil, err
	}
	ans := &types.GetNodeInfoResponse{
		DefaultNodeInfo: resp.DefaultNodeInfo,
		ApplicationVersion: &types.VersionInfo{
			Name:             resp.ApplicationVersion.Name,
			AppName:          resp.ApplicationVersion.AppName,
			Version:          resp.ApplicationVersion.Version,
			GitCommit:        resp.ApplicationVersion.GitCommit,
			BuildTags:        resp.ApplicationVersion.BuildTags,
			GoVersion:        resp.ApplicationVersion.GoVersion,
			BuildDeps:        make([]*types.Module, 0, 0),
			CosmosSdkVersion: resp.ApplicationVersion.CosmosSdkVersion,
		},
	}
	for _, v := range resp.ApplicationVersion.BuildDeps {
		module := &types.Module{
			Path:    v.Path,
			Version: v.Version,
			Sum:     v.Sum,
		}
		ans.ApplicationVersion.BuildDeps = append(ans.ApplicationVersion.BuildDeps, module)
	}
	return ans, nil
}

func (s *server) GetSyncing(ctx context.Context, req *types.GetSyncingRequest) (*types.GetSyncingResponse, error) {
	grpcConn, _ := grpc.Dial(
		GRPC_SERVER_ADDRESS, // your gRPC server address.
		grpc.WithInsecure(),
	)
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			fmt.Println("grpcConn.Close() err:", err)
		}
	}(grpcConn)
	client := tmservice.NewServiceClient(grpcConn)
	issync, err := client.GetSyncing(ctx, &tmservice.GetSyncingRequest{})
	if err != nil {
		return nil, err
	}

	return &types.GetSyncingResponse{Syncing: issync.Syncing}, nil
}

func (s *server) GetLatestBlock(ctx context.Context, req *types.GetLatestBlockRequest) (*types.GetLatestBlockResponse, error) {
	grpcConn, _ := grpc.Dial(
		GRPC_SERVER_ADDRESS, // your gRPC server address.
		grpc.WithInsecure(),
	)
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			fmt.Println("grpcConn.Close() err:", err)
		}
	}(grpcConn)
	client := tmservice.NewServiceClient(grpcConn)
	latestBlock, err := client.GetLatestBlock(ctx, &tmservice.GetLatestBlockRequest{})
	if err != nil {
		return nil, err
	}

	return &types.GetLatestBlockResponse{
		BlockId: latestBlock.BlockId,
		Block:   latestBlock.Block,
	}, nil
}

func (s *server) GetBlockByHeight(ctx context.Context, req *types.GetBlockByHeightRequest) (*types.GetBlockByHeightResponse, error) {
	grpcConn, _ := grpc.Dial(
		GRPC_SERVER_ADDRESS, // your gRPC server address.
		grpc.WithInsecure(),
	)
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			fmt.Println("grpcConn.Close() err:", err)
		}
	}(grpcConn)
	client := tmservice.NewServiceClient(grpcConn)
	block, err := client.GetBlockByHeight(ctx, &tmservice.GetBlockByHeightRequest{
		Height: req.Height,
	})
	if err != nil {
		return nil, err
	}
	return &types.GetBlockByHeightResponse{
		BlockId: block.BlockId,
		Block:   block.Block,
	}, nil
}

func (s *server) GetLatestValidatorSet(ctx context.Context, req *types.GetLatestValidatorSetRequest) (*types.GetLatestValidatorSetResponse, error) {
	grpcConn, _ := grpc.Dial(
		GRPC_SERVER_ADDRESS, // your gRPC server address.
		grpc.WithInsecure(),
	)
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			fmt.Println("grpcConn.Close() err:", err)
		}
	}(grpcConn)
	client := tmservice.NewServiceClient(grpcConn)
	valSet, err := client.GetLatestValidatorSet(ctx, &tmservice.GetLatestValidatorSetRequest{Pagination: req.Pagination})
	if err != nil {
		return nil, err
	}
	validators := make([]*types.Validator, 0, 0)
	for _, v := range valSet.Validators {
		validator := &types.Validator{
			Address: v.Address,
			PubKey: &anypb.Any{
				TypeUrl: v.PubKey.TypeUrl,
				Value:   v.PubKey.Value,
			},
			VotingPower:      v.VotingPower,
			ProposerPriority: v.ProposerPriority,
		}
		validators = append(validators, validator)
	}
	return &types.GetLatestValidatorSetResponse{
		BlockHeight: valSet.BlockHeight,
		Validators:  validators,
		Pagination:  valSet.Pagination,
	}, nil

}

func (s *server) GetValidatorSetByHeight(ctx context.Context, req *types.GetValidatorSetByHeightRequest) (*types.GetValidatorSetByHeightResponse, error) {
	grpcConn, _ := grpc.Dial(
		GRPC_SERVER_ADDRESS, // your gRPC server address.
		grpc.WithInsecure(),
	)
	defer func(grpcConn *grpc.ClientConn) {
		err := grpcConn.Close()
		if err != nil {
			fmt.Println("grpcConn.Close() err:", err)
		}
	}(grpcConn)
	client := tmservice.NewServiceClient(grpcConn)
	valSet, err := client.GetValidatorSetByHeight(ctx, &tmservice.GetValidatorSetByHeightRequest{Pagination: req.Pagination, Height: req.Height})
	if err != nil {
		return nil, err
	}
	validators := make([]*types.Validator, 0, 0)
	for _, v := range valSet.Validators {
		validator := &types.Validator{
			Address: v.Address,
			PubKey: &anypb.Any{
				TypeUrl: v.PubKey.TypeUrl,
				Value:   v.PubKey.Value,
			},
			VotingPower:      v.VotingPower,
			ProposerPriority: v.ProposerPriority,
		}
		validators = append(validators, validator)
	}
	return &types.GetValidatorSetByHeightResponse{
		BlockHeight: valSet.BlockHeight,
		Validators:  validators,
		Pagination:  valSet.Pagination,
	}, nil
}

func Serve() {

	listen, err := net.Listen("tcp", "localhost:9090")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	gs := grpc.NewServer()
	reflection.Register(gs)
	types.RegisterGrpcQueryServiceServer(gs, &server{})
	fmt.Println("grpc server is started")
	err = gs.Serve(listen)
	if err != nil {
		fmt.Printf("failed to serve: %v", err)
		return
	}
	log.Println("OsmosisServiceServer is closed!")
}

func main() {
	InitCcontext()
	Serve()
}