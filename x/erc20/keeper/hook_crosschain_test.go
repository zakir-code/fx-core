package keeper_test

import (
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"testing"

	polygontypes "github.com/functionx/fx-core/v2/x/polygon/types"

	"github.com/functionx/fx-core/v2/app/helpers"

	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ibcclienttypes "github.com/cosmos/ibc-go/v3/modules/core/02-client/types"
	"github.com/cosmos/ibc-go/v3/modules/core/exported"
	ibctmtypes "github.com/cosmos/ibc-go/v3/modules/light-clients/07-tendermint/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/evmos/ethermint/crypto/ethsecp256k1"
	"github.com/stretchr/testify/require"

	"github.com/functionx/fx-core/v2/app"
	fxtypes "github.com/functionx/fx-core/v2/types"
	bsctypes "github.com/functionx/fx-core/v2/x/bsc/types"
	"github.com/functionx/fx-core/v2/x/crosschain"
	crosschainkeeper "github.com/functionx/fx-core/v2/x/crosschain/keeper"
	crosschaintypes "github.com/functionx/fx-core/v2/x/crosschain/types"
	"github.com/functionx/fx-core/v2/x/erc20/types"
)

func (suite *KeeperTestSuite) TestHookChainBSC() {
	suite.purseBalance = sdk.NewInt(100000).Mul(sdk.NewInt(1e18))
	suite.SetupTest()

	signer1, addr1 := privateSigner()
	_, addr2 := privateSigner()

	suite.ctx = testInitBscCrossChain(suite.T(), suite.ctx, suite.app, suite.address.Bytes(), addr1.Bytes(), addr2)

	purseID := suite.app.Erc20Keeper.GetDenomMap(suite.ctx, PurseDenom)
	suite.Require().NotEmpty(purseID)

	tokenPair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, purseID)
	suite.Require().True(found)
	suite.Require().NotNil(tokenPair)
	suite.Require().NotEmpty(tokenPair.GetERC20Contract())

	require.Equal(suite.T(), types.TokenPair{
		Erc20Address:  tokenPair.GetErc20Address(),
		Denom:         PurseDenom,
		Enabled:       true,
		ContractOwner: types.OWNER_MODULE,
	}, tokenPair)

	fip20, err := suite.app.Erc20Keeper.QueryERC20(suite.ctx, tokenPair.GetERC20Contract())
	suite.Require().NoError(err)
	suite.Require().Equal("PURSE", fip20.Symbol)

	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100))
	err = suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), addr1.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amt), sdk.NewCoin(PurseDenom, amt)))
	suite.Require().NoError(err)

	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1.Bytes())
	_ = balances

	err = suite.app.Erc20Keeper.RelayConvertCoin(suite.ctx, addr1.Bytes(), addr1, sdk.NewCoin(PurseDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))))
	suite.Require().NoError(err)

	balanceOf, err := suite.app.Erc20Keeper.BalanceOf(suite.ctx, tokenPair.GetERC20Contract(), addr1)
	suite.Require().NoError(err)
	_ = balanceOf

	token := tokenPair.GetERC20Contract()
	crossChainTarget := fmt.Sprintf("%s%s", fxtypes.FIP20TransferToChainPrefix, bsctypes.ModuleName)
	transferChainData := packTransferCrossData(suite.T(), addr2.String(), big.NewInt(1e18), big.NewInt(1e18), crossChainTarget)
	sendEthTx(suite.T(), suite.ctx, suite.app, signer1, addr1, token, transferChainData)

	transactions := suite.app.BscKeeper.GetUnbatchedTransactions(suite.ctx)
	require.Equal(suite.T(), 1, len(transactions))
	require.Equal(suite.T(), transactions[0].DestAddress, addr2.String())
	require.Equal(suite.T(), transactions[0].Token.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transactions[0].Fee.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transactions[0].Sender, sdk.AccAddress(addr1.Bytes()).String())
}

