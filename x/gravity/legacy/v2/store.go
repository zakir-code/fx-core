// nolint:staticcheck
package v2

import (
	"fmt"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
	stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

	fxtypes "github.com/functionx/fx-core/v3/types"
	crosschaintypes "github.com/functionx/fx-core/v3/x/crosschain/types"
	ethtypes "github.com/functionx/fx-core/v3/x/eth/types"
	"github.com/functionx/fx-core/v3/x/gravity/types"
)

// MigrateStore performs in-place store migrations from v1 to v2.
// migrate data from gravity module
func MigrateStore(cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore) {

	// gravity 0x2 -> eth 0x13
	// key                                        			value
	// prefix     external-address                			oracle-address
	// [0x13][0xd98F9E3B1Bc6927700ce4A963429DC157dD4EBDf]   [0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9]
	// gravity 0xe -> eth 0x14
	// key                                        			value
	// prefix     bridger-address                			oracle-address
	// [0x14][0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9]    	[0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9]
	// migrate on MigrateValidatorToOracle

	// gravity 0x3 -> eth 0x15
	// key                                        			value
	// prefix     nonce  			                		OracleSet
	// [0x15][0 0 0 0 0 0 0 1]                           	[object marshal bytes]
	migratePrefix(gravityStore, ethStore, types.ValsetRequestKey, crosschaintypes.OracleSetRequestKey)

	// gravity 0x4 -> eth 0x16
	// key                                       				 			value
	// prefix     nonce  		oracle-address		               			MsgOracleSetConfirm
	// [0x16][0 0 0 0 0 0 0 1][0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9]     [object marshal bytes]
	migrateOracleSetConfirm(cdc, gravityStore, ethStore)

	// gravity 0x5 -> eth 0x17
	// key                                       				 									value
	// prefix     nonce                             claim-details-hash								Attestation
	// [0x5][0 0 0 0 0 0 0 1][fd1af8cec6c67fcf156f1b61fdf91ebc04d05484d007436e75342fc05bbff35a]		[object marshal bytes]
	migrateAttestation(cdc, gravityStore, ethStore)

	// gravity 0x6 and 0x7 -> eth 0x18
	// key                                       				 				value
	// prefix            token-address            		 fee_amount(byte32) 	IDSet -delete-
	// [0x7][0xb4fA5979babd8Bb7e427157d0d353Cf205F43752][1000000000000]			[object marshal bytes] delete
	// prefix            id											 			OutgoingTransferTx
	// [0x18][0 0 0 0 0 0 0 1]													[object marshal bytes]
	migrateOutgoingTxPool(cdc, gravityStore, ethStore)

	// gravity 0x8 -> eth 0x20
	// key                                       				 			value
	// prefix            token-address            		 nonce		 		OutgoingTxBatch
	// [0x20][0xb4fA5979babd8Bb7e427157d0d353Cf205F43752][0 0 0 0 0 0 0 1]	[object marshal bytes]
	migratePrefix(gravityStore, ethStore, types.OutgoingTxBatchKey, crosschaintypes.OutgoingTxBatchKey)

	// gravity 0x9 -> eth 0x21
	// key                                  value
	// prefix	block-height		 		OutgoingTxBatch
	// [0x21][0 0 0 0 0 0 0 1]				[object marshal bytes]
	migratePrefix(gravityStore, ethStore, types.OutgoingTxBatchBlockKey, crosschaintypes.OutgoingTxBatchBlockKey)

	// gravity 0xa -> eth 0x22
	// key                                  																			value
	// prefix           token-address                		batch-nonce            oracle-address						MsgConfirmBatch
	// [0x22][0xb4fA5979babd8Bb7e427157d0d353Cf205F43752][0 0 0 0 0 0 0 1][fx1mx8euwcmc6f8wqxwf2trg2wuz47af67lads8yg]	[object marshal bytes]
	migrateConfirmBatch(cdc, gravityStore, ethStore)

	// gravity 0xb -> eth 0x23
	// key                                  					value
	// prefix           oracle-address                			event-nonce
	// [0x23][0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9]			[0 0 0 0 0 0 0 1]
	migratePrefix(gravityStore, ethStore, types.LastEventNonceByValidatorKey, crosschaintypes.LastEventNonceByOracleKey)

	// gravity 0xc -> eth 0x24
	// key         		value
	// prefix           event-nonce
	// [0x24]			[0 0 0 0 0 0 0 1]
	migratePrefix(gravityStore, ethStore, types.LastObservedEventNonceKey, crosschaintypes.LastObservedEventNonceKey)

	// gravity 0xd+"lastTxPoolId" -> eth 0x25+"lastTxPoolId"
	// gravity 0xd+"lastBatchId"  -> eth 0x25+"lastBatchId"
	migratePrefix(gravityStore, ethStore, types.SequenceKeyPrefix, crosschaintypes.SequenceKeyPrefix)

	// gravity 0xf -> eth 0x26
	// key         														value
	// prefix  denom		   											BridgeToken
	// [0x26][eth0xb4fA5979babd8Bb7e427157d0d353Cf205F43752]			[object marshal bytes]

	// gravity 0x10 -> eth 0x27
	// key         														value
	// prefix  token-address   											BridgeToken
	// [0x27][0xb4fA5979babd8Bb7e427157d0d353Cf205F43752]				[object marshal bytes]
	migrateBridgeToken(cdc, gravityStore, ethStore)

	// gravity 0x11 -> eth 0x28
	// key         		value
	// prefix           oracle-set-nonce
	// [0x28]			[0 0 0 0 0 0 0 1]
	migratePrefix(gravityStore, ethStore, types.LastSlashedValsetNonce, crosschaintypes.LastSlashedOracleSetNonce)

	// gravity 0x12 -> eth 0x29
	// key         		value
	// prefix           oracle-set-nonce
	// [0x29]			[0 0 0 0 0 0 0 1]
	migratePrefix(gravityStore, ethStore, types.LatestValsetNonce, crosschaintypes.LatestOracleSetNonce)

	// gravity 0x13 -> eth 0x30
	// key         		value
	// prefix           block-height
	// [0x30]			[0 0 0 0 0 0 0 1]
	migratePrefix(gravityStore, ethStore, types.LastSlashedBatchBlock, crosschaintypes.LastSlashedBatchBlock)

	// gravity 0x14 -delete-
	// key         		value
	// prefix           block-height
	// [0x14]			[0 0 0 0 0 0 0 1]
	deletePrefixKey(gravityStore, types.LastUnBondingBlockHeight)

	// gravity 0x15 -> eth 0x32
	// key         		value
	// prefix           LastObservedBlockHeight
	// [0x32]			[object marshal bytes]
	migrateLastObservedBlockHeight(cdc, gravityStore, ethStore)

	// gravity 0x16 -> eth 0x33
	// key         		value
	// prefix           OracleSet
	// [0x33]			[object marshal bytes]
	migratePrefix(gravityStore, ethStore, types.LastObservedValsetKey, crosschaintypes.LastObservedOracleSetKey)

	// gravity 0x17 delete
	deletePrefixKey(gravityStore, types.IbcSequenceHeightKey)

	// gravity 0x18 -> eth 0x35
	// key                                  					value
	// prefix           oracle-address                			block-height
	// [0x35][0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9]			[0 0 0 0 0 0 0 1]
	migratePrefix(gravityStore, ethStore, types.LastEventBlockHeightByValidatorKey, crosschaintypes.LastEventBlockHeightByOracleKey)
}

