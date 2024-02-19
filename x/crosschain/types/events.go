package types

const (
	EventTypeContractEvent = "observation"
	AttributeKeyClaimType  = "claim_type"
	AttributeKeyEventNonce = "event_nonce"
	AttributeKeyClaimHash  = "claim_hash"

	AttributeKeyBlockHeight  = "block_height"
	AttributeKeyStateSuccess = "state_success"

	EventTypeOracleSetUpdate   = "oracle_set_update"
	AttributeKeyOracleSetNonce = "oracle_set_nonce"
	AttributeKeyOracleSetLen   = "oracle_set_len"

	EventTypeSendToExternal         = "send_to_external"
	EventTypeSendToExternalCanceled = "send_to_external_canceled"
	EventTypeIncreaseBridgeFee      = "increase_bridge_fee"
	AttributeKeyOutgoingTxID        = "outgoing_tx_id"
	AttributeKeyIncreaseFee         = "increase_fee"

	EventTypeOutgoingBatch           = "outgoing_batch"
	EventTypeOutgoingBatchCanceled   = "outgoing_batch_canceled"
	AttributeKeyOutgoingTxIds        = "outgoing_tx_ids"
	AttributeKeyOutgoingBatchNonce   = "batch_nonce"
	AttributeKeyOutgoingBatchTimeout = "outgoing_batch_timeout"

	EventTypeIbcTransfer         = "ibc_transfer"
	AttributeKeyIbcSendSequence  = "ibc_send_sequence"
	AttributeKeyIbcSourcePort    = "ibc_source_port"
	AttributeKeyIbcSourceChannel = "ibc_source_channel"

	EventTypeEvmTransfer = "evm_transfer"

	EventTypeBridgeCallEvm    = "bridge_call_evm"
	AttributeKeyEvmCallResult = "evm_call_result"
	AttributeKeyEvmCallError  = "evm_call_error"
)
