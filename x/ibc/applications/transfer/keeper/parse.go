package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	transfertypes "github.com/cosmos/ibc-go/v6/modules/apps/transfer/types"
	channeltypes "github.com/cosmos/ibc-go/v6/modules/core/04-channel/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/ibc/applications/transfer/types"
)

func parseIBCCoinDenom(packet channeltypes.Packet, packetDenom string) string {
	// This is the prefix that would have been prefixed to the denomination
	// on sender chain IF and only if the token originally came from the
	// receiving chain.
	//
	// NOTE: We use SourcePort and SourceChannel here, because the counterparty
	// chain would have prefixed with DestPort and DestChannel when originally
	// receiving this coin as seen in the "sender chain is the source" condition.

	var receiveDenom string
	if transfertypes.ReceiverChainIsSource(packet.GetSourcePort(), packet.GetSourceChannel(), packetDenom) {
		// sender chain is not the source, unescrow tokens

		// remove prefix added by sender chain
		voucherPrefix := transfertypes.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())
		unprefixedDenom := packetDenom[len(voucherPrefix):]

		// coin denomination used in sending from the escrow address
		denom := unprefixedDenom

		// The denomination used to send the coins is either the native denom or the hash of the path
		// if the denomination is not native.
		denomTrace := transfertypes.ParseDenomTrace(unprefixedDenom)
		if denomTrace.Path != "" {
			denom = denomTrace.IBCDenom()
		}
		receiveDenom = denom
	} else {
		// sender chain is the source, mint vouchers

		// since SendPacket did not prefix the denomination, we must prefix denomination here
		sourcePrefix := transfertypes.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel())
		// NOTE: sourcePrefix contains the trailing "/"
		prefixedDenom := sourcePrefix + packetDenom

		// construct the denomination trace from the full raw denomination
		denomTrace := transfertypes.ParseDenomTrace(prefixedDenom)

		receiveDenom = denomTrace.IBCDenom()
	}
	return receiveDenom
}

func parseReceiveAndAmountByPacket(data types.FungibleTokenPacketData) (sdk.AccAddress, bool, sdkmath.Int, sdkmath.Int, error) {
	// parse the transfer amount
	transferAmount, ok := sdkmath.NewIntFromString(data.Amount)
	if !ok {
		return nil, false, sdkmath.Int{}, sdkmath.Int{}, errorsmod.Wrapf(transfertypes.ErrInvalidAmount, "unable to parse transfer amount (%s) into sdkmath.Int", data.Amount)
	}

	if data.Router != "" {
		addressBytes, _, err := parsePacketAddress(data.Sender)
		if err != nil {
			return nil, false, sdkmath.Int{}, sdkmath.Int{}, err
		}
		feeAmount, ok := sdkmath.NewIntFromString(data.Fee)
		if !ok || feeAmount.IsNegative() {
			return nil, false, sdkmath.Int{}, sdkmath.Int{}, errorsmod.Wrapf(transfertypes.ErrInvalidAmount, "fee amount is invalid:%s", data.Fee)
		}
		return addressBytes, false, transferAmount, feeAmount, nil
	}

	// decode the receiver address
	receiverAddr, isEvmAddr, err := parsePacketAddress(data.Receiver)
	return receiverAddr, isEvmAddr, transferAmount, sdkmath.ZeroInt(), err
}

func parseAmountAndFeeByPacket(data types.FungibleTokenPacketData) (sdkmath.Int, sdkmath.Int, error) {
	// parse the transfer amount
	transferAmount, ok := sdkmath.NewIntFromString(data.Amount)
	if !ok {
		return sdkmath.Int{}, sdkmath.Int{}, errorsmod.Wrapf(transfertypes.ErrInvalidAmount, "unable to parse transfer amount (%s) into sdkmath.Int", data.Amount)
	}

	feeAmount := sdkmath.ZeroInt()
	if data.Router != "" {
		fee, ok := sdkmath.NewIntFromString(data.Fee)
		if !ok || fee.IsNegative() {
			return sdkmath.Int{}, sdkmath.Int{}, errorsmod.Wrapf(transfertypes.ErrInvalidAmount, "fee amount is invalid:%s", data.Fee)
		}
		feeAmount = fee
	}
	return transferAmount, feeAmount, nil
}

func parsePacketAddress(ibcSender string) (addr sdk.AccAddress, isEvmAddr bool, err error) {
	_, bytes, decodeErr := bech32.DecodeAndConvert(ibcSender)
	if decodeErr == nil {
		return bytes, false, nil
	}
	ethAddrError := fxtypes.ValidateEthereumAddress(ibcSender)
	if ethAddrError == nil {
		return common.HexToAddress(ibcSender).Bytes(), true, nil
	}
	return nil, false, decodeErr
}
