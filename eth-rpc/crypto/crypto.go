package crypto

import (
	"fmt"

	"github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/common"
	"github.com/hexoul/eth-rpc-on-aws-lambda/eth-rpc/db"

	ethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

func getPrivKeyFromDB(propVal string) string {
	//dbHelper := db.New("ap-northeast-2")
	dbHelper := db.New("aws-region")
	if dbHelper == nil {
		return ""
	}

	ret := dbHelper.GetItem(config.DbConfigTblName, config.DbConfigPropName, propVal, config.DbConfigValName)
	item := common.DbConfigResult{}
	for _, elem := range ret.Items {
		dbHelper.UnmarshalMap(elem, &item)
		return item.Value
	}
	return ""
}

func Sign() {
	privKey := getPrivKeyFromDB("priv_key1")
	if privKey == "" {
		return
	}
}

// signHash is a helper function that calculates a hash for the given message that can be
// safely used to calculate a signature from.
//
// The hash is calulcated as
//   keccak256("\x19Ethereum Signed Message:\n"${message length}${message}).
//
// This gives context to the signed message and prevents signing of transactions.
func signHash(data []byte) []byte {
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data)
	return crypto.Keccak256([]byte(msg))
}

// EcRecover returns the address for the Account that was used to create the signature.
// Note, this function is compatible with eth_sign and personal_sign. As such it recovers
// the address of:
// hash = keccak256("\x19Ethereum Signed Message:\n"${message length}${message})
// addr = ecrecover(hash, signature)
//
// Note, the signature must conform to the secp256k1 curve R, S and V values, where
// the V value must be be 27 or 28 for legacy reasons.
//
// https://github.com/ethereum/go-ethereum/wiki/Management-APIs#personal_ecRecover
func EcRecover(dataStr, sigStr string) (string, error) {
	data := hexutil.MustDecode(dataStr)
	sig := hexutil.MustDecode(sigStr)
	if len(sig) != 65 {
		return "", fmt.Errorf("signature must be 65 bytes long")
	}
	if sig[64] != 27 && sig[64] != 28 {
		return "", fmt.Errorf("invalid Ethereum signature (V is not 27 or 28)")
	}
	sig[64] -= 27 // Transform yellow paper V from 27/28 to 0/1

	rpk, err := crypto.Ecrecover(signHash(data), sig)
	if err != nil {
		return "", err
	}
	pubKey := crypto.ToECDSAPub(rpk)
	recoveredAddr := crypto.PubkeyToAddress(*pubKey)
	return fmt.Sprintf("0x%x", recoveredAddr), nil
}

func EcRecoverToPubkey(hash, sig string) ([]byte, error) {
	return crypto.Ecrecover(hexutil.MustDecode(hash), hexutil.MustDecode(sig))
}

func PubkeyToAddress(p []byte) ethcommon.Address {
	return ethcommon.BytesToAddress(crypto.Keccak256(p[1:])[12:])
}
