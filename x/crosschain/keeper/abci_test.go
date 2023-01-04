package keeper_test

import (
	"encoding/hex"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ethereum/go-ethereum/crypto"
	tronAddress "github.com/fbsobreira/gotron-sdk/pkg/address"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/functionx/fx-core/v3/app/helpers"
	fxtypes "github.com/functionx/fx-core/v3/types"
	"github.com/functionx/fx-core/v3/x/crosschain"
	"github.com/functionx/fx-core/v3/x/crosschain/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

func (suite *KeeperTestSuite) TestABCIEndBlockDepositClaim() {
	normalMsg := &types.MsgBondedOracle{
		OracleAddress:    suite.oracles[0].String(),
		BridgerAddress:   suite.bridgers[0].String(),
		ExternalAddress:  crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex(),
		ValidatorAddress: suite.validator[0].String(),
		DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(10 * 1e3).MulRaw(1e18)},
		ChainName:        suite.chainName,
	}
	if trontypes.ModuleName == suite.chainName {
		normalMsg.ExternalAddress = tronAddress.PubkeyToAddress(suite.externals[0].PublicKey).String()
	}
	_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), normalMsg)
	require.NoError(suite.T(), err)

	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)

	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	bridgeToken := helpers.GenerateAddress().String()
	sendToFxSendAddr := helpers.GenerateAddress().String()
	if trontypes.ModuleName == suite.chainName {
		bridgeToken = trontypes.AddressFromHex(bridgeToken)
		sendToFxSendAddr = tronAddress.PubkeyToAddress(suite.externals[0].PublicKey).String()
	}
	addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
		EventNonce:     1,
		BlockHeight:    1000,
		TokenContract:  bridgeToken,
		Name:           "Test Token",
		Symbol:         "TEST",
		Decimals:       18,
		BridgerAddress: suite.bridgers[0].String(),
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
		Amount:         sdk.NewInt(1234),
		Sender:         sendToFxSendAddr,
		Receiver:       sdk.AccAddress(helpers.GenerateAddress().Bytes()).String(),
		TargetIbc:      hex.EncodeToString([]byte("px/transfer/channel-0")),
		BridgerAddress: suite.bridgers[0].String(),
		ChainName:      suite.chainName,
	}
	_, err = suite.MsgServer().SendToFxClaim(sdk.WrapSDKContext(suite.ctx), sendToFxClaim)
	require.NoError(suite.T(), err)

	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	allBalances := suite.app.BankKeeper.GetAllBalances(suite.ctx, sdk.MustAccAddressFromBech32(sendToFxClaim.Receiver))
	tokenName := fmt.Sprintf("%s%s", suite.chainName, bridgeToken)
	require.EqualValues(suite.T(), fmt.Sprintf("%s%s", sendToFxClaim.Amount.String(), tokenName), allBalances.String())
}

