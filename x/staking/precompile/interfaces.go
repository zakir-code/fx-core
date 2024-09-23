package precompile

import (
	"context"
	"math/big"
	"time"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	distrtypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	"github.com/ethereum/go-ethereum/common"
	evmtypes "github.com/evmos/ethermint/x/evm/types"
)

type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx context.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx context.Context, recipientAddr sdk.AccAddress, senderModule string, amt sdk.Coins) error
	SendCoinsFromModuleToModule(ctx context.Context, senderModule, recipientModule string, amt sdk.Coins) error
	GetBalance(ctx context.Context, addr sdk.AccAddress, denom string) sdk.Coin
}

type StakingKeeper interface {
	GetDelegation(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress) (delegation stakingtypes.Delegation, err error)
	SetDelegation(ctx context.Context, delegation stakingtypes.Delegation) error
	RemoveDelegation(ctx context.Context, delegation stakingtypes.Delegation) error
	GetValidator(ctx context.Context, addr sdk.ValAddress) (validator stakingtypes.Validator, err error)
	Delegate(
		ctx context.Context, delAddr sdk.AccAddress, bondAmt sdkmath.Int, tokenSrc stakingtypes.BondStatus,
		validator stakingtypes.Validator, subtractAccount bool,
	) (newShares sdkmath.LegacyDec, err error)
	HasMaxUnbondingDelegationEntries(ctx context.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress) (bool, error)
	Unbond(ctx context.Context, delAddr sdk.AccAddress, valAddr sdk.ValAddress, shares sdkmath.LegacyDec) (amount sdkmath.Int, err error)
	UnbondingTime(ctx context.Context) (time.Duration, error)
	SetUnbondingDelegationEntry(
		ctx context.Context, delegatorAddr sdk.AccAddress, validatorAddr sdk.ValAddress,
		creationHeight int64, minTime time.Time, balance sdkmath.Int,
	) (stakingtypes.UnbondingDelegation, error)
	InsertUBDQueue(ctx context.Context, ubd stakingtypes.UnbondingDelegation, completionTime time.Time) error
	GetAllowance(ctx sdk.Context, valAddr sdk.ValAddress, owner, spender sdk.AccAddress) *big.Int
	SetAllowance(ctx sdk.Context, valAddr sdk.ValAddress, owner, spender sdk.AccAddress, shares *big.Int)
	HasReceivingRedelegation(ctx context.Context, delAddr sdk.AccAddress, valDstAddr sdk.ValAddress) (bool, error)
	BeginRedelegation(
		ctx context.Context, delAddr sdk.AccAddress, valSrcAddr, valDstAddr sdk.ValAddress, sharesAmount sdkmath.LegacyDec,
	) (completionTime time.Time, err error)
	GetLastValidators(ctx context.Context) (validators []stakingtypes.Validator, err error)
}

type DistrKeeper interface {
	GetDelegatorWithdrawAddr(ctx context.Context, delAddr sdk.AccAddress) (sdk.AccAddress, error)
	IncrementValidatorPeriod(ctx context.Context, val stakingtypes.ValidatorI) (uint64, error)
	CalculateDelegationRewards(ctx context.Context, val stakingtypes.ValidatorI, del stakingtypes.DelegationI, endingPeriod uint64) (rewards sdk.DecCoins, err error)
	GetDelegatorStartingInfo(ctx context.Context, val sdk.ValAddress, del sdk.AccAddress) (period distrtypes.DelegatorStartingInfo, err error)
	SetDelegatorStartingInfo(ctx context.Context, val sdk.ValAddress, del sdk.AccAddress, period distrtypes.DelegatorStartingInfo) error
	DeleteDelegatorStartingInfo(ctx context.Context, val sdk.ValAddress, del sdk.AccAddress) error
	GetValidatorCurrentRewards(ctx context.Context, val sdk.ValAddress) (rewards distrtypes.ValidatorCurrentRewards, err error)
	GetValidatorHistoricalRewards(ctx context.Context, val sdk.ValAddress, period uint64) (rewards distrtypes.ValidatorHistoricalRewards, err error)
	SetValidatorHistoricalRewards(ctx context.Context, val sdk.ValAddress, period uint64, rewards distrtypes.ValidatorHistoricalRewards) error
	DeleteValidatorHistoricalReward(ctx context.Context, val sdk.ValAddress, period uint64) error
}

type EvmKeeper interface {
	GetParams(ctx sdk.Context) evmtypes.Params
}

type GovKeeper interface {
	CheckDisabledPrecompiles(ctx sdk.Context, contractAddress common.Address, methodId []byte) error
}

type SlashingKeeper interface {
	GetValidatorSigningInfo(ctx context.Context, address sdk.ConsAddress) (slashingtypes.ValidatorSigningInfo, error)
}
