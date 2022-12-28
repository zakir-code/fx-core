package tests

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"strconv"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gethcommon "github.com/ethereum/go-ethereum/common"
	ethcrypto "github.com/ethereum/go-ethereum/crypto"
	tronAddress "github.com/fbsobreira/gotron-sdk/pkg/address"

	"github.com/functionx/fx-core/v3/app/helpers"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	trontypes "github.com/functionx/fx-core/v3/x/tron/types"
)

type CrosschainTestSuite struct {
	*TestSuite
	params          crosschaintypes.Params
	chainName       string
	oraclePrivKey   cryptotypes.PrivKey
	bridgerPrivKey  cryptotypes.PrivKey
	externalPrivKey *ecdsa.PrivateKey
	privKey         cryptotypes.PrivKey
}

func NewCrosschainWithTestSuite(chainName string, ts *TestSuite) CrosschainTestSuite {
	key, err := ethcrypto.GenerateKey()
	if err != nil {
		panic(err.Error())
	}
	return CrosschainTestSuite{
		TestSuite:       ts,
		chainName:       chainName,
		oraclePrivKey:   helpers.NewPriKey(),
		bridgerPrivKey:  helpers.NewPriKey(),
		externalPrivKey: key,
		privKey:         helpers.NewEthPrivKey(),
	}
}

func (suite *CrosschainTestSuite) Init() {
	suite.Send(suite.OracleAddr(), suite.NewCoin(sdk.NewInt(10_100).MulRaw(1e18)))
	suite.Send(suite.BridgerAddr(), suite.NewCoin(sdk.NewInt(1_000).MulRaw(1e18)))
	suite.Send(suite.AccAddress(), suite.NewCoin(sdk.NewInt(1_000).MulRaw(1e18)))
	suite.params = suite.QueryParams()
}

func (suite *CrosschainTestSuite) OracleAddr() sdk.AccAddress {
	return suite.oraclePrivKey.PubKey().Address().Bytes()
}

func (suite *CrosschainTestSuite) ExternalAddr() string {
	if suite.chainName == trontypes.ModuleName {
		return tronAddress.PubkeyToAddress(suite.externalPrivKey.PublicKey).String()
	}
	return ethcrypto.PubkeyToAddress(suite.externalPrivKey.PublicKey).String()
}

func (suite *CrosschainTestSuite) BridgerAddr() sdk.AccAddress {
	return suite.bridgerPrivKey.PubKey().Address().Bytes()
}

func (suite *CrosschainTestSuite) AccAddress() sdk.AccAddress {
	return suite.privKey.PubKey().Address().Bytes()
}

func (suite *CrosschainTestSuite) HexAddress() gethcommon.Address {
	return gethcommon.BytesToAddress(suite.privKey.PubKey().Address())
}

func (suite *CrosschainTestSuite) CrosschainQuery() crosschaintypes.QueryClient {
	return suite.GRPCClient().CrosschainQuery()
}

func (suite *CrosschainTestSuite) QueryParams() crosschaintypes.Params {
	response, err := suite.CrosschainQuery().Params(suite.ctx,
		&crosschaintypes.QueryParamsRequest{ChainName: suite.chainName})
	suite.NoError(err)
	return response.Params
}

func (suite *CrosschainTestSuite) queryFxLastEventNonce() uint64 {
	lastEventNonce, err := suite.CrosschainQuery().LastEventNonceByAddr(suite.ctx,
		&crosschaintypes.QueryLastEventNonceByAddrRequest{
			ChainName:      suite.chainName,
			BridgerAddress: suite.BridgerAddr().String(),
		},
	)
	suite.NoError(err)
	return lastEventNonce.EventNonce + 1
}

func (suite *CrosschainTestSuite) queryObserverExternalBlockHeight() uint64 {
	response, err := suite.CrosschainQuery().LastObservedBlockHeight(suite.ctx,
		&crosschaintypes.QueryLastObservedBlockHeightRequest{
			ChainName: suite.chainName,
		},
	)
	suite.NoError(err)
	return response.ExternalBlockHeight
}

