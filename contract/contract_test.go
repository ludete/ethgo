package contract

import (
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"

	"github.com/ethereum/go-ethereum/common"

	"github.com/stretchr/testify/assert"
	"github.com/umbracle/ethgo"
	"github.com/umbracle/ethgo/abi"
	"github.com/umbracle/ethgo/jsonrpc"
	"github.com/umbracle/ethgo/testutil"
)

var (
	addr0  = "0x0000000000000000000000000000000000000000"
	addr0B = ethgo.HexToAddress(addr0)
)

func TestContract_NoInput(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	cc := &testutil.Contract{}
	cc.AddOutputCaller("set")

	contract, addr := s.DeployContract(cc)

	abi0, err := abi.NewABI(contract.Abi)
	assert.NoError(t, err)

	p, _ := jsonrpc.NewClient(s.HTTPAddr())
	c := NewContract(common.Address(addr), abi0, WithJsonRPC(p.Eth()))

	vals, err := c.Call("set", ethgo.Latest)
	assert.NoError(t, err)
	assert.Equal(t, vals["0"], big.NewInt(1))

	abi1, err := abi.NewABIFromList([]string{
		"function set() view returns (uint256)",
	})
	assert.NoError(t, err)

	c1 := NewContract(common.Address(addr), abi1, WithJsonRPC(p.Eth()))
	vals, err = c1.Call("set", ethgo.Latest)
	assert.NoError(t, err)
	assert.Equal(t, vals["0"], big.NewInt(1))
}

func TestContract_IO(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	cc := &testutil.Contract{}
	cc.AddDualCaller("setA", "address", "uint256")

	contract, addr := s.DeployContract(cc)

	abi, err := abi.NewABI(contract.Abi)
	assert.NoError(t, err)

	c := NewContract(common.Address(addr), abi, WithJsonRPCEndpoint(s.HTTPAddr()))

	resp, err := c.Call("setA", ethgo.Latest, addr0B, 1000)
	assert.NoError(t, err)

	assert.Equal(t, resp["0"], addr0B)
	assert.Equal(t, resp["1"], big.NewInt(1000))
}

func TestContract_From(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	cc := &testutil.Contract{}
	cc.AddCallback(func() string {
		return `function example() public view returns (address) {
			return msg.sender;	
		}`
	})

	contract, addr := s.DeployContract(cc)

	abi, err := abi.NewABI(contract.Abi)
	assert.NoError(t, err)

	from := ethgo.Address{0x1}
	c := NewContract(common.Address(addr), abi, WithSender(nil), WithJsonRPCEndpoint(s.HTTPAddr()))

	resp, err := c.Call("example", ethgo.Latest)
	assert.NoError(t, err)
	assert.Equal(t, resp["0"], from)
}

func TestContract_Deploy(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	// create an address and fund it
	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	s.Transfer(ethgo.Address(addr), big.NewInt(1000000000000000000))

	//p, _ := jsonrpc.NewClient(s.HTTPAddr())

	cc := &testutil.Contract{}
	cc.AddConstructor("address", "uint256")

	artifact, err := cc.Compile()
	assert.NoError(t, err)

	abi, err := abi.NewABI(artifact.Abi)
	assert.NoError(t, err)

	bin, err := hex.DecodeString(artifact.Bin)
	assert.NoError(t, err)

	_, err = DeployContract(abi, bin, []interface{}{ethgo.Address{0x1}, 1000}, nil, WithSender(key))
	assert.NoError(t, err)

}

func TestContract_Transaction(t *testing.T) {
	s := testutil.NewTestServer(t, nil)
	defer s.Close()

	// create an address and fund it
	key, _ := crypto.GenerateKey()
	s.Transfer(ethgo.Address(crypto.PubkeyToAddress(key.PublicKey)), big.NewInt(1000000000000000000))

	cc := &testutil.Contract{}
	cc.AddEvent(testutil.NewEvent("A").Add("uint256", true))
	cc.EmitEvent("setA", "A", "1")

	artifact, addr := s.DeployContract(cc)

	abi, err := abi.NewABI(artifact.Abi)
	assert.NoError(t, err)

	// create a transaction
	i := NewContract(common.Address(addr), abi, WithJsonRPCEndpoint(s.HTTPAddr()), WithSender(key))
	_, err = i.Txn("setA")
	assert.NoError(t, err)

}
