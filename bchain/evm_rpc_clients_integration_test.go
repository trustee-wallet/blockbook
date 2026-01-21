//go:build integration

package bchain_test

import (
	"context"
	"testing"
	"time"

	"github.com/trezor/blockbook/bchain"
	"github.com/trezor/blockbook/bchain/coins/avalanche"
	"github.com/trezor/blockbook/bchain/coins/eth"
)

type openRPCFunc func(string, string) (bchain.EVMRPCClient, bchain.EVMClient, error)

func TestEVMRPCClients(t *testing.T) {
	openRPCOverrides := map[string]openRPCFunc{
		"avalanche": avalanche.OpenRPC,
	}
	aliases := []string{"ethereum", "avalanche", "arbitrum", "base", "bsc", "optimism", "polygon"}

	for _, alias := range aliases {
		alias := alias
		openRPC := eth.OpenRPC
		if override, ok := openRPCOverrides[alias]; ok {
			openRPC = override
		}
		t.Run(alias, func(t *testing.T) {
			runEVMRPCClientIntegrationTest(t, alias, openRPC)
		})
	}
}

func runEVMRPCClientIntegrationTest(t *testing.T, coinAlias string, openRPC openRPCFunc) {
	t.Helper()

	cfg := bchain.LoadBlockchainCfg(t, coinAlias)
	if cfg.RpcUrl == "" {
		t.Fatalf("empty rpc_url for %s", coinAlias)
	}
	if cfg.RpcUrlWs == "" {
		t.Fatalf("empty rpc_url_ws for %s", coinAlias)
	}

	rpcClient, _, err := openRPC(cfg.RpcUrl, cfg.RpcUrlWs)
	if err != nil {
		t.Fatalf("open rpc clients: %v", err)
	}
	defer rpcClient.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var version string
	if err := rpcClient.CallContext(ctx, &version, "web3_clientVersion"); err != nil {
		t.Fatalf("CallContext web3_clientVersion failed: %v", err)
	}
	if version == "" {
		t.Fatalf("empty web3_clientVersion")
	}

	subCtx, subCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer subCancel()
	sub, err := rpcClient.EthSubscribe(subCtx, make(chan interface{}, 1), "newHeads")
	if err != nil {
		t.Fatalf("EthSubscribe newHeads failed: %v", err)
	}
	sub.Unsubscribe()
}
