// Package predefined manages predefined functions for RPC request
package predefined

import (
	"fmt"

	"bitbucket.org/coinplugin/proxy/json"
	"bitbucket.org/coinplugin/proxy/rpc"
	"bitbucket.org/coinplugin/proxy/web3"
)

// Sample
func foo(req json.RPCRequest) (json.RPCResponse, error) {
	fmt.Println("foo")
	return json.RPCResponse{}, nil
}

// getBalance is a wrapper to support fromWei for eth_getBalance
func getBalance(req json.RPCRequest) (json.RPCResponse, error) {
	// Preprocessing
	var unit string
	if len(req.Params) > 2 {
		unit = req.Params[2].(string)
		req.Params = req.Params[:2]
	}

	// RPC
	var resp json.RPCResponse
	respBody, err := rpc.GetInstance().DoRPC(req)
	if err == nil {
		resp = json.GetRPCResponseFromJSON(respBody)
		// Postprocessing
		if unit != "" {
			if val, err := web3.FromWei(resp.Result.(string), unit); err == nil {
				resp.Result = val
			}
		}
	}

	return resp, err
}

// Forward delivers RPCRequest to predefined function and returns that
func Forward(req json.RPCRequest) (json.RPCResponse, error) {
	for k, v := range predefinedPaths {
		if k == req.Method {
			return v.(func(json.RPCRequest) (json.RPCResponse, error))(req)
		}
	}
	return json.RPCResponse{}, fmt.Errorf("predefined NOT FOUND")
}

// Contains check if given path is in predefined or not
func Contains(path string) bool {
	for k := range predefinedPaths {
		if k == path {
			return true
		}
	}
	return false
}

var predefinedPaths = map[string]interface{}{
	"foo":              foo,
	"eth_getBalance":   getBalance,
	"register_meta_id": registerMetaID,
	"revoke_meta_id":   revokeMetaID,
}
