//go:build integration

package coins

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/trezor/blockbook/bchain"
	"github.com/trezor/blockbook/bchain/coins/eth"
)

const defaultBatchSize = 200

type ERC20BatchCase struct {
	Name            string
	RPCURL          string
	Addr            common.Address
	Contracts       []common.Address
	BatchSize       int
	SkipUnavailable bool
}

func RunERC20BatchBalanceTest(t *testing.T, tc ERC20BatchCase) {
	t.Helper()
	if tc.BatchSize <= 0 {
		tc.BatchSize = defaultBatchSize
	}
	rc, _, err := eth.OpenRPC(tc.RPCURL)
	if err != nil {
		handleRPCError(t, tc, fmt.Errorf("rpc dial error: %w", err))
		return
	}
	t.Cleanup(func() { rc.Close() })

	rpcClient := &eth.EthereumRPC{
		RPC:         rc,
		Timeout:     15 * time.Second,
		ChainConfig: &eth.Configuration{Erc20BatchSize: tc.BatchSize},
	}
	if err := verifyBatchBalances(rpcClient, tc.Addr, tc.Contracts); err != nil {
		handleRPCError(t, tc, err)
		return
	}
	chunkedContracts := expandContracts(tc.Contracts, tc.BatchSize+1)
	if err := verifyBatchBalances(rpcClient, tc.Addr, chunkedContracts); err != nil {
		handleRPCError(t, tc, err)
		return
	}
}

func handleRPCError(t *testing.T, tc ERC20BatchCase, err error) {
	t.Helper()
	if tc.SkipUnavailable && isRPCUnavailable(err) {
		t.Skipf("WARN: %s RPC not available: %v", tc.Name, err)
		return
	}
	t.Fatalf("%v", err)
}

func expandContracts(contracts []common.Address, minLen int) []common.Address {
	if len(contracts) >= minLen {
		return contracts
	}
	out := make([]common.Address, 0, minLen)
	for len(out) < minLen {
		out = append(out, contracts...)
	}
	if len(out) > minLen {
		out = out[:minLen]
	}
	return out
}

func verifyBatchBalances(rpcClient *eth.EthereumRPC, addr common.Address, contracts []common.Address) error {
	if len(contracts) == 0 {
		return errors.New("no contracts to query")
	}
	contractDescs := make([]bchain.AddressDescriptor, len(contracts))
	for i, c := range contracts {
		contractDescs[i] = bchain.AddressDescriptor(c.Bytes())
	}
	addrDesc := bchain.AddressDescriptor(addr.Bytes())
	balances, err := rpcClient.EthereumTypeGetErc20ContractBalances(addrDesc, contractDescs)
	if err != nil {
		return fmt.Errorf("batch balances error: %w", err)
	}
	if len(balances) != len(contractDescs) {
		return fmt.Errorf("expected %d balances, got %d", len(contractDescs), len(balances))
	}
	for i, contractDesc := range contractDescs {
		single, err := rpcClient.EthereumTypeGetErc20ContractBalance(addrDesc, contractDesc)
		if err != nil {
			return fmt.Errorf("single balance error for %s: %w", contracts[i].Hex(), err)
		}
		if balances[i] == nil {
			return fmt.Errorf("batch balance missing for %s", contracts[i].Hex())
		}
		if balances[i].Cmp(single) != 0 {
			return fmt.Errorf("balance mismatch for %s: batch=%s single=%s", contracts[i].Hex(), balances[i].String(), single.String())
		}
	}
	return nil
}

func isRPCUnavailable(err error) bool {
	if err == nil {
		return false
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	msg := strings.ToLower(err.Error())
	switch {
	case strings.Contains(msg, "context deadline exceeded"),
		strings.Contains(msg, "connection refused"),
		strings.Contains(msg, "no such host"),
		strings.Contains(msg, "i/o timeout"),
		strings.Contains(msg, "timeout"):
		return true
	}
	return false
}