func migratePrefix(gravityStore, ethStore sdk.KVStore, oldPrefix, newPrefix []byte) {
	oldStore := prefix.NewStore(gravityStore, oldPrefix)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()
	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		key := oldStoreIter.Key()
		ethStore.Set(append(newPrefix, key...), oldStoreIter.Value())
		oldStore.Delete(key)
	}
}

func MigrateValidatorToOracle(ctx sdk.Context, cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore, stakingKeeper StakingKeeper, bankKeeper BankKeeper) {

	chainOracle := new(crosschaintypes.ProposalOracle)
	totalPower := sdk.ZeroInt()

	ethOracles := EthInitOracles(ctx.ChainID())
	index := 0
	minDelegateAmount := sdk.DefaultPowerReduction.MulRaw(100)

	oldStore := prefix.NewStore(gravityStore, types.ValidatorAddressByOrchestratorAddress)

	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()
	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		bridgerAddr := sdk.AccAddress(oldStoreIter.Key())
		oldOracleAddress := sdk.AccAddress(oldStoreIter.Value())
		externalAddress := string(gravityStore.Get(append(types.EthAddressByValidatorKey, oldOracleAddress.Bytes()...)))
		validator, found := stakingKeeper.GetValidator(ctx, oldOracleAddress.Bytes())
		if !found {
			//ctx.Logger().Error("no found validator", "address", sdk.ValAddress(oldOracleAddress))
			//continue
			panic(fmt.Sprintf("no found validator: %s", sdk.ValAddress(oldOracleAddress).String()))
		}
		oracle := crosschaintypes.Oracle{
			BridgerAddress:    bridgerAddr.String(),
			ExternalAddress:   externalAddress,
			StartHeight:       0,
			DelegateValidator: oldOracleAddress.String(),
			DelegateAmount:    sdk.ZeroInt(),
			Online:            false,
			OracleAddress:     oldOracleAddress.String(),
			SlashTimes:        0,
		}
		if len(ethOracles) > index {
			oracle.OracleAddress = ethOracles[index]
			oracleAddr := oracle.GetOracle()
			balances := bankKeeper.GetAllBalances(ctx, oracleAddr)
			if balances.AmountOf(fxtypes.DefaultDenom).GTE(minDelegateAmount) {
				delegateAddr := oracle.GetDelegateAddress(ethtypes.ModuleName)
				if err := bankKeeper.SendCoins(ctx, oracleAddr, delegateAddr,
					sdk.NewCoins(sdk.NewCoin(fxtypes.DefaultDenom, minDelegateAmount))); err != nil {
					panic("send to coins error: " + err.Error())
				}
				newShares, err := stakingKeeper.Delegate(ctx,
					delegateAddr, minDelegateAmount, stakingtypes.Unbonded, validator, true)
				if err != nil {
					panic("gravity migrate to eth error: " + err.Error())
				}
				oracle.StartHeight = ctx.BlockHeight()
				oracle.DelegateAmount = minDelegateAmount
				oracle.Online = true
				ctx.EventManager().EmitEvents(sdk.Events{
					sdk.NewEvent(
						stakingtypes.EventTypeDelegate,
						sdk.NewAttribute(stakingtypes.AttributeKeyValidator, oracle.DelegateValidator),
						sdk.NewAttribute(sdk.AttributeKeyAmount, oracle.DelegateAmount.String()),
						sdk.NewAttribute(stakingtypes.AttributeKeyNewShares, newShares.String()),
					),
				})
			}
		}
		index = index + 1

		if oracle.Online {
			totalPower = totalPower.Add(oracle.GetPower())
		}
		oracleAddress := oracle.GetOracle()
		ethStore.Set(append(crosschaintypes.OracleAddressByExternalKey, []byte(oracle.ExternalAddress)...), oracleAddress.Bytes())
		ethStore.Set(append(crosschaintypes.OracleAddressByBridgerKey, bridgerAddr.Bytes()...), oracleAddress.Bytes())
		// SetOracle
		ethStore.Set(crosschaintypes.GetOracleKey(oracleAddress), cdc.MustMarshal(&oracle))
		oldStore.Delete(oldStoreIter.Key())

		chainOracle.Oracles = append(chainOracle.Oracles, oracle.OracleAddress)
	}

	// SetProposalOracle eth 0x38
	if len(chainOracle.Oracles) > 0 {
		ethStore.Set(crosschaintypes.ProposalOracleKey, cdc.MustMarshal(chainOracle))
	}
	// setLastTotalPower eth 0x39
	ethStore.Set(crosschaintypes.LastTotalPowerKey, cdc.MustMarshal(&sdk.IntProto{Int: totalPower}))

	// gravity 0x1 -> eth 0x12
	deletePrefixKey(gravityStore, types.EthAddressByValidatorKey)
}

