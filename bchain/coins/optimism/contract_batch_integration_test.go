//go:build integration

package optimism_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/trezor/blockbook/bchain"
	"github.com/trezor/blockbook/bchain/coins/eth"
)

func TestOptimismErc20ContractBalancesIntegration(t *testing.T) {
	bchain.RunERC20BatchBalanceTest(t, bchain.ERC20BatchCase{
		Name:   "optimism",
		RPCURL: bchain.LoadBlockchainCfg(t, "optimism").RpcUrl,
		Addr:   common.HexToAddress("0xDF90C9B995a3b10A5b8570a47101e6c6a29eb945"),
		Contracts: []common.Address{
			common.HexToAddress("0x4200000000000000000000000000000000000006"), // WETH
			common.HexToAddress("0x7F5c764cBc14f9669B88837ca1490cCa17c31607"), // USDC
			common.HexToAddress("0x94b008aa00579c1307b0ef2c499ad98a8ce58e58"), // USDT
			common.HexToAddress("0xDA10009cBd5D07dd0CeCc66161FC93D7c9000da1"), // DAI
			common.HexToAddress("0x4200000000000000000000000000000000000042"), // OP
		},
		BatchSize:       200,
		SkipUnavailable: true,
		NewClient:       eth.NewERC20BatchIntegrationClient,
	})
}
