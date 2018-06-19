package abi

import (
	"math/big"

	"bitbucket.org/coinplugin/proxy/crypto"
	"bitbucket.org/coinplugin/proxy/json"
	"bitbucket.org/coinplugin/proxy/rpc"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

func DummySendTransaction(abi abi.ABI, targetNet, to, name string, inputs []interface{}, gas int) (resp json.RpcResponse, err error) {
	data, err := Pack(abi, name, inputs...)
	if err != nil {
		return
	}

	c := crypto.GetDummy()
	r := rpc.GetInstance(targetNet)
	respStr, err := r.SendTransaction(c.Address, to, data, gas)
	if err != nil {
		return
	}

	resp = json.GetRpcResponseFromJson(respStr)
	return
}

func DummySendTransactionWithSign(abi abi.ABI, targetNet, to, name string, inputs []interface{}, gasLimit, gasPrice uint64) (resp json.RpcResponse, err error) {
	data, err := abi.Pack(name, inputs...)
	if err != nil {
		return
	}

	c := crypto.GetDummy()
	tx := types.NewTransaction(0, common.HexToAddress(to), big.NewInt(0), uint64(gasLimit), big.NewInt(int64(gasPrice)), data)
	tx, err = c.SignTx(tx)
	if err != nil {
		return
	}

	rlpTx, err := rlp.EncodeToBytes(tx)
	if err != nil {
		return
	}

	r := rpc.GetInstance(targetNet)
	respStr, err := r.SendRawTransaction(rlpTx)
	if err != nil {
		return
	}

	resp = json.GetRpcResponseFromJson(respStr)
	return
}