func (suite *KeeperTestSuite) TestHookChainUSDT() {
	suite.supportManyToOneBlock = true
	suite.bscUSDTBalance = sdk.NewInt(100000).Mul(sdk.NewInt(1e18))
	suite.polygonUSDTBalance = sdk.NewInt(100000).Mul(sdk.NewInt(1e18))
	suite.SetupTest()

	signer1, addr1 := privateSigner()
	_, addr2 := privateSigner()
	_ = signer1

	signer3, addr3 := privateSigner()
	_, addr4 := privateSigner()
	_ = signer3

	suite.ctx = testInitBscCrossChain(suite.T(), suite.ctx, suite.app, suite.address.Bytes(), addr1.Bytes(), addr2)
	suite.ctx = testInitPolygonCrossChain(suite.T(), suite.ctx, suite.app, suite.address.Bytes(), addr3.Bytes(), addr4)

	usdtCopy := usdtMatedata
	usdtCopy.DenomUnits[0].Aliases = []string{bscUSDT, polygonUSDT}

	tokenPair, err := suite.app.Erc20Keeper.RegisterCoin(suite.ctx, usdtCopy)
	suite.Require().NoError(err)
	require.Equal(suite.T(), types.TokenPair{
		Erc20Address:  tokenPair.GetErc20Address(),
		Denom:         "usdt",
		Enabled:       true,
		ContractOwner: types.OWNER_MODULE,
	}, *tokenPair)

	denomBytes := suite.app.Erc20Keeper.GetAliasDenom(suite.ctx, bscUSDT)
	require.Equal(suite.T(), usdtCopy.Base, string(denomBytes))

	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100))
	err = suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), addr1.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amt), sdk.NewCoin(bscUSDT, amt)))
	suite.Require().NoError(err)

	err = suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), addr3.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amt), sdk.NewCoin(polygonUSDT, amt)))
	suite.Require().NoError(err)

	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1.Bytes())
	_ = balances

	balances = suite.app.BankKeeper.GetAllBalances(suite.ctx, addr3.Bytes())
	_ = balances

	usdtCoin, err := suite.app.Erc20Keeper.RelayConvertDenom(suite.ctx, addr1.Bytes(), sdk.NewCoin(bscUSDT, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))))
	suite.Require().NoError(err)
	err = suite.app.Erc20Keeper.RelayConvertCoin(suite.ctx, addr1.Bytes(), addr1, usdtCoin)
	suite.Require().NoError(err)
	balanceOf, err := suite.app.Erc20Keeper.BalanceOf(suite.ctx, tokenPair.GetERC20Contract(), addr1)
	suite.Require().NoError(err)
	_ = balanceOf

	usdtCoin, err = suite.app.Erc20Keeper.RelayConvertDenom(suite.ctx, addr3.Bytes(), sdk.NewCoin(polygonUSDT, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))))
	suite.Require().NoError(err)
	err = suite.app.Erc20Keeper.RelayConvertCoin(suite.ctx, addr3.Bytes(), addr3, usdtCoin)
	suite.Require().NoError(err)
	balanceOf, err = suite.app.Erc20Keeper.BalanceOf(suite.ctx, tokenPair.GetERC20Contract(), addr3)
	suite.Require().NoError(err)
	_ = balanceOf

	// addr1 transfer usdt to bsc
	token := tokenPair.GetERC20Contract()
	crossChainTarget := fmt.Sprintf("%s%s", fxtypes.FIP20TransferToChainPrefix, bsctypes.ModuleName)
	transferChainData := packTransferCrossData(suite.T(), addr2.String(), big.NewInt(1e18), big.NewInt(1e18), crossChainTarget)
	sendEthTx(suite.T(), suite.ctx, suite.app, signer1, addr1, token, transferChainData)

	transactions := suite.app.BscKeeper.GetUnbatchedTransactions(suite.ctx)
	require.Equal(suite.T(), 1, len(transactions))
	require.Equal(suite.T(), transactions[0].DestAddress, addr2.String())
	require.Equal(suite.T(), transactions[0].Token.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transactions[0].Fee.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transactions[0].Sender, sdk.AccAddress(addr1.Bytes()).String())

	// addr1 transfer usdt to polygon
	crossChainTarget = fmt.Sprintf("%s%s", fxtypes.FIP20TransferToChainPrefix, polygontypes.ModuleName)
	transferChainData = packTransferCrossData(suite.T(), addr2.String(), big.NewInt(1e18), big.NewInt(1e18), crossChainTarget)
	sendEthTx(suite.T(), suite.ctx, suite.app, signer1, addr1, token, transferChainData)

	transactions = suite.app.PolygonKeeper.GetUnbatchedTransactions(suite.ctx)
	require.Equal(suite.T(), 1, len(transactions))
	require.Equal(suite.T(), transactions[0].DestAddress, addr2.String())
	require.Equal(suite.T(), transactions[0].Token.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transactions[0].Fee.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transactions[0].Sender, sdk.AccAddress(addr1.Bytes()).String())

	// addr3 transfer usdt to bsc
	crossChainTarget = fmt.Sprintf("%s%s", fxtypes.FIP20TransferToChainPrefix, bsctypes.ModuleName)
	transferChainData = packTransferCrossData(suite.T(), addr4.String(), big.NewInt(1e18), big.NewInt(1e18), crossChainTarget)
	sendEthTx(suite.T(), suite.ctx, suite.app, signer3, addr3, token, transferChainData)

	transactions = suite.app.BscKeeper.GetUnbatchedTransactions(suite.ctx)
	require.Equal(suite.T(), 2, len(transactions))
	transaction := getOutgoingTransferTxById(transactions, 2)
	require.NotNil(suite.T(), transaction)
	require.Equal(suite.T(), transaction.DestAddress, addr4.String())
	require.Equal(suite.T(), transaction.Token.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transaction.Fee.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transaction.Sender, sdk.AccAddress(addr3.Bytes()).String())

	// addr2 transfer usdt to polygon
	crossChainTarget = fmt.Sprintf("%s%s", fxtypes.FIP20TransferToChainPrefix, polygontypes.ModuleName)
	transferChainData = packTransferCrossData(suite.T(), addr4.String(), big.NewInt(1e18), big.NewInt(1e18), crossChainTarget)
	sendEthTx(suite.T(), suite.ctx, suite.app, signer3, addr3, token, transferChainData)

	transactions = suite.app.PolygonKeeper.GetUnbatchedTransactions(suite.ctx)
	require.Equal(suite.T(), 2, len(transactions))
	transaction = getOutgoingTransferTxById(transactions, 2)
	require.NotNil(suite.T(), transaction)
	require.Equal(suite.T(), transaction.DestAddress, addr4.String())
	require.Equal(suite.T(), transaction.Token.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transaction.Fee.Amount.BigInt(), big.NewInt(1e18))
	require.Equal(suite.T(), transaction.Sender, sdk.AccAddress(addr3.Bytes()).String())

	suite.supportManyToOneBlock = false
}

