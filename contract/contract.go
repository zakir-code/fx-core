package contract

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

const (
	FIP20LogicAddress = "0x0000000000000000000000000000000000001001"
	WFXLogicAddress   = "0x0000000000000000000000000000000000001002"

	StakingAddress    = "0x0000000000000000000000000000000000001003"
	CrossChainAddress = "0x0000000000000000000000000000000000001004"
)

var (
	initFIP20Code = MustDecodeHex("0x60806040526004361061011f5760003560e01c8063715018a6116100a0578063b86d529811610064578063b86d529814610306578063c5cb9b5114610324578063dd62ed3e14610344578063de7ea79d1461038a578063f2fde38b146103aa5761011f565b8063715018a61461026a5780638da5cb5b1461027f57806395d89b41146102b15780639dc29fac146102c6578063a9059cbb146102e65761011f565b80633659cfe6116100e75780633659cfe6146101e057806340c10f19146102025780634f1ef2861461022257806352d1902d1461023557806370a082311461024a5761011f565b806306fdde0314610124578063095ea7b31461014f57806318160ddd1461017f57806323b872dd1461019e578063313ce567146101be575b600080fd5b34801561013057600080fd5b506101396103ca565b6040516101469190611b5b565b60405180910390f35b34801561015b57600080fd5b5061016f61016a3660046118df565b61045c565b6040519015158152602001610146565b34801561018b57600080fd5b5060cc545b604051908152602001610146565b3480156101aa57600080fd5b5061016f6101b9366004611845565b6104b2565b3480156101ca57600080fd5b5060cb5460405160ff9091168152602001610146565b3480156101ec57600080fd5b506102006101fb3660046117f9565b61055f565b005b34801561020e57600080fd5b5061020061021d3660046118df565b61063f565b610200610230366004611880565b610655565b34801561024157600080fd5b50610190610722565b34801561025657600080fd5b506101906102653660046117f9565b6107d5565b34801561027657600080fd5b506102006107f4565b34801561028b57600080fd5b506097546001600160a01b03165b6040516001600160a01b039091168152602001610146565b3480156102bd57600080fd5b50610139610808565b3480156102d257600080fd5b506102006102e13660046118df565b610817565b3480156102f257600080fd5b5061016f6103013660046118df565b610829565b34801561031257600080fd5b5060cf546001600160a01b0316610299565b34801561033057600080fd5b5061016f61033f366004611a3c565b61083f565b34801561035057600080fd5b5061019061035f366004611813565b6001600160a01b03918216600090815260ce6020908152604080832093909416825291909152205490565b34801561039657600080fd5b506102006103a53660046119b3565b6108f6565b3480156103b657600080fd5b506102006103c53660046117f9565b610a65565b606060c980546103d990611d34565b80601f016020809104026020016040519081016040528092919081815260200182805461040590611d34565b80156104525780601f1061042757610100808354040283529160200191610452565b820191906000526020600020905b81548152906001019060200180831161043557829003601f168201915b5050505050905090565b6000610469338484610adb565b6040518281526001600160a01b0384169033907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259060200160405180910390a350600192915050565b6001600160a01b038316600090815260ce60209081526040808320338452909152812054828110156105355760405162461bcd60e51b815260206004820152602160248201527f7472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636044820152606560f81b60648201526084015b60405180910390fd5b61054985336105448685611cf1565b610adb565b610554858585610b5d565b506001949350505050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000010011614156105a85760405162461bcd60e51b815260040161052c90611b9d565b7f00000000000000000000000000000000000000000000000000000000000010016001600160a01b03166105f1600080516020611d9c833981519152546001600160a01b031690565b6001600160a01b0316146106175760405162461bcd60e51b815260040161052c90611be9565b61062081610d0c565b6040805160008082526020820190925261063c91839190610d14565b50565b610647610e98565b6106518282610ef2565b5050565b306001600160a01b037f000000000000000000000000000000000000000000000000000000000000100116141561069e5760405162461bcd60e51b815260040161052c90611b9d565b7f00000000000000000000000000000000000000000000000000000000000010016001600160a01b03166106e7600080516020611d9c833981519152546001600160a01b031690565b6001600160a01b03161461070d5760405162461bcd60e51b815260040161052c90611be9565b61071682610d0c565b61065182826001610d14565b6000306001600160a01b037f000000000000000000000000000000000000000000000000000000000000100116146107c25760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c0000000000000000606482015260840161052c565b50600080516020611d9c83398151915290565b6001600160a01b038116600090815260cd60205260409020545b919050565b6107fc610e98565b6108066000610fd1565b565b606060ca80546103d990611d34565b61081f610e98565b6106518282611023565b6000610836338484610b5d565b50600192915050565b600063ffffffff333b16156108965760405162461bcd60e51b815260206004820152601960248201527f63616c6c65722063616e6e6f7420626520636f6e747261637400000000000000604482015260640161052c565b6108a33386868686611165565b336001600160a01b03167f282dd1817b996776123a00596764d4d54cc16460c9854f7a23f6be020ba0463d868686866040516108e29493929190611b6e565b60405180910390a25060015b949350505050565b600054610100900460ff16158080156109165750600054600160ff909116105b806109305750303b158015610930575060005460ff166001145b6109935760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b606482015260840161052c565b6000805460ff1916600117905580156109b6576000805461ff0019166101001790555b84516109c99060c99060208801906116ec565b5083516109dd9060ca9060208701906116ec565b5060cb805460ff191660ff851617905560cf80546001600160a01b0319166001600160a01b038416179055610a10611284565b610a186112b3565b8015610a5e576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050505050565b610a6d610e98565b6001600160a01b038116610ad25760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b606482015260840161052c565b61063c81610fd1565b6001600160a01b038316610b315760405162461bcd60e51b815260206004820152601d60248201527f617070726f76652066726f6d20746865207a65726f2061646472657373000000604482015260640161052c565b6001600160a01b03928316600090815260ce602090815260408083209490951682529290925291902055565b6001600160a01b038316610bb35760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f20616464726573730000604482015260640161052c565b6001600160a01b038216610c095760405162461bcd60e51b815260206004820152601c60248201527f7472616e7366657220746f20746865207a65726f206164647265737300000000604482015260640161052c565b6001600160a01b038316600090815260cd602052604090205481811015610c725760405162461bcd60e51b815260206004820152601f60248201527f7472616e7366657220616d6f756e7420657863656564732062616c616e636500604482015260640161052c565b610c7c8282611cf1565b6001600160a01b03808616600090815260cd60205260408082209390935590851681529081208054849290610cb2908490611cd9565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610cfe91815260200190565b60405180910390a350505050565b61063c610e98565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff1615610d4c57610d47836112da565b610e93565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b815260040160206040518083038186803b158015610d8557600080fd5b505afa925050508015610db5575060408051601f3d908101601f19168201909252610db291810190611928565b60015b610e185760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b606482015260840161052c565b600080516020611d9c8339815191528114610e875760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b606482015260840161052c565b50610e93838383611376565b505050565b6097546001600160a01b031633146108065760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e6572604482015260640161052c565b6001600160a01b038216610f485760405162461bcd60e51b815260206004820152601860248201527f6d696e7420746f20746865207a65726f20616464726573730000000000000000604482015260640161052c565b8060cc6000828254610f5a9190611cd9565b90915550506001600160a01b038216600090815260cd602052604081208054839290610f87908490611cd9565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a35050565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6001600160a01b0382166110795760405162461bcd60e51b815260206004820152601a60248201527f6275726e2066726f6d20746865207a65726f2061646472657373000000000000604482015260640161052c565b6001600160a01b038216600090815260cd6020526040902054818110156110e25760405162461bcd60e51b815260206004820152601b60248201527f6275726e20616d6f756e7420657863656564732062616c616e63650000000000604482015260640161052c565b6110ec8282611cf1565b6001600160a01b038416600090815260cd602052604081209190915560cc805484929061111a908490611cf1565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b6001600160a01b0385166111bb5760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f20616464726573730000604482015260640161052c565b60008451116112005760405162461bcd60e51b81526020600482015260116024820152701a5b9d985b1a59081c9958da5c1a595b9d607a1b604482015260640161052c565b8061123e5760405162461bcd60e51b815260206004820152600e60248201526d1a5b9d985b1a59081d185c99d95d60921b604482015260640161052c565b60cf5461125f9086906001600160a01b031661125a8587611cd9565b610b5d565b61127c8585858585604051806020016040528060008152506113a1565b505050505050565b600054610100900460ff166112ab5760405162461bcd60e51b815260040161052c90611c35565b610806611459565b600054610100900460ff166108065760405162461bcd60e51b815260040161052c90611c35565b6001600160a01b0381163b6113475760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b606482015260840161052c565b600080516020611d9c83398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b61137f83611489565b60008251118061138c5750805b15610e935761139b83836114c9565b50505050565b600080806110046113b68a8a8a8a8a8a6114f5565b6040516113c39190611aba565b6000604051808303816000865af19150503d8060008114611400576040519150601f19603f3d011682016040523d82523d6000602084013e611405565b606091505b5091509150611443828260405180604001604052806016815260200175199a5c0b58dc9bdcdccb58da185a5b8819985a5b195960521b815250611548565b61144c816115c2565b9998505050505050505050565b600054610100900460ff166114805760405162461bcd60e51b815260040161052c90611c35565b61080633610fd1565b611492816112da565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b60606114ee8383604051806060016040528060278152602001611dbc602791396115d9565b9392505050565b606086868686868660405160240161151296959493929190611b13565b60408051601f198184030181529190526020810180516001600160e01b0316633c3e7d7760e01b17905290509695505050505050565b82610e93576000828060200190518101906115639190611940565b9050600182511015611589578060405162461bcd60e51b815260040161052c9190611b5b565b818160405160200161159c929190611ad6565b60408051601f198184030181529082905262461bcd60e51b825261052c91600401611b5b565b600080828060200190518101906114ee9190611908565b6060600080856001600160a01b0316856040516115f69190611aba565b600060405180830381855af49150503d8060008114611631576040519150601f19603f3d011682016040523d82523d6000602084013e611636565b606091505b509150915061164786838387611651565b9695505050505050565b606083156116bd5782516116b6576001600160a01b0385163b6116b65760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e7472616374000000604482015260640161052c565b50816108ee565b6108ee83838151156116d25781518083602001fd5b8060405162461bcd60e51b815260040161052c9190611b5b565b8280546116f890611d34565b90600052602060002090601f01602090048101928261171a5760008555611760565b82601f1061173357805160ff1916838001178555611760565b82800160010185558215611760579182015b82811115611760578251825591602001919060010190611745565b5061176c929150611770565b5090565b5b8082111561176c5760008155600101611771565b600061179861179384611cb1565b611c80565b90508281528383830111156117ac57600080fd5b828260208301376000602084830101529392505050565b80356001600160a01b03811681146107ef57600080fd5b600082601f8301126117ea578081fd5b6114ee83833560208501611785565b60006020828403121561180a578081fd5b6114ee826117c3565b60008060408385031215611825578081fd5b61182e836117c3565b915061183c602084016117c3565b90509250929050565b600080600060608486031215611859578081fd5b611862846117c3565b9250611870602085016117c3565b9150604084013590509250925092565b60008060408385031215611892578182fd5b61189b836117c3565b9150602083013567ffffffffffffffff8111156118b6578182fd5b8301601f810185136118c6578182fd5b6118d585823560208401611785565b9150509250929050565b600080604083850312156118f1578182fd5b6118fa836117c3565b946020939093013593505050565b600060208284031215611919578081fd5b815180151581146114ee578182fd5b600060208284031215611939578081fd5b5051919050565b600060208284031215611951578081fd5b815167ffffffffffffffff811115611967578182fd5b8201601f81018413611977578182fd5b805161198561179382611cb1565b818152856020838501011115611999578384fd5b6119aa826020830160208601611d08565b95945050505050565b600080600080608085870312156119c8578081fd5b843567ffffffffffffffff808211156119df578283fd5b6119eb888389016117da565b95506020870135915080821115611a00578283fd5b50611a0d878288016117da565b935050604085013560ff81168114611a23578182fd5b9150611a31606086016117c3565b905092959194509250565b60008060008060808587031215611a51578384fd5b843567ffffffffffffffff811115611a67578485fd5b611a73878288016117da565b97602087013597506040870135966060013595509350505050565b60008151808452611aa6816020860160208601611d08565b601f01601f19169290920160200192915050565b60008251611acc818460208701611d08565b9190910192915050565b60008351611ae8818460208801611d08565b6101d160f51b9083019081528351611b07816002840160208801611d08565b01600201949350505050565b6001600160a01b038716815260c060208201819052600090611b3790830188611a8e565b86604084015285606084015284608084015282810360a084015261144c8185611a8e565b6000602082526114ee6020830184611a8e565b600060808252611b816080830187611a8e565b6020830195909552506040810192909252606090910152919050565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b604051601f8201601f1916810167ffffffffffffffff81118282101715611ca957611ca9611d85565b604052919050565b600067ffffffffffffffff821115611ccb57611ccb611d85565b50601f01601f191660200190565b60008219821115611cec57611cec611d6f565b500190565b600082821015611d0357611d03611d6f565b500390565b60005b83811015611d23578181015183820152602001611d0b565b8381111561139b5750506000910152565b600281046001821680611d4857607f821691505b60208210811415611d6957634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052604160045260246000fdfe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a26469706673582212208b5e915fa08dc1ea750a0ba5395d33e136da48ff45b03242745cc1516b80f05164736f6c63430008020033")
	initWFXCode   = MustDecodeHex("0x6080604052600436106101445760003560e01c8063715018a6116100b6578063c5cb9b511161006f578063c5cb9b511461038f578063d0e30db014610153578063dd62ed3e146103a2578063de7ea79d146103e8578063f2fde38b14610408578063f3fef3a31461042857610153565b8063715018a6146102d55780638da5cb5b146102ea57806395d89b411461031c5780639dc29fac14610331578063a9059cbb14610351578063b86d52981461037157610153565b8063313ce56711610108578063313ce567146102155780633659cfe61461023757806340c10f19146102575780634f1ef2861461027757806352d1902d1461028a57806370a082311461029f57610153565b806306fdde031461015b578063095ea7b31461018657806318160ddd146101b657806323b872dd146101d55780632e1a7d4d146101f557610153565b3661015357610151610448565b005b610151610448565b34801561016757600080fd5b50610170610489565b60405161017d9190611d1e565b60405180910390f35b34801561019257600080fd5b506101a66101a1366004611a9f565b61051b565b604051901515815260200161017d565b3480156101c257600080fd5b5060cc545b60405190815260200161017d565b3480156101e157600080fd5b506101a66101f03660046119fe565b610571565b34801561020157600080fd5b50610151610210366004611c39565b61061e565b34801561022157600080fd5b5060cb5460405160ff909116815260200161017d565b34801561024357600080fd5b5061015161025236600461197f565b61068f565b34801561026357600080fd5b50610151610272366004611a9f565b61076f565b610151610285366004611a3e565b610785565b34801561029657600080fd5b506101c7610852565b3480156102ab57600080fd5b506101c76102ba36600461197f565b6001600160a01b0316600090815260cd602052604090205490565b3480156102e157600080fd5b50610151610905565b3480156102f657600080fd5b506097546001600160a01b03165b6040516001600160a01b03909116815260200161017d565b34801561032857600080fd5b50610170610919565b34801561033d57600080fd5b5061015161034c366004611a9f565b610928565b34801561035d57600080fd5b506101a661036c366004611a9f565b61093a565b34801561037d57600080fd5b5060cf546001600160a01b0316610304565b6101a661039d366004611be7565b610950565b3480156103ae57600080fd5b506101c76103bd3660046119c6565b6001600160a01b03918216600090815260ce6020908152604080832093909416825291909152205490565b3480156103f457600080fd5b50610151610403366004611b5c565b610a15565b34801561041457600080fd5b5061015161042336600461197f565b610b84565b34801561043457600080fd5b5061015161044336600461199b565b610bfa565b6104523334610c7f565b60405134815233907fe1fffcc4923d04b559f4d29a8bfc6cda04eb5b0d3c460751c2402c5c5cc9109c9060200160405180910390a2565b606060c9805461049890611ef7565b80601f01602080910402602001604051908101604052809291908181526020018280546104c490611ef7565b80156105115780601f106104e657610100808354040283529160200191610511565b820191906000526020600020905b8154815290600101906020018083116104f457829003601f168201915b5050505050905090565b6000610528338484610d57565b6040518281526001600160a01b0384169033907f8c5be1e5ebec7d5bd14f71427d1e84f3dd0314c0f7b2291e5b200ac8c7c3b9259060200160405180910390a350600192915050565b6001600160a01b038316600090815260ce60209081526040808320338452909152812054828110156105f45760405162461bcd60e51b815260206004820152602160248201527f7472616e7366657220616d6f756e74206578636565647320616c6c6f77616e636044820152606560f81b60648201526084015b60405180910390fd5b61060885336106038685611eb4565b610d57565b610613858585610dd9565b506001949350505050565b610629335b82610f88565b604051339082156108fc029083906000818181858888f19350505050158015610656573d6000803e3d6000fd5b5060405181815233907f884edad9ce6fa2440d8a54cc123490eb96d2768479d49ff9c7366125a94243649060200160405180910390a250565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000010021614156106d85760405162461bcd60e51b81526004016105eb90611d60565b7f00000000000000000000000000000000000000000000000000000000000010026001600160a01b0316610721600080516020611f74833981519152546001600160a01b031690565b6001600160a01b0316146107475760405162461bcd60e51b81526004016105eb90611dac565b610750816110ca565b6040805160008082526020820190925261076c918391906110d2565b50565b610777611256565b6107818282610c7f565b5050565b306001600160a01b037f00000000000000000000000000000000000000000000000000000000000010021614156107ce5760405162461bcd60e51b81526004016105eb90611d60565b7f00000000000000000000000000000000000000000000000000000000000010026001600160a01b0316610817600080516020611f74833981519152546001600160a01b031690565b6001600160a01b03161461083d5760405162461bcd60e51b81526004016105eb90611dac565b610846826110ca565b610781828260016110d2565b6000306001600160a01b037f000000000000000000000000000000000000000000000000000000000000100216146108f25760405162461bcd60e51b815260206004820152603860248201527f555550535570677261646561626c653a206d757374206e6f742062652063616c60448201527f6c6564207468726f7567682064656c656761746563616c6c000000000000000060648201526084016105eb565b50600080516020611f7483398151915290565b61090d611256565b61091760006112b0565b565b606060ca805461049890611ef7565b610930611256565b6107818282610f88565b6000610947338484610dd9565b50600192915050565b600063ffffffff333b16156109a75760405162461bcd60e51b815260206004820152601960248201527f63616c6c65722063616e6e6f7420626520636f6e74726163740000000000000060448201526064016105eb565b34156109b5576109b5610448565b6109c23386868686611302565b336001600160a01b03167f282dd1817b996776123a00596764d4d54cc16460c9854f7a23f6be020ba0463d86868686604051610a019493929190611d31565b60405180910390a25060015b949350505050565b600054610100900460ff1615808015610a355750600054600160ff909116105b80610a4f5750303b158015610a4f575060005460ff166001145b610ab25760405162461bcd60e51b815260206004820152602e60248201527f496e697469616c697a61626c653a20636f6e747261637420697320616c72656160448201526d191e481a5b9a5d1a585b1a5e995960921b60648201526084016105eb565b6000805460ff191660011790558015610ad5576000805461ff0019166101001790555b8451610ae89060c9906020880190611889565b508351610afc9060ca906020870190611889565b5060cb805460ff191660ff851617905560cf80546001600160a01b0319166001600160a01b038416179055610b2f611421565b610b37611450565b8015610b7d576000805461ff0019169055604051600181527f7f26b83ff96e1f2b6a682f133852f6798a09c465da95921460cefb38474024989060200160405180910390a15b5050505050565b610b8c611256565b6001600160a01b038116610bf15760405162461bcd60e51b815260206004820152602660248201527f4f776e61626c653a206e6577206f776e657220697320746865207a65726f206160448201526564647265737360d01b60648201526084016105eb565b61076c816112b0565b610c0333610623565b6040516001600160a01b0383169082156108fc029083906000818181858888f19350505050158015610c39573d6000803e3d6000fd5b506040518181526001600160a01b0383169033907f9b1bfa7fa9ee420a16e124f794c35ac9f90472acc99140eb2f6447c714cad8eb906020015b60405180910390a35050565b6001600160a01b038216610cd55760405162461bcd60e51b815260206004820152601860248201527f6d696e7420746f20746865207a65726f2061646472657373000000000000000060448201526064016105eb565b8060cc6000828254610ce79190611e9c565b90915550506001600160a01b038216600090815260cd602052604081208054839290610d14908490611e9c565b90915550506040518181526001600160a01b038316906000907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef90602001610c73565b6001600160a01b038316610dad5760405162461bcd60e51b815260206004820152601d60248201527f617070726f76652066726f6d20746865207a65726f206164647265737300000060448201526064016105eb565b6001600160a01b03928316600090815260ce602090815260408083209490951682529290925291902055565b6001600160a01b038316610e2f5760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f2061646472657373000060448201526064016105eb565b6001600160a01b038216610e855760405162461bcd60e51b815260206004820152601c60248201527f7472616e7366657220746f20746865207a65726f20616464726573730000000060448201526064016105eb565b6001600160a01b038316600090815260cd602052604090205481811015610eee5760405162461bcd60e51b815260206004820152601f60248201527f7472616e7366657220616d6f756e7420657863656564732062616c616e63650060448201526064016105eb565b610ef88282611eb4565b6001600160a01b03808616600090815260cd60205260408082209390935590851681529081208054849290610f2e908490611e9c565b92505081905550826001600160a01b0316846001600160a01b03167fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef84604051610f7a91815260200190565b60405180910390a350505050565b6001600160a01b038216610fde5760405162461bcd60e51b815260206004820152601a60248201527f6275726e2066726f6d20746865207a65726f206164647265737300000000000060448201526064016105eb565b6001600160a01b038216600090815260cd6020526040902054818110156110475760405162461bcd60e51b815260206004820152601b60248201527f6275726e20616d6f756e7420657863656564732062616c616e6365000000000060448201526064016105eb565b6110518282611eb4565b6001600160a01b038416600090815260cd602052604081209190915560cc805484929061107f908490611eb4565b90915550506040518281526000906001600160a01b038516907fddf252ad1be2c89b69c2b068fc378daa952ba7f163c4a11628f55a4df523b3ef9060200160405180910390a3505050565b61076c611256565b7f4910fdfa16fed3260ed0e7147f7cc6da11a60208b5b9406d12a635614ffd91435460ff161561110a5761110583611477565b611251565b826001600160a01b03166352d1902d6040518163ffffffff1660e01b815260040160206040518083038186803b15801561114357600080fd5b505afa925050508015611173575060408051601f3d908101601f1916820190925261117091810190611ad1565b60015b6111d65760405162461bcd60e51b815260206004820152602e60248201527f45524331393637557067726164653a206e657720696d706c656d656e7461746960448201526d6f6e206973206e6f74205555505360901b60648201526084016105eb565b600080516020611f7483398151915281146112455760405162461bcd60e51b815260206004820152602960248201527f45524331393637557067726164653a20756e737570706f727465642070726f786044820152681a58589b195555525160ba1b60648201526084016105eb565b50611251838383611513565b505050565b6097546001600160a01b031633146109175760405162461bcd60e51b815260206004820181905260248201527f4f776e61626c653a2063616c6c6572206973206e6f7420746865206f776e657260448201526064016105eb565b609780546001600160a01b038381166001600160a01b0319831681179093556040519116919082907f8be0079c531659141344cd1fd0a4f28419497f9722a3daafe3b4186f6b6457e090600090a35050565b6001600160a01b0385166113585760405162461bcd60e51b815260206004820152601e60248201527f7472616e736665722066726f6d20746865207a65726f2061646472657373000060448201526064016105eb565b600084511161139d5760405162461bcd60e51b81526020600482015260116024820152701a5b9d985b1a59081c9958da5c1a595b9d607a1b60448201526064016105eb565b806113db5760405162461bcd60e51b815260206004820152600e60248201526d1a5b9d985b1a59081d185c99d95d60921b60448201526064016105eb565b60cf546113fc9086906001600160a01b03166113f78587611e9c565b610dd9565b61141985858585856040518060200160405280600081525061153e565b505050505050565b600054610100900460ff166114485760405162461bcd60e51b81526004016105eb90611df8565b6109176115f6565b600054610100900460ff166109175760405162461bcd60e51b81526004016105eb90611df8565b6001600160a01b0381163b6114e45760405162461bcd60e51b815260206004820152602d60248201527f455243313936373a206e657720696d706c656d656e746174696f6e206973206e60448201526c1bdd08184818dbdb9d1c9858dd609a1b60648201526084016105eb565b600080516020611f7483398151915280546001600160a01b0319166001600160a01b0392909216919091179055565b61151c83611626565b6000825111806115295750805b15611251576115388383611666565b50505050565b600080806110046115538a8a8a8a8a8a611692565b6040516115609190611c7d565b6000604051808303816000865af19150503d806000811461159d576040519150601f19603f3d011682016040523d82523d6000602084013e6115a2565b606091505b50915091506115e0828260405180604001604052806016815260200175199a5c0b58dc9bdcdccb58da185a5b8819985a5b195960521b8152506116e5565b6115e98161175f565b9998505050505050505050565b600054610100900460ff1661161d5760405162461bcd60e51b81526004016105eb90611df8565b610917336112b0565b61162f81611477565b6040516001600160a01b038216907fbc7cd75a20ee27fd9adebab32041f755214dbc6bffa90cc0225b39da2e5c2d3b90600090a250565b606061168b8383604051806060016040528060278152602001611f9460279139611776565b9392505050565b60608686868686866040516024016116af96959493929190611cd6565b60408051601f198184030181529190526020810180516001600160e01b0316633c3e7d7760e01b17905290509695505050505050565b82611251576000828060200190518101906117009190611ae9565b9050600182511015611726578060405162461bcd60e51b81526004016105eb9190611d1e565b8181604051602001611739929190611c99565b60408051601f198184030181529082905262461bcd60e51b82526105eb91600401611d1e565b6000808280602001905181019061168b9190611ab1565b6060600080856001600160a01b0316856040516117939190611c7d565b600060405180830381855af49150503d80600081146117ce576040519150601f19603f3d011682016040523d82523d6000602084013e6117d3565b606091505b50915091506117e4868383876117ee565b9695505050505050565b6060831561185a578251611853576001600160a01b0385163b6118535760405162461bcd60e51b815260206004820152601d60248201527f416464726573733a2063616c6c20746f206e6f6e2d636f6e747261637400000060448201526064016105eb565b5081610a0d565b610a0d838381511561186f5781518083602001fd5b8060405162461bcd60e51b81526004016105eb9190611d1e565b82805461189590611ef7565b90600052602060002090601f0160209004810192826118b757600085556118fd565b82601f106118d057805160ff19168380011785556118fd565b828001600101855582156118fd579182015b828111156118fd5782518255916020019190600101906118e2565b5061190992915061190d565b5090565b5b80821115611909576000815560010161190e565b600061193561193084611e74565b611e43565b905082815283838301111561194957600080fd5b828260208301376000602084830101529392505050565b600082601f830112611970578081fd5b61168b83833560208501611922565b600060208284031215611990578081fd5b813561168b81611f5e565b600080604083850312156119ad578081fd5b82356119b881611f5e565b946020939093013593505050565b600080604083850312156119d8578182fd5b82356119e381611f5e565b915060208301356119f381611f5e565b809150509250929050565b600080600060608486031215611a12578081fd5b8335611a1d81611f5e565b92506020840135611a2d81611f5e565b929592945050506040919091013590565b60008060408385031215611a50578182fd5b8235611a5b81611f5e565b9150602083013567ffffffffffffffff811115611a76578182fd5b8301601f81018513611a86578182fd5b611a9585823560208401611922565b9150509250929050565b600080604083850312156119ad578182fd5b600060208284031215611ac2578081fd5b8151801515811461168b578182fd5b600060208284031215611ae2578081fd5b5051919050565b600060208284031215611afa578081fd5b815167ffffffffffffffff811115611b10578182fd5b8201601f81018413611b20578182fd5b8051611b2e61193082611e74565b818152856020838501011115611b42578384fd5b611b53826020830160208601611ecb565b95945050505050565b60008060008060808587031215611b71578081fd5b843567ffffffffffffffff80821115611b88578283fd5b611b9488838901611960565b95506020870135915080821115611ba9578283fd5b50611bb687828801611960565b935050604085013560ff81168114611bcc578182fd5b91506060850135611bdc81611f5e565b939692955090935050565b60008060008060808587031215611bfc578384fd5b843567ffffffffffffffff811115611c12578485fd5b611c1e87828801611960565b97602087013597506040870135966060013595509350505050565b600060208284031215611c4a578081fd5b5035919050565b60008151808452611c69816020860160208601611ecb565b601f01601f19169290920160200192915050565b60008251611c8f818460208701611ecb565b9190910192915050565b60008351611cab818460208801611ecb565b6101d160f51b9083019081528351611cca816002840160208801611ecb565b01600201949350505050565b6001600160a01b038716815260c060208201819052600090611cfa90830188611c51565b86604084015285606084015284608084015282810360a08401526115e98185611c51565b60006020825261168b6020830184611c51565b600060808252611d446080830187611c51565b6020830195909552506040810192909252606090910152919050565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b19195b1959d85d1958d85b1b60a21b606082015260800190565b6020808252602c908201527f46756e6374696f6e206d7573742062652063616c6c6564207468726f7567682060408201526b6163746976652070726f787960a01b606082015260800190565b6020808252602b908201527f496e697469616c697a61626c653a20636f6e7472616374206973206e6f74206960408201526a6e697469616c697a696e6760a81b606082015260800190565b604051601f8201601f1916810167ffffffffffffffff81118282101715611e6c57611e6c611f48565b604052919050565b600067ffffffffffffffff821115611e8e57611e8e611f48565b50601f01601f191660200190565b60008219821115611eaf57611eaf611f32565b500190565b600082821015611ec657611ec6611f32565b500390565b60005b83811015611ee6578181015183820152602001611ece565b838111156115385750506000910152565b600281046001821680611f0b57607f821691505b60208210811415611f2c57634e487b7160e01b600052602260045260246000fd5b50919050565b634e487b7160e01b600052601160045260246000fd5b634e487b7160e01b600052604160045260246000fd5b6001600160a01b038116811461076c57600080fdfe360894a13ba1a3210667c828492db98dca3e2076cc3735a920a3ca505d382bbc416464726573733a206c6f772d6c6576656c2064656c65676174652063616c6c206661696c6564a2646970667358221220bfbf873314a53bcd61631f3197ba5b4ea2d93847a8b61d091291a4ed060e644164736f6c63430008020033")
)

