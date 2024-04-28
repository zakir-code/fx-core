package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (k Keeper) GetBridgeTokenDenom(ctx sdk.Context, tokenContract string) *types.BridgeToken {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.GetDenomToTokenKey(tokenContract))
	if len(data) == 0 {
		return nil
	}
	return &types.BridgeToken{
		Denom: string(data),
		Token: tokenContract,
	}
}

func (k Keeper) GetDenomBridgeToken(ctx sdk.Context, denom string) *types.BridgeToken {
	store := ctx.KVStore(k.storeKey)
	data := store.Get(types.GetTokenToDenomKey(denom))
	if len(data) == 0 {
		return nil
	}
	return &types.BridgeToken{
		Denom: denom,
		Token: string(data),
	}
}

func (k Keeper) HasBridgeToken(ctx sdk.Context, tokenContract string) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetDenomToTokenKey(tokenContract))
}

func (k Keeper) AddBridgeToken(ctx sdk.Context, token, denom string) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetTokenToDenomKey(denom), []byte(token))
	store.Set(types.GetDenomToTokenKey(token), []byte(denom))
}

func (k Keeper) IterateBridgeTokenToDenom(ctx sdk.Context, cb func(*types.BridgeToken) bool) {
	store := ctx.KVStore(k.storeKey)
	iter := sdk.KVStorePrefixIterator(store, types.TokenToDenomKey)
	defer iter.Close()

	for ; iter.Valid(); iter.Next() {
		bridgeToken := &types.BridgeToken{
			Denom: string(iter.Key()[len(types.TokenToDenomKey):]),
			Token: string(iter.Value()),
		}
		// cb returns true to stop early
		if cb(bridgeToken) {
			break
		}
	}
}

func (k Keeper) SetIbcDenomTrace(ctx sdk.Context, token, channelIBC string) (string, error) {
	denom := types.NewBridgeDenom(k.moduleName, token)
	denomTrace, err := fxtypes.GetIbcDenomTrace(denom, channelIBC)
	if err != nil {
		return denom, err
	}
	if denomTrace.Path != "" {
		k.ibcTransferKeeper.SetDenomTrace(ctx, denomTrace)
		return denomTrace.IBCDenom(), nil
	}
	return denomTrace.BaseDenom, nil
}
