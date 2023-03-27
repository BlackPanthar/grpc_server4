package test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"testing"

	types "grpc_server4/proto/generated"

	"google.golang.org/grpc"
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
		t.Fatalf("dial err: %v", err)
	}
	defer grpcConn.Close()

	c := types.NewGrpcQueryServiceClient(grpcConn)
	ctx := context.Background()

	// basic Comparison of result atributes with Osmosis direct call attributes

	abciInfo, err := c.GetABCIInfo(ctx, &types.GetABCIInfoRequest{})
	if err != nil {
		t.Fatalf("GetABCIInfo err: %v", err)
	}

	abciJSON := map[string]interface{}{
		"jsonrpc":     abciInfo.Jsonrpc,
		"data":        abciInfo.Response.Data,
		"version":     abciInfo.Response.Version,
		"app_version": abciInfo.Response.AppVersion,
	}
	abciBytes, err := json.Marshal(abciJSON)
	if err != nil {
		t.Fatalf("GetABCIInfo err: %v", err)
	}
	fmt.Println(string(abciBytes))

	// Direct call to Osmosis endpoint to compare with server results
	abciInfoDirect, err := getABCIInfoDirect()
	if err != nil {
		t.Fatalf("getABCIInfoDirect err: %v", err)
	}

	abciDirectJSON := map[string]interface{}{
		"jsonrpc":     abciInfoDirect.Jsonrpc,
		"data":        abciInfoDirect.Response.Data,
		"version":     abciInfoDirect.Response.Version,
		"app_version": abciInfoDirect.Response.AppVersion,
	}
	abciDirectBytes, err := json.Marshal(abciDirectJSON)
	if err != nil {
		t.Fatalf("getABCIInfoDirect err: %v", err)
	}
	fmt.Println(string(abciDirectBytes))

	// Compare server results with direct results
	if string(abciDirectBytes) == "" || string(abciBytes) == "" {
		t.Errorf("One or both GetABCIInfo results are empty: server=%v, direct=%v", string(abciBytes), string(abciDirectBytes))
	}
	if string(abciDirectBytes) == string(abciBytes) {
		fmt.Printf("ABCIInfo grpc response and server response equal: their formats are exactly same/equal %s %s"+string(abciDirectBytes), string(abciBytes))
	}
	statusInfo, err := c.GetStatusInfo(ctx, &types.GetStatusInfoRequest{})
	if err != nil {
		t.Fatalf("GetStatusInfo err: %v", err)
	}

	// Direct call to Osmosis endpoint to compare with server results
	statusInfoDirect, err := getStatusInfoDirect()
	if err != nil {
		t.Fatalf("getStatusInfoDirect err: %v", err)
	}

	// Compare server results with direct results
	if statusInfoDirect.ResponseString == "" || statusInfo.ResponseString == "" {
		t.Errorf("One or both GetStatusInfo results are empty: server=%v, direct=%v", statusInfo.ResponseString, statusInfoDirect.ResponseString)
	}

	if statusInfoDirect.ResponseString == statusInfo.ResponseString {
		fmt.Printf("StatusInfo grpc response and server response equal: their formats are exactly same/equal %s %s"+statusInfo.ResponseString, statusInfoDirect.ResponseString)
	}

	// Basic function tests to check errors - no comparison with Osmosis in this section
	nodeInfo, err := c.GetNodeInfo(ctx, &types.GetNodeInfoRequest{})
	if err != nil {
		t.Fatalf("GetNodeInfo err: %v", err)
	}
	nodeOut, err := json.Marshal(nodeInfo)
	if err != nil {
		t.Fatalf("GetNodeInfo err: %v", err)
	}
	fmt.Println(string(nodeOut))

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
	out, err := json.Marshal(latestBlock)
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
	heightOut, err := json.Marshal(heightBlock)
	if err != nil {
		log.Fatalf("GetBlockByHeight err: %v", err)
		return
	}
	fmt.Println(string(heightOut))

	latestValidators, err := c.GetLatestValidatorSet(ctx, &types.GetLatestValidatorSetRequest{})
	if err != nil {
		log.Fatalf("GetLatestValidatorSet err: %v", err)
		return
	}
	lValOut, err := json.Marshal(latestValidators)
	if err != nil {
		log.Fatalf("GetLatestValidatorSet err: %v", err)
		return
	}
	fmt.Println(string(lValOut))

	heightValidators, err := c.GetValidatorSetByHeight(ctx, &types.GetValidatorSetByHeightRequest{Height: 8658239})
	if err != nil {
		log.Fatalf("GetValidatorSetByHeight err: %v", err)
	}
	hValOut, err := json.Marshal(heightValidators)
	if err != nil {
		log.Fatalf("GetValidatorSetByHeight err: %v", err)
		return
	}
	fmt.Println(string(hValOut))
}
func getABCIInfoDirect() (*types.GetABCIInfoResponse, error) {
	resp, err := http.Get("https://rpc.osmosis.zone/abci_info")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Print response before decoding
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(bodyBytes))

	var response struct {
		Result struct {
			Jsonrpc      string             `json:"jsonrpc"`
			Id           interface{}        `json:"id"`
			Response     types.ABCIResponse `json:"response"`
			ABCIResponse struct {
				Data       string `json:"data"`
				Version    string `json:"version"`
				AppVersion string `json:"app_version"`
			} `json:"abci_info"`
		} `json:"result"`
	}

	err = json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&response)
	if err != nil {
		return nil, err
	}

	// Print response after decoding
	fmt.Printf("%+v\n", response)

	id, ok := response.Result.Id.(string)
	if !ok || id == "" {
		id = "0"
		//return nil, fmt.Errorf("invalid id value")
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("failed to convert id value to int32: %v", err)
	}

	return &types.GetABCIInfoResponse{
		Jsonrpc: response.Result.Jsonrpc,
		Id:      int32(idInt),
		Response: &types.ABCIResponse{
			Data:       response.Result.ABCIResponse.Data,
			Version:    response.Result.ABCIResponse.Version,
			AppVersion: response.Result.ABCIResponse.AppVersion,
		},
	}, nil
}

func getStatusInfoDirect() (*types.GetStatusInfoResponse, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://rpc.osmosis.zone/status", nil)
	if err != nil {
		return nil, err
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	jsonResp, err := json.Marshal(result["result"])
	if err != nil {
		return nil, err
	}

	responseStr := string(jsonResp)
	fmt.Println(responseStr)

	return &types.GetStatusInfoResponse{
		ResponseString: responseStr,
	}, nil
}