var (
	fip20Init = Contract{
		Address: common.HexToAddress(FIP20LogicAddress),
		ABI:     MustABIJson(FIP20UpgradableMetaData.ABI),
		Bin:     MustDecodeHex(FIP20UpgradableMetaData.Bin),
		Code:    initFIP20Code,
	}
	wfxInit = Contract{
		Address: common.HexToAddress(WFXLogicAddress),
		ABI:     MustABIJson(WFXUpgradableMetaData.ABI),
		Bin:     MustDecodeHex(WFXUpgradableMetaData.Bin),
		Code:    initWFXCode,
	}
	erc1967Proxy = Contract{
		Address: common.Address{},
		ABI:     MustABIJson(ERC1967ProxyMetaData.ABI),
		Bin:     MustDecodeHex(ERC1967ProxyMetaData.Bin),
		Code:    []byte{},
	}

	fxBridgeABI = MustABIJson(IFxBridgeLogicMetaData.ABI)
)

type Contract struct {
	Address common.Address
	ABI     abi.ABI
	Bin     []byte
	Code    []byte
}

func (c Contract) CodeHash() common.Hash {
	return crypto.Keccak256Hash(c.Code)
}

func GetFIP20() Contract {
	return fip20Init
}

func GetWFX() Contract {
	return wfxInit
}

func GetERC1967Proxy() Contract {
	return erc1967Proxy
}

func GetFxBridgeABI() abi.ABI {
	return fxBridgeABI
}

func MustDecodeHex(str string) []byte {
	bz, err := hexutil.Decode(str)
	if err != nil {
		panic(err)
	}
	return bz
}

func MustABIJson(str string) abi.ABI {
	j, err := abi.JSON(strings.NewReader(str))
	if err != nil {
		panic(err)
	}
	return j
}