type IBCTransferSimulate struct {
	T *testing.T
}

func (it *IBCTransferSimulate) SendTransfer(ctx sdk.Context, sourcePort, sourceChannel string, token sdk.Coin, sender sdk.AccAddress,
	receiver string, timeoutHeight ibcclienttypes.Height, timeoutTimestamp uint64, router string, fee sdk.Coin) error {
	require.Equal(it.T, token.Amount.BigInt(), ibcTransferAmount)
	return nil
}

type IBCChannelSimulate struct {
}

func (ic *IBCChannelSimulate) GetChannelClientState(ctx sdk.Context, portID, channelID string) (string, exported.ClientState, error) {
	return "", &ibctmtypes.ClientState{
		ChainId:         "fxcore",
		TrustLevel:      ibctmtypes.Fraction{},
		TrustingPeriod:  0,
		UnbondingPeriod: 0,
		MaxClockDrift:   0,
		FrozenHeight: ibcclienttypes.Height{
			RevisionHeight: 1000,
			RevisionNumber: 1000,
		},
		LatestHeight: ibcclienttypes.Height{
			RevisionHeight: 10,
			RevisionNumber: 10,
		},
		ProofSpecs:                   nil,
		UpgradePath:                  nil,
		AllowUpdateAfterExpiry:       false,
		AllowUpdateAfterMisbehaviour: false,
	}, nil
}
func (ic *IBCChannelSimulate) GetNextSequenceSend(ctx sdk.Context, portID, channelID string) (uint64, bool) {
	return 1, true
}

var (
	ibcTransferAmount = big.NewInt(1e18)
)

func (suite *KeeperTestSuite) TestHookIBC() {
	suite.SetupTest()

	pairId := suite.app.Erc20Keeper.GetDenomMap(suite.ctx, "FX")
	suite.Require().Greater(len(pairId), 0)

	pair, found := suite.app.Erc20Keeper.GetTokenPair(suite.ctx, pairId)
	suite.Require().True(found)

	signer1, addr1 := privateSigner()
	amt := sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(100))
	err := suite.app.BankKeeper.SendCoins(suite.ctx, suite.address.Bytes(), addr1.Bytes(), sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, amt)))
	suite.Require().NoError(err)

	balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, addr1.Bytes())
	suite.Require().False(balances.IsZero())

	err = suite.app.Erc20Keeper.RelayConvertCoin(suite.ctx, addr1.Bytes(), addr1, sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromUint64(1e18).Mul(sdk.NewInt(10))))
	suite.Require().NoError(err)

	balanceOf, err := suite.app.Erc20Keeper.BalanceOf(suite.ctx, pair.GetERC20Contract(), addr1)
	suite.Require().NoError(err)
	suite.Require().Equal(balanceOf.Cmp(big.NewInt(0)), 1)

	//reset ibc
	suite.app.Erc20Keeper.SetIBCTransferKeeperForTest(&IBCTransferSimulate{T: suite.T()})
	suite.app.Erc20Keeper.SetIBCChannelKeeperForTest(&IBCChannelSimulate{})

	token := pair.GetERC20Contract()
	ibcTarget := fmt.Sprintf("%s%s", fxtypes.FIP20TransferToIBCPrefix, "px/transfer/channel-0")
	transferIBCData := packTransferCrossData(suite.T(), "px16u6kjunrcxkvaln9aetxwjpruply3sgwpr9z8u", ibcTransferAmount, big.NewInt(0), ibcTarget)
	sendEthTx(suite.T(), suite.ctx, suite.app, signer1, addr1, token, transferIBCData)
}

