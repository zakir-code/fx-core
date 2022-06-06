package keeper

import (
	"encoding/hex"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"

	"github.com/functionx/fx-core/x/gravity/types"
)

var targetEvmPrefix = hex.EncodeToString([]byte("module/evm"))

func (a AttestationHandler) handlerRelayTransfer(ctx sdk.Context, claim *types.MsgDepositClaim, receiver sdk.AccAddress, coin sdk.Coin) {
	if claim.TargetIbc == targetEvmPrefix {
		a.handlerEvmTransfer(ctx, claim, receiver, coin)
		return
	}
	a.handleIbcTransfer(ctx, claim, receiver, coin)
}

func (a AttestationHandler) handleIbcTransfer(ctx sdk.Context, claim *types.MsgDepositClaim, receiveAddr sdk.AccAddress, coin sdk.Coin) {
	logger := a.keeper.Logger(ctx)
	targetIBC, ok := fxtypes.ParseHexTargetIBC(claim.TargetIbc)
	if !ok {
		logger.Error("convert target ibc data error!!!", "targetIbc", claim.GetTargetIbc())
		return
	}
	ibcReceiveAddress, err := types.CovertIbcPacketReceiveAddressByPrefix(targetIBC.Prefix, receiveAddr)
	if err != nil {
		logger.Error("convert ibc transfer receive address error!!!", "fxReceive", claim.FxReceiver,
			"ibcPrefix", targetIBC.Prefix, "sourcePort", targetIBC.SourcePort, "sourceChannel", targetIBC.SourceChannel, "error", err)
		return
	}

	_, clientState, err := a.keeper.ibcChannelKeeper.GetChannelClientState(ctx, targetIBC.SourcePort, targetIBC.SourceChannel)
	if err != nil {
		logger.Error("get channel client state error!!!", "sourcePort", targetIBC.SourcePort, "sourceChannel", targetIBC.SourceChannel)
		return
	}
	params := a.keeper.GetParams(ctx)
	clientStateHeight := clientState.GetLatestHeight()
	destTimeoutHeight := clientStateHeight.GetRevisionHeight() + params.IbcTransferTimeoutHeight
	ibcTimeoutHeight := ibcclienttypes.Height{
		RevisionNumber: clientStateHeight.GetRevisionNumber(),
		RevisionHeight: destTimeoutHeight,
	}
	nextSequenceSend, found := a.keeper.ibcChannelKeeper.GetNextSequenceSend(ctx, targetIBC.SourcePort, targetIBC.SourceChannel)
	if !found {
		logger.Error("ibc channel next sequence send not found!!!", "sourcePort", targetIBC.SourcePort, "sourceChannel", targetIBC.SourceChannel)
		return
	}
	logger.Info("gravity start ibc transfer", "sender", receiveAddr, "receive", ibcReceiveAddress, "coin", coin, "destCurrentHeight", clientStateHeight.GetRevisionHeight(), "destTimeoutHeight", destTimeoutHeight, "nextSequenceSend", nextSequenceSend)
	if err = a.keeper.ibcTransferKeeper.SendTransfer(ctx,
		targetIBC.SourcePort, targetIBC.SourceChannel,
		coin, receiveAddr, ibcReceiveAddress,
		ibcTimeoutHeight, 0,
		"", sdk.NewCoin(coin.Denom, sdk.ZeroInt())); err != nil {
		logger.Error("gravity ibc transfer fail. ", "sender", receiveAddr, "receive", ibcReceiveAddress, "coin", coin, "err", err)
		return
	}

	a.keeper.SetIbcSequenceHeight(ctx, targetIBC.SourcePort, targetIBC.SourceChannel, nextSequenceSend, uint64(ctx.BlockHeight()))

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeIbcTransfer,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(claim.EventNonce)),
		sdk.NewAttribute(types.AttributeKeyIbcSendSequence, fmt.Sprint(nextSequenceSend)),
		sdk.NewAttribute(types.AttributeKeyIbcSourcePort, targetIBC.SourcePort),
		sdk.NewAttribute(types.AttributeKeyIbcSourceChannel, targetIBC.SourceChannel),
	))
}

func (a AttestationHandler) handlerEvmTransfer(ctx sdk.Context, claim *types.MsgDepositClaim, receiver sdk.AccAddress, coin sdk.Coin) {
	logger := a.keeper.Logger(ctx)
	receiverEthType := common.BytesToAddress(receiver.Bytes())
	logger.Info("convert denom to fip20", "eth sender", claim.EthSender, "receiver", claim.FxReceiver,
		"receiver-eth-type", receiverEthType.String(), "amount", coin.String(), "target", claim.TargetIbc)
	err := a.keeper.erc20Keeper.RelayConvertCoin(ctx, receiver, receiverEthType, coin)
	if err != nil {
		logger.Error("evm transfer, convert denom to fip20 failed", "error", err.Error())
		return
	}
	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeEvmTransfer,
		sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprint(claim.EventNonce)),
	))
}
