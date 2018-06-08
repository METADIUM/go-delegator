package main

import (
	"context"
	"fmt"

	_ "bitbucket.org/coinplugin/proxy/eth-rpc/crypto"
	"bitbucket.org/coinplugin/proxy/eth-rpc/json"
	"bitbucket.org/coinplugin/proxy/eth-rpc/rpc"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	ParamFuncName = "func"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate RPC request
	req := json.GetRpcRequestFromJson(request.Body)
	if method := request.QueryStringParameters[ParamFuncName]; method != "" {
		req.Method = method
	} else if method := request.PathParameters[ParamFuncName]; method != "" {
		req.Method = method
	}
	fmt.Printf("RpcRequest: %#v\n", req)

	// Forward RPC request to Ether node
	respBody, err := rpc.GetInstance(rpc.Testnet).DoRpc(req)

	// Relay a response from the node
	resp := json.GetRpcResponseFromJson(respBody)
	fmt.Printf("RpcResponse: %#v\n", resp)
	retCode := 200
	if err != nil {
		// In case of server-side RPC fail
		resp.Error.Message = err.Error()
		respBody = resp.String()
		retCode = 400
	} else if resp.Error.Code != 0 {
		// In case of ether-node-side RPC fail
		retCode = 400
	}
	return events.APIGatewayProxyResponse{Body: respBody, StatusCode: retCode}, nil
}

func main() {
	lambda.Start(Handler)
}
