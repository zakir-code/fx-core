package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math/big"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	tmrand "github.com/tendermint/tendermint/libs/rand"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/functionx/fx-core/v7/contract"
	"github.com/functionx/fx-core/v7/testutil/helpers"
	fxtypes "github.com/functionx/fx-core/v7/types"
	"github.com/functionx/fx-core/v7/x/crosschain"
	"github.com/functionx/fx-core/v7/x/crosschain/types"
)

func (suite *KeeperTestSuite) TestABCIEndBlockDepositClaim() {
	normalMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracleAddrs[0].String(),
		BridgerAddress:   suite.bridgerAddrs[0].String(),
		ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey),
		ValidatorAddress: suite.valAddrs[0].String(),
		DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
		ChainName:        suite.chainName,
	}
	_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), normalMsg)
	require.NoError(suite.T(), err)

	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	bridgeToken := helpers.GenerateAddressByModule(suite.chainName)
	sendToFxSendAddr := helpers.GenerateAddressByModule(suite.chainName)
	addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
		EventNonce:     1,
		BlockHeight:    1000,
		TokenContract:  bridgeToken,
		Name:           "Test Token",
		Symbol:         "TEST",
		Decimals:       18,
		BridgerAddress: suite.bridgerAddrs[0].String(),
		ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
		ChainName:      suite.chainName,
	}
	_, err = suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), addBridgeTokenClaim)
	require.NoError(suite.T(), err)

	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	sendToFxClaim := &types.MsgSendToFxClaim{
		EventNonce:     2,
		BlockHeight:    1001,
		TokenContract:  bridgeToken,
		Amount:         sdkmath.NewInt(1234),
		Sender:         sendToFxSendAddr,
		Receiver:       sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
		TargetIbc:      hex.EncodeToString([]byte("px/transfer/channel-0")),
		BridgerAddress: suite.bridgerAddrs[0].String(),
		ChainName:      suite.chainName,
	}
	_, err = suite.MsgServer().SendToFxClaim(sdk.WrapSDKContext(suite.ctx), sendToFxClaim)
	require.NoError(suite.T(), err)

	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	allBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(sendToFxClaim.Receiver))
	denom := fmt.Sprintf("%s%s", suite.chainName, bridgeToken)
	trace, err := fxtypes.GetIbcDenomTrace(denom, addBridgeTokenClaim.ChannelIbc)
	suite.NoError(err)
	denom = trace.IBCDenom()
	require.EqualValues(suite.T(), fmt.Sprintf("%s%s", sendToFxClaim.Amount.String(), denom), allBalances.String())
}

