// Code generated by ethgo/abigen. DO NOT EDIT.
// Hash: d644ffae9e5df06d8b503a99da12a04d43a389598849c4f582d8610d8f846672
package ens

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

// ENS is a solidity contract
type ENS struct {
	c *contract.Contract
}

// DeployENS deploys a new ENS contract
func DeployENS(provider *jsonrpc.Client, from ethgo.Address, args ...interface{}) *contract.Txn {
	return contract.DeployContract(provider, from, abiENS, binENS, args...)
}

// NewENS creates a new instance of the contract at a specific address
func NewENS(addr ethgo.Address, provider *jsonrpc.Client) *ENS {
	return &ENS{c: contract.NewContract(addr, abiENS, provider)}
}

// Contract returns the contract object
func (e *ENS) Contract() *contract.Contract {
	return e.c
}

// calls

// Owner calls the owner method in the solidity contract
func (e *ENS) Owner(node [32]byte, block ...ethgo.BlockNumber) (retval0 ethgo.Address, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = e.c.Call("owner", ethgo.EncodeBlock(block...), node)
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(ethgo.Address)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}

	return
}

// Resolver calls the resolver method in the solidity contract
func (e *ENS) Resolver(node [32]byte, block ...ethgo.BlockNumber) (retval0 ethgo.Address, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = e.c.Call("resolver", ethgo.EncodeBlock(block...), node)
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(ethgo.Address)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}

	return
}

// Ttl calls the ttl method in the solidity contract
func (e *ENS) Ttl(node [32]byte, block ...ethgo.BlockNumber) (retval0 uint64, err error) {
	var out map[string]interface{}
	var ok bool

	out, err = e.c.Call("ttl", ethgo.EncodeBlock(block...), node)
	if err != nil {
		return
	}

	// decode outputs
	retval0, ok = out["0"].(uint64)
	if !ok {
		err = fmt.Errorf("failed to encode output at index 0")
		return
	}

	return
}

// txns

// SetOwner sends a setOwner transaction in the solidity contract
func (e *ENS) SetOwner(node [32]byte, owner ethgo.Address) *contract.Txn {
	return e.c.Txn("setOwner", node, owner)
}

// SetResolver sends a setResolver transaction in the solidity contract
func (e *ENS) SetResolver(node [32]byte, resolver ethgo.Address) *contract.Txn {
	return e.c.Txn("setResolver", node, resolver)
}

// SetSubnodeOwner sends a setSubnodeOwner transaction in the solidity contract
func (e *ENS) SetSubnodeOwner(node [32]byte, label [32]byte, owner ethgo.Address) *contract.Txn {
	return e.c.Txn("setSubnodeOwner", node, label, owner)
}

// SetTTL sends a setTTL transaction in the solidity contract
func (e *ENS) SetTTL(node [32]byte, ttl uint64) *contract.Txn {
	return e.c.Txn("setTTL", node, ttl)
}

// events

func (e *ENS) NewOwnerEventSig() ethgo.Hash {
	return e.c.ABI().Events["NewOwner"].ID()
}

func (e *ENS) NewResolverEventSig() ethgo.Hash {
	return e.c.ABI().Events["NewResolver"].ID()
}

func (e *ENS) NewTTLEventSig() ethgo.Hash {
	return e.c.ABI().Events["NewTTL"].ID()
}

func (e *ENS) TransferEventSig() ethgo.Hash {
	return e.c.ABI().Events["Transfer"].ID()
}
