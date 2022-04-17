package contract

import (
	"crypto/ecdsa"

	"github.com/umbracle/ethgo/jsonrpc"
)

type Opts struct {
	JsonRPCEndpoint string
	JsonRPCClient   *jsonrpc.Eth
	Provider        Provider
	Sender          *ecdsa.PrivateKey
	chainId         int64
}

type ContractOption func(*Opts)

func WithJsonRPCEndpoint(endpoint string) ContractOption {
	return func(o *Opts) {
		o.JsonRPCEndpoint = endpoint
	}
}

func WithChainId(chainId int64) ContractOption {
	return func(opts *Opts) {
		opts.chainId = chainId
	}
}

func WithJsonRPC(client *jsonrpc.Eth) ContractOption {
	return func(o *Opts) {
		o.JsonRPCClient = client
	}
}

func WithProvider(provider Provider) ContractOption {
	return func(o *Opts) {
		o.Provider = provider
	}
}

func WithSender(sender *ecdsa.PrivateKey) ContractOption {
	return func(o *Opts) {
		o.Sender = sender
	}
}
