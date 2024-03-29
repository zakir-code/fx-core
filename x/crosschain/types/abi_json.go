package types

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

// The go-ethereum ABI encoder *only* encodes function calls and then it only encodes
// function calls for which you provide an ABI json just like you would get out of the
// solidity compiler with your compiled contract.
// You are supposed to compile your contract, use abigen to generate an ABI , import
// this generated go module and then use for that for all testing and development.
// This abstraction layer is more trouble than it's worth, because we don't want to
// encode a function call at all, but instead we want to emulate a Solidity encode operation
// which has no equal available from go-ethereum.
//
// In order to work around this absurd series of problems we have to manually write the below
// 'function specification' that will encode the same arguments into a function call. We can then
// truncate the first several bytes where the call name is encoded to finally get the equal of the

var (
	outgoingBatchTxCheckpointABI      abi.ABI
	oracleSetCheckpointABI            abi.ABI
	outgoingBridgeCallTxCheckpointABI abi.ABI
)

func init() {
	const (
		// OutgoingBatchTxCheckpointABIJSON checks the ETH ABI for compatibility of the OutgoingBatchTx message
		outgoingBatchTxCheckpointABIJSON = `[{
		"name":"submitBatch",
		"stateMutability":"nonpayable",
		"type":"function",
		"inputs":[
			{ "internalType": "bytes32",   "name": "_fxbridgeId",    "type": "bytes32" },
			{ "internalType": "bytes32",   "name": "_methodName",    "type": "bytes32" },
			{ "internalType": "uint256[]", "name": "_amounts",       "type": "uint256[]" },
			{ "internalType": "address[]", "name": "_destinations",  "type": "address[]" },
			{ "internalType": "uint256[]", "name": "_fees",          "type": "uint256[]" },
			{ "internalType": "uint256",   "name": "_batchNonce",    "type": "uint256" },
			{ "internalType": "address",   "name": "_tokenContract", "type": "address" },
			{ "internalType": "uint256",   "name": "_batchTimeout",  "type": "uint256" },
			{ "internalType": "address",   "name": "_feeReceive",    "type": "address" }
		],
		"outputs": [
			{ "internalType": "bytes32", "name": "", "type": "bytes32" }
		]
	}]`

		// OracleSetCheckpointABIJSON checks the ETH ABI for compatibility of the OracleSet update message
		oracleSetCheckpointABIJSON = `[{
		"name": "checkpoint",
		"stateMutability": "pure",
		"type": "function",
		"inputs": [
			{ "internalType": "bytes32",   "name": "_fxbridgeId",  "type": "bytes32"   },
			{ "internalType": "bytes32",   "name": "_checkpoint",  "type": "bytes32"   },
			{ "internalType": "uint256",   "name": "_valsetNonce", "type": "uint256"   },
			{ "internalType": "address[]", "name": "_validators",  "type": "address[]" },
			{ "internalType": "uint256[]", "name": "_powers",      "type": "uint256[]" }
		],
		"outputs": [
			{ "internalType": "bytes32", "name": "", "type": "bytes32" }
		]
	}]`

		outgoingBridgeCallTxCheckpointABIJSON = `[{
		"name":"submitBridgeCall",
		"stateMutability":"nonpayable",
		"type":"function",
		"inputs":[
			{ "internalType": "bytes32",   "name": "_fxbridgeId",    "type": "bytes32" },
			{ "internalType": "bytes32",   "name": "_methodName",    "type": "bytes32" },
			{ "internalType": "address",   "name": "_receiver",      "type": "address" },
			{ "internalType": "address[]", "name": "_tokenContract", "type": "address[]" },
			{ "internalType": "uint256[]", "name": "_amounts",       "type": "uint256[]" },
			{ "internalType": "uint256",   "name": "_eventNonce",    "type": "uint256" },
			{ "internalType": "uint256",   "name": "_timeout",       "type": "uint256" }
		],
		"outputs": [
			{ "internalType": "bytes32", "name": "", "type": "bytes32" }
		]
	}]`
	)

	var err error
	outgoingBatchTxCheckpointABI, err = abi.JSON(strings.NewReader(outgoingBatchTxCheckpointABIJSON))
	if err != nil {
		panic(err.Error())
	}
	oracleSetCheckpointABI, err = abi.JSON(strings.NewReader(oracleSetCheckpointABIJSON))
	if err != nil {
		panic(err.Error())
	}
	outgoingBridgeCallTxCheckpointABI, err = abi.JSON(strings.NewReader(outgoingBridgeCallTxCheckpointABIJSON))
	if err != nil {
		panic(err.Error())
	}
}