func (suite *KeeperTestSuite) TestOracleUpdate() {
	if len(suite.oracleAddrs) < 10 {
		return
	}
	for i := 0; i < 10; i++ {
		msgBondedOracle := &types.MsgBondedOracle{
			OracleAddress:    suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.valAddrs[i].String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		require.NoError(suite.T(), msgBondedOracle.ValidateBasic())
		_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), msgBondedOracle)

		require.NoError(suite.T(), err)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		oracleSets := suite.Keeper().GetOracleSets(suite.ctx)
		require.NotNil(suite.T(), oracleSets)
		require.EqualValues(suite.T(), i+1, len(oracleSets))

		power := suite.Keeper().GetLastTotalPower(suite.ctx)
		expectPower := sdkmath.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdkmath.NewInt(int64(i + 1))).Quo(sdk.DefaultPowerReduction)
		require.True(suite.T(), expectPower.Equal(power))
	}

	bridgeToken := helpers.GenerateAddressByModule(suite.chainName)

	for i := 0; i < 6; i++ {
		addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     1,
			BlockHeight:    1000,
			TokenContract:  bridgeToken,
			Name:           "Test Token",
			Symbol:         "TEST",
			Decimals:       18,
			BridgerAddress: suite.bridgerAddrs[i].String(),
			ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
			ChainName:      suite.chainName,
		}
		_, err := suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), addBridgeTokenClaim)
		require.NoError(suite.T(), err)
		endBlockBeforeAttestation := suite.Keeper().GetAttestation(suite.ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), endBlockBeforeAttestation)
		require.False(suite.T(), endBlockBeforeAttestation.Observed)
		require.NotNil(suite.T(), endBlockBeforeAttestation.Votes)
		require.EqualValues(suite.T(), i+1, len(endBlockBeforeAttestation.Votes))

		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		endBlockAfterAttestation := suite.Keeper().GetAttestation(suite.ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), endBlockAfterAttestation)
		require.False(suite.T(), endBlockAfterAttestation.Observed)
	}

	addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
		EventNonce:     1,
		BlockHeight:    1000,
		TokenContract:  bridgeToken,
		Name:           "Test Token",
		Symbol:         "TEST",
		Decimals:       18,
		BridgerAddress: suite.bridgerAddrs[6].String(),
		ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
		ChainName:      suite.chainName,
	}
	_, err := suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), addBridgeTokenClaim)
	require.NoError(suite.T(), err)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	attestation := suite.Keeper().GetAttestation(suite.ctx, addBridgeTokenClaim.EventNonce, addBridgeTokenClaim.ClaimHash())

	require.NotNil(suite.T(), attestation)
	require.True(suite.T(), attestation.Observed)

	proposalHandler := crosschain.NewCrosschainProposalHandler(suite.app.CrosschainKeeper)

	var newOracleList []string
	for i := 0; i < 7; i++ {
		newOracleList = append(newOracleList, suite.oracleAddrs[i].String())
	}
	err = proposalHandler(suite.ctx, &types.UpdateChainOraclesProposal{ // nolint:staticcheck
		Title:       "proposal 1: try update chain oracle power >= 30%, expect error",
		Description: "",
		Oracles:     newOracleList,
		ChainName:   suite.chainName,
	})
	require.ErrorIs(suite.T(), types.ErrInvalid, err)

	expectTotalPower := sdkmath.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdkmath.NewInt(10)).Quo(sdk.DefaultPowerReduction)
	actualTotalPower := suite.Keeper().GetLastTotalPower(suite.ctx)
	require.True(suite.T(), expectTotalPower.Equal(actualTotalPower))

	expectMaxChangePower := types.AttestationProposalOracleChangePowerThreshold.Mul(expectTotalPower).Quo(sdkmath.NewInt(100))

	expectDeletePower := sdkmath.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdkmath.NewInt(3)).Quo(sdk.DefaultPowerReduction)
	require.EqualValues(suite.T(), fmt.Sprintf("max change power, maxChangePowerThreshold: %s, deleteTotalPower: %s: %s", expectMaxChangePower.String(), expectDeletePower.String(), types.ErrInvalid), err.Error())

	var newOracleList2 []string
	for i := 0; i < 8; i++ {
		newOracleList2 = append(newOracleList2, suite.oracleAddrs[i].String())
	}
	err = proposalHandler(suite.ctx, &types.UpdateChainOraclesProposal{ // nolint:staticcheck
		Title:       "proposal 2: try update chain oracle power <= 30%, expect success",
		Description: "",
		Oracles:     newOracleList2,
		ChainName:   suite.chainName,
	})
	require.NoError(suite.T(), err)
}

