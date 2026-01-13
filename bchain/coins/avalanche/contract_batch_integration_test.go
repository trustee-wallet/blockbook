//go:build integration

package avalanche

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/trezor/blockbook/bchain/coins"
)

const defaultAvaxRpcURL = "http://localhost:8098/ext/bc/C/rpc"

func TestAvalancheErc20ContractBalancesIntegration(t *testing.T) {
	coins.RunERC20BatchBalanceTest(t, coins.ERC20BatchCase{
		Name:   "avalanche",
		RPCURL: defaultAvaxRpcURL,
		// Token-rich address on Avalanche C-Chain (balanceOf works for any address).
		Addr: common.HexToAddress("0x60aE616a2155Ee3d9A68541Ba4544862310933d4"),
		Contracts: []common.Address{
			common.HexToAddress("0xB31f66AA3C1e785363F0875A1B74E27b85FD66c7"), // WAVAX
			common.HexToAddress("0xA7D7079b0FEAD91F3e65f86E8915Cb59c1a4C664"), // USDC.e
			common.HexToAddress("0xc7198437980c041c805A1EDcbA50c1Ce5db95118"), // USDT.e
			common.HexToAddress("0xd586e7f844cea2f87f50152665bcbc2c279d8d70"), // DAI.e
			common.HexToAddress("0x49D5c2BdFfac6Ce2BFdB6640F4F80f226bc10bAB"), // WETH.e
			common.HexToAddress("0x60781C2586D68229fde47564546784ab3fACA982"), // PNG
		},
		BatchSize:       200,
		SkipUnavailable: true,
	})
}
