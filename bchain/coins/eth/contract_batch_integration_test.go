//go:build integration

package eth_test

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/trezor/blockbook/bchain/coins"
)

const defaultEthRpcURL = "http://localhost:8545"

func TestEthereumTypeGetErc20ContractBalancesIntegration(t *testing.T) {
	coins.RunERC20BatchBalanceTest(t, coins.ERC20BatchCase{
		Name:   "ethereum",
		RPCURL: defaultEthRpcURL,
		// Token-rich EOA (CEX hot wallet) used as a stable address reference.
		Addr: common.HexToAddress("0x28C6c06298d514Db089934071355E5743bf21d60"),
		Contracts: []common.Address{
			common.HexToAddress("0xA0b86991c6218b36c1d19d4a2e9eb0ce3606eb48"), // USDC
			common.HexToAddress("0xdAC17F958D2ee523a2206206994597C13D831ec7"), // USDT
			common.HexToAddress("0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2"), // WETH
			common.HexToAddress("0x6B175474E89094C44Da98b954EedeAC495271d0F"), // DAI
		},
		BatchSize:       200,
		SkipUnavailable: false,
	})
}
