package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	// ModuleName is the name of the module
	ModuleName = "crosschain"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	// RouterKey is the module name router key
	RouterKey = ModuleName

	// QuerierRoute to be used for querier msgs
	QuerierRoute = ModuleName
)

var (
	// OracleTotalDepositKey
	OracleTotalDepositKey = []byte{0x11}

	// OracleKey key oracle address -> Oracle
	OracleKey = []byte{0x12}

	// OracleAddressByExternalKey key external address -> value oracle address
	OracleAddressByExternalKey = []byte{0x13}

	// OracleAddressByOrchestratorKey key orchestrator address -> value oracle address
	OracleAddressByOrchestratorKey = []byte{0x14}

	// OracleSetRequestKey indexes valset requests by nonce
	OracleSetRequestKey = []byte{0x15}

	// OracleSetConfirmKey indexes valset confirmations by nonce and the validator account address
	OracleSetConfirmKey = []byte{0x16}

	// OracleAttestationKey attestation details by nonce and validator address
	// An attestation can be thought of as the 'event to be executed' while
	// the Claims are an individual validator saying that they saw an event
	// occur the Attestation is 'the event' that multiple claims vote on and
	// eventually executes
	OracleAttestationKey = []byte{0x17}

	// OutgoingTxPoolKey indexes the last nonce for the outgoing tx pool
	OutgoingTxPoolKey = []byte{0x18}

	// SecondIndexOutgoingTxFeeKey indexes fee amounts by token contract address
	SecondIndexOutgoingTxFeeKey = []byte{0x19}

	// OutgoingTxBatchKey indexes outgoing tx batches under a nonce and token address
	OutgoingTxBatchKey = []byte{0x20}

	// OutgoingTxBatchBlockKey indexes outgoing tx batches under a block height and token address
	OutgoingTxBatchBlockKey = []byte{0x21}

	// BatchConfirmKey indexes validator confirmations by token contract address
	BatchConfirmKey = []byte{0x22}

	// LastEventNonceByValidatorKey indexes latest event nonce by validator
	LastEventNonceByValidatorKey = []byte{0x23}

	// LastObservedEventNonceKey indexes the latest event nonce
	LastObservedEventNonceKey = []byte{0x24}

	// SequenceKeyPrefix indexes different txIds
	SequenceKeyPrefix = []byte{0x25}

	// KeyLastTXPoolID indexes the lastTxPoolID
	KeyLastTXPoolID = append(SequenceKeyPrefix, []byte("lastTxPoolId")...)

	// KeyLastOutgoingBatchID indexes the lastBatchID
	KeyLastOutgoingBatchID = append(SequenceKeyPrefix, []byte("lastBatchId")...)

	// DenomToTokenKey prefixes the index of asset denom to external token
	DenomToTokenKey = []byte{0x26}

	// TokenToDenomKey prefixes the index of assets external token to denom
	TokenToDenomKey = []byte{0x27}

	// LastSlashedOracleSetNonce indexes the latest slashed oracleSet nonce
	LastSlashedOracleSetNonce = []byte{0x28}

	// LatestOracleSetNonce indexes the latest oracleSet nonce
	LatestOracleSetNonce = []byte{0x29}

	// LastSlashedBatchBlock indexes the latest slashed batch block height
	LastSlashedBatchBlock = []byte{0x30}

	// LastProposalBlockHeight indexes the last validator update block height
	LastProposalBlockHeight = []byte{0x31}

	// LastObservedBlockHeightKey indexes the latest observed external block height
	LastObservedBlockHeightKey = []byte{0x32}

	// LastObservedOracleSetKey indexes the latest observed OracleSet nonce
	LastObservedOracleSetKey = []byte{0x33}

	// KeyIbcSequenceHeight  indexes the gravity -> ibc sequence block height
	// DEPRECATED: delete by v2
	KeyIbcSequenceHeight = []byte{0x34}

	// LastEventBlockHeightByValidatorKey indexes latest event blockHeight by validator
	LastEventBlockHeightByValidatorKey = []byte{0x35}

	// PastExternalSignatureCheckpointKey indexes eth signature checkpoints that have existed
	PastExternalSignatureCheckpointKey = []byte{0x36}

	// LastOracleSlashBlockHeight indexes the last oracle slash block height
	LastOracleSlashBlockHeight = []byte{0x37}

	// KeyChainOracles -> value ChainOracle
	KeyChainOracles = []byte{0x38}

	// LastTotalPowerKey
	LastTotalPowerKey = []byte{0x39}
)

