package keeper

import (
	"encoding/hex"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	errortypes "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ethereum/go-ethereum/common"

	"github.com/functionx/fx-core/v7/contract"
	fxtypes "github.com/functionx/fx-core/v7/types"
)

// TransferAfter
// 1. Hook operation after transfer transaction triggered by IBC module
// 2. Hook operation after transferCrossChain triggered by ERC20 module
func (k Keeper) TransferAfter(ctx sdk.Context, sender sdk.AccAddress, receive string, amount, fee sdk.Coin, originToken bool) error {
	if err := contract.ValidateEthereumAddress(receive); err != nil {
		return errortypes.ErrInvalidAddress.Wrapf("invalid receive address: %s", err)
	}

	txID, err := k.AddToOutgoingPool(ctx, sender, receive, amount, fee)
	if err != nil {
		return err
	}

	if !originToken {
		k.erc20Keeper.SetOutgoingTransferRelation(ctx, k.moduleName, txID)
	}
	return nil
}

func (k Keeper) PrecompileCancelSendToExternal(ctx sdk.Context, txID uint64, sender sdk.AccAddress) (sdk.Coin, error) {
	return k.RemoveFromOutgoingPoolAndRefund(ctx, txID, sender)
}

func (k Keeper) PrecompileIncreaseBridgeFee(ctx sdk.Context, txID uint64, sender sdk.AccAddress, addBridgeFee sdk.Coin) error {
	return k.AddUnbatchedTxBridgeFee(ctx, txID, sender, addBridgeFee)
}

func (k Keeper) PrecompileBridgeCall(
	ctx sdk.Context,
	sender, receiver, to common.Address,
	coins sdk.Coins,
	message []byte,
	value *big.Int,
	gasLimit uint64,
) (eventNonce uint64, err error) {
	tokens, err := k.bridgeCallCoinsHandler(ctx, sender.Bytes(), coins)
	if err != nil {
		return 0, err
	}

	outCall, err := k.AddOutgoingBridgeCall(
		ctx,
		sender.Bytes(),
		fxtypes.AddressToStr(receiver.Bytes(), k.moduleName),
		fxtypes.AddressToStr(to.Bytes(), k.moduleName),
		tokens,
		hex.EncodeToString(message),
		sdkmath.NewIntFromBigInt(value),
		gasLimit,
	)
	if err != nil {
		return 0, err
	}

	return outCall.Nonce, nil
}