func (suite *KeeperTestSuite) TestAttestationAfterOracleUpdate() {
	if len(suite.bridgerAddrs) < 20 {
		return
	}
	for i := 0; i < 20; i++ {
		msgBondedOracle := &types.MsgBondedOracle{
			OracleAddress:    suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.valAddrs[i].String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), msgBondedOracle)
		require.NoError(suite.T(), err)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		oracleSets := suite.Keeper().GetOracleSets(suite.ctx)
		require.NotNil(suite.T(), oracleSets)
		require.EqualValues(suite.T(), i+1, len(oracleSets))

		power := suite.Keeper().GetLastTotalPower(suite.ctx)
		expectPower := sdkmath.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdkmath.NewInt(int64(i + 1))).Quo(sdk.DefaultPowerReduction)
		require.True(suite.T(), expectPower.Equal(power))
	}

	{
		firstBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     1,
			BlockHeight:    1000,
			TokenContract:  helpers.GenerateAddressByModule(suite.chainName),
			Name:           "Test Token",
			Symbol:         "TEST",
			Decimals:       18,
			BridgerAddress: "",
			ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
			ChainName:      suite.chainName,
		}

		for i := 0; i < 13; i++ {
			firstBridgeTokenClaim.BridgerAddress = suite.bridgerAddrs[i].String()
			_, err := suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), firstBridgeTokenClaim)
			require.NoError(suite.T(), err)
			endBlockBeforeAttestation := suite.Keeper().GetAttestation(suite.ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
			require.NotNil(suite.T(), endBlockBeforeAttestation)
			require.False(suite.T(), endBlockBeforeAttestation.Observed)
			require.NotNil(suite.T(), endBlockBeforeAttestation.Votes)
			require.EqualValues(suite.T(), i+1, len(endBlockBeforeAttestation.Votes))

			endBlockAfterAttestation := suite.Keeper().GetAttestation(suite.ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
			require.NotNil(suite.T(), endBlockAfterAttestation)
			require.False(suite.T(), endBlockAfterAttestation.Observed)
		}

		firstBridgeTokenClaim.BridgerAddress = suite.bridgerAddrs[13].String()
		_, err := suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), firstBridgeTokenClaim)
		require.NoError(suite.T(), err)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		attestation := suite.Keeper().GetAttestation(suite.ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())

		require.NotNil(suite.T(), attestation)
		require.True(suite.T(), attestation.Observed)
	}

	{
		secondBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     2,
			BlockHeight:    1001,
			TokenContract:  helpers.GenerateAddressByModule(suite.chainName),
			Name:           "Test Token2",
			Symbol:         "TEST2",
			Decimals:       18,
			BridgerAddress: "",
			ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
			ChainName:      suite.chainName,
		}

		for i := 0; i < 6; i++ {
			secondBridgeTokenClaim.BridgerAddress = suite.bridgerAddrs[i].String()
			_, err := suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), secondBridgeTokenClaim)
			require.NoError(suite.T(), err)
			endBlockBeforeAttestation := suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
			require.NotNil(suite.T(), endBlockBeforeAttestation)
			require.False(suite.T(), endBlockBeforeAttestation.Observed)
			require.NotNil(suite.T(), endBlockBeforeAttestation.Votes)
			require.EqualValues(suite.T(), i+1, len(endBlockBeforeAttestation.Votes))

			suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
			suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
			endBlockAfterAttestation := suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
			require.NotNil(suite.T(), endBlockAfterAttestation)
			require.False(suite.T(), endBlockAfterAttestation.Observed)
		}

		secondClaimAttestation := suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), secondClaimAttestation)
		require.False(suite.T(), secondClaimAttestation.Observed)
		require.NotNil(suite.T(), secondClaimAttestation.Votes)
		require.EqualValues(suite.T(), 6, len(secondClaimAttestation.Votes))

		var newOracleList []string
		for i := 0; i < 15; i++ {
			newOracleList = append(newOracleList, suite.oracleAddrs[i].String())
		}
		_, err := suite.MsgServer().UpdateChainOracles(suite.ctx, &types.MsgUpdateChainOracles{
			Oracles:   newOracleList,
			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			ChainName: suite.chainName,
		})
		require.NoError(suite.T(), err)
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

		secondClaimAttestation = suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), secondClaimAttestation)
		require.False(suite.T(), secondClaimAttestation.Observed)
		require.NotNil(suite.T(), secondClaimAttestation.Votes)
		require.EqualValues(suite.T(), 6, len(secondClaimAttestation.Votes))

		activeOracles := suite.Keeper().GetAllOracles(suite.ctx, true)
		require.NotNil(suite.T(), activeOracles)
		require.EqualValues(suite.T(), 15, len(activeOracles))
		for i := 0; i < 15; i++ {
			require.NotNil(suite.T(), newOracleList[i], activeOracles[i].OracleAddress)
		}

		var newOracleList2 []string
		for i := 0; i < 11; i++ {
			newOracleList2 = append(newOracleList2, suite.oracleAddrs[i].String())
		}
		_, err = suite.MsgServer().UpdateChainOracles(suite.ctx, &types.MsgUpdateChainOracles{
			Oracles:   newOracleList2,
			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			ChainName: suite.chainName,
		})
		require.NoError(suite.T(), err)
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

		secondClaimAttestation = suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), secondClaimAttestation)
		require.False(suite.T(), secondClaimAttestation.Observed)
		require.NotNil(suite.T(), secondClaimAttestation.Votes)
		require.EqualValues(suite.T(), 6, len(secondClaimAttestation.Votes))

		activeOracles = suite.Keeper().GetAllOracles(suite.ctx, true)
		require.NotNil(suite.T(), activeOracles)
		require.EqualValues(suite.T(), 11, len(activeOracles))
		for i := 0; i < 11; i++ {
			require.NotNil(suite.T(), newOracleList2[i], activeOracles[i].OracleAddress)
		}

		var newOracleList3 []string
		for i := 0; i < 10; i++ {
			newOracleList3 = append(newOracleList3, suite.oracleAddrs[i].String())
		}
		_, err = suite.MsgServer().UpdateChainOracles(suite.ctx, &types.MsgUpdateChainOracles{
			Oracles:   newOracleList3,
			Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
			ChainName: suite.chainName,
		})
		require.NoError(suite.T(), err)
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

		secondClaimAttestation = suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), secondClaimAttestation)
		require.False(suite.T(), secondClaimAttestation.Observed)
		require.NotNil(suite.T(), secondClaimAttestation.Votes)
		require.EqualValues(suite.T(), 6, len(secondClaimAttestation.Votes))

		activeOracles = suite.Keeper().GetAllOracles(suite.ctx, true)
		require.NotNil(suite.T(), activeOracles)
		require.EqualValues(suite.T(), 10, len(activeOracles))
		for i := 0; i < 10; i++ {
			require.NotNil(suite.T(), newOracleList3[i], activeOracles[i].OracleAddress)
		}

		secondBridgeTokenClaim.BridgerAddress = suite.bridgerAddrs[6].String()
		_, err = suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), secondBridgeTokenClaim)
		require.NoError(suite.T(), err)

		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

		secondClaimAttestation = suite.Keeper().GetAttestation(suite.ctx, secondBridgeTokenClaim.EventNonce, secondBridgeTokenClaim.ClaimHash())
		require.NotNil(suite.T(), secondClaimAttestation)
		require.True(suite.T(), secondClaimAttestation.Observed)
		require.NotNil(suite.T(), secondClaimAttestation.Votes)
		require.EqualValues(suite.T(), 7, len(secondClaimAttestation.Votes))
	}
}