// GetOracleKey returns the following key format
func GetOracleKey(oracle sdk.AccAddress) []byte {
	return append(OracleKey, oracle.Bytes()...)
}

// GetOracleAddressByOrchestratorKey returns the following key format
func GetOracleAddressByOrchestratorKey(orchestrator sdk.AccAddress) []byte {
	return append(OracleAddressByOrchestratorKey, orchestrator.Bytes()...)
}

// GetOracleAddressByExternalKey returns the following key format
func GetOracleAddressByExternalKey(externalAddress string) []byte {
	return append(OracleAddressByExternalKey, []byte(externalAddress)...)
}

// GetOracleSetKey returns the following key format
func GetOracleSetKey(nonce uint64) []byte {
	return append(OracleSetRequestKey, sdk.Uint64ToBigEndian(nonce)...)
}

// GetOracleSetConfirmKey returns the following key format
func GetOracleSetConfirmKey(nonce uint64, oracleAddr sdk.AccAddress) []byte {
	return append(OracleSetConfirmKey, append(sdk.Uint64ToBigEndian(nonce), oracleAddr.Bytes()...)...)
}

// GetAttestationKey returns the following key format
// An attestation is an event multiple people are voting on, this function needs the claim
// details because each Attestation is aggregating all claims of a specific event, lets say
// validator X and validator y where making different claims about the same event nonce
// Note that the claim hash does NOT include the claimer address and only identifies an event
func GetAttestationKey(eventNonce uint64, claimHash []byte) []byte {
	return append(OracleAttestationKey, append(sdk.Uint64ToBigEndian(eventNonce), claimHash...)...)
}

// GetOutgoingTxPoolContractPrefix returns the following key format
// This prefix is used for iterating over unbatched transactions for a given contract
func GetOutgoingTxPoolContractPrefix(tokenContract string) []byte {
	return append(OutgoingTxPoolKey, []byte(tokenContract)...)
}

// GetOutgoingTxPoolKey returns the following key format
func GetOutgoingTxPoolKey(fee ExternalToken, id uint64) []byte {
	amount := make([]byte, 32)
	amount = fee.Amount.BigInt().FillBytes(amount)
	return append(OutgoingTxPoolKey, append([]byte(fee.Contract), append(amount, sdk.Uint64ToBigEndian(id)...)...)...)
}

// GetOutgoingTxBatchKey returns the following key format
func GetOutgoingTxBatchKey(tokenContract string, nonce uint64) []byte {
	return append(append(OutgoingTxBatchKey, []byte(tokenContract)...), sdk.Uint64ToBigEndian(nonce)...)
}

// GetOutgoingTxBatchBlockKey returns the following key format
func GetOutgoingTxBatchBlockKey(block uint64) []byte {
	return append(OutgoingTxBatchBlockKey, sdk.Uint64ToBigEndian(block)...)
}

// GetBatchConfirmKey returns the following key format
func GetBatchConfirmKey(tokenContract string, batchNonce uint64, oracleAddr sdk.AccAddress) []byte {
	return append(BatchConfirmKey, append([]byte(tokenContract), append(sdk.Uint64ToBigEndian(batchNonce), oracleAddr.Bytes()...)...)...)
}

// GetLastEventNonceByOracleKey returns the following key format
func GetLastEventNonceByOracleKey(validator sdk.AccAddress) []byte {
	return append(LastEventNonceByValidatorKey, validator.Bytes()...)
}

//GetIbcSequenceHeightKey returns the following key format
// DEPRECATED: delete by v2
func GetIbcSequenceHeightKey(sourcePort, sourceChannel string, sequence uint64) []byte {
	key := fmt.Sprintf("%s/%s/%d", sourcePort, sourceChannel, sequence)
	return append(KeyIbcSequenceHeight, []byte(key)...)
}

// GetLastEventBlockHeightByOracleKey returns the following key format
func GetLastEventBlockHeightByOracleKey(validator sdk.AccAddress) []byte {
	return append(LastEventBlockHeightByValidatorKey, validator.Bytes()...)
}

// GetDenomToTokenKey returns the following key format
func GetDenomToTokenKey(token string) []byte {
	return append(DenomToTokenKey, []byte(token)...)
}

// GetTokenToDenomKey returns the following key format
func GetTokenToDenomKey(denom string) []byte {
	return append(TokenToDenomKey, []byte(denom)...)
}

// GetPastExternalSignatureCheckpointKey returns the following key format
func GetPastExternalSignatureCheckpointKey(checkpoint []byte) []byte {
	return append(PastExternalSignatureCheckpointKey, checkpoint...)
}