func migrateOutgoingTxPool(cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore) {

	oldStore := prefix.NewStore(gravityStore, types.OutgoingTxPoolKey)
	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var transact crosschaintypes.OutgoingTransferTx
		cdc.MustUnmarshal(oldStoreIter.Value(), &transact)

		ethStore.Set(crosschaintypes.GetOutgoingTxPoolKey(transact.Fee, transact.Id), oldStoreIter.Value())
		oldStore.Delete(oldStoreIter.Key())
	}

	oldStore2 := prefix.NewStore(gravityStore, types.SecondIndexOutgoingTxFeeKey)
	oldStoreIter2 := oldStore2.Iterator(nil, nil)
	defer oldStoreIter2.Close()
	for ; oldStoreIter2.Valid(); oldStoreIter2.Next() {
		oldStore2.Delete(oldStoreIter2.Key())
	}
}

func migrateOracleSetConfirm(cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore) {
	oldStore := prefix.NewStore(gravityStore, types.ValsetConfirmKey)
	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var msg types.MsgValsetConfirm
		cdc.MustUnmarshal(oldStoreIter.Value(), &msg)

		key := oldStoreIter.Key()
		nonce := sdk.BigEndianToUint64(key[:8])
		if nonce != msg.Nonce {
			panic(fmt.Sprintf("invalid nonce, expect: %d, actual: %d", nonce, msg.Nonce))
		}
		bridgeAddr := key[8:]
		oracleAddr := ethStore.Get(crosschaintypes.GetOracleAddressByBridgerKey(bridgeAddr))
		if len(oracleAddr) >= 20 {
			panic(fmt.Sprintf("invalid oracle address: %v", oracleAddr))
		}

		ethStore.Set(crosschaintypes.GetOracleSetConfirmKey(msg.Nonce, oracleAddr),
			cdc.MustMarshal(&crosschaintypes.MsgOracleSetConfirm{
				Nonce:           msg.Nonce,
				BridgerAddress:  msg.Orchestrator,
				ExternalAddress: msg.EthAddress,
				Signature:       msg.Signature,
				ChainName:       ethtypes.ModuleName,
			}),
		)
		oldStore.Delete(oldStoreIter.Key())
	}
}