func (suite *KeeperTestSuite) TestOracleDelete() {
	for i := 0; i < len(suite.oracleAddrs); i++ {
		msgBondedOracle := &types.MsgBondedOracle{
			OracleAddress:    suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.valAddrs[i].String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		require.NoError(suite.T(), msgBondedOracle.ValidateBasic())
		_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), msgBondedOracle)
		require.NoError(suite.T(), err)
	}
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	allOracles := suite.Keeper().GetAllOracles(suite.ctx, false)
	require.NotNil(suite.T(), allOracles)
	require.EqualValues(suite.T(), len(suite.oracleAddrs), len(allOracles))

	oracle := suite.oracleAddrs[0]
	bridger := suite.bridgerAddrs[0]
	externalAddress := suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey)

	oracleAddr, found := suite.Keeper().GetOracleAddressByBridgerKey(suite.ctx, bridger)
	require.True(suite.T(), found)
	require.EqualValues(suite.T(), oracle.String(), oracleAddr.String())

	oracleAddr, found = suite.Keeper().GetOracleByExternalAddress(suite.ctx, externalAddress)
	require.True(suite.T(), found)
	require.EqualValues(suite.T(), oracle.String(), oracleAddr.String())

	oracleData, found := suite.Keeper().GetOracle(suite.ctx, oracle)
	require.True(suite.T(), found)
	require.NotNil(suite.T(), oracleData)
	require.EqualValues(suite.T(), oracle.String(), oracleData.OracleAddress)
	require.EqualValues(suite.T(), bridger.String(), oracleData.BridgerAddress)
	require.EqualValues(suite.T(), externalAddress, oracleData.ExternalAddress)

	require.True(suite.T(), sdkmath.NewInt(10*1e3).MulRaw(1e18).Equal(oracleData.DelegateAmount))

	newOracleAddressList := make([]string, 0, len(suite.oracleAddrs)-1)
	for _, address := range suite.oracleAddrs[1:] {
		newOracleAddressList = append(newOracleAddressList, address.String())
	}

	_, err := suite.MsgServer().UpdateChainOracles(suite.ctx, &types.MsgUpdateChainOracles{
		Oracles:   newOracleAddressList,
		Authority: authtypes.NewModuleAddress(govtypes.ModuleName).String(),
		ChainName: suite.chainName,
	})
	require.NoError(suite.T(), err)
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	oracleAddr, found = suite.Keeper().GetOracleAddressByBridgerKey(suite.ctx, bridger)
	require.True(suite.T(), found)
	require.Equal(suite.T(), oracleAddr, oracle)

	oracleAddr, found = suite.Keeper().GetOracleByExternalAddress(suite.ctx, externalAddress)
	require.True(suite.T(), found)
	require.Equal(suite.T(), oracleAddr, oracle)

	oracleData, found = suite.Keeper().GetOracle(suite.ctx, oracle)
	require.True(suite.T(), found)
}