func packTransferCrossData(t *testing.T, to string, amount, fee *big.Int, target string) []byte {
	fip20 := fxtypes.GetERC20()
	targetBytes := fxtypes.StringToByte32(target)
	pack, err := fip20.ABI.Pack("transferCrossChain", to, amount, fee, targetBytes)
	require.NoError(t, err)
	return pack
}

func privateSigner() (keyring.Signer, common.Address) {
	// account key
	priKey := NewPriKey()
	//ethsecp256k1.GenerateKey()
	ethPriv := &ethsecp256k1.PrivKey{Key: priKey.Bytes()}

	return helpers.NewSigner(ethPriv), common.BytesToAddress(ethPriv.PubKey().Address())
}

var (
	BSCBridgeTokenContract  = common.HexToAddress("0x29a63F4B209C29B4DC47f06FFA896F32667DAD2C")
	BSCUSDTokenContract     = common.HexToAddress("0x0000000000000000000000000000000000000000")
	PolygonUSDTokenContract = common.HexToAddress("0x0000000000000000000000000000000000000001")
	bscUSDT                 = fmt.Sprintf("bsc%s", BSCUSDTokenContract.String())
	polygonUSDT             = fmt.Sprintf("polygon%s", PolygonUSDTokenContract.String())
)

func testInitBscCrossChain(t *testing.T, ctx sdk.Context, myApp *app.App, oracleAddress, bridgeAddress sdk.AccAddress, externalAddress common.Address) sdk.Context {
	deposit := sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromBigInt(big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18))))
	err := myApp.BankKeeper.SendCoinsFromAccountToModule(ctx, oracleAddress, bsctypes.ModuleName, sdk.NewCoins(deposit))
	require.NoError(t, err)

	testCrossChainParamsProposal(t, ctx, myApp.BscKeeper, oracleAddress, bsctypes.ModuleName)

	oracle := crosschaintypes.Oracle{
		OracleAddress:   oracleAddress.String(),
		BridgerAddress:  bridgeAddress.String(),
		ExternalAddress: externalAddress.String(),
		DelegateAmount:  deposit.Amount,
		StartHeight:     ctx.BlockHeight(),
		Online:          true,
		SlashTimes:      0,
	}
	// save oracle
	myApp.BscKeeper.SetOracle(ctx, oracle)

	myApp.BscKeeper.SetOracleByBridger(ctx, bridgeAddress, oracleAddress)
	// set the ethereum address
	myApp.BscKeeper.SetOracleByExternalAddress(ctx, externalAddress.String(), oracleAddress)

	myApp.BscKeeper.CommonSetOracleTotalPower(ctx)

	testCrossChainOracleSetUpdateClaim(t, ctx, myApp.BscKeeper, bridgeAddress, externalAddress, 1, bsctypes.ModuleName)

	testCrossChainBridgeTokenClaim(t, ctx, myApp.BscKeeper, bridgeAddress, 2,
		BSCBridgeTokenContract, "PURSE Token", "PURSE", 18, bsctypes.ModuleName, hex.EncodeToString([]byte("transfer/channel-0")))
	testCrossChainBridgeTokenClaim(t, ctx, myApp.BscKeeper, bridgeAddress, 3,
		BSCUSDTokenContract, "USDT Token", "USDT", 6, bsctypes.ModuleName, "")

	crosschain.EndBlocker(ctx, myApp.BscKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	return ctx
}

