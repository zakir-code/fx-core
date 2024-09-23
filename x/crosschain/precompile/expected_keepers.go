package precompile

import (
	"context"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	ibctransfertypes "github.com/cosmos/ibc-go/v8/modules/apps/transfer/types"
	"github.com/ethereum/go-ethereum/common"

	fxtypes "github.com/functionx/fx-core/v8/types"
	erc20types "github.com/functionx/fx-core/v8/x/erc20/types"
)

type Erc20Keeper interface {
	ModuleAddress() common.Address
	GetTokenPairByAddress(ctx sdk.Context, address common.Address) (erc20types.TokenPair, bool)
	ConvertDenomToTarget(ctx sdk.Context, from sdk.AccAddress, coin sdk.Coin, fxTarget fxtypes.FxTarget) (sdk.Coin, error)
	GetIbcTimeout(ctx sdk.Context) time.Duration
	SetIBCTransferRelation(ctx sdk.Context, channel string, sequence uint64)
	HasOutgoingTransferRelation(ctx sdk.Context, moduleName string, txID uint64) bool
	ToTargetDenom(ctx sdk.Context, denom, base string, aliases []string, fxTarget fxtypes.FxTarget) string
	GetTokenPair(ctx sdk.Context, tokenOrDenom string) (erc20types.TokenPair, bool)
	IsOriginDenom(ctx sdk.Context, denom string) bool
	HasDenomAlias(ctx sdk.Context, denom string) (banktypes.Metadata, bool)
}

type BankKeeper interface {
	GetSupply(ctx context.Context, denom string) sdk.Coin
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
	MintCoins(ctx context.Context, moduleName string, amounts sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	GetDenomMetaData(ctx context.Context, denom string) (banktypes.Metadata, bool)
}

type IBCTransferKeeper interface {
	Transfer(goCtx context.Context, msg *ibctransfertypes.MsgTransfer) (*ibctransfertypes.MsgTransferResponse, error)
}

type AccountKeeper interface {
	GetModuleAddress(moduleName string) sdk.AccAddress
}

type CrosschainKeeper interface {
	TransferAfter(ctx sdk.Context, sender sdk.AccAddress, receive string, coins, fee sdk.Coin, originToken, insufficientLiquidity bool) error
	PrecompileCancelSendToExternal(ctx sdk.Context, txID uint64, sender sdk.AccAddress) (sdk.Coin, error)
	PrecompileIncreaseBridgeFee(ctx sdk.Context, txID uint64, sender sdk.AccAddress, addBridgeFee sdk.Coin) error
	PrecompileBridgeCall(ctx sdk.Context, sender, refund common.Address, coins sdk.Coins, to common.Address, data, memo []byte) (uint64, error)
	PrecompileCancelPendingBridgeCall(ctx sdk.Context, nonce uint64, sender sdk.AccAddress) (sdk.Coins, error)
	PrecompileAddPendingPoolRewards(ctx sdk.Context, txID uint64, sender sdk.AccAddress, reward sdk.Coin) error
}

type GovKeeper interface {
	CheckDisabledPrecompiles(ctx sdk.Context, contractAddress common.Address, methodId []byte) error
}
