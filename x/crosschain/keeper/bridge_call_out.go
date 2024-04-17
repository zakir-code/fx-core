package keeper

import (
	"fmt"
	"math"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/gogo/protobuf/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (k Keeper) AddOutgoingBridgeCall(
	ctx sdk.Context,
	sender sdk.AccAddress, receiver, to string,
	tokens []types.ERC20Token, message string, value sdkmath.Int,
	gasLimit uint64,
) (*types.OutgoingBridgeCall, error) {
	params := k.GetParams(ctx)
	bridgeCallTimeout := k.CalExternalTimeoutHeight(ctx, params, params.BridgeCallTimeout)
	if bridgeCallTimeout <= 0 {
		return nil, errorsmod.Wrap(types.ErrInvalid, "bridge call timeout height")
	}

	oracleSet := k.GetLatestOracleSet(ctx)
	if oracleSet == nil {
		return nil, errorsmod.Wrap(types.ErrInvalid, "no oracle set")
	}

	nextID := k.autoIncrementID(ctx, types.KeyLastBridgeCallID)

	bridgeCall := &types.OutgoingBridgeCall{
		Nonce:          nextID,
		Timeout:        bridgeCallTimeout,
		Sender:         fxtypes.AddressToStr(sender, k.moduleName),
		Receiver:       receiver,
		To:             to,
		Tokens:         tokens,
		Message:        message,
		Value:          value,
		GasLimit:       gasLimit,
		OracleSetNonce: oracleSet.Nonce,
		BlockHeight:    uint64(ctx.BlockHeight()),
	}
	k.SetOutgoingBridgeCall(ctx, bridgeCall)

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBridgeCall,
		sdk.NewAttribute(sdk.AttributeKeyModule, k.moduleName),
		sdk.NewAttribute(sdk.AttributeKeySender, bridgeCall.Sender),
		sdk.NewAttribute(types.AttributeKeyBridgeCallNonce, fmt.Sprint(bridgeCall.Nonce)),
	))

	return bridgeCall, nil
}

func (k Keeper) SetOutgoingBridgeCall(ctx sdk.Context, out *types.OutgoingBridgeCall) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetOutgoingBridgeCallNonceKey(out.Nonce), k.cdc.MustMarshal(out))
	// value is just a placeholder
	store.Set(types.GetOutgoingBridgeCallAddressAndNonceKey(out.Sender, out.Nonce), k.cdc.MustMarshal(&gogotypes.BoolValue{Value: true}))
}

func (k Keeper) GetOutgoingBridgeCallByNonce(ctx sdk.Context, nonce uint64) (*types.OutgoingBridgeCall, bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.GetOutgoingBridgeCallNonceKey(nonce))
	if bz == nil {
		return nil, false
	}
	var out types.OutgoingBridgeCall
	k.cdc.MustUnmarshal(bz, &out)
	return &out, true
}

func (k Keeper) DeleteOutgoingBridgeCall(ctx sdk.Context, nonce uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetOutgoingBridgeCallNonceKey(nonce))
}

func (k Keeper) IterateOutgoingBridgeCalls(ctx sdk.Context, cb func(*types.OutgoingBridgeCall) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.OutgoingBridgeCallNonceKey)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var value types.OutgoingBridgeCall
		k.cdc.MustUnmarshal(iterator.Value(), &value)
		if cb(&value) {
			break
		}
	}
}

func (k Keeper) IterateOutgoingBridgeCallsByAddress(ctx sdk.Context, addr string, cb func(record *types.OutgoingBridgeCall) bool) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.GetOutgoingBridgeCallAddressKey(addr))
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		nonce := types.ParseOutgoingBridgeCallNonce(iterator.Key(), addr)
		record, found := k.GetOutgoingBridgeCallByNonce(ctx, nonce)
		if !found {
			continue
		}
		if cb(record) {
			break
		}
	}
}

func (k Keeper) HandleOutgoingBridgeCallRefund(ctx sdk.Context, data *types.OutgoingBridgeCall) {
	receiveAddr := types.ExternalAddressToAccAddress(k.moduleName, data.GetSender())
	if err := k.bridgeCallAssetRefundHandler(ctx, receiveAddr, data.Tokens); err != nil {
		panic(err)
	}
}

func (k Keeper) bridgeCallAssetRefundHandler(ctx sdk.Context, receive sdk.AccAddress, tokens []types.ERC20Token) error {
	coins, err := k.bridgeCallTransferToSender(ctx, receive, tokens)
	if err != nil {
		return err
	}
	return k.bridgeCallTransferToReceiver(ctx, receive, receive, coins)
}

func (k Keeper) IterateBridgeCallByNonce(ctx sdk.Context, startNonce uint64, cb func(bridgeCall *types.OutgoingBridgeCall) bool) {
	store := ctx.KVStore(k.storeKey)
	startKey := append(types.OutgoingBridgeCallNonceKey, sdk.Uint64ToBigEndian(startNonce)...)
	endKey := append(types.OutgoingBridgeCallNonceKey, sdk.Uint64ToBigEndian(math.MaxUint64)...)
	iter := store.Iterator(startKey, endKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		value := new(types.OutgoingBridgeCall)
		k.cdc.MustUnmarshal(iter.Value(), value)
		if cb(value) {
			break
		}
	}
}
