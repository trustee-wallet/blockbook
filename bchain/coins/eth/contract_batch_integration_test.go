//go:build integration

package eth

import (
	"os"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/trezor/blockbook/bchain"
)

const defaultEthRpcURL = "http://naked:8545"

func integrationRpcURL() string {
	if v := os.Getenv("BLOCKBOOK_ETH_RPC_URL"); v != "" {
		return v
	}
	return defaultEthRpcURL
}

func TestEthereumTypeGetErc20ContractBalancesIntegration(t *testing.T) {
	rpcURL := integrationRpcURL()
	rc, _, err := OpenRPC(rpcURL)
	if err != nil {
		t.Skipf("skipping: cannot connect to RPC at %s: %v", rpcURL, err)
		return
	}
	defer rc.Close()

	// Use stable mainnet ERC20 contracts and a well-known EOA.
	addr := common.HexToAddress("0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045")
	contracts := []common.Address{
		common.HexToAddress("0xA0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"), // USDC
		common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), // USDT
		common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"), // WETH
	}
	contractDescs := make([]bchain.AddressDescriptor, len(contracts))
	for i, c := range contracts {
		contractDescs[i] = bchain.AddressDescriptor(c.Bytes())
	}

	rpcClient := &EthereumRPC{
		RPC:     rc,
		Timeout: 15 * time.Second,
	}
	addrDesc := bchain.AddressDescriptor(addr.Bytes())
	balances, err := rpcClient.EthereumTypeGetErc20ContractBalances(addrDesc, contractDescs)
	if err != nil {
		t.Fatalf("batch balances error: %v", err)
	}
	if len(balances) != len(contractDescs) {
		t.Fatalf("expected %d balances, got %d", len(contractDescs), len(balances))
	}
	for i, contractDesc := range contractDescs {
		single, err := rpcClient.EthereumTypeGetErc20ContractBalance(addrDesc, contractDesc)
		if err != nil {
			t.Fatalf("single balance error for %s: %v", contracts[i].Hex(), err)
		}
		if balances[i] == nil {
			t.Fatalf("batch balance missing for %s", contracts[i].Hex())
		}
		if balances[i].Cmp(single) != 0 {
			t.Fatalf("balance mismatch for %s: batch=%s single=%s", contracts[i].Hex(), balances[i].String(), single.String())
		}
	}
}