func testInitPolygonCrossChain(t *testing.T, ctx sdk.Context, myApp *app.App, oracleAddress, bridgeAddress sdk.AccAddress, externalAddress common.Address) sdk.Context {
	deposit := sdk.NewCoin(fxtypes.DefaultDenom, sdk.NewIntFromBigInt(big.NewInt(0).Mul(big.NewInt(10), big.NewInt(1e18))))
	err := myApp.BankKeeper.SendCoinsFromAccountToModule(ctx, oracleAddress, polygontypes.ModuleName, sdk.NewCoins(deposit))
	require.NoError(t, err)

	testCrossChainParamsProposal(t, ctx, myApp.PolygonKeeper, oracleAddress, polygontypes.ModuleName)

	oracle := crosschaintypes.Oracle{
		OracleAddress:   oracleAddress.String(),
		BridgerAddress:  bridgeAddress.String(),
		ExternalAddress: externalAddress.String(),
		DelegateAmount:  deposit.Amount,
		StartHeight:     ctx.BlockHeight(),
		Online:          true,
		SlashTimes:      0,
	}
	// save oracle
	myApp.PolygonKeeper.SetOracle(ctx, oracle)

	myApp.PolygonKeeper.SetOracleByBridger(ctx, bridgeAddress, oracleAddress)
	// set the ethereum address
	myApp.PolygonKeeper.SetOracleByExternalAddress(ctx, externalAddress.String(), oracleAddress)

	myApp.PolygonKeeper.CommonSetOracleTotalPower(ctx)

	testCrossChainOracleSetUpdateClaim(t, ctx, myApp.PolygonKeeper, bridgeAddress, externalAddress, 1, polygontypes.ModuleName)

	testCrossChainBridgeTokenClaim(t, ctx, myApp.PolygonKeeper, bridgeAddress, 2,
		PolygonUSDTokenContract, "USDT Token", "USDT", 6, polygontypes.ModuleName, "")

	crosschain.EndBlocker(ctx, myApp.PolygonKeeper)

	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + 1)

	return ctx
}

func testCrossChainParamsProposal(t *testing.T, ctx sdk.Context, cck crosschainkeeper.Keeper, oracles sdk.AccAddress, chain string) {
	proposal := &crosschaintypes.UpdateChainOraclesProposal{
		Title:       fmt.Sprintf("%s cross chain", chain),
		Description: fmt.Sprintf("%s cross chain oracles init", chain),
		Oracles:     []string{oracles.String()},
		ChainName:   chain,
	}

	k := &crosschainkeeper.EthereumMsgServer{Keeper: cck}
	err := crosschain.HandleUpdateChainOraclesProposal(ctx, k, proposal)
	require.NoError(t, err)
}

func testCrossChainBridgeTokenClaim(t *testing.T, ctx sdk.Context, cck crosschainkeeper.Keeper,
	orchAddr sdk.AccAddress, eventNonce uint64, contract common.Address,
	name, symbol string, decimals uint64, chain, channelIBC string) {
	msg := &crosschaintypes.MsgBridgeTokenClaim{
		EventNonce:     eventNonce,
		BlockHeight:    uint64(ctx.BlockHeight()),
		TokenContract:  contract.String(),
		Name:           name,
		Symbol:         symbol,
		Decimals:       decimals,
		BridgerAddress: orchAddr.String(),
		ChannelIbc:     channelIBC,
		ChainName:      chain,
	}

	any, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)

	// Add the claim to the store
	_, err = cck.Attest(ctx, msg, any)
	require.NoError(t, err)
}

func testCrossChainOracleSetUpdateClaim(t *testing.T, ctx sdk.Context, cck crosschainkeeper.Keeper,
	orch sdk.AccAddress, addr common.Address, eventNonce uint64, chain string) {
	msg := &crosschaintypes.MsgOracleSetUpdatedClaim{
		EventNonce:     eventNonce,
		BlockHeight:    uint64(ctx.BlockHeight()),
		OracleSetNonce: 0,
		Members: crosschaintypes.BridgeValidators{
			{
				Power:           uint64(math.MaxUint32),
				ExternalAddress: addr.String(),
			},
		},
		BridgerAddress: orch.String(),
		ChainName:      chain,
	}
	for _, member := range msg.Members {
		_, found := cck.GetOracleByExternalAddress(ctx, member.ExternalAddress)
		require.True(t, found)
	}

	any, err := codectypes.NewAnyWithValue(msg)
	require.NoError(t, err)

	// Add the claim to the store
	_, err = cck.Attest(ctx, msg, any)
	require.NoError(t, err)
}

func getOutgoingTransferTxById(txs []*crosschaintypes.OutgoingTransferTx, id uint64) *crosschaintypes.OutgoingTransferTx {
	for _, tx := range txs {
		if tx.Id == id {
			return tx
		}
	}
	return nil
}
