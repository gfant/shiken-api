package handler

import (
	"fmt"

	"github.com/gnolang/gno/gno.land/pkg/gnoclient"
	rpcclient "github.com/gnolang/gno/tm2/pkg/bft/rpc/client"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys"
	"github.com/gnolang/gno/tm2/pkg/std"
)

func SetupRegisterEnvironment(
	keyPath,
	address,
	chainId,
	httpClientPath string,
) (*std.BaseAccount, gnoclient.Client) {
	// Initialize keybase from a directory
	keybase, _ := keys.NewKeyBaseFromDir(keyPath)

	// Create signer
	signer := gnoclient.SignerFromKeybase{
		Keybase:  keybase,
		Account:  address, // Name of your keypair in keybase
		Password: "",      // Password to decrypt your keypair
		ChainID:  chainId, // id of gno.land chain
	}

	// Initialize the RPC client
	rpc, err := rpcclient.NewHTTPClient(httpClientPath)
	if err != nil {
		panic(err)
	}

	// Initialize the gnoclient
	client := gnoclient.Client{
		Signer:    signer,
		RPCClient: rpc,
	}

	// Convert Gno address string to `crypto.Address`
	addr, err := crypto.AddressFromBech32(address) // your Gno address
	if err != nil {
		panic(err)
	}

	accountRes, _, err := client.QueryAccount(addr)
	if err != nil {
		panic(err)
	}

	return accountRes, client
}

func MakeTx(realmPath, gasFee, memo, funcName string, gasWanted int64, args []string, accountRes *std.BaseAccount, client gnoclient.Client) error {
	txCfg := gnoclient.BaseTxCfg{
		GasFee:         gasFee,                        // gas price
		GasWanted:      gasWanted,                     // gas limit
		AccountNumber:  accountRes.GetAccountNumber(), // account ID
		SequenceNumber: accountRes.GetSequence(),      // account nonce
		Memo:           memo,                          // transaction memo
	}

	msg := gnoclient.MsgCall{
		PkgPath:  realmPath, // wrapped ugnot realm path
		FuncName: funcName,  // function to call
		Args:     args,      // arguments in string format
	}
	res, err := client.Call(txCfg, msg)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	return nil
}
