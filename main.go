package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"bitbucket.org/coinplugin/proxy/crypto"
	"bitbucket.org/coinplugin/proxy/json"
	"bitbucket.org/coinplugin/proxy/predefined"
	"bitbucket.org/coinplugin/proxy/rpc"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

const (
	// ParamFuncName is a name indicating function's
	ParamFuncName = "func"
	// Targetnet indicates target network
	Targetnet = rpc.Testnet
)

func handler(req json.RPCRequest) (body string, statusCode int) {
	var resp json.RPCResponse
	var err error
	if predefined.Contains(req.Method) {
		// Forward RPC request to predefined function
		resp, err = predefined.Forward(req)
	} else {
		// Forward RPC request to Ether node
		var respBody string
		if respBody, err = rpc.GetInstance().DoRPC(req); err == nil {
			// Relay a response from the node
			resp = json.GetRPCResponseFromJSON(respBody)
		}
	}

	statusCode = 200
	if err != nil {
		// In case of server-side RPC fail
		fmt.Println(err.Error())
		resp.Error = &json.RPCError{
			Code:    -1,
			Message: err.Error(),
		}
		statusCode = 400
	} else if resp.Error != nil && resp.Error.Code != 0 {
		// In case of ether-node-side RPC fail
		statusCode = 400
	}
	body = resp.String()
	return
}

// lambdaHandler handles APIGatewayProxyRequest as JSON-RPC request
func lambdaHandler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	// Validate RPC request
	req := json.GetRPCRequestFromJSON(request.Body)
	if method := request.QueryStringParameters[ParamFuncName]; method != "" {
		req.Method = method
	} else if method := request.PathParameters[ParamFuncName]; method != "" {
		req.Method = method
	}

	respBody, statusCode := handler(req)
	return events.APIGatewayProxyResponse{Body: respBody, StatusCode: statusCode}, nil
}

// httpHandler handles http.Request as JSON-RPC request
func httpHandler(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	req := json.GetRPCRequestFromJSON(string(b))
	respBody, statusCode := handler(req)
	w.WriteHeader(statusCode)
	w.Write([]byte(respBody))
}

func help() {
	fmt.Println("====== Need arguments")
	fmt.Println("======== Option 1. key path only as argument")
	fmt.Println("========== $> proxy [path]")
	fmt.Println("======== Option 2. key path and passphrase as argument")
	fmt.Println("========== $> proxy [path] [passphrase]")
	fmt.Println("======== Option 3. key path and passphrase as environment variable")
	fmt.Println("========== $> export KEY_PATH=[path]")
	fmt.Println("========== $> export KEY_PASSPHRASE=[passphrase]")
	fmt.Println("========== $> proxy")
}

func main() {
	rpc.NetType = Targetnet

	// Initalize Crypto with arguments
	var path, passphrase string
	if path = os.Getenv(crypto.Path); path != "" {
		passphrase = os.Getenv(crypto.Passphrase)
		os.Setenv(crypto.Path, "")
		os.Setenv(crypto.Passphrase, "")
	} else if len(os.Args) > 1 {
		path = os.Args[1]
		if len(os.Args) > 2 {
			passphrase = os.Args[2]
		} else {
			fmt.Printf("Passphrase: ")
			fmt.Scanln(&passphrase)
		}
	} else {
		help()
		return
	}
	go func() {
		crypto.PathChan <- path
		crypto.PassphraseChan <- passphrase
	}()
	crypto.GetInstance()

	// Run
	if os.Getenv(crypto.IsAwsLambda) != "" {
		fmt.Println("Ready to start Lambda")
		lambda.Start(lambdaHandler)
	} else {
		fmt.Println("Ready to start HTTP/HTTPS")
		http.HandleFunc("/", httpHandler)
		http.ListenAndServe(":8545", nil)
		// http.ListenAndServeTLS()
	}
}
