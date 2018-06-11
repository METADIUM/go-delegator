package abi

import (
	"encoding/hex"
	"strings"

	"bitbucket.org/coinplugin/proxy/json"
	"bitbucket.org/coinplugin/proxy/rpc"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

func Pack(abi abi.ABI, name string, args ...interface{}) (string, error) {
	data, err := abi.Pack(name, args...)
	if err != nil {
		return "", err
	} else {
		return hexutil.Encode(data), nil
	}
}

func Unpack(abi abi.ABI, v interface{}, name string, output string) error {
	var data []byte
	var err error
	if output[:2] == "0x" {
		data, err = hex.DecodeString(output[2:])
	} else {
		data, err = hex.DecodeString(output)
	}

	if err != nil {
		return err
	}
	return abi.Unpack(v, name, data)
}

func Call(abi abi.ABI, to, name string, inputs []interface{}, outputs interface{}) error {
	data, err := Pack(abi, name, inputs...)
	if err != nil {
		return err
	}

	r := rpc.GetInstance(rpc.Testnet)
	respStr, err := r.Call(to, data)
	if err != nil {
		return err
	}

	resp := json.GetRpcResponseFromJson(respStr)
	return Unpack(abi, outputs, name, resp.Result.(string))
}

func GetAbiFromJson(raw string) (abi.ABI, error) {
	return abi.JSON(strings.NewReader(raw))
}

// getAbiFromAddress is NOT YET SUPPORTED
// TODO: use eth.compile.solidity?
func getAbiFromAddress(addr string) (abi abi.ABI) {
	r := rpc.GetInstance(rpc.Testnet)
	respStr, err := r.GetCode(addr)
	if err != nil {
		return
	}

	json.GetRpcResponseFromJson(respStr)
	return
}
