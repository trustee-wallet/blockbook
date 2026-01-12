package eth

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/trezor/blockbook/bchain"
)

type mockBatchRPC struct {
	results   map[string]string
	perErr    map[string]error
	lastBatch []rpc.BatchElem
}

func (m *mockBatchRPC) EthSubscribe(ctx context.Context, channel interface{}, args ...interface{}) (bchain.EVMClientSubscription, error) {
	return nil, errors.New("not implemented")
}

func (m *mockBatchRPC) CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error {
	return errors.New("not implemented")
}

func (m *mockBatchRPC) Close() {}

func (m *mockBatchRPC) BatchCallContext(ctx context.Context, batch []rpc.BatchElem) error {
	m.lastBatch = batch
	for i := range batch {
		elem := &batch[i]
		if elem.Method != "eth_call" {
			elem.Error = errors.New("unexpected method")
			continue
		}
		if len(elem.Args) < 2 {
			elem.Error = errors.New("missing args")
			continue
		}
		args, ok := elem.Args[0].(map[string]interface{})
		if !ok {
			elem.Error = errors.New("bad args")
			continue
		}
		to, _ := args["to"].(string)
		if err, ok := m.perErr[to]; ok {
			elem.Error = err
			continue
		}
		res, ok := m.results[to]
		if !ok {
			elem.Error = errors.New("missing result")
			continue
		}
		out, ok := elem.Result.(*string)
		if !ok {
			elem.Error = errors.New("bad result type")
			continue
		}
		*out = res
	}
	return nil
}

func TestEthereumTypeGetErc20ContractBalances(t *testing.T) {
	addr := common.HexToAddress("0x0000000000000000000000000000000000000011")
	contractA := common.HexToAddress("0x00000000000000000000000000000000000000aa")
	contractB := common.HexToAddress("0x00000000000000000000000000000000000000bb")
	contractAKey := hexutil.Encode(contractA.Bytes())
	contractBKey := hexutil.Encode(contractB.Bytes())
	mock := &mockBatchRPC{
		results: map[string]string{
			contractAKey: fmt.Sprintf("0x%064x", 123),
			contractBKey: fmt.Sprintf("0x%064x", 0),
		},
	}
	rpcClient := &EthereumRPC{
		RPC:     mock,
		Timeout: time.Second,
	}
	balances, err := rpcClient.EthereumTypeGetErc20ContractBalances(
		bchain.AddressDescriptor(addr.Bytes()),
		[]bchain.AddressDescriptor{
			bchain.AddressDescriptor(contractA.Bytes()),
			bchain.AddressDescriptor(contractB.Bytes()),
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(balances) != 2 {
		t.Fatalf("expected 2 balances, got %d", len(balances))
	}
	if balances[0] == nil || balances[0].Cmp(big.NewInt(123)) != 0 {
		t.Fatalf("unexpected balance[0]: %v", balances[0])
	}
	if balances[1] == nil || balances[1].Sign() != 0 {
		t.Fatalf("unexpected balance[1]: %v", balances[1])
	}
}

func TestEthereumTypeGetErc20ContractBalancesPartialError(t *testing.T) {
	addr := common.HexToAddress("0x0000000000000000000000000000000000000011")
	contractA := common.HexToAddress("0x00000000000000000000000000000000000000aa")
	contractB := common.HexToAddress("0x00000000000000000000000000000000000000bb")
	contractAKey := hexutil.Encode(contractA.Bytes())
	contractBKey := hexutil.Encode(contractB.Bytes())
	mock := &mockBatchRPC{
		results: map[string]string{
			contractAKey: fmt.Sprintf("0x%064x", 42),
		},
		perErr: map[string]error{
			contractBKey: errors.New("boom"),
		},
	}
	rpcClient := &EthereumRPC{
		RPC:     mock,
		Timeout: time.Second,
	}
	balances, err := rpcClient.EthereumTypeGetErc20ContractBalances(
		bchain.AddressDescriptor(addr.Bytes()),
		[]bchain.AddressDescriptor{
			bchain.AddressDescriptor(contractA.Bytes()),
			bchain.AddressDescriptor(contractB.Bytes()),
		},
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if balances[0] == nil || balances[0].Cmp(big.NewInt(42)) != 0 {
		t.Fatalf("unexpected balance[0]: %v", balances[0])
	}
	if balances[1] != nil {
		t.Fatalf("expected balance[1] to be nil, got %v", balances[1])
	}
}
