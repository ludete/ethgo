package contract

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
	"github.com/umbracle/ethgo/jsonrpc"

	//"github.com/umbracle/ethgo/wallet"
	eTypes "github.com/ethereum/go-ethereum/core/types"
)

// Provider handles the interactions with the Ethereum 1x node
type Provider interface {
	Call(ethgo.Address, []byte, *CallOpts) ([]byte, error)
	SendEIP1559Tx(to common.Address, nonce uint64,
		key *ecdsa.PrivateKey, input []byte, opts *TxnOpts) (ethgo.Hash, error)
	Wait(txHash ethgo.Hash) (*ethgo.Receipt, error)
}

type JonRPCNodeProvider struct {
	chainId       *big.Int
	client        *jsonrpc.Eth
	queryInterval time.Duration
}

func NewJonRPCNodeProvider(client *jsonrpc.Eth, chainId int64) *JonRPCNodeProvider {
	return &JonRPCNodeProvider{
		client:        client,
		chainId:       big.NewInt(chainId),
		queryInterval: time.Millisecond,
	}
}

func (j *JonRPCNodeProvider) Call(addr ethgo.Address, input []byte, opts *CallOpts) ([]byte, error) {
	msg := &ethgo.CallMsg{
		To:   &addr,
		Data: input,
	}
	if opts.From != ethgo.ZeroAddress {
		msg.From = opts.From
	}
	rawStr, err := j.client.Call(msg, opts.Block)
	if err != nil {
		return nil, err
	}
	raw, err := hex.DecodeString(rawStr[2:])
	if err != nil {
		return nil, err
	}
	return raw, nil
}

func (j *JonRPCNodeProvider) SendEIP1559Tx(to common.Address,
	nonce uint64, key *ecdsa.PrivateKey, input []byte, opts *TxnOpts) (ethgo.Hash, error) {

	var (
		err  error
		hash ethgo.Hash
	)
	dynamicFeeTx := &eTypes.DynamicFeeTx{
		Nonce:     nonce,
		GasTipCap: opts.MaxPriorityFeePerGas,
		GasFeeCap: opts.MaxFeePerGas,
		Gas:       opts.GasLimit,
		To:        &to,
		Value:     opts.Value,
		Data:      input,
	}
	signedTx, err := eTypes.SignNewTx(key, eTypes.LatestSignerForChainID(j.chainId), dynamicFeeTx)
	if err != nil {
		return hash, err
	}
	txData, err := signedTx.MarshalBinary()
	if err != nil {
		return hash, err
	}
	hash, err = j.client.SendRawTransaction(txData)
	return hash, err
}

func (j *JonRPCNodeProvider) Wait(txHash ethgo.Hash) (*ethgo.Receipt, error) {
	if (txHash == ethgo.Hash{}) {
		return nil, nil
	}

	for {
		receipt, err := j.client.GetTransactionReceipt(txHash)
		if err != nil {
			if err.Error() != "not found" {
				return nil, err
			}
		}
		if receipt != nil {
			return receipt, nil
		}
		time.Sleep(j.queryInterval)
	}
}

type Opts struct {
	JsonRPCEndpoint string
	JsonRPCClient   *jsonrpc.Eth
	Provider        Provider
	Sender          ethgo.Key
}

type ContractOption func(*Opts)

func WithJsonRPCEndpoint(endpoint string) ContractOption {
	return func(o *Opts) {
		o.JsonRPCEndpoint = endpoint
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

func WithSender(sender ethgo.Key) ContractOption {
	return func(o *Opts) {
		o.Sender = sender
	}
}

func DeployContract(abi *abi.ABI, bin []byte, args []interface{}, txOpts *TxnOpts, opts ...ContractOption) (Txn, error) {
	a := NewContract(ethgo.Address{}, abi, opts...)
	a.bin = bin

	return a.Txn("constructor", append(args, txOpts))
}

func NewContract(addr ethgo.Address, abi *abi.ABI, opts ...ContractOption) *Contract {
	opt := &Opts{
		JsonRPCEndpoint: "http://localhost:8545",
	}
	for _, c := range opts {
		c(opt)
	}

	var provider Provider
	if opt.Provider != nil {
		provider = opt.Provider
	} else if opt.JsonRPCClient != nil {
		provider = &JonRPCNodeProvider{client: opt.JsonRPCClient}
	} else {
		client, _ := jsonrpc.NewClient(opt.JsonRPCEndpoint)
		provider = &JonRPCNodeProvider{client: client.Eth()}
	}

	a := &Contract{
		addr:     addr,
		abi:      abi,
		provider: provider,
		key:      opt.Sender,
	}

	return a
}

// Contract is a wrapper to make abi calls to contract with a state provider
type Contract struct {
	addr     ethgo.Address
	abi      *abi.ABI
	bin      []byte
	provider Provider
	key      ethgo.Key
}

func (a *Contract) GetABI() *abi.ABI {
	return a.abi
}

type TxnOpts struct {
	Value                *big.Int
	GasPrice             uint64
	GasLimit             uint64
	Nonce                uint64
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
}

func (a *Contract) Txn(method string, args ...interface{}) (Txn, error) {
	if a.key == nil {
		return nil, fmt.Errorf("no key selected")
	}
	length := len(args)
	isContractDeployment := method == "constructor"

	var input []byte
	if isContractDeployment {
		input = append(input, a.bin...)
	}

	var abiMethod *abi.Method
	if isContractDeployment {
		if a.abi.Constructor != nil {
			abiMethod = a.abi.Constructor
		}
	} else {
		if abiMethod = a.abi.GetMethod(method); abiMethod == nil {
			return nil, fmt.Errorf("method %s not found", method)
		}
	}
	if abiMethod != nil {
		data, err := abi.Encode(args[:length-1], abiMethod.Inputs)
		if err != nil {
			return nil, fmt.Errorf("failed to encode arguments: %v", err)
		}
		if isContractDeployment {
			input = append(input, data...)
		} else {
			input = append(abiMethod.ID(), data...)
		}
	}
	txCfg := args[length-1].(*TxnOpts)
	txn, err := a.provider.Txn(a.addr, a.key, input, txCfg)
	if err != nil {
		return nil, err
	}
	return txn, nil
}

type CallOpts struct {
	Block ethgo.BlockNumber
	From  ethgo.Address
}

func (a *Contract) Call(method string, block ethgo.BlockNumber, args ...interface{}) (map[string]interface{}, error) {
	m := a.abi.GetMethod(method)
	if m == nil {
		return nil, fmt.Errorf("method %s not found", method)
	}

	data, err := m.Encode(args)
	if err != nil {
		return nil, err
	}

	opts := &CallOpts{
		Block: ethgo.Latest,
	}
	if a.key != nil {
		opts.From = a.key.Address()
	}
	rawOutput, err := a.provider.Call(a.addr, data, opts)
	if err != nil {
		return nil, err
	}

	resp, err := m.Decode(rawOutput)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