func (suite *CrosschainTestSuite) AddBridgeTokenClaim(name, symbol string, decimals uint64, token, channelIBC string) string {
	bridgeToken, err := suite.CrosschainQuery().TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
		ChainName: suite.chainName,
		Token:     token,
	})
	suite.ErrorContains(err, "code = NotFound desc = bridge token")
	suite.Nil(bridgeToken)

	suite.BroadcastTx(suite.bridgerPrivKey, &crosschaintypes.MsgBridgeTokenClaim{
		EventNonce:     suite.queryFxLastEventNonce(),
		BlockHeight:    suite.queryObserverExternalBlockHeight() + 1,
		TokenContract:  token,
		Name:           name,
		Symbol:         symbol,
		Decimals:       decimals,
		BridgerAddress: suite.BridgerAddr().String(),
		ChannelIbc:     hex.EncodeToString([]byte(channelIBC)),
		ChainName:      suite.chainName,
	})

	bridgeToken, err = suite.CrosschainQuery().TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
		ChainName: suite.chainName,
		Token:     token,
	})
	suite.NoError(err)
	suite.T().Log("bridge token", bridgeToken)
	return bridgeToken.Denom
}

func (suite *CrosschainTestSuite) GetBridgeDenoms() (denoms []string) {
	response, err := suite.CrosschainQuery().BridgeTokens(suite.ctx, &crosschaintypes.QueryBridgeTokensRequest{
		ChainName: suite.chainName,
	})
	suite.NoError(err)
	for _, token := range response.BridgeTokens {
		denoms = append(denoms, token.Denom)
	}
	return denoms
}

func (suite *CrosschainTestSuite) GetBridgeToken(denom string) string {
	response, err := suite.CrosschainQuery().DenomToToken(suite.ctx, &crosschaintypes.QueryDenomToTokenRequest{
		ChainName: suite.chainName,
		Denom:     denom,
	})
	suite.NoError(err)
	return response.Token

}

func (suite *CrosschainTestSuite) BondedOracle() {
	response, err := suite.CrosschainQuery().GetOracleByBridgerAddr(suite.ctx,
		&crosschaintypes.QueryOracleByBridgerAddrRequest{
			BridgerAddress: suite.BridgerAddr().String(),
			ChainName:      suite.chainName,
		},
	)
	suite.Error(err, crosschaintypes.ErrNoFoundOracle)
	suite.Nil(response)

	suite.BroadcastTx(suite.oraclePrivKey, &crosschaintypes.MsgBondedOracle{
		OracleAddress:    suite.OracleAddr().String(),
		BridgerAddress:   suite.BridgerAddr().String(),
		ExternalAddress:  suite.ExternalAddr(),
		ValidatorAddress: suite.GetFirstValiAddr().String(),
		DelegateAmount:   suite.params.DelegateThreshold,
		ChainName:        suite.chainName,
	})

	response, err = suite.CrosschainQuery().GetOracleByBridgerAddr(suite.ctx,
		&crosschaintypes.QueryOracleByBridgerAddrRequest{
			BridgerAddress: suite.BridgerAddr().String(),
			ChainName:      suite.chainName,
		},
	)
	suite.NoError(err)
	suite.T().Log("oracle", response.Oracle)
}

func (suite *CrosschainTestSuite) SendUpdateChainOraclesProposal() (proposalId uint64) {
	content := &crosschaintypes.UpdateChainOraclesProposal{
		Title:       fmt.Sprintf("Update %s cross chain oracle", suite.chainName),
		Description: "foo",
		Oracles:     []string{suite.OracleAddr().String()},
		ChainName:   suite.chainName,
	}
	return suite.BroadcastProposalTx(content)
}

