package examples

import (
	"fmt"

	"github.com/umbracle/ethgo/abi"
	"github.com/umbracle/ethgo/contract"
)

func contractDeploy() {
	abiContract, err := abi.NewABIFromList([]string{})
	handleErr(err)

	// bytecode of the contract
	bin := []byte{}

	txHash, err := contract.DeployContract(abiContract, bin, []interface{}{}, nil)
	handleErr(err)
	fmt.Printf("Contract: %s", txHash)
}
