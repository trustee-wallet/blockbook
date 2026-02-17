//go:build unittest

package server

import (
	"errors"
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

func TestGetEthereumInternalTransfersMissingData(t *testing.T) {
	tx := bchain.Tx{}

	transfers := getEthereumInternalTransfers(&tx)

	if len(transfers) != 0 {
		t.Fatalf("len(transfers) = %d, want 0", len(transfers))
	}
}

func TestGetEthereumInternalTransfersReturnsTransfers(t *testing.T) {
	expected := []bchain.EthereumInternalTransfer{
		{From: "0x111", To: "0x222"},
	}
	tx := bchain.Tx{
		CoinSpecificData: bchain.EthereumSpecificData{
			InternalData: &bchain.EthereumInternalData{
				Transfers: expected,
			},
		},
	}

	transfers := getEthereumInternalTransfers(&tx)

	if len(transfers) != len(expected) {
		t.Fatalf("len(transfers) = %d, want %d", len(transfers), len(expected))
	}
	if transfers[0].From != expected[0].From || transfers[0].To != expected[0].To {
		t.Fatalf("transfers[0] = %+v, want %+v", transfers[0], expected[0])
	}
}

func TestSetEthereumReceiptIfAvailableKeepsTxWhenReceiptFails(t *testing.T) {
	tx := bchain.Tx{
		Txid: "0xabc",
		CoinSpecificData: bchain.EthereumSpecificData{
			Tx: &bchain.RpcTransaction{Hash: "0xabc"},
		},
	}

	setEthereumReceiptIfAvailable(&tx, func(string) (*bchain.RpcReceipt, error) {
		return nil, errors.New("rpc failure")
	})

	csd, ok := tx.CoinSpecificData.(bchain.EthereumSpecificData)
	if !ok {
		t.Fatal("CoinSpecificData has unexpected type")
	}
	if csd.Receipt != nil {
		t.Fatalf("Receipt = %+v, want nil", csd.Receipt)
	}
}

func TestSetEthereumReceiptIfAvailableSetsReceipt(t *testing.T) {
	tx := bchain.Tx{
		Txid: "0xdef",
		CoinSpecificData: bchain.EthereumSpecificData{
			Tx: &bchain.RpcTransaction{Hash: "0xdef"},
		},
	}
	wantReceipt := &bchain.RpcReceipt{GasUsed: "0x5208"}

	setEthereumReceiptIfAvailable(&tx, func(string) (*bchain.RpcReceipt, error) {
		return wantReceipt, nil
	})

	csd, ok := tx.CoinSpecificData.(bchain.EthereumSpecificData)
	if !ok {
		t.Fatal("CoinSpecificData has unexpected type")
	}
	if csd.Receipt != wantReceipt {
		t.Fatalf("Receipt = %+v, want %+v", csd.Receipt, wantReceipt)
	}
}
