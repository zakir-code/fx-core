package keeper

import (
	"time"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (k Keeper) HandleRefundTokenClaim(ctx sdk.Context, claim *types.MsgRefundTokenClaim) {
	record, found := k.GetRefundRecord(ctx, claim.RefundNonce)
	if !found {
		return
	}
	// todo: If need be to slash unsigned oracles, can't delete refund record and refund confirm here
	// 1. delete refund record
	k.DeleteRefundRecord(ctx, record)

	// 2. delete record confirm
	k.DeleteRefundConfirm(ctx, claim.RefundNonce)

	// 3. delete snapshot oracle event nonce or snapshot oracle
	oracle, found := k.GetSnapshotOracle(ctx, record.OracleSetNonce)
	if !found {
		return
	}

	for i, nonce := range oracle.EventNonces {
		if nonce == claim.RefundNonce {
			oracle.EventNonces = append(oracle.EventNonces[:i], oracle.EventNonces[i+1:]...)
			break
		}
	}
	if len(oracle.EventNonces) == 0 {
		k.DeleteSnapshotOracle(ctx, record.OracleSetNonce)
	} else {
		k.SetSnapshotOracle(ctx, oracle)
	}
}

func (k Keeper) AddRefundRecord(ctx sdk.Context, receiver string, eventNonce uint64, tokens []types.ERC20Token) error {
	oracleSet := k.GetLatestOracleSet(ctx)
	if oracleSet == nil {
		return errorsmod.Wrap(types.ErrInvalid, "no oracle set")
	}
	snapshotOracle, found := k.GetSnapshotOracle(ctx, oracleSet.Nonce)
	if !found {
		snapshotOracle = &types.SnapshotOracle{
			OracleSetNonce: oracleSet.Nonce,
			Members:        oracleSet.Members,
			EventNonces:    []uint64{},
		}
	}
	snapshotOracle.EventNonces = append(snapshotOracle.EventNonces, eventNonce)
	k.SetSnapshotOracle(ctx, snapshotOracle)

	k.SetRefundRecord(ctx, &types.RefundRecord{
		EventNonce:     eventNonce,
		Receiver:       receiver,
		Timeout:        uint64(ctx.BlockTime().Add(time.Hour * 24 * 7).Second()), // TODO need to be configurable
		OracleSetNonce: oracleSet.Nonce,
		Tokens:         tokens,
	})
	return nil
}

func (k Keeper) SetRefundRecord(ctx sdk.Context, record *types.RefundRecord) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetBridgeCallRefundEventNonceKey(record.EventNonce), k.cdc.MustMarshal(record))
	store.Set(types.GetBridgeCallRefundKey(record.Receiver, record.EventNonce), sdk.Uint64ToBigEndian(record.OracleSetNonce))
}

func (k Keeper) DeleteRefundRecord(ctx sdk.Context, record *types.RefundRecord) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetBridgeCallRefundEventNonceKey(record.EventNonce))
	store.Delete(types.GetBridgeCallRefundKey(record.Receiver, record.EventNonce))
}

func (k Keeper) GetRefundRecord(ctx sdk.Context, eventNonce uint64) (*types.RefundRecord, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetBridgeCallRefundEventNonceKey(eventNonce))
	if bz == nil {
		return nil, false
	}
	refundRecord := new(types.RefundRecord)
	k.cdc.MustUnmarshal(bz, refundRecord)
	return refundRecord, true
}

func (k Keeper) SetSnapshotOracle(ctx sdk.Context, snapshotOracleKey *types.SnapshotOracle) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetSnapshotOracleKey(snapshotOracleKey.OracleSetNonce), k.cdc.MustMarshal(snapshotOracleKey))
}

func (k Keeper) GetSnapshotOracle(ctx sdk.Context, oracleSetNonce uint64) (*types.SnapshotOracle, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetSnapshotOracleKey(oracleSetNonce))
	if bz == nil {
		return nil, false
	}
	snapshotOracle := new(types.SnapshotOracle)
	k.cdc.MustUnmarshal(bz, snapshotOracle)
	return snapshotOracle, true
}

func (k Keeper) DeleteSnapshotOracle(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetSnapshotOracleKey(nonce))
}

func (k Keeper) GetRefundConfirm(ctx sdk.Context, nonce uint64, addr sdk.AccAddress) (*types.MsgConfirmRefund, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetRefundConfirmKey(nonce, addr))
	if bz == nil {
		return nil, false
	}
	var msg types.MsgConfirmRefund
	k.cdc.MustUnmarshal(bz, &msg)
	return &msg, true
}

func (k Keeper) SetRefundConfirm(ctx sdk.Context, addr sdk.AccAddress, msg *types.MsgConfirmRefund) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetRefundConfirmKey(msg.Nonce, addr), k.cdc.MustMarshal(msg))
}

func (k Keeper) DeleteRefundConfirm(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetRefundConfirmKeyByNonce(nonce))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		store.Delete(iterator.Key())
	}
}