func (suite *KeeperTestSuite) TestOracleUpdate() {
	if len(suite.oracles) < 10 {
		return
	}
	for i := 0; i < 10; i++ {
		msgBondedOracle := &types.MsgBondedOracle{
			OracleAddress:    suite.oracles[i].String(),
			BridgerAddress:   suite.bridgers[i].String(),
			ExternalAddress:  crypto.PubkeyToAddress(suite.externals[i].PublicKey).Hex(),
			ValidatorAddress: suite.validator[i].String(),
			DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(10 * 1e3).MulRaw(1e18)},
			ChainName:        suite.chainName,
		}
		if trontypes.ModuleName == suite.chainName {
			msgBondedOracle.ExternalAddress = tronAddress.PubkeyToAddress(suite.externals[i].PublicKey).String()
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
		expectPower := sdk.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdk.NewInt(int64(i + 1))).Quo(sdk.DefaultPowerReduction)
		require.True(suite.T(), expectPower.Equal(power))
	}

	bridgeToken := helpers.GenerateAddress().String()
	if trontypes.ModuleName == suite.chainName {
		bridgeToken = trontypes.AddressFromHex(bridgeToken)
	}

	for i := 0; i < 6; i++ {
		addBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     1,
			BlockHeight:    1000,
			TokenContract:  bridgeToken,
			Name:           "Test Token",
			Symbol:         "TEST",
			Decimals:       18,
			BridgerAddress: suite.bridgers[i].String(),
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
		BridgerAddress: suite.bridgers[6].String(),
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

	proposalHandler := crosschain.NewChainProposalHandler(suite.app.CrosschainKeeper)

	var newOracleList []string
	for i := 0; i < 7; i++ {
		newOracleList = append(newOracleList, suite.oracles[i].String())
	}
	err = proposalHandler(suite.ctx, &types.UpdateChainOraclesProposal{
		Title:       "proposal 1: try update chain oracle power >= 30%, expect error",
		Description: "",
		Oracles:     newOracleList,
		ChainName:   suite.chainName,
	})
	require.ErrorIs(suite.T(), types.ErrInvalid, err)

	expectTotalPower := sdk.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdk.NewInt(10)).Quo(sdk.DefaultPowerReduction)
	actualTotalPower := suite.Keeper().GetLastTotalPower(suite.ctx)
	require.True(suite.T(), expectTotalPower.Equal(actualTotalPower))

	expectMaxChangePower := types.AttestationProposalOracleChangePowerThreshold.Mul(expectTotalPower).Quo(sdk.NewInt(100))

	expectDeletePower := sdk.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdk.NewInt(3)).Quo(sdk.DefaultPowerReduction)
	require.EqualValues(suite.T(), fmt.Sprintf("max change power, maxChangePowerThreshold: %s, deleteTotalPower: %s: %s", expectMaxChangePower.String(), expectDeletePower.String(), types.ErrInvalid), err.Error())

	var newOracleList2 []string
	for i := 0; i < 8; i++ {
		newOracleList2 = append(newOracleList2, suite.oracles[i].String())
	}
	err = proposalHandler(suite.ctx, &types.UpdateChainOraclesProposal{
		Title:       "proposal 2: try update chain oracle power <= 30%, expect success",
		Description: "",
		Oracles:     newOracleList2,
		ChainName:   suite.chainName,
	})
	require.NoError(suite.T(), err)
}

