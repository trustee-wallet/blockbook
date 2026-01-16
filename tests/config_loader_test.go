//go:build integration

package tests

import (
	"testing"

	"github.com/trezor/blockbook/bchain"
)

// TestLoadBlockchainCfgEnvOverride verifies env-based overrides land in blockchaincfg.json.
func TestLoadBlockchainCfgEnvOverride(t *testing.T) {
	const want = "ws://backend_hostname:1234"
	t.Setenv("BB_RPC_URL_ethereum_archive", want)

	cfg := bchain.LoadBlockchainCfg(t, "ethereum_archive")
	if cfg.RpcUrl != want {
		t.Fatalf("expected rpc_url %q, got %q", want, cfg.RpcUrl)
	}
}