func (suite *CrosschainTestSuite) SendOracleSetConfirm() {
	queryResponse, err := suite.CrosschainQuery().LastPendingOracleSetRequestByAddr(suite.ctx,
		&crosschaintypes.QueryLastPendingOracleSetRequestByAddrRequest{
			BridgerAddress: suite.BridgerAddr().String(),
			ChainName:      suite.chainName,
		},
	)
	suite.NoError(err)

	for _, oracleSet := range queryResponse.OracleSets {
		var signature []byte
		if suite.chainName == trontypes.ModuleName {
			checkpoint, err := trontypes.GetCheckpointOracleSet(oracleSet, suite.params.GravityId)
			suite.NoError(err)
			signature, err = trontypes.NewTronSignature(checkpoint, suite.externalPrivKey)
			suite.NoError(err)
			err = trontypes.ValidateTronSignature(checkpoint, signature, suite.ExternalAddr())
			suite.NoError(err)
		} else {
			checkpoint, err := oracleSet.GetCheckpoint(suite.params.GravityId)
			suite.NoError(err)
			signature, err = crosschaintypes.NewEthereumSignature(checkpoint, suite.externalPrivKey)
			suite.NoError(err)
			err = crosschaintypes.ValidateEthereumSignature(checkpoint, signature, suite.ExternalAddr())
			suite.NoError(err)
		}

		suite.BroadcastTx(suite.bridgerPrivKey, &crosschaintypes.MsgOracleSetConfirm{
			Nonce:           oracleSet.Nonce,
			BridgerAddress:  suite.BridgerAddr().String(),
			ExternalAddress: suite.ExternalAddr(),
			Signature:       hex.EncodeToString(signature),
			ChainName:       suite.chainName,
		})
	}
}

func (suite *CrosschainTestSuite) SendToFxClaim(token string, amount sdk.Int, targetIbc string) {
	sender := suite.HexAddress().Hex()
	if suite.chainName == trontypes.ModuleName {
		sender = trontypes.AddressFromHex(sender)
	}
	suite.BroadcastTx(suite.bridgerPrivKey, &crosschaintypes.MsgSendToFxClaim{
		EventNonce:     suite.queryFxLastEventNonce(),
		BlockHeight:    suite.queryObserverExternalBlockHeight() + 1,
		TokenContract:  token,
		Amount:         amount,
		Sender:         sender,
		Receiver:       suite.AccAddress().String(),
		TargetIbc:      hex.EncodeToString([]byte(targetIbc)),
		BridgerAddress: suite.BridgerAddr().String(),
		ChainName:      suite.chainName,
	})
	bridgeToken, err := suite.CrosschainQuery().TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
		ChainName: suite.chainName,
		Token:     token,
	})
	suite.NoError(err)
	if len(targetIbc) <= 0 {
		balances := suite.QueryBalances(suite.AccAddress())
		suite.True(balances.IsAllGTE(sdk.NewCoins(sdk.NewCoin(bridgeToken.Denom, amount))))
	}
}

func (suite *CrosschainTestSuite) SendToExternal(count int, amount sdk.Coin) uint64 {
	msgList := make([]sdk.Msg, 0, count)
	for i := 0; i < count; i++ {
		dest := suite.HexAddress().Hex()
		if suite.chainName == trontypes.ModuleName {
			dest = trontypes.AddressFromHex(dest)
		}
		msgList = append(msgList, &crosschaintypes.MsgSendToExternal{
			Sender:    suite.AccAddress().String(),
			Dest:      dest,
			Amount:    amount.SubAmount(sdk.NewInt(1)),
			BridgeFee: sdk.NewCoin(amount.Denom, sdk.NewInt(1)),
			ChainName: suite.chainName,
		})
	}
	txResponse := suite.BroadcastTx(suite.privKey, msgList...)
	for _, eventLog := range txResponse.Logs {
		for _, event := range eventLog.Events {
			if event.Type != crosschaintypes.EventTypeSendToExternal {
				continue
			}
			for _, attribute := range event.Attributes {
				if attribute.Key != crosschaintypes.AttributeKeyOutgoingTxID {
					continue
				}
				txId, err := strconv.ParseUint(attribute.Value, 10, 64)
				suite.NoError(err)
				return txId
			}
		}
	}
	return 0
}

func (suite *CrosschainTestSuite) SendToExternalAndCancel(coin sdk.Coin) {
	txId := suite.SendToExternal(1, coin)
	suite.Greater(txId, uint64(0))

	suite.SendCancelSendToExternal(txId)
}

