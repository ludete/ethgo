package examples

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
	"github.com/umbracle/ethgo/contract"
	"github.com/umbracle/ethgo/jsonrpc"
)

// call a contract and use a custom `from` address
func contractCallFrom() {
	var functions = []string{
		"function totalSupply() view returns (uint256)",
	}

	abiContract, err := abi.NewABIFromList(functions)
	handleErr(err)

	// Matic token
	addr := common.HexToAddress("0x7d1afa7b718fb893db30a3abc0cfc608aacfebb0")

	client, err := jsonrpc.NewClient("https://mainnet.infura.io")
	handleErr(err)

	// from address (msg.sender in solidity)
	//from := common.Address{0x1}

	c := contract.NewContract(addr, abiContract, contract.WithSender(nil), contract.WithJsonRPC(client.Eth()))
	res, err := c.Call("totalSupply", ethgo.Latest)
	handleErr(err)

	fmt.Printf("TotalSupply: %s", res["totalSupply"].(*big.Int))
}
