package keeper

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/hashicorp/go-metrics"

	"github.com/pundiai/fx-core/v8/contract"
	fxtypes "github.com/pundiai/fx-core/v8/types"
	"github.com/pundiai/fx-core/v8/x/crosschain/types"
)

func (k Keeper) BridgeCallExecuted(ctx sdk.Context, caller contract.Caller, msg *types.MsgBridgeCallClaim) error {
	k.CreateBridgeAccount(ctx, msg.TxOrigin)
	if senderAccount := k.ak.GetAccount(ctx, msg.GetSenderAddr().Bytes()); senderAccount != nil {
		if _, ok := senderAccount.(sdk.ModuleAccountI); ok {
			return types.ErrInvalid.Wrap("sender is module account")
		}
	}
	isMemoSendCallTo := msg.IsMemoSendCallTo()
	receiverAddr := msg.GetToAddr()
	if isMemoSendCallTo {
		receiverAddr = msg.GetSenderAddr()
	}

	baseCoins := sdk.NewCoins()
	for i, tokenAddr := range msg.TokenContracts {
		baseCoin, err := k.DepositBridgeTokenToBaseCoin(ctx, receiverAddr.Bytes(), msg.Amounts[i], tokenAddr)
		if err != nil {
			return err
		}
		baseCoins = baseCoins.Add(baseCoin)
	}

	if err := k.HandlerBridgeCallInFee(ctx, caller, msg.GetSenderAddr(), msg.QuoteId.BigInt(), msg.GasLimit.Uint64()); err != nil {
		return err
	}

	cacheCtx, commit := sdk.UnwrapSDKContext(ctx).CacheContext()
	err := k.BridgeCallEvm(cacheCtx, caller, msg.GetSenderAddr(), msg.GetRefundAddr(), msg.GetToAddr(),
		receiverAddr, baseCoins, msg.MustData(), msg.MustMemo(), isMemoSendCallTo, msg.GetGasLimit())
	if !ctx.IsCheckTx() {
		telemetry.IncrCounterWithLabels(
			[]string{types.ModuleName, "bridge_call_in"},
			float32(1),
			[]metrics.Label{
				telemetry.NewLabel("module", k.moduleName),
				telemetry.NewLabel("success", strconv.FormatBool(err == nil)),
			},
		)
	}
	if err == nil {
		commit()
		return nil
	}
	revertMsg := err.Error()

	// refund bridge-call case of error
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventTypeBridgeCallEvent, sdk.NewAttribute(types.AttributeKeyErrCause, err.Error())))

	if !baseCoins.Empty() {
		if err = k.bankKeeper.SendCoins(ctx, receiverAddr.Bytes(), msg.GetRefundAddr().Bytes(), baseCoins); err != nil {
			return err
		}

		for _, coin := range baseCoins {
			if coin.Denom == fxtypes.DefaultDenom {
				continue
			}
			if _, err = k.erc20Keeper.BaseCoinToEvm(ctx, caller, msg.GetRefundAddr(), coin); err != nil {
				return err
			}
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(
		types.EventTypeBridgeCallFailed,
		sdk.NewAttribute(types.AttributeKeyEventNonce, fmt.Sprintf("%d", msg.EventNonce)),
		sdk.NewAttribute(types.AttributeKeyBridgeCallFailedRefundAddr, msg.GetRefundAddr().Hex()),
	))

	// onRevert bridgeCall
	_, err = k.AddOutgoingBridgeCall(ctx, msg.GetToAddr(), common.Address{}, sdk.NewCoins(),
		msg.GetSenderAddr(), []byte(revertMsg), []byte{}, 0, msg.EventNonce)
	return err
}

func (k Keeper) BridgeCallEvm(ctx sdk.Context, caller contract.Caller, sender, refundAddr, to, receiverAddr common.Address, baseCoins sdk.Coins, data, memo []byte, isMemoSendCallTo bool, gasLimit uint64) error {
	tokens := make([]common.Address, 0, baseCoins.Len())
	amounts := make([]*big.Int, 0, baseCoins.Len())
	for _, coin := range baseCoins {
		tokenContract, err := k.erc20Keeper.BaseCoinToEvm(ctx, caller, receiverAddr, coin)
		if err != nil {
			return err
		}
		tokens = append(tokens, common.HexToAddress(tokenContract))
		amounts = append(amounts, coin.Amount.BigInt())
	}

	if !k.evmKeeper.IsContract(ctx, to) {
		return nil
	}
	var callEvmSender common.Address
	var args []byte

	if isMemoSendCallTo {
		args = data
		callEvmSender = sender
	} else {
		var err error
		args, err = contract.PackOnBridgeCall(sender, refundAddr, tokens, amounts, data, memo)
		if err != nil {
			return err
		}
		callEvmSender = k.GetCallbackFrom()
	}

	if gasLimit == 0 {
		gasLimit = k.GetBridgeCallMaxGasLimit(ctx)
	}
	txResp, err := caller.ExecuteEVM(ctx, callEvmSender, &to, nil, gasLimit, args)
	if err != nil {
		return err
	}
	if txResp.Failed() {
		return types.ErrInvalid.Wrap(txResp.VmError)
	}
	return nil
}

func (k Keeper) CreateBridgeAccount(ctx sdk.Context, address string) {
	accAddress := fxtypes.ExternalAddrToAccAddr(k.moduleName, address)
	if account := k.ak.GetAccount(ctx, accAddress); account != nil {
		return
	}
	k.ak.NewAccountWithAddress(ctx, accAddress)
}
