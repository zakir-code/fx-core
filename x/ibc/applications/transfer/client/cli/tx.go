package cli

import (
	"fmt"
	"strings"

	"github.com/functionx/fx-core/x/ibc/applications/transfer/types"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/version"
	clienttypes "github.com/cosmos/cosmos-sdk/x/ibc/core/02-client/types"
	channelutils "github.com/cosmos/cosmos-sdk/x/ibc/core/04-channel/client/utils"
)

const (
	flagPacketTimeoutHeight    = "packet-timeout-height"
	flagPacketTimeoutTimestamp = "packet-timeout-timestamp"
	flagAbsoluteTimeouts       = "absolute-timeouts"
	flagIbcRouter              = "ibc-router"
	flagIbcFee                 = "ibc-fee"
)

// NewTransferTxCmd returns the command to create a NewMsgTransfer transaction
func NewTransferTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer [src-port] [src-channel] [receiver] [amount]",
		Short: "Transfer a fungible token through IBC",
		Long: strings.TrimSpace(`Transfer a fungible token through IBC. Timeouts can be specified
as absolute or relative using the "absolute-timeouts" flag. Timeout height can be set by passing in the height string
in the form {revision}-{height} using the "packet-timeout-height" flag. Relative timeouts are added to
the block height and block timestamp queried from the latest consensus state corresponding
to the counterparty channel. Any timeout set to 0 is disabled.`),
		Example: fmt.Sprintf("%s tx ibc-transfer transfer [src-port] [src-channel] [receiver] [amount]", version.AppName),
		Args:    cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}
			sender := clientCtx.GetFromAddress()
			srcPort := args[0]
			srcChannel := args[1]
			receiver := args[2]

			coin, err := sdk.ParseCoinNormalized(args[3])
			if err != nil {
				return err
			}

			if !strings.HasPrefix(coin.Denom, "ibc/") {
				denomTrace := types.ParseDenomTrace(coin.Denom)
				coin.Denom = denomTrace.IBCDenom()
			}

			timeoutHeightStr, err := cmd.Flags().GetString(flagPacketTimeoutHeight)
			if err != nil {
				return err
			}
			timeoutHeight, err := clienttypes.ParseHeight(timeoutHeightStr)
			if err != nil {
				return err
			}

			timeoutTimestamp, err := cmd.Flags().GetUint64(flagPacketTimeoutTimestamp)
			if err != nil {
				return err
			}

			absoluteTimeouts, err := cmd.Flags().GetBool(flagAbsoluteTimeouts)
			if err != nil {
				return err
			}

			// if the timeouts are not absolute, retrieve latest block height and block timestamp
			// for the consensus state connected to the destination port/channel
			if !absoluteTimeouts {
				consensusState, height, _, err := channelutils.QueryLatestConsensusState(clientCtx, srcPort, srcChannel)
				if err != nil {
					return err
				}

				if !timeoutHeight.IsZero() {
					absoluteHeight := height
					absoluteHeight.RevisionNumber += timeoutHeight.RevisionNumber
					absoluteHeight.RevisionHeight += timeoutHeight.RevisionHeight
					timeoutHeight = absoluteHeight
				}

				if timeoutTimestamp != 0 {
					timeoutTimestamp = consensusState.GetTimestamp() + timeoutTimestamp
				}
			}
			router, err := cmd.Flags().GetString(flagIbcRouter)
			if err != nil {
				return err
			}
			var ibcFee sdk.Coin
			if router != "" {
				ibcFeeString, err := cmd.Flags().GetString(flagIbcFee)
				if err != nil {
					return err
				}
				if ibcFeeAmount, ok := sdk.NewIntFromString(ibcFeeString); !ok {
					return fmt.Errorf("ibc-fee invalid!!!input:%v", ibcFeeString)
				} else {
					ibcFee = sdk.NewCoin(coin.Denom, ibcFeeAmount)
				}
			}

			msg := types.NewMsgTransfer(
				srcPort, srcChannel, coin, sender, receiver, timeoutHeight, timeoutTimestamp, router, ibcFee,
			)
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	cmd.Flags().String(flagPacketTimeoutHeight, "0-0", "Packet timeout block height. Example: 0-20000, The timeout is disabled when set to 0-0.")
	cmd.Flags().Uint64(flagPacketTimeoutTimestamp, types.DefaultRelativePacketTimeoutTimestamp, "Packet timeout timestamp in nanoseconds. Default is 24 hours. The timeout is disabled when set to 0.")
	cmd.Flags().String(flagIbcRouter, "", "Ibc transfer after router module, Default is nothing")
	cmd.Flags().String(flagIbcFee, "", "Ibc transfer after crosschain chain fee, Default is nothing")
	cmd.Flags().Bool(flagAbsoluteTimeouts, false, "Timeout flags are used as absolute timeouts.")
	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