func (suite *KeeperTestSuite) TestOracleSetSlash() {
	for i := 0; i < len(suite.oracleAddrs); i++ {
		msgBondedOracle := &types.MsgBondedOracle{
			OracleAddress:    suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.valAddrs[i].String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		require.NoError(suite.T(), msgBondedOracle.ValidateBasic())
		_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), msgBondedOracle)
		require.NoError(suite.T(), err)
	}
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.Keeper().EndBlocker(suite.ctx)

	allOracles := suite.Keeper().GetAllOracles(suite.ctx, false)
	require.NotNil(suite.T(), allOracles)
	require.Equal(suite.T(), len(suite.oracleAddrs), len(allOracles))

	oracleSets := suite.Keeper().GetOracleSets(suite.ctx)
	require.NotNil(suite.T(), oracleSets)
	require.EqualValues(suite.T(), 1, len(oracleSets))

	for i := 0; i < len(suite.oracleAddrs)-1; i++ {
		externalAddress, signature := suite.SignOracleSetConfirm(suite.externalPris[i], oracleSets[0])
		oracleSetConfirm := &types.MsgOracleSetConfirm{
			Nonce:           oracleSets[0].Nonce,
			BridgerAddress:  suite.bridgerAddrs[i].String(),
			ExternalAddress: externalAddress,
			Signature:       hex.EncodeToString(signature),
			ChainName:       suite.chainName,
		}
		require.NoError(suite.T(), oracleSetConfirm.ValidateBasic())
		_, err := suite.MsgServer().OracleSetConfirm(sdk.WrapSDKContext(suite.ctx), oracleSetConfirm)
		require.NoError(suite.T(), err)
	}

	suite.Keeper().EndBlocker(suite.ctx)
	oracleSetHeight := int64(oracleSets[0].Height)
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	oracle, found := suite.Keeper().GetOracle(suite.ctx, suite.oracleAddrs[len(suite.oracleAddrs)-1])
	require.True(suite.T(), found)
	require.True(suite.T(), oracle.Online)
	require.Equal(suite.T(), int64(0), oracle.SlashTimes)

	suite.ctx = suite.ctx.WithBlockHeight(oracleSetHeight + int64(suite.Keeper().GetParams(suite.ctx).SignedWindow) + 1)
	suite.Keeper().EndBlocker(suite.ctx)

	oracle, found = suite.Keeper().GetOracle(suite.ctx, suite.oracleAddrs[len(suite.oracleAddrs)-1])
	require.True(suite.T(), found)
	require.False(suite.T(), oracle.Online)
	require.Equal(suite.T(), int64(1), oracle.SlashTimes)
}