func (suite *KeeperTestSuite) TestAttestationAfterOracleUpdate() {
	if len(suite.bridgers) < 20 {
		return
	}
	for i := 0; i < 20; i++ {
		msgBondedOracle := &types.MsgBondedOracle{
			OracleAddress:    suite.oracles[i].String(),
			BridgerAddress:   suite.bridgers[i].String(),
			ExternalAddress:  crypto.PubkeyToAddress(suite.externals[i].PublicKey).Hex(),
			ValidatorAddress: suite.validator[i].String(),
			DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(10 * 1e3).MulRaw(1e18)},
			ChainName:        suite.chainName,
		}
		if trontypes.ModuleName == suite.chainName {
			msgBondedOracle.ExternalAddress = tronAddress.PubkeyToAddress(suite.externals[i].PublicKey).String()
		}
		_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), msgBondedOracle)
		require.NoError(suite.T(), err)
		suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
		suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
		oracleSets := suite.Keeper().GetOracleSets(suite.ctx)
		require.NotNil(suite.T(), oracleSets)
		require.EqualValues(suite.T(), i+1, len(oracleSets))

		power := suite.Keeper().GetLastTotalPower(suite.ctx)
		expectPower := sdk.NewInt(10 * 1e3).MulRaw(1e18).Mul(sdk.NewInt(int64(i + 1))).Quo(sdk.DefaultPowerReduction)
		require.True(suite.T(), expectPower.Equal(power))
	}

	bridgeToken := helpers.GenerateAddress().String()
	if trontypes.ModuleName == suite.chainName {
		bridgeToken = trontypes.AddressFromHex(bridgeToken)
	}

	{
		firstBridgeTokenClaim := &types.MsgBridgeTokenClaim{
			EventNonce:     1,
			BlockHeight:    1000,
			TokenContract:  bridgeToken,
			Name:           "Test Token",
			Symbol:         "TEST",
			Decimals:       18,
			BridgerAddress: "",
			ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
			ChainName:      suite.chainName,
		}

		for i := 0; i < 13; i++ {
			firstBridgeTokenClaim.BridgerAddress = suite.bridgers[i].String()
			_, err := suite.MsgServer().BridgeTokenClaim(sdk.WrapSDKContext(suite.ctx), firstBridgeTokenClaim)
			require.NoError(suite.T(), err)
			endBlockBeforeAttestation := suite.Keeper().GetAttestation(suite.ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
			require.NotNil(suite.T(), endBlockBeforeAttestation)
			require.False(suite.T(), endBlockBeforeAttestation.Observed)
			require.NotNil(suite.T(), endBlockBeforeAttestation.Votes)
			require.EqualValues(suite.T(), i+1, len(endBlockBeforeAttestation.Votes))

			suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
			suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
			endBlockAfterAttestation := suite.Keeper().GetAttestation(suite.ctx, firstBridgeTokenClaim.EventNonce, firstBridgeTokenClaim.ClaimHash())
			require.NotNil(suite.T(), endBlockAfterAttestation)
			require.False(suite.T(), endBlockAfterAttestation.Observed)
		}

		firstBridgeTokenClaim.BridgerAddress = suite.bridgers[13].String()
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
			TokenContract:  bridgeToken,
			Name:           "Test Token2",
			Symbol:         "TEST2",
			Decimals:       18,
			BridgerAddress: "",
			ChannelIbc:     hex.EncodeToString([]byte("transfer/channel-0")),
			ChainName:      suite.chainName,
		}

		for i := 0; i < 6; i++ {
			secondBridgeTokenClaim.BridgerAddress = suite.bridgers[i].String()
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
			newOracleList = append(newOracleList, suite.oracles[i].String())
		}
		err := suite.Keeper().UpdateChainOraclesProposal(suite.ctx, &types.UpdateChainOraclesProposal{
			Title:       "proposal 1: try update chain oracle save top 15 oracle, expect success",
			Description: "",
			Oracles:     newOracleList,
			ChainName:   suite.chainName,
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
			newOracleList2 = append(newOracleList2, suite.oracles[i].String())
		}
		err = suite.Keeper().UpdateChainOraclesProposal(suite.ctx, &types.UpdateChainOraclesProposal{
			Title:       "proposal 2: try update chain oracle save top 11 oracle, expect success",
			Description: "",
			Oracles:     newOracleList2,
			ChainName:   suite.chainName,
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
			newOracleList3 = append(newOracleList3, suite.oracles[i].String())
		}
		err = suite.Keeper().UpdateChainOraclesProposal(suite.ctx, &types.UpdateChainOraclesProposal{
			Title:       "proposal 3: try update chain oracle save top 10 oracle, expect success",
			Description: "",
			Oracles:     newOracleList3,
			ChainName:   suite.chainName,
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

		secondBridgeTokenClaim.BridgerAddress = suite.bridgers[6].String()
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
	for i := 0; i < len(suite.oracles); i++ {
		msgBondedOracle := &types.MsgBondedOracle{
			OracleAddress:    suite.oracles[i].String(),
			BridgerAddress:   suite.bridgers[i].String(),
			ExternalAddress:  crypto.PubkeyToAddress(suite.externals[i].PublicKey).Hex(),
			ValidatorAddress: suite.validator[i].String(),
			DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(10 * 1e3).MulRaw(1e18)},
			ChainName:        suite.chainName,
		}
		if trontypes.ModuleName == suite.chainName {
			msgBondedOracle.ExternalAddress = tronAddress.PubkeyToAddress(suite.externals[i].PublicKey).String()
		}
		require.NoError(suite.T(), msgBondedOracle.ValidateBasic())
		_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), msgBondedOracle)
		require.NoError(suite.T(), err)
	}
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	allOracles := suite.Keeper().GetAllOracles(suite.ctx, false)
	require.NotNil(suite.T(), allOracles)
	require.EqualValues(suite.T(), len(suite.oracles), len(allOracles))

	oracle := suite.oracles[0]
	bridger := suite.bridgers[0]
	externalAddress := crypto.PubkeyToAddress(suite.externals[0].PublicKey).Hex()
	if trontypes.ModuleName == suite.chainName {
		externalAddress = tronAddress.PubkeyToAddress(suite.externals[0].PublicKey).String()
	}

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

	require.True(suite.T(), sdk.NewInt(10*1e3).MulRaw(1e18).Equal(oracleData.DelegateAmount))

	var newOracleAddressList []string
	for _, address := range suite.oracles[1:] {
		newOracleAddressList = append(newOracleAddressList, address.String())
	}

	err := suite.Keeper().UpdateChainOraclesProposal(suite.ctx, &types.UpdateChainOraclesProposal{
		Title:       "proposal 1: try update chain oracle remove first oracle, expect success",
		Description: "",
		Oracles:     newOracleAddressList,
		ChainName:   suite.chainName,
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
	for i := 0; i < len(suite.oracles); i++ {
		msgBondedOracle := &types.MsgBondedOracle{
			OracleAddress:    suite.oracles[i].String(),
			BridgerAddress:   suite.bridgers[i].String(),
			ExternalAddress:  crypto.PubkeyToAddress(suite.externals[i].PublicKey).Hex(),
			ValidatorAddress: suite.validator[i].String(),
			DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(10 * 1e3).MulRaw(1e18)},
			ChainName:        suite.chainName,
		}
		if trontypes.ModuleName == suite.chainName {
			msgBondedOracle.ExternalAddress = tronAddress.PubkeyToAddress(suite.externals[i].PublicKey).String()
		}
		require.NoError(suite.T(), msgBondedOracle.ValidateBasic())
		_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), msgBondedOracle)
		require.NoError(suite.T(), err)
	}
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.Keeper().EndBlocker(suite.ctx)

	allOracles := suite.Keeper().GetAllOracles(suite.ctx, false)
	require.NotNil(suite.T(), allOracles)
	require.Equal(suite.T(), len(suite.oracles), len(allOracles))

	oracleSets := suite.Keeper().GetOracleSets(suite.ctx)
	require.NotNil(suite.T(), oracleSets)
	require.EqualValues(suite.T(), 1, len(oracleSets))

	gravityId := suite.Keeper().GetGravityID(suite.ctx)
	for i := 0; i < len(suite.oracles)-1; i++ {
		externalAddress := crypto.PubkeyToAddress(suite.externals[i].PublicKey).Hex()

		checkpoint, err := oracleSets[0].GetCheckpoint(gravityId)
		require.NoError(suite.T(), err)
		signature, err := types.NewEthereumSignature(checkpoint, suite.externals[i])
		require.NoError(suite.T(), err)

		if trontypes.ModuleName == suite.chainName {
			externalAddress = tronAddress.PubkeyToAddress(suite.externals[i].PublicKey).String()

			checkpoint, err = trontypes.GetCheckpointOracleSet(oracleSets[0], gravityId)
			require.NoError(suite.T(), err)

			signature, err = trontypes.NewTronSignature(checkpoint, suite.externals[i])
			require.NoError(suite.T(), err)
		}
		oracleSetConfirm := &types.MsgOracleSetConfirm{
			Nonce:           oracleSets[0].Nonce,
			BridgerAddress:  suite.bridgers[i].String(),
			ExternalAddress: externalAddress,
			Signature:       hex.EncodeToString(signature),
			ChainName:       suite.chainName,
		}
		require.NoError(suite.T(), oracleSetConfirm.ValidateBasic())
		_, err = suite.MsgServer().OracleSetConfirm(sdk.WrapSDKContext(suite.ctx), oracleSetConfirm)
		require.NoError(suite.T(), err)
	}

	suite.Keeper().EndBlocker(suite.ctx)
	oracleSetHeight := int64(oracleSets[0].Height)
	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 1)
	suite.app.EndBlock(abci.RequestEndBlock{Height: suite.ctx.BlockHeight()})

	oracle, found := suite.Keeper().GetOracle(suite.ctx, suite.oracles[len(suite.oracles)-1])
	require.True(suite.T(), found)
	require.True(suite.T(), oracle.Online)
	require.Equal(suite.T(), int64(0), oracle.SlashTimes)

	suite.ctx = suite.ctx.WithBlockHeight(oracleSetHeight + int64(suite.Keeper().GetParams(suite.ctx).SignedWindow) + 1)
	suite.Keeper().EndBlocker(suite.ctx)

	oracle, found = suite.Keeper().GetOracle(suite.ctx, suite.oracles[len(suite.oracles)-1])
	require.True(suite.T(), found)
	require.False(suite.T(), oracle.Online)
	require.Equal(suite.T(), int64(1), oracle.SlashTimes)
}