func (suite *CrosschainTestSuite) SendCancelSendToExternal(txId uint64) {
	suite.BroadcastTx(suite.privKey, &crosschaintypes.MsgCancelSendToExternal{
		TransactionId: txId,
		Sender:        suite.AccAddress().String(),
		ChainName:     suite.chainName,
	})
}

func (suite *CrosschainTestSuite) SendBatchRequest(minTxs uint64) {
	msgList := make([]sdk.Msg, 0)
	batchFeeResponse, err := suite.CrosschainQuery().BatchFees(suite.ctx, &crosschaintypes.QueryBatchFeeRequest{ChainName: suite.chainName})
	suite.NoError(err)
	suite.True(len(batchFeeResponse.BatchFees) >= 1)
	for _, batchToken := range batchFeeResponse.BatchFees {
		suite.Equal(batchToken.TotalTxs, minTxs)

		denomResponse, err := suite.CrosschainQuery().TokenToDenom(suite.ctx, &crosschaintypes.QueryTokenToDenomRequest{
			Token:     batchToken.TokenContract,
			ChainName: suite.chainName,
		})
		suite.NoError(err)

		feeReceive := suite.HexAddress().String()
		if suite.chainName == trontypes.ModuleName {
			feeReceive = trontypes.AddressFromHex(feeReceive)
		}
		msgList = append(msgList, &crosschaintypes.MsgRequestBatch{
			Sender:     suite.BridgerAddr().String(),
			Denom:      denomResponse.Denom,
			MinimumFee: batchToken.TotalFees,
			FeeReceive: feeReceive,
			ChainName:  suite.chainName,
		})
	}
	suite.BroadcastTx(suite.bridgerPrivKey, msgList...)
}

func (suite *CrosschainTestSuite) SendConfirmBatch() {
	response, err := suite.CrosschainQuery().LastPendingBatchRequestByAddr(
		suite.ctx,
		&crosschaintypes.QueryLastPendingBatchRequestByAddrRequest{
			BridgerAddress: suite.BridgerAddr().String(),
			ChainName:      suite.chainName,
		},
	)
	suite.NoError(err)
	suite.NotNil(response.Batch)

	outgoingTxBatch := response.Batch
	var signatureBytes []byte
	if suite.chainName == trontypes.ModuleName {
		checkpoint, err := trontypes.GetCheckpointConfirmBatch(outgoingTxBatch, suite.params.GravityId)
		suite.NoError(err)
		signatureBytes, err = trontypes.NewTronSignature(checkpoint, suite.externalPrivKey)
		suite.NoError(err)
		err = trontypes.ValidateTronSignature(checkpoint, signatureBytes, suite.ExternalAddr())
		suite.NoError(err)
	} else {
		checkpoint, err := outgoingTxBatch.GetCheckpoint(suite.params.GravityId)
		suite.NoError(err)
		signatureBytes, err = crosschaintypes.NewEthereumSignature(checkpoint, suite.externalPrivKey)
		suite.NoError(err)
		err = crosschaintypes.ValidateEthereumSignature(checkpoint, signatureBytes, suite.ExternalAddr())
		suite.NoError(err)
	}

	suite.BroadcastTx(suite.bridgerPrivKey,
		&crosschaintypes.MsgConfirmBatch{
			Nonce:           outgoingTxBatch.BatchNonce,
			TokenContract:   outgoingTxBatch.TokenContract,
			BridgerAddress:  suite.BridgerAddr().String(),
			ExternalAddress: suite.ExternalAddr(),
			Signature:       hex.EncodeToString(signatureBytes),
			ChainName:       suite.chainName,
		},
		&crosschaintypes.MsgSendToExternalClaim{
			EventNonce:     suite.queryFxLastEventNonce(),
			BlockHeight:    suite.queryObserverExternalBlockHeight() + 1,
			BatchNonce:     outgoingTxBatch.BatchNonce,
			TokenContract:  outgoingTxBatch.TokenContract,
			BridgerAddress: suite.BridgerAddr().String(),
			ChainName:      suite.chainName,
		},
	)
}