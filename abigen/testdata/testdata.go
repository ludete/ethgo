// Code generated by ethgo/abigen. DO NOT EDIT.
// Hash: 3f1af52b391dcf1991b5cee7468a69f382cfa0f819eaff85474464c969fe7ea9
package testdata

import (
	"fmt"
	"math/big"

	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/contract"
	"github.com/umbracle/ethgo/jsonrpc"
)

var (
	_ = big.NewInt
)

// Testdata is a solidity contract
type Testdata struct {
	c *contract.Contract
}

// NewTestdata creates a new instance of the contract at a specific address
func NewTestdata(addr ethgo.Address, provider *jsonrpc.Client) *Testdata {
	return &Testdata{c: contract.NewContract(addr, abiTestdata, provider)}
}

// Contract returns the contract object
func (t *Testdata) Contract() *contract.Contract {
	return t.c
}

// calls

// CallBasicInput calls the callBasicInput method in the solidity contract
func (t *Testdata) CallBasicInput(block ...ethgo.BlockNumber) (retval0 *big.Int, retval1 ethgo.Address, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = t.c.Call("callBasicInput", ethgo.EncodeBlock(block...))
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(*big.Int)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}
	retval1, ok = out["1"].(ethgo.Address)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 1")
		return
	}

	return
}

// txns

// TxnBasicInput sends a txnBasicInput transaction in the solidity contract
func (t *Testdata) TxnBasicInput(val1 ethgo.Address, val2 *big.Int) *contract.Txn {
	return t.c.Txn("txnBasicInput", val1, val2)
}

// events

func (t *Testdata) EventBasicEventSig() ethgo.Hash {
	return t.c.ABI().Events["EventBasic"].ID()
}
