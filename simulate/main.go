package main

import (
	"fmt"

	"github.com/gnolang/gno/gno.land/pkg/gnoclient"
	rpcclient "github.com/gnolang/gno/tm2/pkg/bft/rpc/client"
	"github.com/gnolang/gno/tm2/pkg/crypto"
	"github.com/gnolang/gno/tm2/pkg/crypto/keys"
	"github.com/gnolang/gno/tm2/pkg/std"
)

func main() {
	// Initialize keybase from a directory
	keybase, _ := keys.NewKeyBaseFromDir("/Users/iam-agf/Library/Application Support/gno")

	// Create signer
	signer := gnoclient.SignerFromKeybase{
		Keybase:  keybase,
		Account:  "g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5", // Name of your keypair in keybase
		Password: "",                                         // Password to decrypt your keypair
		ChainID:  "dev",                                      // id of gno.land chain
	}

	// Initialize the RPC client
	rpc, err := rpcclient.NewHTTPClient("http://127.0.0.1:26657")
	if err != nil {
		panic(err)
	}

	// Initialize the gnoclient
	client := gnoclient.Client{
		Signer:    signer,
		RPCClient: rpc,
	}

	// Convert Gno address string to `crypto.Address`
	addr, err := crypto.AddressFromBech32("g1jg8mtutu9khhfwc4nxmuhcpftf0pajdhfvsqf5") // your Gno address
	if err != nil {
		panic(err)
	}

	accountRes, _, err := client.QueryAccount(addr)
	if err != nil {
		panic(err)
	}

	fmt.Println(accountRes)

	makeTx(
		"AddNewProblem",
		[]string{"Revert the strings", "Create a function 'Problem' that reverses the given string.", "gno->ong,dao->oad,very long sentence->ecnetnes gnol yrev"},
		accountRes,
		client,
	)
	makeTx(
		"AddNewProblem",
		[]string{
			"Primes under n",
			"Create a function 'Problem' that counts the number of primes under n. (A prime is a number that its only positive divisors are 1 and itself)",
			"2->0,20->8,1000000->78498",
		},
		accountRes,
		client)
	makeTx("AddNewProblem",
		[]string{"Remove the vowels",
			"Create a function 'Problem' that removes the vowels of a given string.",
			"gnolang->gnlng,dao->d,very long sentence->vry lng sntnc",
		},
		accountRes,
		client,
	)
	return
}

func makeTx(funcName string, args []string, accountRes *std.BaseAccount, client gnoclient.Client) error {
	txCfg := gnoclient.BaseTxCfg{
		GasFee:         "1000000ugnot",                 // gas price
		GasWanted:      1000000,                        // gas limit
		AccountNumber:  accountRes.GetAccountNumber(),  // account ID
		SequenceNumber: accountRes.GetSequence(),       // account nonce
		Memo:           "This is a cool how-to guide!", // transaction memo
	}

	msg := gnoclient.MsgCall{
		PkgPath:  "gno.land/r/demo/shiken", // wrapped ugnot realm path
		FuncName: funcName,                 // function to call
		Args:     nil,                      // arguments in string format
		Send:     "1000000ugnot",           // coins to send along with transaction
	}
	res, err := client.Call(txCfg, msg)
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
	return nil
}
