package examples

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
	"github.com/umbracle/ethgo/contract"
	"github.com/umbracle/ethgo/jsonrpc"
)

func contractTransaction() {
	var functions = []string{
		"function transferFrom(address from, address to, uint256 value)",
	}

	abiContract, err := abi.NewABIFromList(functions)
	handleErr(err)

	// Matic token
	addr := common.HexToAddress("0x7d1afa7b718fb893db30a3abc0cfc608aacfebb0")

	client, err := jsonrpc.NewClient("https://mainnet.infura.io")
	handleErr(err)

	// wallet signer
	//key, err := wallet.GenerateKey()
	//handleErr(err)

	opts := []contract.ContractOption{
		contract.WithJsonRPC(client.Eth()),
		contract.WithSender(nil),
	}
	c := contract.NewContract(addr, abiContract, opts...)
	txHash, err := c.Txn("transferFrom", ethgo.Latest)
	handleErr(err)

	fmt.Printf("Transaction mined at: %x", txHash)
}
