//go:build integration

package eth

import (
	"testing"
	"time"

	"github.com/trezor/blockbook/bchain"
)

const (
	blockHeightLag = 1000
)

func newTestEthereumRPC(t *testing.T) *EthereumRPC {
	t.Helper()
	rpcURL := bchain.LoadBlockchainCfg(t, "ethereum").RpcUrl
	rc, ec, err := OpenRPC(rpcURL)
	if err != nil {
		t.Skipf("skipping: cannot connect to RPC at %s: %v", rpcURL, err)
		return nil
	}
	t.Cleanup(func() { rc.Close() })

	return &EthereumRPC{
		BaseChain:   &bchain.BaseChain{},
		Client:      ec,
		RPC:         rc,
		Timeout:     30 * time.Second,
		Parser:      NewEthereumParser(100, false),
		ChainConfig: &Configuration{},
	}
}

func assertBlockBasics(t *testing.T, block *bchain.Block, hash string, height uint32) {
	t.Helper()
	if block.Hash != hash {
		t.Fatalf("hash mismatch: got %s want %s", block.Hash, hash)
	}
	if block.Height != height {
		t.Fatalf("height mismatch: got %d want %d", block.Height, height)
	}
	if block.Confirmations <= 0 {
		t.Fatalf("expected confirmations > 0, got %d", block.Confirmations)
	}
	if block.Time <= 0 {
		t.Fatalf("expected block time > 0, got %d", block.Time)
	}
}

func assertTxSpecificData(t *testing.T, block *bchain.Block) {
	t.Helper()
	for i := range block.Txs {
		tx := &block.Txs[i]
		csd, ok := tx.CoinSpecificData.(bchain.EthereumSpecificData)
		if !ok {
			t.Fatalf("tx %d missing ethereum specific data", i)
		}
		if csd.Tx == nil {
			t.Fatalf("tx %d missing raw tx data", i)
		}
		if csd.Receipt == nil {
			t.Fatalf("tx %d missing receipt data", i)
		}
		if csd.Tx.Hash != tx.Txid {
			t.Fatalf("tx %d hash mismatch: raw %s vs txid %s", i, csd.Tx.Hash, tx.Txid)
		}
	}
}

// TestEthereumRPCGetBlockIntegration validates GetBlock by hash/height and exercises
// the parallel log + internal trace
func TestEthereumRPCGetBlockIntegration(t *testing.T) {
	prev := ProcessInternalTransactions
	t.Cleanup(func() { ProcessInternalTransactions = prev })

	rpcClient := newTestEthereumRPC(t)
	if rpcClient == nil {
		return
	}
	best, err := rpcClient.GetBestBlockHeight()
	if err != nil {
		t.Fatalf("GetBestBlockHeight: %v", err)
	}
	if best <= blockHeightLag {
		t.Skipf("best height %d too low for lag %d", best, blockHeightLag)
		return
	}
	height := best - blockHeightLag

	hash, err := rpcClient.GetBlockHash(height)
	if err != nil {
		t.Fatalf("GetBlockHash height %d: %v", height, err)
	}

	// Baseline: no internal tracing.
	ProcessInternalTransactions = false
	blockByHash, err := rpcClient.GetBlock(hash, 0)
	if err != nil {
		t.Fatalf("GetBlock by hash: %v", err)
	}
	assertBlockBasics(t, blockByHash, hash, height)
	assertTxSpecificData(t, blockByHash)

	blockByHeight, err := rpcClient.GetBlock("", height)
	if err != nil {
		t.Fatalf("GetBlock by height: %v", err)
	}
	assertBlockBasics(t, blockByHeight, hash, height)
	if len(blockByHeight.Txs) != len(blockByHash.Txs) {
		t.Fatalf("tx count mismatch: by hash %d vs by height %d", len(blockByHash.Txs), len(blockByHeight.Txs))
	}

	// Internal tracing enabled: should return the same block/tx set while logs and traces run in parallel.
	ProcessInternalTransactions = true
	blockWithTraces, err := rpcClient.GetBlock(hash, 0)
	if err != nil {
		t.Fatalf("GetBlock with internal traces: %v", err)
	}
	assertBlockBasics(t, blockWithTraces, hash, height)
	if len(blockWithTraces.Txs) != len(blockByHash.Txs) {
		t.Fatalf("tx count mismatch with traces: %d vs %d", len(blockWithTraces.Txs), len(blockByHash.Txs))
	}
	for i := range blockWithTraces.Txs {
		if blockWithTraces.Txs[i].Txid != blockByHash.Txs[i].Txid {
			t.Fatalf("txid mismatch at %d: %s vs %s", i, blockWithTraces.Txs[i].Txid, blockByHash.Txs[i].Txid)
		}
	}
	assertTxSpecificData(t, blockWithTraces)
}
