package contract

import (
	"crypto/ecdsa"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
	"github.com/umbracle/ethgo/jsonrpc"
)

type Txn interface {
	Hash() ethgo.Hash
	EstimatedGas() uint64
	GasPrice() uint64
	Do() error
	Wait() (*ethgo.Receipt, error)
}

type TxnOpts struct {
	Value                *big.Int
	GasPrice             uint64
	GasLimit             uint64
	Nonce                uint64
	MaxFeePerGas         *big.Int
	MaxPriorityFeePerGas *big.Int
}

// Contract is a wrapper to make abi calls to contract with a state provider
type Contract struct {
	addr     common.Address
	abi      *abi.ABI
	bin      []byte
	provider Provider
	key      *ecdsa.PrivateKey
	fromAddr common.Address
}

func DeployContract(abi *abi.ABI, bin []byte,
	args []interface{}, txOpts *TxnOpts, opts ...ContractOption) (ethgo.Hash, error) {

	a := NewContract(common.Address{}, abi, opts...)
	a.bin = bin
	return a.Txn("constructor", append(args, txOpts))
}

func NewContract(addr common.Address, abi *abi.ABI, opts ...ContractOption) *Contract {
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
		provider = NewJonRPCNodeProvider(opt.JsonRPCClient, opt.chainId)
	} else {
		client, _ := jsonrpc.NewClient(opt.JsonRPCEndpoint)
		provider = NewJonRPCNodeProvider(client.Eth(), opt.chainId)
	}

	a := &Contract{
		addr:     addr,
		abi:      abi,
		provider: provider,
		key:      opt.Sender,
		fromAddr: crypto.PubkeyToAddress(opt.Sender.PublicKey),
	}
	return a
}

func (a *Contract) GetABI() *abi.ABI {
	return a.abi
}

func (a *Contract) Txn(method string, args ...interface{}) (ethgo.Hash, error) {
	if a.key == nil {
		return ethgo.ZeroHash, fmt.Errorf("no key selected")
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
			return ethgo.ZeroHash, fmt.Errorf("method %s not found", method)
		}
	}
	if abiMethod != nil {
		data, err := abi.Encode(args[:length-1], abiMethod.Inputs)
		if err != nil {
			return ethgo.ZeroHash, fmt.Errorf("failed to encode arguments: %v", err)
		}
		if isContractDeployment {
			input = append(input, data...)
		} else {
			input = append(abiMethod.ID(), data...)
		}
	}
	txCfg := args[length-1].(*TxnOpts)
	txHash, err := a.provider.SendEIP1559Tx(a.addr, a.key, input, txCfg)
	return txHash, err
}

//addr := crypto.PubkeyToAddress(privKey.PublicKey)

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
		Block: block,
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
