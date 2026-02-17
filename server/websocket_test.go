//go:build unittest

package server

import (
	"testing"

	"github.com/trezor/blockbook/bchain"
)

func TestSetConfirmedBlockTxMetadataSetsConfirmedFields(t *testing.T) {
	tx := bchain.Tx{
		Confirmations: 0,
		Blocktime:     0,
		Time:          0,
	}

	setConfirmedBlockTxMetadata(&tx, 123456)

	if tx.Confirmations != 1 {
		t.Fatalf("Confirmations = %d, want 1", tx.Confirmations)
	}
	if tx.Blocktime != 123456 {
		t.Fatalf("Blocktime = %d, want 123456", tx.Blocktime)
	}
	if tx.Time != 123456 {
		t.Fatalf("Time = %d, want 123456", tx.Time)
	}
}

func TestSetConfirmedBlockTxMetadataLeavesConfirmedTxUnchanged(t *testing.T) {
	tx := bchain.Tx{
		Confirmations: 3,
		Blocktime:     100,
		Time:          200,
	}

	setConfirmedBlockTxMetadata(&tx, 123456)

	if tx.Confirmations != 3 {
		t.Fatalf("Confirmations = %d, want 3", tx.Confirmations)
	}
	if tx.Blocktime != 100 {
		t.Fatalf("Blocktime = %d, want 100", tx.Blocktime)
	}
	if tx.Time != 200 {
		t.Fatalf("Time = %d, want 200", tx.Time)
	}
}
