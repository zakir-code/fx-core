package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/functionx/fx-core/x/intrarelayer/types"
)

// MintingEnabled checks that:
//  - the global parameter for intrarelaying is enabled
//  - minting is enabled for the given (erc20,coin) token pair
//  - recipient address is not on the blocked list
//  - bank module transfers are enabled for the Cosmos coin
func (k Keeper) MintingEnabled(ctx sdk.Context, sender, receiver sdk.AccAddress, token string) (types.TokenPair, error) {
	params := k.GetParams(ctx)
	if !params.EnableIntrarelayer {
		return types.TokenPair{}, sdkerrors.Wrap(types.ErrInternalTokenPair, "intrarelaying is currently disabled by governance")
	}
	id := k.GetTokenPairID(ctx, token)
	if len(id) == 0 {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrInternalTokenPair, "token %s not registered", token)
	}

	pair, found := k.GetTokenPair(ctx, id)
	if !found {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrInternalTokenPair, "not registered")
	}

	if !pair.Enabled {
		return types.TokenPair{}, sdkerrors.Wrapf(types.ErrNotAllowedBridge, "minting token %s is not enabled by governance", token)
	}

	if k.bankKeeper.BlockedAddr(receiver.Bytes()) {
		return types.TokenPair{}, sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "%s is not allowed to receive transactions", receiver)
	}

	// NOTE: ignore amount as only denom is checked on IsSendEnabledCoin
	coin := sdk.Coin{Denom: pair.Denom}

	// check if minting to a recipient address other than the sender is enabled for for the given coin denom
	if !sender.Equals(receiver) && !k.bankKeeper.SendEnabledCoin(ctx, coin) {
		return types.TokenPair{}, sdkerrors.Wrapf(banktypes.ErrSendDisabled, "minting %s coins to an external address is currently disabled", token)
	}

	return pair, nil
}