func (suite *KeeperTestSuite) TestSlashOracle() {
	for i := 0; i < len(suite.oracleAddrs); i++ {
		msgBondedOracle := &types.MsgBondedOracle{
			OracleAddress:    suite.oracleAddrs[i].String(),
			BridgerAddress:   suite.bridgerAddrs[i].String(),
			ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[i].PublicKey),
			ValidatorAddress: suite.valAddrs[i].String(),
			DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
			ChainName:        suite.chainName,
		}
		require.NoError(suite.T(), msgBondedOracle.ValidateBasic())
		_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), msgBondedOracle)
		require.NoError(suite.T(), err)
	}

	params := suite.Keeper().GetParams(suite.ctx)
	err := suite.Keeper().SetParams(suite.ctx, &params)
	suite.Require().NoError(err)
	for i := 0; i < len(suite.oracleAddrs); i++ {
		oracle, found := suite.Keeper().GetOracle(suite.ctx, suite.oracleAddrs[i])
		require.True(suite.T(), found)
		require.True(suite.T(), oracle.Online)
		require.Equal(suite.T(), int64(0), oracle.SlashTimes)

		suite.Keeper().SlashOracle(suite.ctx, oracle.OracleAddress)

		oracle, found = suite.Keeper().GetOracle(suite.ctx, suite.oracleAddrs[i])
		require.True(suite.T(), found)
		require.False(suite.T(), oracle.Online)
		require.Equal(suite.T(), int64(1), oracle.SlashTimes)
	}

	// repeat slash test.
	for i := 0; i < len(suite.oracleAddrs); i++ {
		oracle, found := suite.Keeper().GetOracle(suite.ctx, suite.oracleAddrs[i])
		require.True(suite.T(), found)
		require.False(suite.T(), oracle.Online)
		require.Equal(suite.T(), int64(1), oracle.SlashTimes)

		suite.Keeper().SlashOracle(suite.ctx, oracle.OracleAddress)

		oracle, found = suite.Keeper().GetOracle(suite.ctx, suite.oracleAddrs[i])
		require.True(suite.T(), found)
		require.False(suite.T(), oracle.Online)
		require.Equal(suite.T(), int64(1), oracle.SlashTimes)
	}
}

