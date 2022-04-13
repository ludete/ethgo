// Code generated by ethgo/abigen. DO NOT EDIT.
// Hash: 3f1af52b391dcf1991b5cee7468a69f382cfa0f819eaff85474464c969fe7ea9
// Version: 0.1.0
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
	_ = jsonrpc.NewClient
)

// Testdata is a solidity contract
type Testdata struct {
	c *contract.Contract
}

// NewTestdata creates a new instance of the contract at a specific address
func NewTestdata(addr ethgo.Address, opts ...contract.ContractOption) *Testdata {
	return &Testdata{c: contract.NewContract(addr, abiTestdata, opts...)}
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
func (t *Testdata) TxnBasicInput(val1 ethgo.Address, val2 *big.Int, opts *contract.TxnOpts) (contract.Txn, error) {
	return t.c.Txn("txnBasicInput", val1, val2, opts)
}

// events

func (t *Testdata) EventBasicEventSig() ethgo.Hash {
	return t.c.GetABI().Events["EventBasic"].ID()
}
