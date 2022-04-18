package contract

import (
	"crypto/ecdsa"
	"encoding/hex"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	eTypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/jsonrpc"
)

// Provider handles the interactions with the Ethereum 1x node
type Provider interface {
	Client() *jsonrpc.Eth
	Call(common.Address, []byte, *CallOpts) ([]byte, error)
	SendEIP1559Tx(to common.Address,
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

func (j *JonRPCNodeProvider) Client() *jsonrpc.Eth {
	return j.client
}

func (j *JonRPCNodeProvider) Call(addr common.Address, input []byte, opts *CallOpts) ([]byte, error) {
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
	key *ecdsa.PrivateKey, input []byte, opts *TxnOpts) (ethgo.Hash, error) {

	var (
		err  error
		hash ethgo.Hash
	)
	dynamicFeeTx := &eTypes.DynamicFeeTx{
		Nonce:     opts.Nonce,
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