func migrateConfirmBatch(cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore) {
	oldStore := prefix.NewStore(gravityStore, types.BatchConfirmKey)
	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var msg types.MsgConfirmBatch
		cdc.MustUnmarshal(oldStoreIter.Value(), &msg)

		key := oldStoreIter.Key()
		token := string(key[:len(msg.TokenContract)])
		if token != msg.TokenContract {
			panic(fmt.Sprintf("invalid token contract, expect: %s, actual: %s", token, msg.TokenContract))
		}
		nonce := sdk.BigEndianToUint64(key[len(msg.TokenContract) : len(msg.TokenContract)+8])
		if nonce != msg.Nonce {
			panic(fmt.Sprintf("invalid nonce, expect: %d, actual: %d", nonce, msg.Nonce))
		}
		bridgeAddr := key[len(msg.TokenContract)+8:]
		oracleAddr := ethStore.Get(crosschaintypes.GetOracleAddressByBridgerKey(bridgeAddr))
		if len(oracleAddr) >= 20 {
			panic(fmt.Sprintf("invalid oracle address: %v", oracleAddr))
		}

		ethStore.Set(crosschaintypes.GetBatchConfirmKey(msg.TokenContract, msg.Nonce, oracleAddr),
			cdc.MustMarshal(&crosschaintypes.MsgConfirmBatch{
				Nonce:           msg.Nonce,
				TokenContract:   msg.TokenContract,
				BridgerAddress:  msg.Orchestrator,
				ExternalAddress: msg.EthSigner,
				Signature:       msg.Signature,
				ChainName:       ethtypes.ModuleName,
			}),
		)

		oldStore.Delete(oldStoreIter.Key())
	}
}

func migrateLastObservedBlockHeight(cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore) {
	var msg types.LastObservedEthereumBlockHeight
	cdc.MustUnmarshal(gravityStore.Get(types.LastObservedEthereumBlockHeightKey), &msg)

	ethStore.Set(crosschaintypes.LastObservedBlockHeightKey,
		cdc.MustMarshal(&crosschaintypes.LastObservedBlockHeight{
			ExternalBlockHeight: msg.EthBlockHeight,
			BlockHeight:         msg.FxBlockHeight,
		}),
	)
	gravityStore.Delete(types.LastObservedEthereumBlockHeightKey)
}