func (suite *KeeperTestSuite) TestSlashOracle() {
	for i := 0; i < len(suite.oracles); i++ {
		msgBondedOracle := &types.MsgBondedOracle{
			OracleAddress:    suite.oracles[i].String(),
			BridgerAddress:   suite.bridgers[i].String(),
			ExternalAddress:  crypto.PubkeyToAddress(suite.externals[i].PublicKey).Hex(),
			ValidatorAddress: suite.validator[i].String(),
			DelegateAmount:   sdk.Coin{Denom: fxtypes.DefaultDenom, Amount: sdk.NewInt(10 * 1e3).MulRaw(1e18)},
			ChainName:        suite.chainName,
		}
		if trontypes.ModuleName == suite.chainName {
			msgBondedOracle.ExternalAddress = tronAddress.PubkeyToAddress(suite.externals[i].PublicKey).String()
		}
		require.NoError(suite.T(), msgBondedOracle.ValidateBasic())
		_, err := suite.MsgServer().BondedOracle(sdk.WrapSDKContext(suite.ctx), msgBondedOracle)
		require.NoError(suite.T(), err)
	}

	params := suite.Keeper().GetParams(suite.ctx)
	suite.Keeper().SetParams(suite.ctx, &params)

	for i := 0; i < len(suite.oracles); i++ {
		oracle, found := suite.Keeper().GetOracle(suite.ctx, suite.oracles[i])
		require.True(suite.T(), found)
		require.True(suite.T(), oracle.Online)
		require.Equal(suite.T(), int64(0), oracle.SlashTimes)

		suite.Keeper().SlashOracle(suite.ctx, oracle.OracleAddress)

		oracle, found = suite.Keeper().GetOracle(suite.ctx, suite.oracles[i])
		require.True(suite.T(), found)
		require.False(suite.T(), oracle.Online)
		require.Equal(suite.T(), int64(1), oracle.SlashTimes)
	}

	// repeat slash test.
	for i := 0; i < len(suite.oracles); i++ {
		oracle, found := suite.Keeper().GetOracle(suite.ctx, suite.oracles[i])
		require.True(suite.T(), found)
		require.False(suite.T(), oracle.Online)
		require.Equal(suite.T(), int64(1), oracle.SlashTimes)

		suite.Keeper().SlashOracle(suite.ctx, oracle.OracleAddress)

		oracle, found = suite.Keeper().GetOracle(suite.ctx, suite.oracles[i])
		require.True(suite.T(), found)
		require.False(suite.T(), oracle.Online)
		require.Equal(suite.T(), int64(1), oracle.SlashTimes)
	}
}