func (suite *KeeperTestSuite) TestCleanUpRefundTimeout() {
	normalMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracleAddrs[0].String(),
		BridgerAddress:   suite.bridgerAddrs[0].String(),
		ExternalAddress:  suite.PubKeyToExternalAddr(suite.externalPris[0].PublicKey),
		ValidatorAddress: suite.valAddrs[0].String(),
		DelegateAmount:   types.NewDelegateAmount(sdkmath.NewInt(10 * 1e3).MulRaw(1e18)),
		ChainName:        suite.chainName,
	}
	_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), normalMsg)
	require.NoError(suite.T(), err)

	suite.Commit()

	bridgeToken := helpers.GenerateAddressByModule(suite.chainName)
	addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
		EventNonce:     1,
		BlockHeight:    1000,
		TokenContract:  bridgeToken,
		Name:           "Test Token",
		Symbol:         "TEST",
		Decimals:       18,
		BridgerAddress: suite.bridgerAddrs[0].String(),
		ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
		ChainName:      suite.chainName,
	}
	_, err = suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), addBridgeTokenClaim)
	require.NoError(suite.T(), err)

	denomResp, err := suite.QueryClient().TokenToDenom(sdk.WrapSDKContext(suite.ctx), &types.QueryTokenToDenomRequest{
		ChainName: suite.chainName,
		Token:     bridgeToken,
	})
	suite.NoError(err)

	_, err = suite.app.Erc20Keeper.RegisterNativeCoin(suite.ctx, banktypes.Metadata{
		Description: "Function X cross chain token",
		DenomUnits: []*banktypes.DenomUnit{
			{
				Denom:    "test",
				Exponent: 0,
				Aliases:  []string{denomResp.Denom},
			},
			{
				Denom:    "TEST",
				Exponent: 18,
				Aliases:  nil,
			},
		},
		Base:    "test",
		Display: "TEST",
		Name:    "Test Token",
		Symbol:  "TEST",
	})
	suite.NoError(err)

	suite.Commit()

	tokenAddr := common.BytesToAddress(types.ExternalAddressToAccAddress(suite.chainName, bridgeToken).Bytes())
	asset, err := contract.PackERC20AssetWithType([]common.Address{tokenAddr}, []*big.Int{big.NewInt(1)})
	suite.NoError(err)

	bridgeCallClaim := &types.MsgBridgeCallClaim{
		DstChainId:     types.FxcoreChainID,
		EventNonce:     2,
		Sender:         helpers.GenerateAddressByModule(suite.chainName),
		Receiver:       helpers.GenerateAddressByModule(suite.chainName),
		Asset:          asset,
		To:             helpers.GenerateAddressByModule(suite.chainName),
		Message:        hex.EncodeToString([]byte{0x1}),
		Value:          sdkmath.NewInt(1),
		GasLimit:       3000000,
		BlockHeight:    1001,
		BridgerAddress: suite.bridgerAddrs[0].String(),
		ChainName:      suite.chainName,
	}
	suite.ctx = suite.ctx.WithEventManager(sdk.NewEventManager())
	_, err = suite.MsgServer().BridgeCallClaim(sdk.WrapSDKContext(suite.ctx), bridgeCallClaim)
	suite.NoError(err)

	recordExist := false
	for _, event := range suite.ctx.EventManager().Events() {
		if event.Type == types.EventTypeBridgeCallRefund {
			recordExist = true
			break
		}
	}
	suite.True(recordExist)
	refundRecord, err := suite.QueryClient().RefundRecordByNonce(sdk.WrapSDKContext(suite.ctx), &types.QueryRefundRecordByNonceRequest{ChainName: suite.chainName, EventNonce: 2})
	suite.NoError(err)
	suite.Equal(uint64(2), refundRecord.Record.EventNonce)

	suite.Commit()

	sendToFxSendAddr := helpers.GenerateAddressByModule(suite.chainName)
	sendToFxClaim := &types.MsgSendToFxClaim{
		EventNonce:     3,
		BlockHeight:    refundRecord.Record.Timeout + 1,
		TokenContract:  bridgeToken,
		Amount:         sdkmath.NewInt(1234),
		Sender:         sendToFxSendAddr,
		Receiver:       sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
		TargetIbc:      hex.EncodeToString([]byte("px/transfer/channel-0")),
		BridgerAddress: suite.bridgerAddrs[0].String(),
		ChainName:      suite.chainName,
	}
	_, err = suite.MsgServer().SendToFxClaim(sdk.WrapSDKContext(suite.ctx), sendToFxClaim)
	require.NoError(suite.T(), err)

	_, err = suite.QueryClient().RefundRecordByNonce(sdk.WrapSDKContext(suite.ctx), &types.QueryRefundRecordByNonceRequest{ChainName: suite.chainName, EventNonce: 2})
	suite.ErrorIs(err, status.Error(codes.NotFound, "refund record"), suite.chainName)
}

func (suite *KeeperTestSuite) TestRefundSlashing() {
	suite.bondedOracle()
	suite.Commit()

	eventNonce := tmrand.Uint64()

	err := suite.Keeper().AddRefundRecord(suite.ctx, helpers.GenerateZeroAddressByModule(suite.chainName), eventNonce, []types.ERC20Token{})
	suite.NoError(err)

	params := suite.Keeper().GetParams(suite.ctx)
	signedWindow := uint64(tmrand.Int63n(10) + 2)
	params.SignedWindow = signedWindow
	suite.NoError(suite.Keeper().SetParams(suite.ctx, &params))

	slashedRefundNonce := suite.Keeper().GetLastSlashedRefundNonce(suite.ctx)
	suite.EqualValues(0, slashedRefundNonce)
	suite.Commit(int64(signedWindow))

	slashedRefundNonce = suite.Keeper().GetLastSlashedRefundNonce(suite.ctx)
	suite.EqualValues(eventNonce, slashedRefundNonce)

	oracle, found := suite.Keeper().GetOracle(suite.ctx, suite.oracleAddrs[0])
	suite.True(found)
	suite.False(oracle.Online)
	suite.EqualValues(1, oracle.SlashTimes)
}