func migrateAttestation(cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore) {
	oldStore := prefix.NewStore(gravityStore, types.OracleAttestationKey)
	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()

	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		var att types.Attestation
		cdc.MustUnmarshal(oldStoreIter.Value(), &att)

		claim, err := types.UnpackAttestationClaim(cdc, &att)
		if err != nil {
			panic(err.Error())
		}
		var newClaim crosschaintypes.ExternalClaim
		switch c := claim.(type) {
		case *types.MsgDepositClaim:
			newClaim = &crosschaintypes.MsgSendToFxClaim{
				EventNonce:     c.EventNonce,
				BlockHeight:    c.BlockHeight,
				TokenContract:  c.TokenContract,
				Amount:         c.Amount,
				Sender:         c.EthSender,
				Receiver:       c.FxReceiver,
				TargetIbc:      c.TargetIbc,
				BridgerAddress: c.Orchestrator,
				ChainName:      ethtypes.ModuleName,
			}
		case *types.MsgWithdrawClaim:
			newClaim = &crosschaintypes.MsgSendToExternalClaim{
				EventNonce:     c.EventNonce,
				BlockHeight:    c.BlockHeight,
				BatchNonce:     c.BatchNonce,
				TokenContract:  c.TokenContract,
				BridgerAddress: c.Orchestrator,
				ChainName:      ethtypes.ModuleName,
			}
		case *types.MsgFxOriginatedTokenClaim:
			//newClaim = &crosschaintypes.MsgBridgeTokenClaim{
			//	EventNonce:     c.EventNonce,
			//	BlockHeight:    c.BlockHeight,
			//	TokenContract:  c.TokenContract,
			//	Name:           c.Name,
			//	Symbol:         c.Symbol,
			//	Decimals:       c.Decimals,
			//	BridgerAddress: c.Orchestrator,
			//	ChannelIbc:     "",
			//	ChainName:      ethtypes.ModuleName,
			//}
		case *types.MsgValsetUpdatedClaim:
			myClaim := &crosschaintypes.MsgOracleSetUpdatedClaim{
				EventNonce:     c.EventNonce,
				BlockHeight:    c.BlockHeight,
				OracleSetNonce: c.ValsetNonce,
				Members:        make([]crosschaintypes.BridgeValidator, len(c.Members)),
				BridgerAddress: c.Orchestrator,
				ChainName:      ethtypes.ModuleName,
			}
			for i := 0; i < len(c.Members); i++ {
				myClaim.Members[i] = crosschaintypes.BridgeValidator{
					Power:           c.Members[i].Power,
					ExternalAddress: c.Members[i].EthAddress,
				}
			}
			newClaim = myClaim
		}

		anyMsg, err := codectypes.NewAnyWithValue(newClaim)
		if err != nil {
			panic(err.Error())
		}

		// new claim hash
		ethStore.Set(crosschaintypes.GetAttestationKey(claim.GetEventNonce(), claim.ClaimHash()),
			cdc.MustMarshal(&crosschaintypes.Attestation{
				Observed: att.Observed,
				Votes:    att.Votes,
				Height:   att.Height,
				Claim:    anyMsg,
			}),
		)
		oldStore.Delete(oldStoreIter.Key())
	}
}

func migrateBridgeToken(cdc codec.BinaryCodec, gravityStore, ethStore sdk.KVStore) {
	token := gravityStore.Get(append(types.DenomToERC20Key, []byte(fxtypes.DefaultDenom)...))
	if len(token) > 0 {
		ethStore.Set(crosschaintypes.GetTokenToDenomKey(fxtypes.DefaultDenom), token)
		ethStore.Set(crosschaintypes.GetDenomToTokenKey(string(token)), []byte(fxtypes.DefaultDenom))
	}
	gravityStore.Delete(types.DenomToERC20Key)
	gravityStore.Delete(types.ERC20ToDenomKey)
}

func MigrateBridgeTokenFromMetaDatas(cdc codec.BinaryCodec, metaDatas []banktypes.Metadata, ethStore sdk.KVStore) {
	for _, data := range metaDatas {
		if len(data.DenomUnits) > 0 && len(data.DenomUnits[0].Aliases) > 0 {
			for i := 0; i < len(data.DenomUnits[0].Aliases); i++ {
				denom := data.DenomUnits[0].Aliases[i]
				if strings.HasPrefix(denom, ethtypes.ModuleName) {
					token := strings.TrimPrefix(denom, ethtypes.ModuleName)
					ethStore.Set(crosschaintypes.GetTokenToDenomKey(denom), []byte(token))
					ethStore.Set(crosschaintypes.GetDenomToTokenKey(token), []byte(denom))
				}
			}
		}
	}
}

func deletePrefixKey(gravityStore sdk.KVStore, prefixKey []byte) {
	oldStore := prefix.NewStore(gravityStore, prefixKey)
	oldStoreIter := oldStore.Iterator(nil, nil)
	defer oldStoreIter.Close()
	for ; oldStoreIter.Valid(); oldStoreIter.Next() {
		oldStore.Delete(oldStoreIter.Key())
	}
}