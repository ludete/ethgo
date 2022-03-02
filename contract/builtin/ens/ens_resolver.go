package ens

import (
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/jsonrpc"
)

type ENSResolver struct {
	e        *ENS
	provider *jsonrpc.Client
}

func NewENSResolver(addr ethgo.Address, provider *jsonrpc.Client) *ENSResolver {
	return &ENSResolver{NewENS(addr, provider), provider}
}

func (e *ENSResolver) Resolve(addr string, block ...ethgo.BlockNumber) (res ethgo.Address, err error) {
	addrHash := NameHash(addr)
	resolverAddr, err := e.e.Resolver(addrHash, block...)
	if err != nil {
		return
	}

	resolver := NewResolver(resolverAddr, e.provider)
	res, err = resolver.Addr(addrHash, block...)
	return
}
