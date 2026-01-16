//go:build integration

package base_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/trezor/blockbook/bchain"
	"github.com/trezor/blockbook/bchain/coins/eth"
)

func TestBaseErc20ContractBalancesIntegration(t *testing.T) {
	bchain.RunERC20BatchBalanceTest(t, bchain.ERC20BatchCase{
		Name:   "base",
		RPCURL: bchain.LoadBlockchainCfg(t, "base").RpcUrl,
		Addr:   common.HexToAddress("0x242E2d70d3AdC00a9eF23CeD6E88811fCefCA788"),
		Contracts: []common.Address{
			common.HexToAddress("0x4200000000000000000000000000000000000006"), // WETH
			common.HexToAddress("0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"), // USDC
			common.HexToAddress("0x50c5725949A6F0c72E6C4a641F24049A917DB0Cb"), // DAI
			common.HexToAddress("0x2ae3f1ec7f1f5012cfeab0185bfc7aa3cf0dec22"), // cbETH
		},
		BatchSize:       200,
		SkipUnavailable: true,
		NewClient:       eth.NewERC20BatchIntegrationClient,
	})
}
