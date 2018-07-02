package abi

import (
	"fmt"
	"math/big"

	"bitbucket.org/coinplugin/proxy/crypto"
	"bitbucket.org/coinplugin/proxy/json"
	"bitbucket.org/coinplugin/proxy/rpc"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
)

// DummySendTransaction invokes abi.SendTransaction with dummy of Crypto struct
func DummySendTransaction(abi abi.ABI, targetNet, to, name string, inputs []interface{}, gas int) (resp json.RPCResponse, err error) {
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

	resp = json.GetRPCResponseFromJSON(respStr)
	return
}

// DummySendTransactionWithSign invokes abi.SendTransactionWithSign with dummy of Crypto struct
func DummySendTransactionWithSign(abi abi.ABI, targetNet, to, name string, inputs []interface{}, gasLimit, gasPrice uint64) (resp json.RPCResponse, err error) {
	data, err := abi.Pack(name, inputs...)
	if err != nil {
		return
	}

	c := crypto.GetDummy()
	r := rpc.GetInstance(targetNet)

	// Make TX function to get nonce
	tx := func(nonce uint64) error {
		tx := types.NewTransaction(nonce, common.HexToAddress(to), zero, uint64(gasLimit), big.NewInt(int64(gasPrice)), data)
		tx, err = c.SignTx(tx)
		if err != nil {
			return err
		}

		rlpTx, err := rlp.EncodeToBytes(tx)
		if err != nil {
			return err
		}

		respStr, err := r.SendRawTransaction(rlpTx)
		if err != nil {
			return err
		}

		resp = json.GetRPCResponseFromJSON(respStr)
		if resp.Error == nil {
			return fmt.Errorf("%s", resp.Error.Message)
		}
		return nil
	}

	c.ApplyNonce(tx)
	return
}
