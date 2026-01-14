//go:build integration

package bsc_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/trezor/blockbook/bchain"
	"github.com/trezor/blockbook/bchain/coins/eth"
)

func TestBNBSmartChainErc20ContractBalancesIntegration(t *testing.T) {
	bchain.RunERC20BatchBalanceTest(t, bchain.ERC20BatchCase{
		Name:   "bsc",
		RPCURL: bchain.RPCURLFromConfig(t, "bsc"),
		Addr:   common.HexToAddress("0x21d45650db732cE5dF77685d6021d7D5d1da807f"),
		Contracts: []common.Address{
			common.HexToAddress("0xBB4CdB9CBd36B01bD1cBaEBF2De08d9173bc095c"), // WBNB
			common.HexToAddress("0x55d398326f99059fF775485246999027B3197955"), // USDT
			common.HexToAddress("0xe9e7CEA3Dedca5984780Bafc599bd69ADd087d56"), // BUSD
			common.HexToAddress("0x8ac76a51cc950d9822d68b83fe1ad97b32cd580d"), // USDC
			common.HexToAddress("0x1AF3F329e8BE154074D8769D1FFa4eE058B1DBc3"), // DAI
		},
		BatchSize:       200,
		SkipUnavailable: true,
		NewClient:       eth.NewERC20BatchIntegrationClient,
	})
}
